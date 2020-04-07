package command

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

var (
	GET *cli.Command = &cli.Command{
		Name:  "get",
		Usage: "send get request",
		Action: func(c *cli.Context) error {
			fmt.Println("GET")
			return nil
		},
	}
)
