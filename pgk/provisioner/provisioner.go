package provisioner

import (
	"github.com/minectl/pgk/automation"
	"github.com/minectl/pgk/cloud/civo"
	"github.com/minectl/pgk/cloud/do"
	"github.com/minectl/pgk/cloud/scaleway"
	"github.com/minectl/pgk/common"
	"github.com/minectl/pgk/manifest"
	"os"
)

type PulumiProvisioner struct {
	auto     automation.Automation
	Manifest *manifest.MinecraftServer
	args     automation.ServerArgs
}

type Provisioner interface {
	CreateServer() (*automation.RessourceResults, error)
	DeleteServer() error
	UpdateServer() (*automation.RessourceResults, error)
}

func (p PulumiProvisioner) UpdateServer() (*automation.RessourceResults, error) {
	return p.auto.UpdateServer(p.args)
}

func (p PulumiProvisioner) CreateServer() (*automation.RessourceResults, error) {
	return p.auto.CreateServer(p.args)
}

func (p PulumiProvisioner) DeleteServer() error {
	return p.auto.DeleteServer(p.args.ID, p.args)
}

func NewCreateServerProvisioner(manifestPath string) (*PulumiProvisioner, error) {
	return newProvisioner(manifestPath, "")
}

func NewDeleteServerProvisioner(id, manifestPath string) (*PulumiProvisioner, error) {
	return newProvisioner(manifestPath, id)
}

func newProvisioner(manifestPath, id string) (*PulumiProvisioner, error) {
	manifest, err := manifest.NewMinecraftServer(manifestPath)
	if err != nil {
		return nil, err
	}
	args := automation.ServerArgs{
		StackName:  manifest.GetName(),
		SSH:        manifest.GetSSH(),
		Region:     manifest.GetRegion(),
		Size:       manifest.GetSize(),
		Properties: manifest.GetProperties(),
		VolumeSize: manifest.GetVolumeSize(),
		ID:         id,
	}
	var cloudProvider automation.Automation
	if manifest.GetCloud() == "do" {
		common.PrintMixedGreen("Using cloud provider %s\n", "DigitalOcean")
		cloudProvider, err = do.NewDigitalOcean(os.Getenv("DIGITALOCEAN_TOKEN"))
		if err != nil {
			return nil, err
		}
	} else if manifest.GetCloud() == "civo" {
		common.PrintMixedGreen("Using cloud provider %s\n", "Civo")
		cloudProvider, err = civo.NewCivo(os.Getenv("CIVO_TOKEN"), args.Region)
		if err != nil {
			return nil, err
		}
	} else if manifest.GetCloud() == "scaleway" {
		common.PrintMixedGreen("Using cloud provider %s\n", "Scaleway")
		cloudProvider, err = scaleway.NewScaleway(os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), os.Getenv("ORGANISATION_ID"), args.Region)
		if err != nil {
			return nil, err
		}
	}

	p := &PulumiProvisioner{
		auto:     cloudProvider,
		Manifest: manifest,
		args:     args,
	}
	return p, nil
}
