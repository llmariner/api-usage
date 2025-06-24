package main

import (
	"context"
	"log"

	"github.com/go-logr/stdr"
	"github.com/llmariner/api-usage/server/internal/cleaner"
	"github.com/llmariner/api-usage/server/internal/config"
	"github.com/llmariner/api-usage/server/internal/store"
	"github.com/llmariner/common/pkg/db"
	"github.com/spf13/cobra"
)

func cleanerCmd() *cobra.Command {
	var path string
	var logLevel int
	cmd := &cobra.Command{
		Use:   "cleaner",
		Short: "Run the cleaner to delete old records",
		Long:  "Run the cleaner as a standalone process to delete records outside the retention period",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.ParseCleaner(path)
			if err != nil {
				return err
			}
			if err := c.Validate(); err != nil {
				return err
			}
			stdr.SetVerbosity(logLevel)

			if err := runCleaner(cmd.Context(), c); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&path, "config", "", "Path to the config file")
	cmd.Flags().IntVar(&logLevel, "v", 0, "Log level")
	_ = cmd.MarkFlagRequired("config")
	return cmd
}

func runCleaner(ctx context.Context, c *config.CleanerConfig) error {
	logger := stdr.New(log.Default())
	log := logger.WithName("api-usage-cleaner")

	log.Info("Setting up the database...")
	dbInst, err := db.OpenDB(c.Database)
	if err != nil {
		return err
	}
	st := store.New(dbInst)

	log.Info("Setting up cleaner...")
	cleaner := cleaner.NewCleaner(st, c.RetentionPeriod, c.PollInterval, logger)

	log.Info("Running cleaner...")
	return cleaner.Run(ctx)
}
