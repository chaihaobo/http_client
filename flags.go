package main

import "github.com/urfave/cli"

var (
	BODY = &cli.StringFlag{
		Name:    "body",
		Usage:   "the request body",
		Aliases: []string{"b"},
	}
	HEADER = &cli.StringFlag{
		Name:    "header",
		Usage:   "the request header",
		Aliases: []string{"h"},
	}
	COOKIE = &cli.StringFlag{
		Name:    "cookie",
		Usage:   "the request cookie",
		Aliases: []string{"c"},
	}
)
