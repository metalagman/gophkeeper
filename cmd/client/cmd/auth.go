package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "gophkeeper/api/proto"
	"gophkeeper/pkg/logger"
	"os"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authorization",
	Long:  `Choose one of the command to do with authorization`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.CheckErr(cmd.Help())
	},
}

var authRegisterCmd = &cobra.Command{
	Use:   "register [email] [password]",
	Short: "Register on the remote server",
	Long:  `Allows you to register on the remote server`,
	Args:  cobra.ExactArgs(2),
	Run:   register,
}

var authLoginCmd = &cobra.Command{
	Use:   "login [email] [password]",
	Short: "Login on the remote server",
	Long:  `Allows you to login on the remote server`,
	Args:  cobra.ExactArgs(2),
	Run:   login,
}

var authForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget current authorization",
	Long:  `Allows you to forget current authorization`,
	Run:   forgetAuth,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authRegisterCmd)
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authForgetCmd)
}

func register(cmd *cobra.Command, args []string) {
	email, password := args[0], args[1]
	ctx := context.Background()

	cl, stop := getUserClient()
	defer stop()

	resp, err := cl.Register(ctx, &pb.RegisterRequest{
		Email:    email,
		Password: password,
	})

	logger.CheckErr(err)
	vAuth.Set("email", email)
	vAuth.Set("token", resp.GetToken())
	logger.CheckErr(vAuth.WriteConfig())
}

func login(cmd *cobra.Command, args []string) {
	email, password := args[0], args[1]
	ctx := context.Background()

	cl, stop := getUserClient()
	defer stop()

	resp, err := cl.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})

	logger.CheckErr(err)
	vAuth.Set("email", email)
	vAuth.Set("token", resp.GetToken())
	logger.CheckErr(vAuth.WriteConfig())
}

func forgetAuth(cmd *cobra.Command, args []string) {
	l := logger.Global()

	if vAuth.GetString("token") == "" {
		l.Warn().Msg("Auth is already empty")
		os.Exit(0)
	}

	vAuth.Set("token", "")
	if err := vAuth.WriteConfig(); err != nil {
		l.Fatal().Err(err)
	}

	l.Info().Msg("Done")
}

func getUserClient() (pb.UserClient, func()) {
	// real client for mocked service
	conn, err := grpc.Dial(
		cfg.Server.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	logger.CheckErr(err)

	stop := func() {
		_ = conn.Close()
	}

	cl := pb.NewUserClient(conn)

	return cl, stop
}
