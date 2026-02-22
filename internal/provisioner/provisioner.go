package provisioner

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/dirien/minectl-sdk/automation"
	"github.com/dirien/minectl-sdk/cloud"
	"github.com/dirien/minectl-sdk/cloud/akamai"
	"github.com/dirien/minectl-sdk/cloud/aws"
	"github.com/dirien/minectl-sdk/cloud/azure"
	"github.com/dirien/minectl-sdk/cloud/civo"
	"github.com/dirien/minectl-sdk/cloud/do"
	"github.com/dirien/minectl-sdk/cloud/exoscale"
	"github.com/dirien/minectl-sdk/cloud/fuga"
	"github.com/dirien/minectl-sdk/cloud/gce"
	"github.com/dirien/minectl-sdk/cloud/hetzner"
	"github.com/dirien/minectl-sdk/cloud/multipass"
	"github.com/dirien/minectl-sdk/cloud/oci"
	"github.com/dirien/minectl-sdk/cloud/ovh"
	"github.com/dirien/minectl-sdk/cloud/scaleway"
	"github.com/dirien/minectl-sdk/cloud/vexxhost"
	"github.com/dirien/minectl-sdk/cloud/vultr"
	"github.com/dirien/minectl-sdk/common"
	"github.com/dirien/minectl-sdk/model"
	"github.com/dirien/minectl/internal/manifest"
	"github.com/dirien/minectl/internal/rcon"
	"github.com/dirien/minectl/internal/ui"
	"github.com/pkg/errors"
)

const (
	minecraftProxyTitle                 = "Minecraft %s Proxy"
	minecraftServerTitle                = "Minecraft %s edition"
	minecraftSelectedCloudProviderTitle = "Using cloud provider %s"
	minecraftListServersTitle           = "Listing all servers"

	minecraftServerDeletingTitle  = "Deleting server (%s)..."
	minecraftServerDeleteTitle    = "Server (%s) deleted."
	minecraftServerNotDeleteTitle = "Server (%s) not deleted."

	minecraftServerCreatingTitle  = "Creating server (%s)..."
	minecraftServerCreateTitle    = "Server (%s) created."
	minecraftServerNotCreateTitle = "Server (%s) not created."

	minecraftServerStartingTitle = "Starting server..."
	minecraftServerStartTitle    = "Server successfully started."
	minecraftServerNotStartTitle = "Server failed to start."

	minecraftServerUpdatingTitle  = "Updating server (%s)..."
	minecraftServerUpdateTitle    = "Server (%s) updated."
	minecraftServerNotUpdateTitle = "Server (%s) update failed."

	startCheckCount = 50
)

type MinectlProvisionerOpts struct {
	ManifestPath      string
	ID                string
	SSHPrivateKeyPath string
}

type MinectlProvisionerListOpts struct {
	Provider string
	Region   string
}

type MinectlProvisioner struct {
	auto automation.Automation
	args automation.ServerArgs
	ui   *ui.UI
}

type Provisioner interface {
	CreateServer(wait bool) (*automation.ResourceResults, error)
	DeleteServer() error
	UpdateServer() error
	UploadPlugin(plugin, destination string) error
	ListServer() ([]automation.ResourceResults, error)
	GetServer() (*automation.ResourceResults, error)
	DoRCON() error
}

func (p *MinectlProvisioner) GetServer() (*automation.ResourceResults, error) {
	return p.auto.GetServer(p.args.ID, p.args)
}

func (p *MinectlProvisioner) DoRCON() error {
	server, err := p.GetServer()
	if err != nil {
		return err
	}
	r := rcon.NewRCON(server.PublicIP, p.args.MinecraftResource.GetRCONPassword(), p.args.MinecraftResource.GetRCONPort())
	r.RunPrompt()
	return nil
}

func (p *MinectlProvisioner) UploadPlugin(plugin, destination string) error {
	p.ui.Warn("Note: Plugins feature is still in beta.")
	spinner := ui.NewSpinner(fmt.Sprintf("Uploading plugin to server (%s)...", common.Green(p.args.MinecraftResource.GetName())), p.ui)
	spinner.FinalMessage = fmt.Sprintf("Plugin (%s) uploaded.", common.Green(p.args.MinecraftResource.GetName()))
	spinner.ErrorMessage = fmt.Sprintf("Plugin (%s) not uploaded.", common.Green(p.args.MinecraftResource.GetName()))
	spinner.Start()
	err := p.auto.UploadPlugin(p.args.ID, p.args, plugin, destination)
	spinner.Stop(err)
	return err
}

func (p *MinectlProvisioner) UpdateServer() error {
	spinner := ui.NewSpinner(fmt.Sprintf(minecraftServerUpdatingTitle, common.Green(p.args.MinecraftResource.GetName())), p.ui)
	spinner.FinalMessage = fmt.Sprintf(minecraftServerUpdateTitle, common.Green(p.args.MinecraftResource.GetName()))
	spinner.ErrorMessage = fmt.Sprintf(minecraftServerNotUpdateTitle, common.Green(p.args.MinecraftResource.GetName()))
	spinner.Start()
	err := p.auto.UpdateServer(p.args.ID, p.args)
	spinner.Stop(err)
	return err
}

// wait that server is ready... Currently, only for Java based Editions (TCP), as Bedrock is UDP
func (p *MinectlProvisioner) waitForMinecraftServerReady(server *automation.ResourceResults) error {
	if p.args.MinecraftResource.GetEdition() != "bedrock" && p.args.MinecraftResource.GetEdition() != "nukkit" && p.args.MinecraftResource.GetEdition() != "powernukkit" {
		spinner := ui.NewSpinner(minecraftServerStartingTitle, p.ui)
		defer spinner.Stop(nil)
		spinner.FinalMessage = minecraftServerStartTitle
		spinner.ErrorMessage = minecraftServerNotStartTitle
		spinner.Start()
		check := net.JoinHostPort(server.PublicIP, strconv.Itoa(p.args.MinecraftResource.GetPort()))
		checkCounter := 0
		dialer := net.Dialer{Timeout: 15 * time.Second}

		for checkCounter < startCheckCount {
			timeout, err := dialer.Dial("tcp", check)
			if err != nil {
				time.Sleep(15 * time.Second)
				checkCounter++
			}
			if timeout != nil {
				err = timeout.Close()
				if err != nil {
					return errors.Errorf("timeout error: %s", err)
				}
				break
			}
		}
	}
	return nil
}

func (p *MinectlProvisioner) CreateServer(wait bool) (*automation.ResourceResults, error) {
	spinner := ui.NewSpinner(fmt.Sprintf(minecraftServerCreatingTitle, common.Green(p.args.MinecraftResource.GetName())), p.ui)
	spinner.FinalMessage = fmt.Sprintf(minecraftServerCreateTitle, common.Green(p.args.MinecraftResource.GetName()))
	spinner.ErrorMessage = fmt.Sprintf(minecraftServerNotCreateTitle, common.Green(p.args.MinecraftResource.GetName()))
	spinner.Start()
	server, err := p.auto.CreateServer(p.args)
	spinner.Stop(err)
	if err != nil {
		return nil, err
	}

	if wait {
		err := p.waitForMinecraftServerReady(server)
		if err != nil {
			return nil, err
		}
	}
	return server, err
}

func (p *MinectlProvisioner) ListServer() ([]automation.ResourceResults, error) {
	return p.auto.ListServer()
}

func (p *MinectlProvisioner) DeleteServer() error {
	spinner := ui.NewSpinner(fmt.Sprintf(minecraftServerDeletingTitle, common.Green(p.args.MinecraftResource.GetName())), p.ui)
	spinner.FinalMessage = fmt.Sprintf(minecraftServerDeleteTitle, common.Green(p.args.MinecraftResource.GetName()))
	spinner.ErrorMessage = fmt.Sprintf(minecraftServerNotDeleteTitle, common.Green(p.args.MinecraftResource.GetName()))
	spinner.Start()
	err := p.auto.DeleteServer(p.args.ID, p.args)
	spinner.Stop(err)
	return err
}

func ListProvisioner(options *MinectlProvisionerListOpts, u *ui.UI) (*MinectlProvisioner, error) {
	u.Info(minecraftListServersTitle)
	cloudProvider, err := getProvisioner(options.Provider, options.Region)
	u.Info(fmt.Sprintf(minecraftSelectedCloudProviderTitle, cloud.GetCloudProviderFullName(options.Provider)))
	if err != nil {
		return nil, err
	}
	p := &MinectlProvisioner{
		auto: cloudProvider,
		ui:   u,
	}
	return p, nil
}

func getProvisioner(provider, region string) (automation.Automation, error) { //nolint:gocyclo // large switch is inherent to provider selection
	switch provider {
	case model.PROVIDER_HETZNER:
		cloudProvider, err := hetzner.NewHetzner(os.Getenv("HCLOUD_TOKEN"))
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_DIGITALOCEAN:
		cloudProvider, err := do.NewDigitalOcean(os.Getenv("DIGITALOCEAN_TOKEN"))
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_CIVO:
		cloudProvider, err := civo.NewCivo(os.Getenv("CIVO_TOKEN"), region)
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_SCALEWAY:
		cloudProvider, err := scaleway.NewScaleway(os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), os.Getenv("ORGANISATION_ID"), region)
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_AKAMAI:
		cloudProvider, err := akamai.NewAkamai(os.Getenv("LINODE_TOKEN"))
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_OVH:
		cloudProvider, err := ovh.NewOVHcloud(os.Getenv("OVH_ENDPOINT"), os.Getenv("APPLICATION_KEY"), os.Getenv("APPLICATION_SECRET"), os.Getenv("CONSUMER_KEY"), os.Getenv("SERVICENAME"), region)
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_GCE:
		cloudProvider, err := gce.NewGCE(region)
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_VULTR:
		cloudProvider, err := vultr.NewVultr(os.Getenv("VULTR_API_KEY"))
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_AZURE:
		cloudProvider, err := azure.NewAzure()
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_OCI:
		cloudProvider, err := oci.NewOCI()
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_AWS:
		cloudProvider, err := aws.NewAWS(region)
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_VEXXHOST:
		cloudProvider, err := vexxhost.NewVEXXHOST()
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_FUGA:
		cloudProvider, err := fuga.NewFuga()
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_EXOSCALE:
		cloudProvider, err := exoscale.NewExoscale(os.Getenv("EXOSCALE_API_KEY"), os.Getenv("EXOSCALE_API_SECRET"))
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case model.PROVIDER_MULTIPASS:
		cloudProvider, err := multipass.NewMultipass()
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil

	default:
		return nil, errors.Errorf("Could not find provider %s", provider)
	}
}

func NewProvisioner(options *MinectlProvisionerOpts, u *ui.UI) (*MinectlProvisioner, error) {
	var cloudProvider automation.Automation

	minecraftResource, err := manifest.NewMinecraftResource(options.ManifestPath)
	if err != nil {
		return nil, err
	}
	args := automation.ServerArgs{
		MinecraftResource: minecraftResource,
		ID:                options.ID,
		SSHPrivateKeyPath: options.SSHPrivateKeyPath,
	}
	cloudProvider, err = getProvisioner(args.MinecraftResource.GetCloud(), args.MinecraftResource.GetRegion())
	if err != nil {
		return nil, err
	}

	u.Info(fmt.Sprintf(minecraftSelectedCloudProviderTitle, cloud.GetCloudProviderFullName(args.MinecraftResource.GetCloud())))

	if args.MinecraftResource.IsProxyServer() {
		u.Info(fmt.Sprintf(minecraftProxyTitle, args.MinecraftResource.GetEdition()))
	} else {
		u.Info(fmt.Sprintf(minecraftServerTitle, args.MinecraftResource.GetEdition()))
	}

	p := &MinectlProvisioner{
		auto: cloudProvider,
		args: args,
		ui:   u,
	}
	return p, nil
}
