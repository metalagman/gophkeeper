package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "gophkeeper/api/proto"
	"gophkeeper/pkg/logger"
	"os"
	"path/filepath"
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

	// real client for mocked service
	conn, err := grpc.Dial(cfg.Server.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	logger.CheckErr(err)

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	ctx := context.Background()

	cl := pb.NewAuthClient(conn)
	resp, err := cl.Register(ctx, &pb.RegisterRequest{
		Email:    email,
		Password: password,
	})

	logger.CheckErr(err)
	logger.Global().Info().Msgf("token: %s", resp.GetToken())
}

func login(cmd *cobra.Command, args []string) {
	email, password := args[0], args[1]
	fmt.Println(email, password)
}

func forgetAuth(cmd *cobra.Command, args []string) {
	cfgDir := appDir("gophkeeper")

	authCfg := viper.New()
	authCfg.SetConfigType("toml")
	var defaultConfig = []byte(`
[auth]
token=""
`)
	logger.CheckErr(authCfg.ReadConfig(bytes.NewBuffer(defaultConfig)))
	logger.CheckErr(authCfg.WriteConfigAs(filepath.Join(cfgDir, "auth.toml")))
}

func ensureDir(dirName string) error {
	if _, err := os.Stat(dirName); err == nil {
		return nil
	}
	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	return nil
}

func appDir(app string) string {
	cfgDir, err := os.UserConfigDir()
	logger.CheckErr(err)
	cfgDir = filepath.Join(cfgDir, app)
	logger.CheckErr(ensureDir(cfgDir))
	return cfgDir
}
