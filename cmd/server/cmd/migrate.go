package cmd

import (
	"database/sql"
	"fmt"
	"gophkeeper/internal/server/migrate"
	"gophkeeper/pkg/logger"

	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migration tool",
	Long:  `Choose one of the command to do with database migrations`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.CheckErr(cmd.Help())
	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply migrations up",
	Long:  `Allows you to apply all the missing migrations`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := getDb()
		logger.CheckErr(err)

		err = migrate.Up(db)
		logger.CheckErr(err)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Apply migrations down",
	Long:  `Allows you to revert one last migration`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := getDb()
		logger.CheckErr(err)

		err = migrate.Down(db)
		logger.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
}

func getDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DB.DSN)
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return db, nil
}
