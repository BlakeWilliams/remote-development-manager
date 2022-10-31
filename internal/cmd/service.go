package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/brasic/launchd"
	"github.com/brasic/launchd/state"
	"github.com/spf13/cobra"
)

func newServiceCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	short := "Manage this program as a launchd system service."
	svc := rdmService()
	cmd := &cobra.Command{
		Use:   "service [subcommand]",
		Short: short,
		Long:  short + "\n" + prettyStatus(svc),
	}
	cmd.AddCommand(serviceInstallCmd(ctx, logger))
	cmd.AddCommand(serviceUninstallCmd(ctx, logger))
	cmd.AddCommand(serviceStartCmd(ctx, logger))
	cmd.AddCommand(serviceStopCmd(ctx, logger))
	return cmd
}

func serviceInstallCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Configures rdm to run on boot as a MacOS LaunchAgent.",
		Run: func(cmd *cobra.Command, args []string) {
			svc := rdmService()

			if svc.IsHealthy() {
				fmt.Println("service is already installed and running, nothing to do!")
				return
			}
			// Configure launchagent to run `rdm server` at boot
			if err := svc.Install(); err != nil {
				fmt.Printf("Problem installing: %v\n", err)
				return
			}

			// In case the service is not running, start it.
			if err := svc.Start(); err != nil {
				fmt.Printf("Problem starting: %v\n", err)
				return
			}

			fmt.Printf("Configured to start at boot. Uninstall using:\n\t%s service uninstall\n", currentExecutableName())
		},
	}
	return cmd
}

func serviceUninstallCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Removes a previously installed LaunchAgent installation.",
		Run: func(cmd *cobra.Command, args []string) {
			svc := rdmService()

			if !svc.InstallState().Is(state.Installed) {
				fmt.Println("Service is not installed.")
				return
			}
			if err := svc.Bootout(true); err != nil {
				fmt.Printf("Problem uninstalling: %v\n", err)
				return
			}
			fmt.Println("Service uninstalled.")
		},
	}
}

func serviceStartCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Starts the launchd service.",
		Run: func(cmd *cobra.Command, args []string) {
			svc := rdmService()

			if svc.RunState().Is(state.Running) {
				fmt.Println("Service is already running.")
				return
			}
			if err := svc.Start(); err != nil {
				fmt.Printf("Problem starting: %v\n", err)
				return
			}
			finalState, timedOut := svc.PollUntil(state.Running, 5*time.Second)
			if timedOut {
				fmt.Println("Service failed to start. Currently:", finalState.Pretty())
				return
			}
			fmt.Println("Service started.")
		},
	}
}

func serviceStopCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stops the launchd service.",
		Run: func(cmd *cobra.Command, args []string) {
			svc := rdmService()

			runState := svc.RunState()

			if !runState.Is(state.Running) {
				fmt.Println("Service is not running, nothing to do.")
				return
			}
			if err := svc.Stop(); err != nil {
				fmt.Printf("Problem stopping: %v\n", err)
				return
			}
			finalState, timedOut := svc.PollUntil(state.NotRunning, 5*time.Second)
			if timedOut {
				fmt.Println("Service failed to stop. Currently:", finalState.Pretty())
				return
			}
			fmt.Println("Service stopped.")
		},
	}
}

// NewRdmService returns a launchd.Service for the currently running program.
func rdmService() *launchd.Service {
	return launchd.ForRunningProgram("com.blakewilliams.rdm", []string{"server"})
}

func currentExecutableName() string {
	return filepath.Base(os.Args[0])
}

func prettyStatus(svc *launchd.Service) string {
	return fmt.Sprintf("  Status of %s:\n    %s\n    %s",
		svc.UserSpecifier(),
		svc.InstallState().Pretty(),
		svc.RunState().Pretty(),
	)
}
