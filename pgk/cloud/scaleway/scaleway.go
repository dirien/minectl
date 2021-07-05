package scaleway

import (
	"fmt"
	"github.com/minectl/pgk/automation"
	"github.com/minectl/pgk/cloud"
	"github.com/scaleway/scaleway-sdk-go/api/account/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"io/ioutil"
	"log"
	"strings"
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
	log.Printf("Provisioning host with Scaleway\n")

	pubKeyFile, err := ioutil.ReadFile(args.SSH)
	log.Printf("CreateSSHKeyn")
	_, err = s.accountAPI.CreateSSHKey(&account.CreateSSHKeyRequest{
		Name:      fmt.Sprintf("%s-ssh", args.StackName),
		PublicKey: string(pubKeyFile),
	})
	if err != nil {
		return nil, err
	}

	log.Printf("CreateServer")
	server, err := s.instanceAPI.CreateServer(&instance.CreateServerRequest{
		Name:              args.StackName,
		CommercialType:    args.Size,
		Image:             "ubuntu_focal",
		Tags:              []string{"minecraft"},
		DynamicIPRequired: scw.BoolPtr(true),
	})

	if err != nil {
		return nil, err
	}

	cloud.CloudConfig = strings.Replace(cloud.CloudConfig, "vdc", "sda", -1)
	cloud.CloudConfig = cloud.ReplaceServerProperties(cloud.CloudConfig, args.Properties)

	log.Printf("SetServerUserData")
	err = s.instanceAPI.SetServerUserData(&instance.SetServerUserDataRequest{
		ServerID: server.Server.ID,
		Key:      "cloud-init",
		Content:  strings.NewReader(cloud.CloudConfig),
	})

	if err != nil {
		return nil, err
	}

	log.Printf("CreateVolume")
	volume, err := s.instanceAPI.CreateVolume(&instance.CreateVolumeRequest{
		Name:       fmt.Sprintf("%s-vol", args.StackName),
		VolumeType: instance.VolumeVolumeTypeBSSD,
		Size:       scw.SizePtr(scw.Size(args.VolumeSize) * scw.GB),
	})
	if err != nil {
		return nil, err
	}
	log.Printf("AttachVolume")
	s.instanceAPI.AttachVolume(&instance.AttachVolumeRequest{
		VolumeID: volume.Volume.ID,
		ServerID: server.Server.ID,
	})

	/*
		log.Printf("ServerAction")
		_, err = s.instanceAPI.ServerAction(&instance.ServerActionRequest{
			ServerID: server.Server.ID,
			Action:   instance.ServerActionPoweron,
		})

		if err != nil {
			return nil, err
		}*/
	return nil, nil
}

func (s Scaleway) DeleteServer(id string, args automation.ServerArgs) error {
	panic("implement me")
}

func (s Scaleway) ListServer(args automation.ServerArgs) (*[]automation.RessourceResults, error) {
	panic("implement me")
}

func (s Scaleway) UpdateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	panic("implement me")
}
