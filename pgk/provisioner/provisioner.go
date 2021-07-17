package provisioner

import (
	"fmt"
	"os"
	"time"

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
	CreateServer() (*automation.RessourceResults, error)
	DeleteServer() error
	UpdateServer() (*automation.RessourceResults, error)
	ListServer() ([]automation.RessourceResults, error)
}

func (p PulumiProvisioner) startSpinner(prefix, finalMSG string) {
	p.spinner.Prefix = prefix
	p.spinner.FinalMSG = finalMSG
	p.spinner.Start()
}

func (p PulumiProvisioner) stopSpinner() {
	p.spinner.Stop()
}

func (p PulumiProvisioner) UpdateServer() (*automation.RessourceResults, error) {
	return p.auto.UpdateServer(p.args)
}

func (p PulumiProvisioner) CreateServer() (*automation.RessourceResults, error) {
	p.startSpinner(
		fmt.Sprintf("🏗 Creating server (%s)... ", common.Green(p.args.MinecraftServer.GetName())),
		fmt.Sprintf("\n✅ Server (%s) created\n", common.Green(p.args.MinecraftServer.GetName())))
	server, err := p.auto.CreateServer(p.args)
	p.stopSpinner()
	return server, err
}

func (p PulumiProvisioner) ListServer() ([]automation.RessourceResults, error) {
	return p.auto.ListServer()
}

func (p PulumiProvisioner) DeleteServer() error {
	p.startSpinner(
		fmt.Sprintf("🪓 Deleting server (%s)... ", common.Green(p.args.MinecraftServer.GetName())),
		fmt.Sprintf("\n🗑 Server (%s) deleted\n", common.Green(p.args.MinecraftServer.GetName())))
	err := p.auto.DeleteServer(p.args.ID, p.args)
	p.stopSpinner()
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
