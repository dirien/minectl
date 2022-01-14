package upcloud

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/minectl/internal/update"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/service"
	"github.com/minectl/internal/automation"
	"github.com/minectl/internal/common"
	minctlTemplate "github.com/minectl/internal/template"
)

type Upcloud struct {
	service *service.Service
	tmpl    *minctlTemplate.Template
}

func NewUpcloud(username, password string) (*Upcloud, error) {
	client := client.New(username, password)

	tmpl, err := minctlTemplate.NewTemplateBash()
	if err != nil {
		return nil, err
	}

	s := service.New(client)

	do := &Upcloud{
		service: s,
		tmpl:    tmpl,
	}
	return do, nil
}

func (u *Upcloud) CreateServer(args automation.ServerArgs) (*automation.ResourceResults, error) {
	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSHKeyFolder()))
	if err != nil {
		return nil, err
	}

	script, err := u.tmpl.GetTemplate(args.MinecraftResource, &minctlTemplate.CreateUpdateTemplateArgs{Name: minctlTemplate.GetTemplateBashName(args.MinecraftResource.IsProxyServer())})
	if err != nil {
		return nil, err
	}

	fmt.Println("Creating server")
	details, err := u.service.CreateServer(&request.CreateServerRequest{
		Hostname: args.MinecraftResource.GetName(),
		Title:    args.MinecraftResource.GetName(),
		Zone:     args.MinecraftResource.GetRegion(),
		Plan:     args.MinecraftResource.GetSize(),
		LoginUser: &request.LoginUser{
			Username: "root",
			SSHKeys:  request.SSHKeySlice{string(pubKeyFile)},
		},
		StorageDevices: []request.CreateServerStorageDevice{
			{
				Action:  request.CreateServerStorageDeviceActionClone,
				Storage: "01000000-0000-4000-8000-000030200200",
				Title:   "Ubuntu from a template",
				Size:    args.MinecraftResource.GetVolumeSize(),
				Tier:    upcloud.StorageTierMaxIOPS,
			},
		},
		Networking: &request.CreateServerNetworking{
			Interfaces: []request.CreateServerInterface{
				{
					IPAddresses: []request.CreateServerIPAddress{
						{
							Family: upcloud.IPAddressFamilyIPv4,
						},
					},
					Type: upcloud.IPAddressAccessPublic,
				},
				{
					IPAddresses: []request.CreateServerIPAddress{
						{
							Family: upcloud.IPAddressFamilyIPv4,
						},
					},
					Type: upcloud.NetworkTypeUtility,
				},
			},
		},
		UserData: script,
	})
	if err != nil {
		zap.S().Errorw("Unable to create server", "error", err)
		return nil, err
	}

	if len(details.UUID) == 0 {
		zap.S().Info("UUID missing or invalid")
		return nil, errors.New("UUID too short")
	}
	details, err = u.service.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         details.UUID,
		DesiredState: upcloud.ServerStateStarted,
		Timeout:      1 * time.Minute,
	})
	if err != nil {
		zap.S().Errorw("Unable to wait for server", "error", err)
		return nil, err
	}

	zap.S().Infow("Created server", "details", details)

	zap.S().Info("Creating firewall rules")
	_, err = u.service.CreateFirewallRule(&request.CreateFirewallRuleRequest{
		ServerUUID: details.UUID,
		FirewallRule: upcloud.FirewallRule{
			Direction:            upcloud.FirewallRuleDirectionIn,
			Action:               upcloud.FirewallRuleActionAccept,
			Family:               upcloud.IPAddressFamilyIPv4,
			Protocol:             upcloud.FirewallRuleProtocolUDP,
			Position:             1,
			DestinationPortStart: "19132",
			DestinationPortEnd:   "19133",
			Comment:              "Accept all UDP input on IPv4",
		},
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create firewall rule: %#v", err)
		zap.S().Errorw("unable to create firewall rule", "error", err)
		return nil, err
	}

	zap.S().Infow("Everything created successfully")
	return &automation.ResourceResults{
		ID:       details.UUID,
		Name:     details.Hostname,
		Region:   details.Zone,
		PublicIP: details.IPAddresses[1].Address,
		Tags:     strings.Join([]string{common.InstanceTag, args.MinecraftResource.GetEdition()}, ","),
	}, err
}

func (u *Upcloud) DeleteServer(id string, args automation.ServerArgs) error {
	instance, err := u.service.GetServerDetails(&request.GetServerDetailsRequest{
		UUID: id,
	})
	if err != nil {
		return err
	}

	if instance.State != upcloud.ServerStateStopped {
		_, err = u.service.StopServer(&request.StopServerRequest{
			UUID:     id,
			StopType: upcloud.StopTypeHard,
		})

		if err != nil {
			return err
		}

		_, err = u.service.WaitForServerState(&request.WaitForServerStateRequest{
			UUID:         id,
			DesiredState: upcloud.ServerStateStopped,
			Timeout:      3 * time.Minute,
		})

		if err != nil {
			return err
		}
	}

	err = u.service.DeleteServerAndStorages(&request.DeleteServerAndStoragesRequest{
		UUID: id,
	})
	if err != nil {
		return err
	}
	zap.S().Infow("Upcloud delete instance", "id", id)

	return nil
}

func (u *Upcloud) ListServer() ([]automation.ResourceResults, error) {
	var result []automation.ResourceResults
	instances, err := u.service.GetServers()
	if err != nil {
		return nil, err
	}
	for _, instance := range instances.Servers {
		details, err := u.service.GetServerDetails(&request.GetServerDetailsRequest{
			UUID: instance.UUID,
		})
		if err != nil {
			zap.S().Infow("Upcloud could not retrieve details for instance", "instance", instance.UUID)
			return nil, err
		}
		result = append(result, automation.ResourceResults{
			ID:       instance.UUID,
			PublicIP: details.IPAddresses[1].Address,
			Name:     instance.Hostname,
			Region:   details.Zone,
			Tags:     strings.Join(instance.Tags, ","),
		})
	}
	if len(result) > 0 {
		zap.S().Infow("Upcloud list all instances", "list", result)
	} else {
		zap.S().Infow("No minectl instances found")
	}
	return result, nil
}

func (u *Upcloud) UpdateServer(id string, args automation.ServerArgs) error {
	instance, err := u.service.GetServerDetails(&request.GetServerDetailsRequest{
		UUID: id,
	})
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), instance.IPAddresses[1].Address, "root")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	zap.S().Infow("Upcloud minectl server updated", "instance", instance)
	return nil
}

func (u *Upcloud) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	instance, err := u.service.GetServerDetails(&request.GetServerDetailsRequest{
		UUID: id,
	})
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), instance.IPAddresses[1].Address, "root")
	err = remoteCommand.TransferFile(plugin, filepath.Join(destination, filepath.Base(plugin)), args.MinecraftResource.GetSSHPort())
	if err != nil {
		return err
	}
	_, err = remoteCommand.ExecuteCommand("systemctl restart minecraft.service", args.MinecraftResource.GetSSHPort())
	if err != nil {
		return err
	}
	zap.S().Infow("Minecraft plugin uploaded", "plugin", plugin, "instance", instance)
	return nil
}

func (u *Upcloud) GetServer(id string, _ automation.ServerArgs) (*automation.ResourceResults, error) {
	instance, err := u.service.GetServerDetails(&request.GetServerDetailsRequest{
		UUID: id,
	})
	if err != nil {
		return nil, err
	}

	return &automation.ResourceResults{
		ID:       instance.UUID,
		Name:     instance.Hostname,
		Region:   instance.Zone,
		PublicIP: instance.IPAddresses[1].Address,
		Tags:     strings.Join(instance.Tags, ","),
	}, err
}
