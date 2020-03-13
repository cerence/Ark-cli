package main

import (
	"errors"
	"fmt"
	"github.com/cerence/Ark-cli/cmd"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"os"
)

const (
	BackEndHost = "backend-host"
	APIHost     = "api-host"
)

func main() {
	displayLogo()

	app := &cli.App{
		Name:    Name,
		Usage:   Usage,
		Version: Version,
		Commands: []*cli.Command{
			{
				Name:    "import",
				Aliases: []string{"i"},
				Usage:   "Import new device",
				Action: func(c *cli.Context) error {
					backendHost := extractBackEndHost(c)
					apiHost := extractAPIHost(c)
					uniqueDeviceCode, err := extractFirstArg(c)
					if err != nil {
						return err
					}

					cmd.Info.Println(fmt.Sprintf("Unique device code is %s", uniqueDeviceCode))

					return cmd.ImportNewDevice(apiHost, backendHost, uniqueDeviceCode)
				},
			},
			{
				Name:    "unbind",
				Aliases: []string{"u"},
				Usage:   "unbind device",
				Action: func(c *cli.Context) error {
					backendHost := extractBackEndHost(c)
					uniqueDeviceCode, err := extractFirstArg(c)
					if err != nil {
						return err
					}

					cmd.Info.Println(fmt.Sprintf("Unique device code is %s", uniqueDeviceCode))

					return cmd.UnbindDevice(backendHost, uniqueDeviceCode)
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        BackEndHost,
				DefaultText: DefaultBackendHost,
				Value:       DefaultBackendHost,
				Usage:       "Ark cloud backend host",
				Required:    false,
			},
			&cli.StringFlag{
				Name:        APIHost,
				DefaultText: DefaultAPIHost,
				Value:       DefaultAPIHost,
				Usage:       "Ark cloud api host",
				Required:    false,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		cmd.Error.Fatal(err)
	}
}

func extractBackEndHost(c *cli.Context) string {
	host := c.String(BackEndHost)
	cmd.Info.Println(fmt.Sprintf("Ark cloud backend server is %s", host))
	return host
}

func extractAPIHost(c *cli.Context) string {
	host := c.String(APIHost)
	cmd.Info.Println(fmt.Sprintf("Ark cloud api server is %s", host))
	return host
}

func extractFirstArg(c *cli.Context) (string, error) {
	args := c.Args()
	if args.Len() == 0 {
		return "", errors.New("no args")
	}
	arg := args.First()
	return arg, nil
}

func displayLogo() {
	if len(os.Args) == 1 {
		color.Cyan(logo, Version)
	} else if len(os.Args) == 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "h" || os.Args[1] == "help" {
			color.Cyan(logo, Version)
		}
	}
}
