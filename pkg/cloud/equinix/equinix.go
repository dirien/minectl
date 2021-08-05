package equinix

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/minectl/pkg/update"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/minectl/pkg/automation"
	"github.com/minectl/pkg/common"
	minctlTemplate "github.com/minectl/pkg/template"
	"github.com/packethost/packngo"
)

type Equinix struct {
	client  *packngo.Client
	project string
	tmpl    *minctlTemplate.Template
}

func NewEquinix(APIKey, project string) (*Equinix, error) {

	httpClient := retryablehttp.NewClient().HTTPClient
	tmpl, err := minctlTemplate.NewTemplateBash()
	if err != nil {
		return nil, err
	}
	return &Equinix{
		client:  packngo.NewClientWithAuth("", APIKey, httpClient),
		project: project,
		tmpl:    tmpl,
	}, nil
}

func (e *Equinix) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	pubKeyFile, err := ioutil.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftServer.GetSSH()))
	if err != nil {
		return nil, err
	}
	key, _, err := e.client.SSHKeys.Create(&packngo.SSHKeyCreateRequest{
		Label:     fmt.Sprintf("%s-ssh", args.MinecraftServer.GetName()),
		ProjectID: e.project,
		Key:       string(pubKeyFile),
	})
	if err != nil {
		return nil, err
	}

	userData, err := e.tmpl.GetTemplate(args.MinecraftServer, "", minctlTemplate.TemplateBash)
	if err != nil {
		return nil, err
	}

	server, _, err := e.client.Devices.Create(&packngo.DeviceCreateRequest{
		Hostname:       args.MinecraftServer.GetName(),
		ProjectID:      e.project,
		OS:             "ubuntu_20_04",
		Plan:           args.MinecraftServer.GetSize(),
		Tags:           []string{common.InstanceTag, args.MinecraftServer.GetEdition()},
		ProjectSSHKeys: []string{key.ID},
		UserData:       userData,
		BillingCycle:   "hourly",
		Metro:          args.MinecraftServer.GetRegion(),
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

	return &automation.RessourceResults{
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
		if key.Label == fmt.Sprintf("%s-ssh", args.MinecraftServer.GetName()) {
			_, err := e.client.SSHKeys.Delete(key.ID)
			if err != nil {
				return err
			}
		}
	}
	instances, _, err := e.client.Devices.List(e.project, &packngo.ListOptions{
		Search: args.MinecraftServer.GetName(),
	})
	if err != nil {
		return err
	}
	for _, instance := range instances {
		if instance.Hostname == args.MinecraftServer.GetName() {
			_, err = e.client.Devices.Delete(instance.ID, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Equinix) ListServer() ([]automation.RessourceResults, error) {
	list, _, err := e.client.Devices.List(e.project, &packngo.ListOptions{
		Search: common.InstanceTag,
	})
	if err != nil {
		return nil, err
	}
	var result []automation.RessourceResults
	for _, server := range list {
		result = append(result, automation.RessourceResults{
			ID:       server.ID,
			Name:     server.Hostname,
			Region:   server.Metro.Code,
			PublicIP: getIP4(&server),
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

	remoteCommand := update.NewRemoteServer(args.MinecraftServer.GetSSH(), getIP4(instance), "root")
	err = remoteCommand.UpdateServer(args.MinecraftServer)
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

	remoteCommand := update.NewRemoteServer(args.MinecraftServer.GetSSH(), getIP4(instance), "root")
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

func (e *Equinix) GetServer(id string, _ automation.ServerArgs) (*automation.RessourceResults, error) {
	instance, _, err := e.client.Devices.Get(id, nil)
	if err != nil {
		return nil, err
	}

	return &automation.RessourceResults{
		ID:       instance.ID,
		Name:     instance.Hostname,
		Region:   instance.Metro.Code,
		PublicIP: getIP4(instance),
		Tags:     strings.Join(instance.Tags, ","),
	}, err
}
