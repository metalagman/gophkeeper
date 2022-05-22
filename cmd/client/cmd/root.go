package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gophkeeper/internal/client/config"
	"gophkeeper/pkg/logger"
	"gophkeeper/pkg/userconfig"
	"gophkeeper/pkg/version"
	"io/fs"
	"os"
	"strings"
)

const (
	appName = "gkcli"
)

var (
	cfg   = config.Config{}
	vAuth *viper.Viper
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "gophkeeper client",
	Long:  `gophkeeper client`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.CheckErr(cmd.Help())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Version: version.Info(),
}

func Execute() {
	logger.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initDotEnv)
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initAuth)

	//rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "set high log verbosity")
}

func initDotEnv() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		logger.CheckErr(fmt.Errorf(".env load: %w", err))
	}
}

func initConfig() {
	viper.SetConfigType("toml")

	var defaultConfig = []byte(`
[log]
verbose=0
pretty=1
`)
	logger.CheckErr(viper.ReadConfig(bytes.NewBuffer(defaultConfig)))

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	logger.CheckErr(viper.BindPFlag("log.verbose", rootCmd.PersistentFlags().Lookup("verbose")))

	logger.CheckErr(viper.Unmarshal(&cfg))
}

func initAuth() {
	uc, err := userconfig.New(appName, "toml")
	logger.CheckErr(err)
	vAuth = uc.Viper("auth")

	l := logger.Global()
	l.Debug().
		Str("email", viper.GetString("email")).
		Str("token", viper.GetString("token")).
		Msg("Auth")
}

func initLogger() {
	logger.NewGlobal(cfg.Logger)
}
