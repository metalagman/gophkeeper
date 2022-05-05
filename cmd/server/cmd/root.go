package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gophkeeper/internal/app/client"
	"gophkeeper/pkg/logger"
	"gophkeeper/pkg/version"
	"io/fs"
	"os"
	"strings"
)

var cfg = client.Config{}

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "gophkeeper server",
	Long:  `gophkeeper server`,
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

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Set high log verbosity")
	rootCmd.PersistentFlags().BoolP("pretty", "p", false, "Set pretty log formatting (instead of json)")
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
pretty=0
`)
	logger.CheckErr(viper.ReadConfig(bytes.NewBuffer(defaultConfig)))

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	logger.CheckErr(viper.BindPFlag("log.verbose", rootCmd.PersistentFlags().Lookup("verbose")))
	logger.CheckErr(viper.BindPFlag("log.pretty", rootCmd.PersistentFlags().Lookup("pretty")))

	logger.CheckErr(viper.Unmarshal(&cfg))
}

func initLogger() {
	logger.NewGlobal(cfg.Logger)
}
