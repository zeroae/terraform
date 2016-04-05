package command

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/builtin/providers/aws"
	"github.com/hashicorp/terraform/builtin/providers/azurerm"
	"github.com/hashicorp/terraform/builtin/providers/cloudflare"
	"github.com/hashicorp/terraform/builtin/providers/digitalocean"
	"github.com/hashicorp/terraform/builtin/providers/google"
	"github.com/hashicorp/terraform/builtin/providers/null"
	"github.com/hashicorp/terraform/builtin/providers/template"
	"github.com/hashicorp/terraform/builtin/provisioners/local-exec"
	"github.com/hashicorp/terraform/builtin/provisioners/remote-exec"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kardianos/osext"
)

// InternalPluginCommand is a Command implementation that allows plugins to be
// compiled into the main Terraform binary and executed via a subcommand.
type InternalPluginCommand struct {
	Meta
}

const TFSPACE = "-TFSPACE-"

var Providers = map[string]plugin.ProviderFunc{
	"aws":          aws.Provider,
	"azurerm":      azurerm.Provider,
	"cloudflare":   cloudflare.Provider,
	"digitalocean": digitalocean.Provider,
	"google":       google.Provider,
	"null":         null.Provider,
	"template":     template.Provider,
}

var Provisioners = map[string]plugin.ProvisionerFunc{
	"local-exec":  func() terraform.ResourceProvisioner { return new(localexec.ResourceProvisioner) },
	"remote-exec": func() terraform.ResourceProvisioner { return new(remoteexec.ResourceProvisioner) },
}

var pluginRegexp = regexp.MustCompile("terraform-(provider|provisioner)-(.+)")

// BuildPluginCommandString builds a special string for executing internal
// plugins. It has the following format:
//
// 	/path/to/terraform-TFSPACE-internal-plugin-TFSPACE-terraform-provider-aws
//
// We split the string on -TFSPACE- to build the command executor. The reason we
// use -TFSPACE- is so we can support spaces in the /path/to/terraform part.
func BuildPluginCommandString(pluginType, pluginName string) (string, error) {
	terraformPath, err := osext.Executable()
	if err != nil {
		return "", err
	}
	pluginString := fmt.Sprintf("terraform-%s-%s", pluginType, pluginName)
	commandString := fmt.Sprintf("%s%s%s%s%s", terraformPath, TFSPACE, "internal-plugin", TFSPACE, pluginString)
	return commandString, nil
}

// parsePluginParts reads a string like terraform-provider-aws and breaks it
// into pluginType and pluginName. This format corresponds to the filenames used
// for disk-based plugins that shipped with Terraform < 0.7
func parsePluginParts(input string) (string, string, error) {
	parts := pluginRegexp.FindStringSubmatch(input)
	if len(parts) != 3 {
		return "", "", fmt.Errorf("Error parsing plugin argument [DEBUG]: %#v", parts)
	}
	pluginType := parts[1] // capture group 1 (provider|provisioner)
	pluginName := parts[2] // capture group 2 (.+)
	return pluginType, pluginName, nil
}

func (c *InternalPluginCommand) Run(args []string) int {
	if len(args) != 1 {
		c.Ui.Error("Wrong number of args")
		return 1
	}

	if args[0] == "version" {
		c.Ui.Output(terraform.Version)
		os.Exit(0)
	}

	pluginType, pluginName, err := parsePluginParts(args[0])
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	switch pluginType {
	case "provider":
		pluginFunc, found := Providers[pluginName]
		if !found {
			c.Ui.Error(fmt.Sprintf("Could not load provider: %s", pluginName))
			return 1
		}
		log.Printf("Starting provider plugin %s", pluginName)
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: pluginFunc,
		})
	case "provisioner":
		pluginFunc, found := Provisioners[pluginName]
		if !found {
			c.Ui.Error(fmt.Sprintf("Could not load provisioner: %s", pluginName))
			return 1
		}
		log.Printf("Starting provisioner plugin %s", pluginName)
		plugin.Serve(&plugin.ServeOpts{
			ProvisionerFunc: pluginFunc,
		})
	default:
		return 1
	}

	return 0
}

func (c *InternalPluginCommand) Help() string {
	helpText := `
Usage: terraform internal-plugin PLUGIN

  Runs an internally-compiled version of a plugin from the terraform binary.

  NOTE: this is an internal command and you should not call it yourself.
`

	return strings.TrimSpace(helpText)
}

func (c *InternalPluginCommand) Synopsis() string {
	return "internal plugin command"
}
