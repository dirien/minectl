package ovh

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minectl/pkg/update"

	ovhsdk "github.com/dirien/ovh-go-sdk/pkg/sdk"

	"github.com/minectl/pkg/automation"
	"github.com/minectl/pkg/common"
	minctlTemplate "github.com/minectl/pkg/template"
)

type OVHcloud struct {
	client *ovhsdk.OVHcloud
	tmpl   *minctlTemplate.Template
}

func NewOVHcloud(endpoint, appKey, appSecret, consumerKey, serviceName, region string) (*OVHcloud, error) {
	client, err := ovhsdk.NewOVHClient(endpoint, appKey, appSecret, consumerKey, region, serviceName)
	if err != nil {
		return nil, err
	}
	tmpl, err := minctlTemplate.NewTemplateBash()
	if err != nil {
		return nil, err
	}
	return &OVHcloud{
		client: client,
		tmpl:   tmpl,
	}, nil
}

func createOVHID(instanceName, label string) (id string) {
	return fmt.Sprintf("%s|%s", instanceName, label)
}

func getOVHFieldsFromID(id string) (label string, err error) {
	fields := strings.Split(id, "|")
	if len(fields) == 3 {
		label = strings.Join([]string{fields[1], fields[2]}, ",")
	} else {
		err = fmt.Errorf("could not get fields from custom ID: fields: %v", fields)
		return "", err
	}
	return label, nil
}

func (o *OVHcloud) CreateServer(args automation.ServerArgs) (*automation.ResourceResults, error) {
	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSH()))
	if err != nil {
		return nil, err
	}

	key, err := o.client.CreateSSHKey(context.Background(), ovhsdk.SSHKeyCreateOptions{
		Name:      fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()),
		PublicKey: string(pubKeyFile),
	})
	if err != nil {
		return nil, err
	}

	image, err := o.client.GetImage(context.Background(), "Ubuntu 20.04", args.MinecraftResource.GetRegion())
	if err != nil {
		return nil, err
	}

	flavor, err := o.client.GetFlavor(context.Background(), args.MinecraftResource.GetSize(), args.MinecraftResource.GetRegion())
	if err != nil {
		return nil, err
	}

	var mount string
	if args.MinecraftResource.GetVolumeSize() > 0 {
		mount = "sdb"
	}
	userData, err := o.tmpl.GetTemplate(args.MinecraftResource, &minctlTemplate.CreateUpdateTemplateArgs{Mount: mount, Name: minctlTemplate.GetTemplateBashName(args.MinecraftResource.IsProxyServer())})
	if err != nil {
		return nil, err
	}

	instance, err := o.client.CreateInstance(context.Background(), ovhsdk.InstanceCreateOptions{
		Name:           createOVHID(args.MinecraftResource.GetName(), strings.Join([]string{common.InstanceTag, args.MinecraftResource.GetEdition()}, "|")),
		Region:         args.MinecraftResource.GetRegion(),
		SSHKeyID:       key.ID,
		FlavorID:       flavor.ID,
		ImageID:        image.ID,
		MonthlyBilling: false,
		UserData:       userData,
	})
	if err != nil {
		return nil, err
	}
	stillCreating := true
	for stillCreating {
		instance, err = o.client.GetInstance(context.Background(), instance.ID)
		if err != nil {
			return nil, err
		}
		if instance.Status == ovhsdk.InstanceActive {
			stillCreating = false
			time.Sleep(2 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	if args.MinecraftResource.GetVolumeSize() > 0 {
		volume, err := o.client.CreateVolume(context.Background(), ovhsdk.VolumeCreateOptions{
			Name:   fmt.Sprintf("%s-vol", args.MinecraftResource.GetName()),
			Size:   args.MinecraftResource.GetVolumeSize(),
			Region: args.MinecraftResource.GetRegion(),
			Type:   ovhsdk.VolumeClassic,
		})
		if err != nil {
			return nil, err
		}

		stillCreating = true
		for stillCreating {
			volume, err = o.client.GetVolume(context.Background(), volume.ID)
			if err != nil {
				return nil, err
			}
			if volume.Status == ovhsdk.VolumeAvailable {
				stillCreating = false
			} else {
				time.Sleep(2 * time.Second)
			}
		}

		_, err = o.client.AttachVolume(context.Background(), volume.ID, &ovhsdk.VolumeAttachOptions{
			InstanceID: instance.ID,
		})
		if err != nil {
			return nil, err
		}
		stillAttaching := true
		for stillAttaching {
			volume, err = o.client.GetVolume(context.Background(), volume.ID)
			if err != nil {
				return nil, err
			}
			if volume.Status == ovhsdk.VolumeInUse {
				stillAttaching = false
			} else {
				time.Sleep(2 * time.Second)
			}
		}
	}

	labels, err := getOVHFieldsFromID(instance.Name)
	if err != nil {
		return nil, err
	}
	ip4, err := ovhsdk.IPv4(instance)
	if err != nil {
		return nil, err
	}
	return &automation.ResourceResults{
		ID:       instance.ID,
		Name:     instance.Name,
		Region:   instance.Region,
		PublicIP: ip4,
		Tags:     labels,
	}, err
}

func (o *OVHcloud) DeleteServer(id string, args automation.ServerArgs) error {
	keys, err := o.client.ListSSHKeys(context.Background())
	if err != nil {
		return err
	}
	for _, key := range keys {
		if key.Name == fmt.Sprintf("%s-ssh", args.MinecraftResource.GetName()) {
			err := o.client.DeleteSSHKey(context.Background(), key.ID)
			if err != nil {
				return err
			}
		}
	}
	volumes, err := o.client.ListVolumes(context.Background())
	if err != nil {
		return err
	}
	for _, volume := range volumes {
		for _, attached := range volume.AttachedTo {
			if attached == id {
				detachVolume, err := o.client.DetachVolume(context.Background(), volume.ID, &ovhsdk.VolumeDetachOptions{
					InstanceID: id,
				})
				if err != nil {
					return err
				}
				stillDetaching := true
				for stillDetaching {
					detachedVolume, err := o.client.GetVolume(context.Background(), detachVolume.ID)
					if err != nil {
						return err
					}
					if detachedVolume.Status == ovhsdk.VolumeAvailable {
						stillDetaching = false
					} else {
						time.Sleep(2 * time.Second)
					}
				}
				err = o.client.DeleteVolume(context.Background(), volume.ID)
				if err != nil {
					return err
				}
			}
		}
	}
	err = o.client.DeleteInstance(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}

func (o *OVHcloud) ListServer() ([]automation.ResourceResults, error) {
	instances, err := o.client.ListInstance(context.Background())
	if err != nil {
		return nil, err
	}
	var result []automation.ResourceResults
	for i, instance := range instances {
		// no error checking. could be server in the region which don't belong to minectl
		labels, _ := getOVHFieldsFromID(instance.Name)
		if strings.Contains(labels, common.InstanceTag) {
			ip4, err := ovhsdk.IPv4(&instances[i])
			if err != nil {
				return nil, err
			}
			result = append(result, automation.ResourceResults{
				ID:       instance.ID,
				Name:     instance.Name,
				Region:   instance.Region,
				PublicIP: ip4,
				Tags:     labels,
			})
		}
	}
	return result, nil
}

func (o *OVHcloud) UpdateServer(id string, args automation.ServerArgs) error {
	instance, err := o.client.GetInstance(context.Background(), id)
	if err != nil {
		return err
	}

	ip4, err := ovhsdk.IPv4(instance)
	if err != nil {
		return err
	}
	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSH(), ip4, "ubuntu")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	return nil
}

func (o *OVHcloud) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	instance, err := o.client.GetInstance(context.Background(), id)
	if err != nil {
		return err
	}

	ip4, err := ovhsdk.IPv4(instance)
	if err != nil {
		return err
	}
	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSH(), ip4, "ubuntu")

	// as we are not allowed to login via root user, we need to add sudo to the command
	source := filepath.Join("/tmp", filepath.Base(plugin))
	err = remoteCommand.TransferFile(plugin, source)
	if err != nil {
		return err
	}
	_, err = remoteCommand.ExecuteCommand(fmt.Sprintf("sudo mv %s %s\nsudo systemctl restart minecraft.service", source, destination))
	if err != nil {
		return err
	}
	return nil
}

func (o *OVHcloud) GetServer(id string, args automation.ServerArgs) (*automation.ResourceResults, error) {
	instance, err := o.client.GetInstance(context.Background(), id)
	if err != nil {
		return nil, err
	}

	ip4, err := ovhsdk.IPv4(instance)
	if err != nil {
		return nil, err
	}
	labels, err := getOVHFieldsFromID(instance.Name)
	if err != nil {
		return nil, err
	}
	return &automation.ResourceResults{
		ID:       instance.ID,
		Name:     instance.Name,
		Region:   instance.Region,
		PublicIP: ip4,
		Tags:     labels,
	}, err
}
