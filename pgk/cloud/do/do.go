package do

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/digitalocean/godo"
	"github.com/minectl/pgk/automation"
	"github.com/minectl/pgk/cloud"
	"github.com/minectl/pgk/common"
	"golang.org/x/oauth2"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type DigitalOcean struct {
	client *godo.Client
}

// TokenSource contains an access token
type TokenSource struct {
	AccessToken string
}

// Token returns an oauth2 token
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func NewDigitalOcean(APIKey string) (*DigitalOcean, error) {
	tokenSource := &TokenSource{
		AccessToken: APIKey,
	}
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)

	do := &DigitalOcean{
		client: client,
	}
	return do, nil
}

func (d *DigitalOcean) ListServer(args automation.ServerArgs) (*[]automation.RessourceResults, error) {
	panic("implement me")
}

func (d *DigitalOcean) UpdateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	panic("implement me")
	return nil, nil
}

func (d *DigitalOcean) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	pubKeyFile, err := ioutil.ReadFile(args.SSH)
	if err != nil {
		return nil, err
	}
	keyRequest := &godo.KeyCreateRequest{
		Name:      fmt.Sprintf("%s-ssh", args.StackName),
		PublicKey: string(pubKeyFile),
	}
	key, _, err := d.client.Keys.Create(context.Background(), keyRequest)
	if err != nil {
		return nil, err
	}

	volumeRequest := &godo.VolumeCreateRequest{
		Name:           fmt.Sprintf("%s-vol", args.StackName),
		Region:         args.Region,
		Description:    "volume for storing the minecraft data",
		FilesystemType: "ext4",
		SizeGigaBytes:  int64(args.VolumeSize),
	}
	volume, _, err := d.client.Storage.CreateVolume(context.Background(), volumeRequest)
	if err != nil {
		return nil, err
	}

	cloud.CloudConfig = strings.Replace(cloud.CloudConfig, "vdc", "sda", -1)
	cloud.CloudConfig = cloud.ReplaceServerProperties(cloud.CloudConfig, args.Properties)
	createRequest := &godo.DropletCreateRequest{
		Name: args.StackName,
		SSHKeys: []godo.DropletCreateSSHKey{
			{
				Fingerprint: key.Fingerprint,
			},
		},
		Region: args.Region,
		Size:   args.Size,
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-20-04-x64",
		},
		UserData: cloud.CloudConfig,
		Volumes: []godo.DropletCreateVolume{
			{
				ID: volume.ID,
			},
		},
	}
	droplet, _, err := d.client.Droplets.Create(context.Background(), createRequest)
	if err != nil {
		return nil, err
	}

	stillCreating := true
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Prefix = fmt.Sprintf("Creating droplet (%s)... ", common.Green(droplet.Name))
	s.FinalMSG = fmt.Sprintf("\nDroplet (%s) created\n", common.Green(droplet.Name))
	s.Start()

	for stillCreating {
		droplet, _, err = d.client.Droplets.Get(context.Background(), droplet.ID)
		if err != nil {
			return nil, err
		}
		if droplet.Status == "active" {
			stillCreating = false
			s.Stop()
		} else {
			time.Sleep(2 * time.Second)
		}
	}
	ipv4, _ := droplet.PublicIPv4()

	return &automation.RessourceResults{
		ID:       strconv.Itoa(droplet.ID),
		PublicIP: ipv4,
	}, err
}

func (d *DigitalOcean) DeleteServer(id string, args automation.ServerArgs) error {
	common.PrintMixedGreen("Delete droplet (%s)... ", id)
	list, _, err := d.client.Keys.List(context.Background(), nil)
	if err != nil {
		return err
	}
	for _, key := range list {
		if key.Name == fmt.Sprintf("%s-ssh", args.StackName) {
			_, err := d.client.Keys.DeleteByID(context.Background(), key.ID)
			if err != nil {
				return err
			}
		}
	}
	intID, _ := strconv.Atoi(id)
	_, err = d.client.Droplets.Delete(context.Background(), intID)
	if err != nil {
		return err
	}
	stillDeleting := true

	for stillDeleting {
		_, _, err := d.client.Droplets.Get(context.Background(), intID)
		if err != nil {
			stillDeleting = false
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	volumes, _, err := d.client.Storage.ListVolumes(context.Background(), &godo.ListVolumeParams{
		Name: fmt.Sprintf("%s-vol", args.StackName),
	})
	if err != nil {
		return err
	}
	for _, volume := range volumes {
		_, err = d.client.Storage.DeleteVolume(context.Background(), volume.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
