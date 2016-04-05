package command

import (
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
)

var InternalProviders = map[string]plugin.ProviderFunc{
	"aws":          aws.Provider,
	"azurerm":      azurerm.Provider,
	"cloudflare":   cloudflare.Provider,
	"digitalocean": digitalocean.Provider,
	"google":       google.Provider,
	"null":         null.Provider,
	"template":     template.Provider,
}

var InternalProvisioners = map[string]plugin.ProvisionerFunc{
	"local-exec":  func() terraform.ResourceProvisioner { return new(localexec.ResourceProvisioner) },
	"remote-exec": func() terraform.ResourceProvisioner { return new(remoteexec.ResourceProvisioner) },
}
