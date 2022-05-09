package scaleway

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minectl/internal/update"

	"github.com/minectl/internal/automation"
	"github.com/minectl/internal/common"
	minctlTemplate "github.com/minectl/internal/template"
	account "github.com/scaleway/scaleway-sdk-go/api/account/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Scaleway struct {
	instanceAPI *instance.API
	accountAPI  *account.API
	tmpl        *minctlTemplate.Template
}

func NewScaleway(accessKey, secretKey, organizationID, region string) (*Scaleway, error) {
	zone, err := scw.ParseZone(region)
	if err != nil {
		return nil, err
	}

	client, err := scw.NewClient(
		scw.WithAuth(accessKey, secretKey),
		scw.WithDefaultOrganizationID(organizationID),
		scw.WithDefaultZone(zone),
	)
	if err != nil {
		return nil, err
	}
	tmpl, err := minctlTemplate.NewTemplateCloudConfig()
	if err != nil {
		return nil, err
	}
	return &Scaleway{
		instanceAPI: instance.NewAPI(client),
		accountAPI:  account.NewAPI(client),
		tmpl:        tmpl,
	}, nil
}

func (s *Scaleway) CreateServer(args automation.ServerArgs) (*automation.ResourceResults, error) {
	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSHKeyFolder()))
	if err != nil {
		return nil, err
	}
	_, err = s.accountAPI.CreateSSHKey(&account.CreateSSHKeyRequest{
		Name:      fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()),
		PublicKey: string(pubKeyFile),
	})
	if err != nil {
		return nil, err
	}
	server, err := s.instanceAPI.CreateServer(&instance.CreateServerRequest{
		Name:              args.MinecraftResource.GetName(),
		CommercialType:    args.MinecraftResource.GetSize(),
		Image:             "ubuntu_jammy",
		Tags:              []string{"minectl"},
		DynamicIPRequired: scw.BoolPtr(true),
	})
	if err != nil {
		return nil, err
	}

	var mount string
	if args.MinecraftResource.GetVolumeSize() > 0 {
		volume, err := s.instanceAPI.CreateVolume(&instance.CreateVolumeRequest{
			Name:       fmt.Sprintf("%s-vol", args.MinecraftResource.GetName()),
			VolumeType: instance.VolumeVolumeTypeBSSD,
			Size:       scw.SizePtr(scw.Size(args.MinecraftResource.GetVolumeSize()) * scw.GB),
		})
		if err != nil {
			return nil, err
		}
		_, err = s.instanceAPI.AttachVolume(&instance.AttachVolumeRequest{
			VolumeID: volume.Volume.ID,
			ServerID: server.Server.ID,
		})
		if err != nil {
			return nil, err
		}
		mount = "sda"
	}
	userData, err := s.tmpl.GetTemplate(args.MinecraftResource, &minctlTemplate.CreateUpdateTemplateArgs{Mount: mount, Name: minctlTemplate.GetTemplateCloudConfigName(args.MinecraftResource.IsProxyServer())})
	if err != nil {
		return nil, err
	}
	err = s.instanceAPI.SetServerUserData(&instance.SetServerUserDataRequest{
		ServerID: server.Server.ID,
		Key:      "cloud-init",
		Content:  strings.NewReader(userData),
	})
	if err != nil {
		return nil, err
	}

	duration := 2 * time.Second
	err = s.instanceAPI.ServerActionAndWait(&instance.ServerActionAndWaitRequest{
		ServerID:      server.Server.ID,
		Action:        instance.ServerActionPoweron,
		RetryInterval: &duration,
	})
	if err != nil {
		return nil, err
	}

	getServer, err := s.instanceAPI.GetServer(&instance.GetServerRequest{
		ServerID: server.Server.ID,
	})
	if err != nil {
		return nil, err
	}

	return &automation.ResourceResults{
		ID:       server.Server.ID,
		Name:     server.Server.Name,
		Region:   server.Server.Zone.String(),
		PublicIP: getServer.Server.PublicIP.Address.String(),
		Tags:     strings.Join(server.Server.Tags, ","),
	}, err
}

func (s *Scaleway) DeleteServer(id string, args automation.ServerArgs) error {
	getServer, err := s.instanceAPI.GetServer(&instance.GetServerRequest{
		ServerID: id,
	})
	if err != nil {
		return err
	}
	duration := 2 * time.Second
	err = s.instanceAPI.ServerActionAndWait(&instance.ServerActionAndWaitRequest{
		ServerID:      getServer.Server.ID,
		Action:        instance.ServerActionPoweroff,
		RetryInterval: &duration,
	})
	if err != nil {
		return err
	}
	err = s.instanceAPI.DeleteServer(&instance.DeleteServerRequest{
		ServerID: getServer.Server.ID,
	})
	if err != nil {
		return err
	}
	for _, volume := range getServer.Server.Volumes {
		err := s.instanceAPI.DeleteVolume(&instance.DeleteVolumeRequest{
			VolumeID: volume.ID,
		})
		if err != nil {
			return err
		}
	}
	keys, err := s.accountAPI.ListSSHKeys(&account.ListSSHKeysRequest{
		Name: scw.StringPtr(fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName())),
	})
	if err != nil {
		return err
	}
	for _, key := range keys.SSHKeys {
		err := s.accountAPI.DeleteSSHKey(&account.DeleteSSHKeyRequest{
			SSHKeyID: key.ID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Scaleway) ListServer() ([]automation.ResourceResults, error) {
	servers, err := s.instanceAPI.ListServers(&instance.ListServersRequest{
		Tags: []string{common.InstanceTag},
	})
	if err != nil {
		return nil, err
	}
	var result []automation.ResourceResults
	for _, server := range servers.Servers {
		result = append(result, automation.ResourceResults{
			ID:       server.ID,
			PublicIP: server.PublicIP.Address.String(),
			Name:     server.Name,
			Region:   server.Zone.String(),
			Tags:     strings.Join(server.Tags, ","),
		})
	}
	return result, nil
}

func (s *Scaleway) UpdateServer(id string, args automation.ServerArgs) error {
	instance, err := s.instanceAPI.GetServer(&instance.GetServerRequest{
		ServerID: id,
	})
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), instance.Server.PublicIP.Address.String(), "root")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	return nil
}

func (s *Scaleway) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	instance, err := s.instanceAPI.GetServer(&instance.GetServerRequest{
		ServerID: id,
	})
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), instance.Server.PublicIP.Address.String(), "root")
	err = remoteCommand.TransferFile(plugin, filepath.Join(destination, filepath.Base(plugin)), args.MinecraftResource.GetSSHPort())
	if err != nil {
		return err
	}
	_, err = remoteCommand.ExecuteCommand("systemctl restart minecraft.service", args.MinecraftResource.GetSSHPort())
	if err != nil {
		return err
	}
	return nil
}

func (s *Scaleway) GetServer(id string, args automation.ServerArgs) (*automation.ResourceResults, error) {
	instance, err := s.instanceAPI.GetServer(&instance.GetServerRequest{
		ServerID: id,
	})
	if err != nil {
		return nil, err
	}
	return &automation.ResourceResults{
		ID:       instance.Server.ID,
		Name:     instance.Server.Name,
		Region:   instance.Server.Zone.String(),
		PublicIP: instance.Server.PublicIP.Address.String(),
		Tags:     strings.Join(instance.Server.Tags, ","),
	}, err
}
