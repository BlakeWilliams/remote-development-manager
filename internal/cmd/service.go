package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/blakewilliams/remote-development-manager/internal/cmd/service"
	"github.com/spf13/cobra"
)

func newServiceCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	short := "Manage this program as a launchd system service."
	svc := service.NewRdmService()
	long := fmt.Sprintf("%s\n  Status of %s:\n    %s\n    %s",
		short,
		svc.UserSpecifier(),
		svc.InstallState().Pretty(),
		svc.RunState().Pretty(),
	)
	cmd := &cobra.Command{
		Use:   "service [subcommand]",
		Short: "Manage this program as a launchd system service",
		Long:  long,
	}
	cmd.AddCommand(service.NewInstallCmd(ctx, logger))
	return cmd
}
