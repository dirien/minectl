package scaleway

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/minectl/pgk/automation"
	"github.com/minectl/pgk/common"
	minctlTemplate "github.com/minectl/pgk/template"
	"github.com/scaleway/scaleway-sdk-go/api/account/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"io/ioutil"
	"strings"
	"time"
)

type Scaleway struct {
	instanceAPI *instance.API
	accountAPI  *account.API
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

	return &Scaleway{
		instanceAPI: instance.NewAPI(client),
		accountAPI:  account.NewAPI(client),
	}, nil
}

func (s Scaleway) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {

	pubKeyFile, err := ioutil.ReadFile(args.MinecraftServer.GetSSH())
	_, err = s.accountAPI.CreateSSHKey(&account.CreateSSHKeyRequest{
		Name:      fmt.Sprintf("%s-ssh", args.MinecraftServer.GetName()),
		PublicKey: string(pubKeyFile),
	})
	if err != nil {
		return nil, err
	}
	server, err := s.instanceAPI.CreateServer(&instance.CreateServerRequest{
		Name:              args.MinecraftServer.GetName(),
		CommercialType:    args.MinecraftServer.GetSize(),
		Image:             "ubuntu_focal",
		Tags:              []string{"minectl"},
		DynamicIPRequired: scw.BoolPtr(true),
	})

	if err != nil {
		return nil, err
	}
	tmpl, err := minctlTemplate.NewTemplateCloudConfig(args.MinecraftServer, "sda")
	if err != nil {
		return nil, err
	}
	userData, err := tmpl.GetTemplate()
	err = s.instanceAPI.SetServerUserData(&instance.SetServerUserDataRequest{
		ServerID: server.Server.ID,
		Key:      "cloud-init",
		Content:  strings.NewReader(userData),
	})

	if err != nil {
		return nil, err
	}

	volume, err := s.instanceAPI.CreateVolume(&instance.CreateVolumeRequest{
		Name:       fmt.Sprintf("%s-vol", args.MinecraftServer.GetName()),
		VolumeType: instance.VolumeVolumeTypeBSSD,
		Size:       scw.SizePtr(scw.Size(args.MinecraftServer.GetVolumeSize()) * scw.GB),
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

	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Prefix = fmt.Sprintf("🏗 Creating instance (%s)... ", common.Green(server.Server.Name))
	spinner.FinalMSG = fmt.Sprintf("\n✅ Instance (%s) created\n", common.Green(server.Server.Name))
	spinner.Start()
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
	spinner.Stop()

	return &automation.RessourceResults{
		ID:       server.Server.ID,
		Name:     server.Server.Name,
		Region:   server.Server.Zone.String(),
		PublicIP: getServer.Server.PublicIP.Address.String(),
		Tags:     strings.Join(server.Server.Tags, ","),
	}, err
}

func (s Scaleway) DeleteServer(id string, args automation.ServerArgs) error {
	common.PrintMixedGreen("🗑 Delete instance (%s)... ", id)
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
		Name: scw.StringPtr(fmt.Sprintf("%s-ssh", args.MinecraftServer.GetName())),
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

func (s Scaleway) ListServer() ([]automation.RessourceResults, error) {
	servers, err := s.instanceAPI.ListServers(&instance.ListServersRequest{
		Tags: []string{common.InstanceTag},
	})
	if err != nil {
		return nil, err
	}
	var result []automation.RessourceResults
	for _, server := range servers.Servers {
		result = append(result, automation.RessourceResults{
			ID:       server.ID,
			PublicIP: server.PublicIP.Address.String(),
			Name:     server.Name,
			Region:   server.Zone.String(),
			Tags:     strings.Join(server.Tags, ","),
		})
	}
	return result, nil
}

func (s Scaleway) UpdateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	panic("implement me")
}
