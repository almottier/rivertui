package main

import (
	"fmt"
	"os"
	"time"

	"github.com/almottier/rivertui/config"
	"github.com/almottier/rivertui/internal/client"
	"github.com/almottier/rivertui/monitor"
	"github.com/spf13/cobra"
)

var (
	dbURL           string
	refreshInterval time.Duration
	jobID           int64
	appConfig       *config.Config
	appClient       *client.Client

	rootCmd = &cobra.Command{
		Use:   "rivertui",
		Short: "rivertui is a terminal-based user interface for River Queue",
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if appClient != nil {
				appClient.Close()
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			var err error
			appConfig, err = config.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			config.UpdateConfigFromFlags(appConfig, dbURL, refreshInterval)

			if appConfig.Database.URL == "" {
				return fmt.Errorf("database URL is required. Set it via --database-url flag or RIVER_DATABASE_URL environment variable")
			}

			appClient, err = client.New(cmd.Context(), appConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize client: %w", err)
			}

			monitor := monitor.NewMonitorApp(appClient, appConfig, jobID)
			monitor.StartRefreshLoop()

			if err := monitor.Run(); err != nil {
				return fmt.Errorf("failed to run monitor: %w", err)
			}

			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&dbURL, "database-url", "", "PostgreSQL connection string/URL (env: RIVER_DATABASE_URL)")
	rootCmd.PersistentFlags().DurationVar(&refreshInterval, "refresh", 1*time.Second, "Refresh interval for the monitor")
	rootCmd.PersistentFlags().Int64Var(&jobID, "job-id", 0, "Job ID to view details for (starts in details view if provided)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
