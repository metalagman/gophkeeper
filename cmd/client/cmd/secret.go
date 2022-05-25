package cmd

import (
	"context"
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

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Secret management",
	Long:  `Choose one of the command to do with secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.CheckErr(cmd.Help())
	},
}

var secretListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List own secrets",
	Long:  `Allows you to list own secrets stored on server`,
	Run:   secretList,
}

var secretCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create secret",
	Long:  `Allows you to create secret`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.CheckErr(cmd.Help())
	},
}

var secretCreateRawCmd = &cobra.Command{
	Use:   "raw",
	Short: "Create raw secret",
	Long:  `Allows you to create raw secret`,
	Run:   createRawSecret,
}

func createRawSecret(cmd *cobra.Command, args []string) {
	var file *os.File
	var err error

	if fromFile := viper.GetString("fromFile"); fromFile != "" {
		file, err = os.Open(fromFile)
		if err != nil {
			l.CheckErr(err)
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	} else {
		file = os.Stdin
	}

	data, err := ioutil.ReadAll(file)
	l.CheckErr(err)

	if len(data) == 0 {
		l.Fatal().Msg("Unable to create empty secret")
	}

	cl, stop := getKeeperClient()
	defer stop()

	ctx := context.Background()

	name, err := cmd.Flags().GetString("name")
	l.CheckErr(err)

	if name == "" {
		l.Fatal().Msg("Please specify secret name")
	}

	_, err = cl.CreateSecret(ctx, &pb.CreateSecretRequest{
		Name:    name,
		Content: data,
	})
	switch status.Code(err) {
	case codes.AlreadyExists:
		l.Fatal().Msg("Secret already exists")
	case codes.OK:
		l.Info().Msg("Secret created successfully")
	}
}

func init() {
	rootCmd.AddCommand(secretCmd)

	secretCmd.AddCommand(secretListCmd)
	secretCmd.AddCommand(secretCreateCmd)

	secretCreateCmd.AddCommand(secretCreateRawCmd)
	secretCreateCmd.PersistentFlags().StringP("name", "n", "", "secret name")

	secretCreateRawCmd.Flags().String("from-file", "", "path to file")
	logger.CheckErr(viper.BindPFlag("fromFile", secretCreateRawCmd.Flags().Lookup("from-file")))
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
	logger.CheckErr(err)

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
