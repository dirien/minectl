package provisioner

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/minectl/pgk/cloud/gce"

	"github.com/minectl/pgk/cloud/equinix"

	"github.com/minectl/pgk/cloud/ovh"

	"github.com/minectl/pgk/cloud/linode"

	"github.com/briandowns/spinner"
	"github.com/minectl/pgk/automation"
	"github.com/minectl/pgk/cloud"
	"github.com/minectl/pgk/cloud/civo"
	"github.com/minectl/pgk/cloud/do"
	"github.com/minectl/pgk/cloud/hetzner"
	"github.com/minectl/pgk/cloud/scaleway"
	"github.com/minectl/pgk/common"
	"github.com/minectl/pgk/manifest"
	"github.com/pkg/errors"
)

type PulumiProvisioner struct {
	auto     automation.Automation
	Manifest *manifest.MinecraftServerManifest
	args     automation.ServerArgs
	spinner  *spinner.Spinner
}

type Provisioner interface {
	CreateServer(wait bool) (*automation.RessourceResults, error)
	DeleteServer() error
	UpdateServer() error
	UploadPlugin(plugin, destination string) error
	ListServer() ([]automation.RessourceResults, error)
	GetServer() (*automation.RessourceResults, error)
}

func (p *PulumiProvisioner) startSpinner(prefix string) {
	p.spinner.Prefix = prefix
	//p.spinner.FinalMSG = finalMSG
	p.spinner.HideCursor = true
	p.spinner.Start()
}

func (p *PulumiProvisioner) stopSpinner() {
	p.spinner.Stop()
}

func (p *PulumiProvisioner) GetServer() (*automation.RessourceResults, error) {
	return p.auto.GetServer(p.args.ID, p.args)
}

func (p *PulumiProvisioner) UploadPlugin(plugin, destination string) error {
	fmt.Println("🚧 Plugins feature is still in beta...")
	p.startSpinner(
		fmt.Sprintf("⤴️ Upload plugin to server (%s)... ", common.Green(p.args.MinecraftServer.GetName())))
	err := p.auto.UploadPlugin(p.args.ID, p.args, plugin, destination)
	if err == nil {
		fmt.Printf("\n✅ Plugin (%s) uploaded\n", common.Green(p.args.MinecraftServer.GetName()))
	}
	p.stopSpinner()
	return err
}

func (p *PulumiProvisioner) UpdateServer() error {
	p.startSpinner(
		fmt.Sprintf("🆙 Update server (%s)... ", common.Green(p.args.MinecraftServer.GetName())))
	err := p.auto.UpdateServer(p.args.ID, p.args)
	if err == nil {
		fmt.Printf("\n✅ Server (%s) updated\n", common.Green(p.args.MinecraftServer.GetName()))
	}
	p.stopSpinner()
	return err
}

//wait that server is ready... Currently on for Java based Editions (TCP), as Bedrock is UDP
func (p *PulumiProvisioner) waitForMinecraftServerReady(server *automation.RessourceResults) {
	if p.args.MinecraftServer.GetEdition() != "bedrock" {
		p.startSpinner("🕹 Starting Minecraft server... ")
		check := fmt.Sprintf("%s:%d", server.PublicIP, p.args.MinecraftServer.GetPort())
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
					fmt.Printf("Timeout error: %s\n", err)
					p.stopSpinner()
				}
				break
			}
		}
		p.stopSpinner()
	}
	fmt.Println("\n✅ Minecraft successfully started.")
}

func (p *PulumiProvisioner) CreateServer(wait bool) (*automation.RessourceResults, error) {
	p.startSpinner(
		fmt.Sprintf("🏗 Creating server (%s)... ", common.Green(p.args.MinecraftServer.GetName())))
	server, err := p.auto.CreateServer(p.args)
	if err != nil {
		p.stopSpinner()
		return nil, err
	}
	p.stopSpinner()
	if wait {
		p.waitForMinecraftServerReady(server)
	}
	fmt.Printf("\n✅ Server (%s) created\n", common.Green(p.args.MinecraftServer.GetName()))
	return server, err
}

func (p *PulumiProvisioner) ListServer() ([]automation.RessourceResults, error) {
	return p.auto.ListServer()
}

func (p *PulumiProvisioner) DeleteServer() error {
	p.startSpinner(
		fmt.Sprintf("🪓 Deleting server (%s)... ", common.Green(p.args.MinecraftServer.GetName())))
	err := p.auto.DeleteServer(p.args.ID, p.args)
	p.stopSpinner()
	if err == nil {
		fmt.Printf("\n🗑 Server (%s) deleted\n", common.Green(p.args.MinecraftServer.GetName()))
	}
	return err
}

// NewProvisioner has variable args: only manifest file or manifest file and the id
func NewProvisioner(args ...string) (*PulumiProvisioner, error) {
	if len(args) == 2 {
		return newProvisioner(args[0], args[1])
	}
	return newProvisioner(args[0], "")
}

func ListProvisioner(args ...string) (*PulumiProvisioner, error) {
	fmt.Println("📒 List all server")
	cloudProvider, err := getProvisioner(args[0], args[1])
	common.PrintMixedGreen("🛎 Using cloud provider %s\n", cloud.GetCloudProviderFullName(args[0]))
	if err != nil {
		return nil, err
	}
	p := &PulumiProvisioner{
		auto: cloudProvider,
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
	default:
		return nil, errors.Errorf("Could not find provider %s", provider)
	}
}

func newProvisioner(manifestPath, id string) (*PulumiProvisioner, error) {
	manifest, err := manifest.NewMinecraftServer(manifestPath)
	if err != nil {
		return nil, err
	}
	args := automation.ServerArgs{
		MinecraftServer: manifest.MinecraftServer,
		ID:              id,
	}
	cloudProvider, err := getProvisioner(args.MinecraftServer.GetCloud(), args.MinecraftServer.GetRegion())
	common.PrintMixedGreen("🛎 Using cloud provider %s\n", cloud.GetCloudProviderFullName(args.MinecraftServer.GetCloud()))
	if err != nil {
		return nil, err
	}
	common.PrintMixedGreen("🗺 Minecraft %s edition\n", args.MinecraftServer.GetEdition())

	p := &PulumiProvisioner{
		auto:     cloudProvider,
		Manifest: manifest,
		args:     args,
		spinner:  spinner.New(spinner.CharSets[11], 100*time.Millisecond),
	}
	return p, nil
}
