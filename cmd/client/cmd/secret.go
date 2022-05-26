package cmd

import (
	"context"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	pb "gophkeeper/api/proto"
	"gophkeeper/pkg/logger"
	"io/ioutil"
	"log"
	"os"
)

var (
	secretCmd = &cobra.Command{
		Use:   "secret",
		Short: "Secret management",
		Long:  `Choose one of the command to do with secrets`,
		Run: func(cmd *cobra.Command, args []string) {
			checkErr(cmd.Help())
		},
	}
	secretListCmd = &cobra.Command{
		Use:   "ls",
		Short: "List own secrets",
		Long:  `Allows you to list own secrets stored on server`,
		Run:   secretList,
	}
	secretCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create secret",
		Long:  `Allows you to create secret`,
		Run: func(cmd *cobra.Command, args []string) {
			checkErr(cmd.Help())
		},
	}
	secretCreateLoginPasswordCmd = &cobra.Command{
		Use:   "lp [login] [password]",
		Short: "Create login/password secret",
		Long:  `Allows you to create raw secret`,
		Args:  cobra.ExactArgs(2),
		Run:   createLoginPasswordSecret,
	}
	secretCreateCardCmd = &cobra.Command{
		Use:   "card [number] [expires] [cvv] [holder]",
		Short: "Create card secret",
		Long:  `Allows you to create raw secret`,
		Args:  cobra.ExactArgs(4),
		Run:   createCardSecret,
	}
	secretCreateRawCmd = &cobra.Command{
		Use:   "raw",
		Short: "Create raw secret",
		Long:  `Allows you to create raw secret`,
		Run:   createRawSecret,
	}
	secretRemoveCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove secret",
		Long:  `Allows you to remove secret`,
		Run:   removeSecret,
	}
)

func init() {
	rootCmd.AddCommand(secretCmd)

	secretCmd.AddCommand(secretListCmd)
	secretCmd.AddCommand(secretRemoveCmd)
	secretRemoveCmd.PersistentFlags().StringP("name", "n", "", "secret name")
	checkErr(secretRemoveCmd.MarkPersistentFlagRequired("name"))

	secretCmd.AddCommand(secretCreateCmd)

	secretCreateCmd.AddCommand(secretCreateRawCmd)
	secretCreateRawCmd.Flags().StringP("name", "n", "", "secret name")
	checkErr(secretCreateRawCmd.MarkFlagRequired("name"))
	secretCreateRawCmd.Flags().String("from-file", "-f", "take secret content from this file")
	secretCreateCmd.AddCommand(secretCreateLoginPasswordCmd)
	secretCreateLoginPasswordCmd.Flags().StringP("name", "n", "", "secret name")
	checkErr(secretCreateLoginPasswordCmd.MarkFlagRequired("name"))
	secretCreateCmd.AddCommand(secretCreateCardCmd)
	secretCreateCardCmd.Flags().StringP("name", "n", "", "secret name")
	checkErr(secretCreateCardCmd.MarkFlagRequired("name"))
}

func removeSecret(cmd *cobra.Command, args []string) {
	var err error

	cl, stop := getKeeperClient()
	defer stop()

	ctx := context.Background()

	name, err := cmd.Flags().GetString("name")
	checkErr(err)

	_, err = cl.DeleteSecret(ctx, &pb.DeleteSecretRequest{
		Name: name,
	})
	switch status.Code(err) {
	case codes.NotFound:
		l.Fatal().Msg("Secret not found")
	case codes.OK:
		l.Info().Msg("Secret removed successfully")
	}
}

func createGenericSecret(t, n string, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		l.Fatal().Err(err).Send()
	}

	if len(data) == 0 {
		l.Fatal().Msg("Unable to create empty secret")
	}

	cl, stop := getKeeperClient()
	defer stop()

	ctx := context.Background()

	if n == "" {
		l.Fatal().Msg("Please specify secret name")
	}

	_, err = cl.CreateSecret(ctx, &pb.CreateSecretRequest{
		Type:    t,
		Name:    n,
		Content: data,
	})
	switch status.Code(err) {
	case codes.AlreadyExists:
		l.Fatal().Msg("Secret already exists")
	case codes.OK:
		l.Info().Msg("Secret created successfully")
	}
}

func createRawSecret(cmd *cobra.Command, args []string) {
	var file *os.File
	var err error

	fromFile, err := cmd.Flags().GetString("from-file")
	checkErr(err)

	if fromFile != "" {
		file, err = os.Open(fromFile)
		if err != nil {
			checkErr(err)
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	} else {
		file = os.Stdin
	}

	data, err := ioutil.ReadAll(file)
	checkErr(err)

	name, err := cmd.Flags().GetString("name")
	checkErr(err)

	createGenericSecret("raw", name, data)
}

func createLoginPasswordSecret(cmd *cobra.Command, args []string) {
	in := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		Login:    args[0],
		Password: args[1],
	}

	name, err := cmd.Flags().GetString("name")
	checkErr(err)

	createGenericSecret("lp", name, in)
}

func createCardSecret(cmd *cobra.Command, args []string) {
	in := struct {
		Number  string `json:"number"`
		Expires string `json:"expires"`
		CVV     string `json:"cvv"`
		Holder  string `json:"holder"`
	}{
		Number:  args[0],
		Expires: args[1],
		CVV:     args[2],
		Holder:  args[3],
	}

	name, err := cmd.Flags().GetString("name")
	checkErr(err)

	createGenericSecret("card", name, in)
}

func secretList(cmd *cobra.Command, args []string) {
	l := logger.Global()
	ctx := context.Background()

	cl, stop := getKeeperClient()
	defer stop()

	resp, err := cl.ListSecrets(ctx, &pb.ListSecretsRequest{})
	if err != nil {
		l.Fatal().Err(err).Msg("Server error")
	}

	log.Println(resp.GetSecrets())
}

func getKeeperClient() (pb.KeeperClient, func()) {
	// real client for mocked service
	conn, err := grpc.Dial(
		viper.GetString("server_addr"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(clientAuthInterceptor),
	)
	checkErr(err)

	stop := func() {
		_ = conn.Close()
	}

	cl := pb.NewKeeperClient(conn)

	return cl, stop
}

func clientAuthInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+authViper.GetString("token"))
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
