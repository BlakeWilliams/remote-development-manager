package service

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/blakewilliams/remote-development-manager/internal/service"
	"github.com/blakewilliams/remote-development-manager/internal/service/state"
	"github.com/spf13/cobra"
)

const rdmServiceName = "com.blakewilliams.rdm"
const detailMessage = "Run `launchctl print %s` for more detail."
const noNeedMessage = "launchd service is already installed and running, nothing to do!" + "\n" + detailMessage + "\n"

// NewInstallCmd returns a new "service install" subcommand that installs a launchctl service.
func NewInstallCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Configures rdm to run on boot as a MacOS LaunchAgent.",
		Run: func(cmd *cobra.Command, args []string) {
			svc := NewRdmService()
			if svc.IsHealthy() {
				fmt.Printf(noNeedMessage, svc.UserSpecifier())
				return
			}
			if err := ensureInstalled(svc); err != nil {
				fmt.Printf("Problem installing: %v\n", err)
				return
			}
			status := svc.RunState()
			fmt.Printf(
				"%s\nRun `launchctl print %s` for more detail\n.",
				status.Pretty(), svc.UserSpecifier(),
			)
		},
	}
	return cmd
}

func NewRdmService() *service.Service {
	return service.New(rdmServiceName)
}

func ensureInstalled(svc *service.Service) (err error) {
	if !svc.InstallState().Is(state.Installed) {
		configFile, err := plistContent()
		if err != nil {
			return err
		}

		fmt.Print("Attempting to set up launchd service... ")
		if err = svc.Install(configFile); err == nil {
			fmt.Println("done!")
		}
	}
	return err
}

//go:embed launchagent.tmpl
var plistTemplateSrc string
var plistTemplate = template.Must(
	template.New("launchctl").Parse(plistTemplateSrc),
)

// Execute the embedded template to produce a rendered file suitable for
// use in launchd
func plistContent() ([]byte, error) {
	exe, err := executablePath()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = plistTemplate.Execute(&buf, struct{ Executable string }{exe})
	return buf.Bytes(), err
}

// The full path of the executable that is currently running.
func executablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return ex, nil
}
