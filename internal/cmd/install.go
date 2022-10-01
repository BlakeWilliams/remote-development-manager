package cmd

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/blakewilliams/remote-development-manager/internal/service"
	"github.com/spf13/cobra"
)

const rdmServiceName = "com.blakewilliams.rdm"

func newInstallCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Configures rdm to run on boot as a MacOS LaunchAgent.",
		Run: func(cmd *cobra.Command, args []string) {
			svc := service.New(rdmServiceName)
			if err := ensureInstalled(svc); err != nil {
				die("ensureInstalled", err)
			}
			status, err := svc.RunState()
			if err != nil {
				die("RunState", err)
			}
			fmt.Printf(
				"Run state for %s: %s\nRun `launchctl print %s` for more detail\n.",
				svc.Name, status.String(), svc.UserSpecifier(),
			)
		},
	}
}

func die(msg string, err error) {
	log.Fatalf("error in %s, %v", msg, err)
}

func ensureInstalled(svc *service.Service) error {
	installState, err := svc.InstallState()
	if err != nil {
		return err
	}
	fmt.Printf("Install state for %+v: %s\n", svc.Name, installState.String())
	if installState != service.Installed {
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
