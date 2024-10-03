package main

import (
	"context"
	"log"

	"github.com/go-logr/stdr"
	"github.com/llmariner/api-usage/server/internal/config"
	"github.com/llmariner/api-usage/server/internal/server"
	"github.com/llmariner/api-usage/server/internal/store"
	"github.com/llmariner/common/pkg/db"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func runCmd() *cobra.Command {
	var path string
	var logLevel int
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Parse(path)
			if err != nil {
				return err
			}
			if err := c.Validate(); err != nil {
				return err
			}
			stdr.SetVerbosity(logLevel)

			if err := run(cmd.Context(), c); err != nil {
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

func run(ctx context.Context, c *config.Config) error {
	logger := stdr.New(log.Default())
	log := logger.WithName("boot")

	log.Info("Setting up the database...")
	dbInst, err := db.OpenDB(c.Database)
	if err != nil {
		return err
	}
	st := store.New(dbInst)
	if err := st.AutoMigrate(); err != nil {
		return err
	}

	log.Info("Setting up the server...")
	asrv := server.NewAdmin(st, logger)
	isrv := server.NewInternal(st, logger)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { return asrv.Run(ctx, c.AdminGRPCPort) })
	eg.Go(func() error { return isrv.Run(ctx, c.InternalGRPCPort) })
	return eg.Wait()
}
