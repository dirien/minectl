package provisioner

import (
	"fmt"
	"github.com/minectl/pgk/automation"
	"github.com/minectl/pgk/cloud"
	"github.com/minectl/pgk/cloud/civo"
	"github.com/minectl/pgk/cloud/do"
	"github.com/minectl/pgk/cloud/scaleway"
	"github.com/minectl/pgk/common"
	"github.com/minectl/pgk/manifest"
	"github.com/pkg/errors"
	"os"
)

type PulumiProvisioner struct {
	auto     automation.Automation
	Manifest *manifest.MinecraftServerManifest
	args     automation.ServerArgs
}

type Provisioner interface {
	CreateServer() (*automation.RessourceResults, error)
	DeleteServer() error
	UpdateServer() (*automation.RessourceResults, error)
	ListServer() ([]automation.RessourceResults, error)
}

func (p PulumiProvisioner) UpdateServer() (*automation.RessourceResults, error) {
	return p.auto.UpdateServer(p.args)
}

func (p PulumiProvisioner) CreateServer() (*automation.RessourceResults, error) {
	return p.auto.CreateServer(p.args)
}

func (p PulumiProvisioner) ListServer() ([]automation.RessourceResults, error) {
	return p.auto.ListServer()
}

func (p PulumiProvisioner) DeleteServer() error {
	return p.auto.DeleteServer(p.args.ID, p.args)
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
	}
	return p, nil
}
