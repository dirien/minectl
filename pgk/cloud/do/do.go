package do

import (
	"context"
	_ "embed"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/minectl/pgk/automation"
	"github.com/minectl/pgk/common"
	minctlTemplate "github.com/minectl/pgk/template"
	"golang.org/x/oauth2"
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

func (d *DigitalOcean) ListServer() ([]automation.RessourceResults, error) {
	droplets, _, err := d.client.Droplets.ListByTag(context.Background(), common.InstanceTag, nil)
	if err != nil {
		return nil, err
	}
	var result []automation.RessourceResults
	for _, droplet := range droplets {
		ipv4, _ := droplet.PublicIPv4()
		result = append(result, automation.RessourceResults{
			ID:       strconv.Itoa(droplet.ID),
			PublicIP: ipv4,
			Name:     droplet.Name,
			Region:   droplet.Region.Slug,
			Tags:     strings.Join(droplet.Tags, ","),
		})
	}
	return result, nil
}

func (d *DigitalOcean) UpdateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	panic("implement me")
}

func (d *DigitalOcean) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	pubKeyFile, err := ioutil.ReadFile(args.MinecraftServer.GetSSH())
	if err != nil {
		return nil, err
	}
	keyRequest := &godo.KeyCreateRequest{
		Name:      fmt.Sprintf("%s-ssh", args.MinecraftServer.GetName()),
		PublicKey: string(pubKeyFile),
	}
	key, _, err := d.client.Keys.Create(context.Background(), keyRequest)
	if err != nil {
		return nil, err
	}

	volumeRequest := &godo.VolumeCreateRequest{
		Name:           fmt.Sprintf("%s-vol", args.MinecraftServer.GetName()),
		Region:         args.MinecraftServer.GetRegion(),
		Description:    "volume for storing the minecraft data",
		FilesystemType: "ext4",
		SizeGigaBytes:  int64(args.MinecraftServer.GetVolumeSize()),
	}
	volume, _, err := d.client.Storage.CreateVolume(context.Background(), volumeRequest)
	if err != nil {
		return nil, err
	}

	tmpl, err := minctlTemplate.NewTemplateCloudConfig(args.MinecraftServer, "sda")
	if err != nil {
		return nil, err
	}
	userData, err := tmpl.GetTemplate()
	if err != nil {
		return nil, err
	}

	createRequest := &godo.DropletCreateRequest{
		Name: args.MinecraftServer.GetName(),
		SSHKeys: []godo.DropletCreateSSHKey{
			{
				Fingerprint: key.Fingerprint,
			},
		},
		Region: args.MinecraftServer.GetRegion(),
		Size:   args.MinecraftServer.GetSize(),
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-20-04-x64",
		},
		UserData: userData,
		Tags:     []string{common.InstanceTag, args.MinecraftServer.GetEdition()},
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
	for stillCreating {
		droplet, _, err = d.client.Droplets.Get(context.Background(), droplet.ID)
		if err != nil {
			return nil, err
		}
		if droplet.Status == "active" {
			stillCreating = false
		} else {
			time.Sleep(2 * time.Second)
		}
	}
	ipv4, _ := droplet.PublicIPv4()

	return &automation.RessourceResults{
		ID:       strconv.Itoa(droplet.ID),
		Name:     droplet.Name,
		Region:   droplet.Region.Slug,
		PublicIP: ipv4,
		Tags:     strings.Join(droplet.Tags, ","),
	}, err
}

func (d *DigitalOcean) DeleteServer(id string, args automation.ServerArgs) error {
	list, _, err := d.client.Keys.List(context.Background(), nil)
	if err != nil {
		return err
	}
	for _, key := range list {
		if key.Name == fmt.Sprintf("%s-ssh", args.MinecraftServer.GetName()) {
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
			time.Sleep(15 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	volumes, _, err := d.client.Storage.ListVolumes(context.Background(), &godo.ListVolumeParams{
		Name: fmt.Sprintf("%s-vol", args.MinecraftServer.GetName()),
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
