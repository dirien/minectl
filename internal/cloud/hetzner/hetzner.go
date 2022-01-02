package hetzner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/minectl/internal/update"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/minectl/internal/automation"
	"github.com/minectl/internal/common"
	minctlTemplate "github.com/minectl/internal/template"
)

type Hetzner struct {
	client *hcloud.Client
	tmpl   *minctlTemplate.Template
}

func NewHetzner(apiKey string) (*Hetzner, error) {
	client := hcloud.NewClient(hcloud.WithToken(apiKey))
	tmpl, err := minctlTemplate.NewTemplateCloudConfig()
	if err != nil {
		return nil, err
	}
	hetzner := &Hetzner{
		client: client,
		tmpl:   tmpl,
	}
	return hetzner, nil
}

func (h *Hetzner) CreateServer(args automation.ServerArgs) (*automation.ResourceResults, error) {
	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSH()))
	if err != nil {
		return nil, err
	}
	key, _, err := h.client.SSHKey.Create(context.Background(), hcloud.SSHKeyCreateOpts{
		Name:      fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()),
		PublicKey: string(pubKeyFile),
	})
	if err != nil {
		return nil, err
	}

	location, _, err := h.client.Location.Get(context.Background(), args.MinecraftResource.GetRegion())
	if err != nil {
		return nil, err
	}

	var volume hcloud.VolumeCreateResult
	var mount string
	if args.MinecraftResource.GetVolumeSize() > 0 {
		volume, _, err = h.client.Volume.Create(context.Background(), hcloud.VolumeCreateOpts{
			Name:     fmt.Sprintf("%s-vol", args.MinecraftResource.GetName()),
			Size:     args.MinecraftResource.GetVolumeSize(),
			Location: location,
			Format:   hcloud.String("ext4"),
		})
		if err != nil {
			return nil, err
		}
		mount = "sdb"
	}
	userData, err := h.tmpl.GetTemplate(args.MinecraftResource, &minctlTemplate.CreateUpdateTemplateArgs{Mount: mount, Name: minctlTemplate.GetTemplateCloudConfigName(args.MinecraftResource.IsProxyServer())})
	if err != nil {
		return nil, err
	}
	image, _, err := h.client.Image.GetByName(context.Background(), "ubuntu-20.04")
	if err != nil {
		return nil, err
	}

	plan, _, err := h.client.ServerType.GetByName(context.Background(), args.MinecraftResource.GetSize())
	if err != nil {
		return nil, err
	}

	requestOpts := hcloud.ServerCreateOpts{
		Name:       args.MinecraftResource.GetName(),
		ServerType: plan,
		Image:      image,
		Location:   location,
		SSHKeys:    []*hcloud.SSHKey{key},
		UserData:   userData,
		Labels:     map[string]string{common.InstanceTag: "true", args.MinecraftResource.GetEdition(): "true"},
	}

	if args.MinecraftResource.GetVolumeSize() > 0 {
		requestOpts.Volumes = []*hcloud.Volume{volume.Volume}
		requestOpts.Automount = hcloud.Bool(true)
	}

	serverCreateReq, _, err := h.client.Server.Create(context.Background(), requestOpts)
	if err != nil {
		return nil, err
	}
	server := serverCreateReq.Server
	stillCreating := true

	for stillCreating {
		server, _, err := h.client.Server.GetByID(context.Background(), server.ID)
		if err != nil {
			return nil, err
		}
		if server.Status == hcloud.ServerStatusRunning {
			stillCreating = false
		} else {
			time.Sleep(2 * time.Second)
		}
	}
	return &automation.ResourceResults{
		ID:       strconv.Itoa(server.ID),
		Name:     server.Name,
		Region:   server.Datacenter.Location.Name,
		PublicIP: server.PublicNet.IPv4.IP.String(),
		Tags:     hetznerLabelsToTags(server.Labels),
	}, err
}

func (h *Hetzner) DeleteServer(id string, args automation.ServerArgs) error {
	serverID, _ := strconv.Atoi(id)
	server, _, err := h.client.Server.GetByID(context.Background(), serverID)
	if err != nil {
		return err
	}

	volume, _, err := h.client.Volume.Get(context.Background(), fmt.Sprintf("%s-vol", args.MinecraftResource.GetName()))
	if err != nil {
		return err
	}

	if volume != nil {
		res, _, err := h.client.Volume.Detach(context.Background(), volume)
		if err != nil {
			return err
		}
		stillDetatching := true
		for stillDetatching {
			action, _, err := h.client.Action.GetByID(context.Background(), res.ID)
			if err != nil {
				return err
			}
			if action.Status == "success" {
				stillDetatching = false
			} else {
				time.Sleep(2 * time.Second)
			}
		}
		_, err = h.client.Volume.Delete(context.Background(), volume)
		if err != nil {
			return err
		}
	}
	_, err = h.client.Server.Delete(context.Background(), server)
	if err != nil {
		return err
	}

	key, _, err := h.client.SSHKey.Get(context.Background(), fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()))
	if err != nil {
		return err
	}
	_, err = h.client.SSHKey.Delete(context.Background(), key)
	if err != nil {
		return err
	}
	return nil
}

func hetznerLabelsToTags(label map[string]string) string {
	var tags []string
	for key := range label {
		tags = append(tags, key)
	}
	return strings.Join(tags, ",")
}

func (h *Hetzner) ListServer() ([]automation.ResourceResults, error) {
	servers, err := h.client.Server.All(context.Background())
	if err != nil {
		return nil, err
	}
	var result []automation.ResourceResults
	for _, server := range servers {
		for key := range server.Labels {
			if key == common.InstanceTag {
				result = append(result, automation.ResourceResults{
					ID:       strconv.Itoa(server.ID),
					Name:     server.Name,
					Region:   server.Datacenter.Location.Name,
					PublicIP: server.PublicNet.IPv4.IP.String(),
					Tags:     hetznerLabelsToTags(server.Labels),
				})
			}
		}
	}
	return result, nil
}

func (h *Hetzner) UpdateServer(id string, args automation.ServerArgs) error {
	intID, _ := strconv.Atoi(id)
	instance, _, err := h.client.Server.GetByID(context.Background(), intID)
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSH(), instance.PublicNet.IPv4.IP.String(), "root")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	return nil
}

func (h *Hetzner) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	intID, _ := strconv.Atoi(id)
	instance, _, err := h.client.Server.GetByID(context.Background(), intID)
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSH(), instance.PublicNet.IPv4.IP.String(), "root")
	err = remoteCommand.TransferFile(plugin, filepath.Join(destination, filepath.Base(plugin)))
	if err != nil {
		return err
	}
	_, err = remoteCommand.ExecuteCommand("systemctl restart minecraft.service")
	if err != nil {
		return err
	}
	return nil
}

func (h *Hetzner) GetServer(id string, args automation.ServerArgs) (*automation.ResourceResults, error) {
	intID, _ := strconv.Atoi(id)
	instance, _, err := h.client.Server.GetByID(context.Background(), intID)
	if err != nil {
		return nil, err
	}
	return &automation.ResourceResults{
		ID:       strconv.Itoa(instance.ID),
		Name:     instance.Name,
		Region:   instance.Datacenter.Location.Name,
		PublicIP: instance.PublicNet.IPv4.IP.String(),
		Tags:     hetznerLabelsToTags(instance.Labels),
	}, err
}
