package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"gophkeeper/internal/server/app"
	"gophkeeper/pkg/logger"
	"os/signal"
	"syscall"
)

// migrateCmd represents the migrate command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run server",
	Long:  `Run the GophKeeker server`,
	Run:   serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM)
	defer stop()

	a, err := app.New(cfg)
	logger.CheckErr(err)

	<-ctx.Done()

	a.Stop()
}
