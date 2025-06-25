package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-logr/stdr"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/llmariner/api-usage/api/v1"
	"github.com/llmariner/api-usage/pkg/store"
	"github.com/llmariner/api-usage/server/internal/cache"
	"github.com/llmariner/api-usage/server/internal/config"
	"github.com/llmariner/api-usage/server/internal/server"
	"github.com/llmariner/common/pkg/db"
	"github.com/llmariner/rbac-manager/pkg/auth"
	uv1 "github.com/llmariner/user-manager/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/encoding/protojson"
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

	conn, err := grpc.NewClient(
		c.CacheConfig.UserManagerServerInternalAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	uClient := uv1.NewUsersInternalServiceClient(conn)

	cstore := cache.NewStore(uClient)

	log.Info("Setting up the server...")
	asrv := server.NewAdmin(st, logger)
	isrv := server.NewInternal(st, logger)
	srv := server.New(st, cstore, logger)

	addr := fmt.Sprintf("localhost:%d", c.GRPCPort)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err = grpc.NewClient(addr, opts...)
	if err != nil {
		return err
	}
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:     true,
				EmitDefaultValues: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithIncomingHeaderMatcher(auth.HeaderMatcher),
		runtime.WithHealthzEndpoint(grpc_health_v1.NewHealthClient(conn)),
	)
	if err := v1.RegisterAPIUsageServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
		return err
	}

	errCh := make(chan error)

	go func() {
		log.Info("Starting HTTP server...", "port", c.HTTPPort)
		errCh <- http.ListenAndServe(fmt.Sprintf(":%d", c.HTTPPort), mux)
	}()

	go func() {
		errCh <- srv.Run(ctx, c.GRPCPort, c.AuthConfig)
	}()

	go func() {
		errCh <- asrv.Run(ctx, c.AdminGRPCPort)
	}()

	go func() {
		errCh <- isrv.Run(ctx, c.InternalGRPCPort)
	}()

	go func() {
		errCh <- cstore.Sync(ctx, c.CacheConfig.SyncInterval)
	}()

	return <-errCh
}
