package provisioner

import (
	"fmt"
	"net"
	"os"
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
	"github.com/dirien/minectl/internal/logging"
	"github.com/dirien/minectl/internal/manifest"
	"github.com/dirien/minectl/internal/progress"
	"github.com/dirien/minectl/internal/rcon"
	"github.com/pkg/errors"
)

const (
	minecraftProxyTitle                 = "📡 Minecraft %s Proxy"
	minecraftServerTitle                = "🗺 Minecraft %s edition"
	minecraftSelectedCloudProviderTitle = "🛎 Using cloud provider %s"
	minecraftListServersTitle           = "📒 List all server"

	minecraftServerDeletingTitle  = "🪓 Deleting server (%s)..."
	minecraftServerDeleteTitle    = "🗑 Server (%s) deleted."
	minecraftServerNotDeleteTitle = "❌ Server (%s) not deleted."

	minecraftServerCreatingTitle  = "🏗 Creating server (%s)..."
	minecraftServerCreateTitle    = "✅ Server (%s) created."
	minecraftServerNotCreateTitle = "❌ Server (%s) not created."

	minecraftServerStartingTitle = "🎬 Starting server..."
	minecraftServerStartTitle    = "✅ Server successfully started."
	minecraftServerNotStartTitle = "❌ Server failed starting."

	minecraftServerUpdatingTitle  = "🆙 Update server (%s)..."
	minecraftServerUpdateTitle    = "✅ Server (%s) updated."
	minecraftServerNotUpdateTitle = "❌ Server (%s) update failed."

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
	auto    automation.Automation
	args    automation.ServerArgs
	logging *logging.MinectlLogging
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
	p.logging.RawMessage("🚧 Plugins feature is still in beta...")
	indicator := progress.NewIndicator(fmt.Sprintf("⤴️ Upload plugin to server (%s)...", common.Green(p.args.MinecraftResource.GetName())), p.logging)
	indicator.FinalMessage = fmt.Sprintf("✅ Plugin (%s) uploaded.", common.Green(p.args.MinecraftResource.GetName()))
	indicator.ErrorMessage = fmt.Sprintf("❌ Plugin (%s) not uploaded.", common.Green(p.args.MinecraftResource.GetName()))
	indicator.Start()
	err := p.auto.UploadPlugin(p.args.ID, p.args, plugin, destination)
	indicator.StopE(err)
	return err
}

func (p *MinectlProvisioner) UpdateServer() error {
	indicator := progress.NewIndicator(fmt.Sprintf(minecraftServerUpdatingTitle, common.Green(p.args.MinecraftResource.GetName())), p.logging)
	indicator.FinalMessage = fmt.Sprintf(minecraftServerUpdateTitle, common.Green(p.args.MinecraftResource.GetName()))
	indicator.ErrorMessage = fmt.Sprintf(minecraftServerNotUpdateTitle, common.Green(p.args.MinecraftResource.GetName()))
	indicator.Start()
	err := p.auto.UpdateServer(p.args.ID, p.args)
	indicator.StopE(err)
	return err
}

// wait that server is ready... Currently, on for Java based Editions (TCP), as Bedrock is UDP
func (p *MinectlProvisioner) waitForMinecraftServerReady(server *automation.ResourceResults) error {
	if p.args.MinecraftResource.GetEdition() != "bedrock" && p.args.MinecraftResource.GetEdition() != "nukkit" && p.args.MinecraftResource.GetEdition() != "powernukkit" {
		indicator := progress.NewIndicator(minecraftServerStartingTitle, p.logging)
		defer indicator.StopE(nil)
		indicator.FinalMessage = minecraftServerStartTitle
		indicator.ErrorMessage = minecraftServerNotStartTitle
		indicator.Start()
		check := fmt.Sprintf("%s:%d", server.PublicIP, p.args.MinecraftResource.GetPort())
		checkCounter := 0

		for checkCounter < startCheckCount {
			timeout, err := net.DialTimeout("tcp", check, 15*time.Second)
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
	indicator := progress.NewIndicator(fmt.Sprintf(minecraftServerCreatingTitle, common.Green(p.args.MinecraftResource.GetName())), p.logging)
	indicator.FinalMessage = fmt.Sprintf(minecraftServerCreateTitle, common.Green(p.args.MinecraftResource.GetName()))
	indicator.ErrorMessage = fmt.Sprintf(minecraftServerNotCreateTitle, common.Green(p.args.MinecraftResource.GetName()))
	indicator.Start()
	server, err := p.auto.CreateServer(p.args)
	indicator.StopE(err)
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
	indicator := progress.NewIndicator(fmt.Sprintf(minecraftServerDeletingTitle, common.Green(p.args.MinecraftResource.GetName())), p.logging)
	indicator.FinalMessage = fmt.Sprintf(minecraftServerDeleteTitle, common.Green(p.args.MinecraftResource.GetName()))
	indicator.ErrorMessage = fmt.Sprintf(minecraftServerNotDeleteTitle, common.Green(p.args.MinecraftResource.GetName()))
	indicator.Start()
	err := p.auto.DeleteServer(p.args.ID, p.args)
	indicator.StopE(err)
	return err
}

func ListProvisioner(options *MinectlProvisionerListOpts, logging ...*logging.MinectlLogging) (*MinectlProvisioner, error) {
	logging[0].RawMessage(minecraftListServersTitle)
	cloudProvider, err := getProvisioner(options.Provider, options.Region)
	logging[0].PrintMixedGreen(minecraftSelectedCloudProviderTitle, cloud.GetCloudProviderFullName(options.Provider))
	if err != nil {
		return nil, err
	}
	p := &MinectlProvisioner{
		auto:    cloudProvider,
		logging: logging[0],
	}
	return p, nil
}

func getProvisioner(provider, region string) (automation.Automation, error) { //nolint: gocyclo
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

func NewProvisioner(options *MinectlProvisionerOpts, logging ...*logging.MinectlLogging) (*MinectlProvisioner, error) {
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

	logging[0].PrintMixedGreen(minecraftSelectedCloudProviderTitle, cloud.GetCloudProviderFullName(args.MinecraftResource.GetCloud()))

	if args.MinecraftResource.IsProxyServer() {
		logging[0].PrintMixedGreen(minecraftProxyTitle, args.MinecraftResource.GetEdition())
	} else {
		logging[0].PrintMixedGreen(minecraftServerTitle, args.MinecraftResource.GetEdition())
	}

	p := &MinectlProvisioner{
		auto:    cloudProvider,
		args:    args,
		logging: logging[0],
	}
	return p, nil
}
