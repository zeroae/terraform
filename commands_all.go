// +build !core

package main

import (
	"github.com/hashicorp/terraform/command"
	"github.com/mitchellh/cli"
)

func init() {
	meta := command.Meta{
		Color:       true,
		ContextOpts: &ContextOpts,
		Ui:          Ui,
	}

	Commands["internal-plugin"] = func() (cli.Command, error) {
		return &command.InternalPluginCommand{
			Meta: meta,
		}, nil
	}
}
