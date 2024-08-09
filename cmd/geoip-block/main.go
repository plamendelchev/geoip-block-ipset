package main

import (
	"github.com/plamendelchev/geoip-block-ipset"
	"github.com/alecthomas/kong"
)

type Globals struct {
	Config string `short:"c" help:"Path to configuration file" default:"/etc/geoip-block.conf" type:"path"`
	Debug  bool   `short:"d" help:"Enable debug mode"          default:"false"`
}

type (
	CreateCmd  struct{}
	DeleteCmd struct{}
)

func (c *CreateCmd) Run(globals *Globals) error {
	err := geoip.Create(globals.Config, globals.Debug)
	if err != nil {
		return err
	}
	return nil
}

func (c *DeleteCmd) Run(globals *Globals) error {
	err := geoip.Delete(globals.Config, globals.Debug)
	if err != nil {
		return err
	}
	return nil
}

var cli struct {
	Globals

	Create  CreateCmd  `cmd:"" help:"Create geoip blocking"`
	Delete DeleteCmd `cmd:"" help:"Remove geoip blocking"`
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
