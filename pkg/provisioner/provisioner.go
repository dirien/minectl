package provisioner

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/minectl/pkg/cloud/vexxhost"

	"github.com/minectl/pkg/cloud/ionos"

	"github.com/minectl/pkg/cloud/oci"

	"github.com/minectl/pkg/progress"

	"github.com/minectl/pkg/logging"

	"github.com/minectl/pkg/cloud/azure"

	"github.com/minectl/pkg/rcon"

	"github.com/minectl/pkg/cloud/vultr"

	"github.com/minectl/pkg/cloud/gce"

	"github.com/minectl/pkg/cloud/equinix"

	"github.com/minectl/pkg/cloud/ovh"

	"github.com/minectl/pkg/cloud/linode"
	"github.com/pkg/errors"

	"github.com/minectl/pkg/automation"
	"github.com/minectl/pkg/cloud"
	"github.com/minectl/pkg/cloud/aws"
	"github.com/minectl/pkg/cloud/civo"
	"github.com/minectl/pkg/cloud/do"
	"github.com/minectl/pkg/cloud/hetzner"
	"github.com/minectl/pkg/cloud/scaleway"
	"github.com/minectl/pkg/common"
	"github.com/minectl/pkg/manifest"
)

type MinectlProvisionerOpts struct {
	ManifestPath string
	ID           string
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
	CreateServer(wait bool) (*automation.RessourceResults, error)
	DeleteServer() error
	UpdateServer() error
	UploadPlugin(plugin, destination string) error
	ListServer() ([]automation.RessourceResults, error)
	GetServer() (*automation.RessourceResults, error)
	DoRCON() error
}

func (p *MinectlProvisioner) GetServer() (*automation.RessourceResults, error) {
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
	indicator := progress.NewIndicator(fmt.Sprintf("🆙 Update server (%s)...", common.Green(p.args.MinecraftResource.GetName())), p.logging)
	indicator.FinalMessage = fmt.Sprintf("✅ Server (%s) updated.", common.Green(p.args.MinecraftResource.GetName()))
	indicator.ErrorMessage = fmt.Sprintf("❌ Server (%s) update failed.", common.Green(p.args.MinecraftResource.GetName()))
	indicator.Start()
	err := p.auto.UpdateServer(p.args.ID, p.args)
	indicator.StopE(err)
	return err
}

// wait that server is ready... Currently, on for Java based Editions (TCP), as Bedrock is UDP
func (p *MinectlProvisioner) waitForMinecraftServerReady(server *automation.RessourceResults) error {
	if p.args.MinecraftResource.GetEdition() != "bedrock" && p.args.MinecraftResource.GetEdition() != "nukkit" && p.args.MinecraftResource.GetEdition() != "powernukkit" {
		indicator := progress.NewIndicator("🕹 Starting Minecraft server...", p.logging)
		defer indicator.StopE(nil)
		indicator.FinalMessage = "✅ Minecraft server successfully started."
		indicator.ErrorMessage = "❌ Minecraft server failed starting."
		indicator.Start()
		check := fmt.Sprintf("%s:%d", server.PublicIP, p.args.MinecraftResource.GetPort())
		checkCounter := 0

		for checkCounter < 500 {
			timeout, err := net.DialTimeout("tcp", check, 15*time.Second)
			if err != nil {
				time.Sleep(15 * time.Second)
				checkCounter++
			}
			if timeout != nil {
				err = timeout.Close()
				if err != nil {
					return errors.Errorf("Timeout error: %s\n", err)
				}
				break
			}
		}
	}
	return nil
}

func (p *MinectlProvisioner) CreateServer(wait bool) (*automation.RessourceResults, error) {
	indicator := progress.NewIndicator(fmt.Sprintf("🏗 Creating server (%s)...", common.Green(p.args.MinecraftResource.GetName())), p.logging)
	indicator.FinalMessage = fmt.Sprintf("✅ Server (%s) created.", common.Green(p.args.MinecraftResource.GetName()))
	indicator.ErrorMessage = fmt.Sprintf("❌ Server (%s) not created.", common.Green(p.args.MinecraftResource.GetName()))
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

func (p *MinectlProvisioner) ListServer() ([]automation.RessourceResults, error) {
	return p.auto.ListServer()
}

func (p *MinectlProvisioner) DeleteServer() error {
	indicator := progress.NewIndicator(fmt.Sprintf("🪓 Deleting server (%s)...", common.Green(p.args.MinecraftResource.GetName())), p.logging)
	indicator.FinalMessage = fmt.Sprintf("🗑 Server (%s) deleted.", common.Green(p.args.MinecraftResource.GetName()))
	indicator.ErrorMessage = fmt.Sprintf("❌ Server (%s) not deleted.", common.Green(p.args.MinecraftResource.GetName()))
	indicator.Start()
	err := p.auto.DeleteServer(p.args.ID, p.args)
	indicator.StopE(err)
	return err
}

func ListProvisioner(options *MinectlProvisionerListOpts, logging ...*logging.MinectlLogging) (*MinectlProvisioner, error) {
	logging[0].RawMessage("📒 List all server")
	cloudProvider, err := getProvisioner(options.Provider, options.Region)
	logging[0].PrintMixedGreen("🛎 Using cloud provider %s", cloud.GetCloudProviderFullName(options.Provider))
	if err != nil {
		return nil, err
	}
	p := &MinectlProvisioner{
		auto:    cloudProvider,
		logging: logging[0],
	}
	return p, nil
}

func getProvisioner(provider, region string) (automation.Automation, error) {
	switch provider {
	case "hetzner":
		cloudProvider, err := hetzner.NewHetzner(os.Getenv("HCLOUD_TOKEN"))
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case "do":
		cloudProvider, err := do.NewDigitalOcean(os.Getenv("DIGITALOCEAN_TOKEN"))
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case "civo":
		cloudProvider, err := civo.NewCivo(os.Getenv("CIVO_TOKEN"), region)
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case "scaleway":
		cloudProvider, err := scaleway.NewScaleway(os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), os.Getenv("ORGANISATION_ID"), region)
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case "linode":
		cloudProvider, err := linode.NewLinode(os.Getenv("LINODE_TOKEN"))
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case "ovh":
		cloudProvider, err := ovh.NewOVHcloud(os.Getenv("OVH_ENDPOINT"), os.Getenv("APPLICATION_KEY"), os.Getenv("APPLICATION_SECRET"), os.Getenv("CONSUMER_KEY"), os.Getenv("SERVICENAME"), region)
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case "equinix":
		cloudProvider, err := equinix.NewEquinix(os.Getenv("PACKET_AUTH_TOKEN"), os.Getenv("EQUINIX_PROJECT"))
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case "gce":
		cloudProvider, err := gce.NewGCE(os.Getenv("GCE_KEY"), region)
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case "vultr":
		cloudProvider, err := vultr.NewVultr(os.Getenv("VULTR_API_KEY"))
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case "azure":
		cloudProvider, err := azure.NewAzure(os.Getenv("AZURE_AUTH_LOCATION"))
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case "oci":
		cloudProvider, err := oci.NewOCI()
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case "ionos":
		cloudProvider, err := ionos.NewIONOS(os.Getenv("IONOS_USERNAME"), os.Getenv("IONOS_PASSWORD"), os.Getenv("IONOS_TOKEN"))
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case "aws":
		cloudProvider, err := aws.NewAWS(region, os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_SESSION_TOKEN"))
		if err != nil {
			return nil, err
		}
		return cloudProvider, nil
	case "vexxhost":
		cloudProvider, err := vexxhost.NewVEXXHOST()
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
	}
	cloudProvider, err = getProvisioner(args.MinecraftResource.GetCloud(), args.MinecraftResource.GetRegion())
	if err != nil {
		return nil, err
	}

	logging[0].PrintMixedGreen("🛎 Using cloud provider %s", cloud.GetCloudProviderFullName(args.MinecraftResource.GetCloud()))

	if args.MinecraftResource.IsProxyServer() {
		logging[0].PrintMixedGreen("📡 Minecraft %s Proxy", args.MinecraftResource.GetEdition())
	} else {
		logging[0].PrintMixedGreen("🗺 Minecraft %s edition", args.MinecraftResource.GetEdition())
	}

	p := &MinectlProvisioner{
		auto:    cloudProvider,
		args:    args,
		logging: logging[0],
	}
	return p, nil
}
