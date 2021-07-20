package ovh

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/minectl/pgk/automation"
	"github.com/minectl/pgk/common"
	minctlTemplate "github.com/minectl/pgk/template"
	"github.com/ovh/go-ovh/ovh"
)

type OVHcloud struct {
	client      *ovh.Client
	serviceName string
	region      string
}

func NewOVHcloud(endpoint, appKey, appSecret, consumerKey, region string) (*OVHcloud, error) {
	client, err := ovh.NewClient(endpoint, appKey, appSecret, consumerKey)
	if err != nil {
		return nil, err
	}
	return &OVHcloud{
		client:      client,
		serviceName: "c3878ba251b5478181eab758e6b34d6a",
		region:      region,
	}, nil
}

func createOVHID(instanceName, label string) (id string) {
	return fmt.Sprintf("%s|%s", instanceName, label)
}

func getOVHFieldsFromID(id string) (instanceName, label string, err error) {
	fields := strings.Split(id, "|")
	err = nil
	if len(fields) == 3 {
		instanceName = fields[0]
		label = strings.Join([]string{fields[1], fields[2]}, ",")
	} else {
		err = fmt.Errorf("could not get fields from custom ID: fields: %v", fields)
		return "", "", err
	}
	return instanceName, label, nil
}

func (o *OVHcloud) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	pubKeyFile, err := ioutil.ReadFile(args.MinecraftServer.GetSSH())
	if err != nil {
		return nil, err
	}

	key, err := o.CreateSSHKey(context.Background(), SSHKeyCreateOptions{
		Name:      fmt.Sprintf("%s-ssh", args.MinecraftServer.GetName()),
		PublicKey: string(pubKeyFile),
	})
	if err != nil {
		return nil, err
	}

	image, err := o.GetImage(context.Background(), "Ubuntu 20.04", args.MinecraftServer.GetRegion())
	if err != nil {
		return nil, err
	}

	flavor, err := o.GetFlavor(context.Background(), args.MinecraftServer.GetSize(), args.MinecraftServer.GetRegion())
	if err != nil {
		return nil, err
	}

	tmpl, err := minctlTemplate.NewTemplateBash(args.MinecraftServer, "sdb")
	if err != nil {
		return nil, err
	}
	userData, err := tmpl.GetTemplate()
	if err != nil {
		return nil, err
	}

	instance, err := o.CreateInstance(context.Background(), InstanceCreateOptions{
		Name:           createOVHID(args.MinecraftServer.GetName(), strings.Join([]string{common.InstanceTag, args.MinecraftServer.GetEdition()}, "|")),
		Region:         args.MinecraftServer.GetRegion(),
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
		instance, err = o.GetInstance(context.Background(), instance.ID)
		if err != nil {
			return nil, err
		}
		if instance.Status == InstanceActive {
			stillCreating = false
			time.Sleep(2 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	volume, err := o.CreateVolume(context.Background(), VolumeCreateOptions{
		Name:   fmt.Sprintf("%s-vol", args.MinecraftServer.GetName()),
		Size:   args.MinecraftServer.GetVolumeSize(),
		Region: args.MinecraftServer.GetRegion(),
		Type:   VolumeClassic,
	})
	if err != nil {
		return nil, err
	}

	stillCreating = true
	for stillCreating {
		volume, err = o.GetVolume(context.Background(), volume.ID)
		if err != nil {
			return nil, err
		}
		if volume.Status == VolumeAvailable {
			stillCreating = false
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	_, err = o.AttachVolume(context.Background(), volume.ID, &VolumeAttachOptions{
		InstanceID: instance.ID,
	})
	if err != nil {
		return nil, err
	}
	stillAttaching := true
	for stillAttaching {
		volume, err = o.GetVolume(context.Background(), volume.ID)
		if err != nil {
			return nil, err
		}
		if volume.Status == VolumeInUse {
			stillAttaching = false
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	_, labels, err := getOVHFieldsFromID(instance.Name)
	if err != nil {
		return nil, err
	}
	ip4, err := IPv4(instance)
	if err != nil {
		return nil, err
	}
	return &automation.RessourceResults{
		ID:       instance.ID,
		Name:     instance.Name,
		Region:   instance.Region,
		PublicIP: ip4,
		Tags:     labels,
	}, err
}

func (o *OVHcloud) DeleteServer(id string, args automation.ServerArgs) error {
	keys, err := o.ListSSHKeys(context.Background())
	if err != nil {
		return err
	}
	for _, key := range keys {
		if key.Name == fmt.Sprintf("%s-ssh", args.MinecraftServer.GetName()) {
			err := o.DeleteSSHKey(context.Background(), key.ID)
			if err != nil {
				return err
			}
		}
	}
	volumes, err := o.ListVolumes(context.Background())
	if err != nil {
		return err
	}
	for _, volume := range volumes {
		for _, attached := range volume.AttachedTo {
			if attached == id {
				detachVolume, err := o.DetachVolume(context.Background(), volume.ID, &VolumeDetachOptions{
					InstanceID: id,
				})
				if err != nil {
					return err
				}
				stillDetaching := true
				for stillDetaching {
					detachedVolume, err := o.GetVolume(context.Background(), detachVolume.ID)
					if err != nil {
						return err
					}
					if detachedVolume.Status == VolumeAvailable {
						stillDetaching = false
					} else {
						time.Sleep(2 * time.Second)
					}
				}
				err = o.DeleteVolume(context.Background(), volume.ID)
				if err != nil {
					return err
				}
			}
		}
	}
	err = o.DeleteInstance(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}

func (o *OVHcloud) ListServer() ([]automation.RessourceResults, error) {
	instances, err := o.ListInstance(context.Background())
	if err != nil {
		return nil, err
	}
	var result []automation.RessourceResults
	for _, instance := range instances {
		// no error checking. could be server in the region which don't belong to minectl
		_, labels, _ := getOVHFieldsFromID(instance.Name)
		if strings.Contains(labels, common.InstanceTag) {
			ip4, err := IPv4(&instance)
			if err != nil {
				return nil, err
			}
			result = append(result, automation.RessourceResults{
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

func (o OVHcloud) UpdateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	panic("implement me")
}
