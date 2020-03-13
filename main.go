package main

import (
	"errors"
	"fmt"
	"github.com/cerence/Ark-cli/cmd"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"os"
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
					host := c.String("host")
					cmd.Info.Println(fmt.Sprintf("Ark cloud server is %s", host))

					args := c.Args()
					if args.Len() == 0 {
						return errors.New("no unique device code")
					}
					uniqueDeviceCode := args.First()
					cmd.Info.Println(fmt.Sprintf("Unique device code is %s", uniqueDeviceCode))

					return cmd.ImportNewDevice(host, uniqueDeviceCode)
				},
			},
			{
				Name:    "unbind",
				Aliases: []string{"u"},
				Usage:   "unbind device",
				Action: func(c *cli.Context) error {
					host := c.String("host")
					cmd.Info.Println(fmt.Sprintf("Ark cloud server is %s", host))

					args := c.Args()
					if args.Len() == 0 {
						return errors.New("no unique device code")
					}
					uniqueDeviceCode := args.First()
					cmd.Info.Println(fmt.Sprintf("Unique device code is %s", uniqueDeviceCode))

					return cmd.UnbindDevice(host, uniqueDeviceCode)
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				DefaultText: DefaultServer,
				Value:       DefaultServer,
				Usage:       "Ark cloud server",
				Required:    false,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
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
