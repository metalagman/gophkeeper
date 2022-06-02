package cmd

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gophkeeper/pkg/logger"
	"gophkeeper/pkg/userconfig"
	"gophkeeper/pkg/version"
	"io/fs"
	"os"
	"time"
)

const (
	appName = "gkcli"
)

var (
	authViper *viper.Viper
	l         *logger.Logger
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "gophkeeper client",
	Long:  `gophkeeper client`,
	Run: func(cmd *cobra.Command, args []string) {
		checkErr(cmd.Help())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Version: version.Info(),
}

func Execute() {
	checkErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initDotEnv)
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initAuth)

	//rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "set high log verbosity")
	rootCmd.PersistentFlags().StringP("server", "s", "localhost:50051", "remote server address and port")
}

func initDotEnv() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		checkErr(fmt.Errorf(".env load: %w", err))
	}
}

func initConfig() {
	viper.SetDefault("server_addr", "localhost:50051")
	viper.SetDefault("log_verbose", 0)

	checkErr(viper.BindPFlag("log_verbose", rootCmd.PersistentFlags().Lookup("verbose")))
	checkErr(viper.BindPFlag("server_addr", rootCmd.PersistentFlags().Lookup("server")))
}

func initAuth() {
	uc, err := userconfig.New(appName, "toml")
	checkErr(err)
	authViper = uc.Viper("auth")

	l.Debug().
		Str("email", authViper.GetString("email")).
		Msg("Auth")
}

func initLogger() {
	logger.NewGlobal(logger.Config{
		Pretty:     true,
		Verbose:    viper.GetBool("log_verbose"),
		TimeFormat: time.Kitchen,
	})
	l = logger.Global()
}

func checkErr(err error) {
	if err == nil {
		return
	}
	logger.Global().Fatal().Msg(err.Error())
}
