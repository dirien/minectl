package equinix

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minectl/internal/update"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/minectl/internal/automation"
	"github.com/minectl/internal/common"
	minctlTemplate "github.com/minectl/internal/template"
	"github.com/packethost/packngo"
)

type Equinix struct {
	client  *packngo.Client
	project string
	tmpl    *minctlTemplate.Template
}

func NewEquinix(apiKey, project string) (*Equinix, error) {
	httpClient := retryablehttp.NewClient().HTTPClient
	tmpl, err := minctlTemplate.NewTemplateBash()
	if err != nil {
		return nil, err
	}
	return &Equinix{
		client:  packngo.NewClientWithAuth("", apiKey, httpClient),
		project: project,
		tmpl:    tmpl,
	}, nil
}

func (e *Equinix) CreateServer(args automation.ServerArgs) (*automation.ResourceResults, error) {
	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSHKeyFolder()))
	if err != nil {
		return nil, err
	}
	key, _, err := e.client.SSHKeys.Create(&packngo.SSHKeyCreateRequest{
		Label:     fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()),
		ProjectID: e.project,
		Key:       string(pubKeyFile),
	})
	if err != nil {
		return nil, err
	}

	userData, err := e.tmpl.GetTemplate(args.MinecraftResource, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.GetTemplateBashName(args.MinecraftResource.IsProxyServer())})
	if err != nil {
		return nil, err
	}

	server, _, err := e.client.Devices.Create(&packngo.DeviceCreateRequest{
		Hostname:       args.MinecraftResource.GetName(),
		ProjectID:      e.project,
		OS:             "ubuntu_22_04",
		Plan:           args.MinecraftResource.GetSize(),
		Tags:           []string{common.InstanceTag, args.MinecraftResource.GetEdition()},
		ProjectSSHKeys: []string{key.ID},
		UserData:       userData,
		BillingCycle:   "hourly",
		Metro:          args.MinecraftResource.GetRegion(),
		SpotInstance:   false,
	})
	if err != nil {
		return nil, err
	}
	stillCreating := true
	for stillCreating {
		server, _, err = e.client.Devices.Get(server.ID, nil)
		if err != nil {
			return nil, err
		}

		if server.State == "active" {
			stillCreating = false
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	return &automation.ResourceResults{
		ID:       server.ID,
		Name:     server.Hostname,
		Region:   server.Metro.Code,
		PublicIP: getIP4(server),
		Tags:     strings.Join(server.Tags, ","),
	}, err
}

func (e *Equinix) DeleteServer(id string, args automation.ServerArgs) error {
	keys, _, err := e.client.SSHKeys.ProjectList(e.project)
	if err != nil {
		return err
	}
	for _, key := range keys {
		if key.Label == fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()) {
			_, err := e.client.SSHKeys.Delete(key.ID)
			if err != nil {
				return err
			}
		}
	}
	instances, _, err := e.client.Devices.List(e.project, &packngo.ListOptions{
		Search: args.MinecraftResource.GetName(),
	})
	if err != nil {
		return err
	}
	for _, instance := range instances {
		if instance.Hostname == args.MinecraftResource.GetName() {
			_, err = e.client.Devices.Delete(instance.ID, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Equinix) ListServer() ([]automation.ResourceResults, error) {
	list, _, err := e.client.Devices.List(e.project, &packngo.ListOptions{
		Search: common.InstanceTag,
	})
	if err != nil {
		return nil, err
	}
	var result []automation.ResourceResults
	for i, server := range list {
		result = append(result, automation.ResourceResults{
			ID:       server.ID,
			Name:     server.Hostname,
			Region:   server.Metro.Code,
			PublicIP: getIP4(&list[i]),
			Tags:     strings.Join(server.Tags, ","),
		})
	}
	return result, nil
}

func (e *Equinix) UpdateServer(id string, args automation.ServerArgs) error {
	instance, _, err := e.client.Devices.Get(id, nil)
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), getIP4(instance), "root")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	return nil
}

func getIP4(server *packngo.Device) string {
	ip4 := ""
	for _, network := range server.Network {
		if network.Public {
			ip4 = network.IpAddressCommon.Address
			break
		}
	}
	return ip4
}

func (e *Equinix) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	instance, _, err := e.client.Devices.Get(id, nil)
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), getIP4(instance), "root")
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

func (e *Equinix) GetServer(id string, _ automation.ServerArgs) (*automation.ResourceResults, error) {
	instance, _, err := e.client.Devices.Get(id, nil)
	if err != nil {
		return nil, err
	}

	return &automation.ResourceResults{
		ID:       instance.ID,
		Name:     instance.Hostname,
		Region:   instance.Metro.Code,
		PublicIP: getIP4(instance),
		Tags:     strings.Join(instance.Tags, ","),
	}, err
}
