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

// NewInstallCmd returns a new "service install" subcommand
func NewInstallCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Configures rdm to run on boot as a MacOS LaunchAgent.",
		Run: func(cmd *cobra.Command, args []string) {
			svc := NewRdmService()
			if err := ensureInstalled(svc); err != nil {
				die("ensureInstalled", err)
			}
			status := svc.RunState()
			fmt.Printf(
				"Run state for %s: %s\nRun `launchctl print %s` for more detail\n.",
				svc.Name, status.Pretty(), svc.UserSpecifier(),
			)
		},
	}
	return cmd
}

func NewRdmService() *service.Service {
	return service.New(rdmServiceName)
}

func die(msg string, err error) {
	log.Fatalf("error in %s, %v", msg, err)
}

func ensureInstalled(svc *service.Service) error {
	installState := svc.InstallState()
	fmt.Println(installState.Pretty())

	if !installState.Is(state.Installed) {
		configFile, err := plistContent()
		if err != nil {
			return err
		}

		fmt.Print("Attempting to set up launchd service... ")
		err = svc.Install(configFile)
		if err != nil {
			fmt.Printf("failed: %v\n", err)
		} else {
			fmt.Println("done!")
		}
		return err
	}
	return nil
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
