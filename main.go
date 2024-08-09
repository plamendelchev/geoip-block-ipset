package main

import (
	"github.com/alecthomas/kong"
)

type Globals struct {
	Config string `short:"c" help:"Path to configuration file" default:"/etc/geoip-block.conf" type:"path"`
	Debug  bool   `short:"d" help:"Enable debug mode"          default:"false"`
}

type (
	CreateCmd  struct{}
	DestroyCmd struct{}
)

func (c *CreateCmd) Run(globals *Globals) error {
	err := Create(globals.Config, globals.Debug)
	if err != nil {
		return err
	}
	return nil
}

func (c *DestroyCmd) Run(globals *Globals) error {
	err := Destroy(globals.Config, globals.Debug)
	if err != nil {
		return err
	}
	return nil
}

var cli struct {
	Globals

	Create  CreateCmd  `cmd:"" help:"Create geoip blocking"`
	Destroy DestroyCmd `cmd:"" help:"Remove geoip blocking"`
}

func main() {
	ctx := kong.Parse(
		&cli,
		kong.Name("geoip-block-ipset"),
		kong.Description("Simple tool to whitelist countries in ipset"),
		kong.UsageOnError(),
	)

	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}
