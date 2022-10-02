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
	"github.com/fatih/color"
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
	return cmd
}

func NewRdmService() *service.Service {
	return service.New(rdmServiceName)
}

func die(msg string, err error) {
	log.Fatalf("error in %s, %v", msg, err)
}

func VerboseInstallState(svc *service.Service) string {
	errDetail := ""
	installState, err := svc.InstallState()
	if err != nil {
		errDetail = fmt.Sprintf(" (%v)", err.Error())
	}

	state := installState.String()
	if state == "Installed" {
		state = green(state)
	}

	return fmt.Sprintf("Install state: [%s]%s", state, errDetail)
}

var green = color.New(color.FgGreen).SprintFunc()

func VerboseRunState(svc *service.Service) string {
	errDetail := ""
	runState, err := svc.RunState()
	if err != nil {
		errDetail = fmt.Sprintf(" (%v)", err.Error())
	}

	state := runState.String()
	if state == "Running" {
		state = green(state)
	}

	return fmt.Sprintf("Run state:     [%s]%s", state, errDetail)
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
