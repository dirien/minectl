package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type SSHKeyCreateOptions struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

type SSHKey struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Regions     []string `json:"regions"`
	FingerPrint string   `json:"fingerPrint"`
	PublicKey   string   `json:"publicKey"`
}

func (o *OVHcloud) CreateSSHKey(ctx context.Context, createOpts SSHKeyCreateOptions) (*SSHKey, error) {
	resp := SSHKey{}
	err := o.client.PostWithContext(ctx, fmt.Sprintf("/cloud/project/%s/sshkey", o.serviceName), createOpts, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

func (o *OVHcloud) ListSSHKeys(ctx context.Context) ([]SSHKey, error) {
	url := fmt.Sprintf("/cloud/project/%s/sshkey", o.serviceName)
	var keys []SSHKey
	err := o.client.GetWithContext(ctx, url, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (o *OVHcloud) DeleteSSHKey(ctx context.Context, id string) error {
	url := fmt.Sprintf("/cloud/project/%s/sshkey/%s", o.serviceName, id)
	err := o.client.DeleteWithContext(ctx, url, nil)
	if err != nil {
		return err
	}
	return nil
}

type VolumeType string

const (
	VolumeClassic   VolumeType = "classic"
	VolumeHighSpeed VolumeType = "high-speed"
)

type VolumeStatus string

const (
	VolumeAttaching VolumeStatus = "attaching"
	VolumeCreating  VolumeStatus = "creating"
	VolumeAvailable VolumeStatus = "available"
	VolumeInUse     VolumeStatus = "in-use"
)

type Volume struct {
	ID           string       `json:"id"`
	AttachedTo   []string     `json:"attachedTo"`
	CreationDate time.Time    `json:"creationDate"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Size         int          `json:"size"`
	Status       VolumeStatus `json:"status"`
	Region       string       `json:"region"`
	Bootable     bool         `json:"bootable"`
	PlanCode     string       `json:"planCode"`
	Type         VolumeType   `json:"type"`
}

type VolumeCreateOptions struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Size        int        `json:"size"`
	Region      string     `json:"region"`
	Type        VolumeType `json:"type"`
}

type VolumeAttachOptions struct {
	InstanceID string `json:"instanceId"`
}

type VolumeDetachOptions VolumeAttachOptions

func (o *OVHcloud) CreateVolume(ctx context.Context, createOpts VolumeCreateOptions) (*Volume, error) {
	resp := Volume{}
	err := o.client.PostWithContext(ctx, fmt.Sprintf("/cloud/project/%s/volume", o.serviceName), createOpts, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

func (o *OVHcloud) ListVolumes(ctx context.Context) ([]Volume, error) {
	var volumes []Volume
	err := o.client.GetWithContext(ctx, fmt.Sprintf("/cloud/project/%s/volume", o.serviceName), &volumes)
	if err != nil {
		return nil, err
	}
	return volumes, err
}

func (o *OVHcloud) GetVolume(ctx context.Context, id string) (*Volume, error) {
	url := fmt.Sprintf("/cloud/project/%s/volume/%s", o.serviceName, id)
	resp := Volume{}
	err := o.client.GetWithContext(ctx, url, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (o *OVHcloud) DeleteVolume(ctx context.Context, id string) error {
	url := fmt.Sprintf("/cloud/project/%s/volume/%s", o.serviceName, id)
	err := o.client.DeleteWithContext(ctx, url, nil)
	if err != nil {
		return err
	}
	return nil
}

func (o *OVHcloud) AttachVolume(ctx context.Context, id string, options *VolumeAttachOptions) (*Volume, error) {
	resp := Volume{}
	err := o.client.PostWithContext(ctx, fmt.Sprintf("/cloud/project/%s/volume/%s/attach", o.serviceName, id), options, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

func (o *OVHcloud) DetachVolume(ctx context.Context, id string, options *VolumeDetachOptions) (*Volume, error) {
	resp := Volume{}
	err := o.client.PostWithContext(ctx, fmt.Sprintf("/cloud/project/%s/volume/%s/detach", o.serviceName, id), options, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

type Image struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Region       string    `json:"region"`
	Visibility   string    `json:"visibility"`
	Type         string    `json:"type"`
	MinDisk      int       `json:"minDisk"`
	MinRAM       int       `json:"minRam"`
	Size         float64   `json:"size"`
	CreationDate time.Time `json:"creationDate"`
	Status       string    `json:"status"`
	User         string    `json:"user"`
	FlavorType   string    `json:"flavorType"`
	Tags         []string  `json:"tags"`
	PlanCode     string    `json:"planCode"`
}

func (o *OVHcloud) GetImage(ctx context.Context, name, region string) (*Image, error) {
	url := fmt.Sprintf("/cloud/project/%s/image", o.serviceName)
	var images []Image
	err := o.client.GetWithContext(ctx, url, &images)
	if err != nil {
		return nil, err
	}
	for _, image := range images {
		if image.Region == region && image.Name == name {
			return &image, nil
		}
	}
	return nil, errors.Errorf("image: %s in region: %s not found", name, region)
}

type Flavor struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Region            string `json:"region"`
	RAM               int    `json:"ram"`
	Disk              int    `json:"disk"`
	Vcpus             int    `json:"vcpus"`
	Type              string `json:"type"`
	OsType            string `json:"osType"`
	InboundBandwidth  int    `json:"inboundBandwidth"`
	OutboundBandwidth int    `json:"outboundBandwidth"`
	Available         bool   `json:"available"`
	PlanCodes         struct {
		Monthly string `json:"monthly"`
		Hourly  string `json:"hourly"`
	} `json:"planCodes"`
	Capabilities []struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	} `json:"capabilities"`
	Quota int `json:"quota"`
}

func (o *OVHcloud) GetFlavor(ctx context.Context, name, region string) (*Flavor, error) {
	url := fmt.Sprintf("/cloud/project/%s/flavor", o.serviceName)
	var flavors []Flavor
	err := o.client.GetWithContext(ctx, url, &flavors)
	if err != nil {
		return nil, err
	}
	for _, flavor := range flavors {
		if flavor.Region == region && flavor.Name == name {
			return &flavor, nil
		}
	}
	return nil, errors.Errorf("flavor: %s in region: %s not found", name, region)
}

type IP struct {
	IP        string `json:"ip"`
	Type      string `json:"type"`
	Version   int    `json:"version"`
	NetworkID string `json:"networkId"`
	GatewayIP string `json:"gatewayIp"`
}

type Instance struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	IPAddresses    []IP           `json:"ipAddresses"`
	FlavorID       string         `json:"flavorId"`
	ImageID        string         `json:"imageId"`
	SSHKeyID       string         `json:"sshKeyId"`
	Created        time.Time      `json:"created"`
	Region         string         `json:"region"`
	MonthlyBilling string         `json:"monthlyBilling"`
	Status         InstanceStatus `json:"status"`
	PlanCode       string         `json:"planCode"`
	OperationIds   []string       `json:"operationIds"`
}

type InstanceStatus string

const (
	InstanceActive   InstanceStatus = "ACTIVE"
	InstanceBuilding InstanceStatus = "BUILDING"
	InstanceDeleted  InstanceStatus = "DELETED"
	InstanceDeleting InstanceStatus = "DELETING"
	InstanceError    InstanceStatus = "ERROR"
	InstanceReboot   InstanceStatus = "REBOOT"
	InstanceStopped  InstanceStatus = "STOPPED"
	InstanceUnknown  InstanceStatus = "UNKNOWN"
	InstanceBuild    InstanceStatus = "BUILD"
	InstanceResuming InstanceStatus = "RESUMING"
	InstanceRebuild  InstanceStatus = "REBUILD"
)

type InstanceCreateOptions struct {
	FlavorID       string `json:"flavorId"`
	ImageID        string `json:"imageId"`
	MonthlyBilling bool   `json:"monthlyBilling"`
	Name           string `json:"name"`
	Region         string `json:"region"`
	SSHKeyID       string `json:"sshKeyId"`
	UserData       string `json:"userData"`
}

func (o *OVHcloud) CreateInstance(ctx context.Context, createOpts InstanceCreateOptions) (*Instance, error) {
	resp := Instance{}
	err := o.client.PostWithContext(ctx, fmt.Sprintf("/cloud/project/%s/instance", o.serviceName), createOpts, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

func (o *OVHcloud) GetInstance(ctx context.Context, id string) (*Instance, error) {
	url := fmt.Sprintf("/cloud/project/%s/instance/%s", o.serviceName, id)
	resp := Instance{}
	err := o.client.GetWithContext(ctx, url, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (o *OVHcloud) ListInstance(ctx context.Context) ([]Instance, error) {
	url := fmt.Sprintf("/cloud/project/%s/instance?region=%s", o.serviceName, o.region)
	var resp []Instance
	err := o.client.GetWithContext(ctx, url, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *OVHcloud) DeleteInstance(ctx context.Context, id string) error {
	url := fmt.Sprintf("/cloud/project/%s/instance/%s", o.serviceName, id)
	err := o.client.DeleteWithContext(ctx, url, nil)
	if err != nil {
		return err
	}
	return nil
}

func IPv4(instance *Instance) (string, error) {
	for _, ip := range instance.IPAddresses {
		if ip.Version == 4 {
			return ip.IP, nil
		}
	}
	return "", errors.Errorf("no ip4 address found for instance: %s", instance.Name)
}
