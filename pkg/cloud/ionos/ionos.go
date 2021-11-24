package ionos

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/minectl/pkg/automation"
	"github.com/minectl/pkg/common"
	minctlTemplate "github.com/minectl/pkg/template"
	"github.com/minectl/pkg/update"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type IONOS struct {
	client *ionoscloud.APIClient
	tmpl   *minctlTemplate.Template
}

const defaultAPI = "https://api.ionos.com/cloudapi/v5"

func NewIONOS(username, password, token string) (*IONOS, error) {
	client := ionoscloud.NewAPIClient(ionoscloud.NewConfiguration(username, password, token, defaultAPI))
	tmpl, err := minctlTemplate.NewTemplateCloudConfig()
	if err != nil {
		return nil, err
	}
	return &IONOS{
		client: client,
		tmpl:   tmpl,
	}, nil
}

func (i *IONOS) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	ctx := context.Background()

	datacenters, _, err := i.client.DataCenterApi.DatacentersGet(ctx).Execute()
	if err != nil {
		return nil, err
	}
	for _, datacenter := range *datacenters.Items {
		if *datacenter.GetProperties().GetName() == fmt.Sprintf("%s-dc", args.MinecraftResource.GetName()) {
			return nil, errors.Errorf("Datancenter %s exist, please change the name of your Minecraft server",
				args.MinecraftResource.GetName())
		}
	}

	datacenter := ionoscloud.Datacenter{
		Properties: &ionoscloud.DatacenterProperties{
			Name:        ionoscloud.PtrString(fmt.Sprintf("%s-dc", args.MinecraftResource.GetName())),
			Location:    ionoscloud.PtrString(args.MinecraftResource.GetRegion()),
			Description: ionoscloud.PtrString(common.InstanceTag),
		},
	}

	createdDatacenter, _, err := i.client.DataCenterApi.DatacentersPost(ctx).Datacenter(datacenter).Execute()
	if err != nil {
		return nil, err
	}
	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSH()))
	if err != nil {
		return nil, err
	}

	lanRequest := ionoscloud.LanPost{
		Properties: &ionoscloud.LanPropertiesPost{
			Name:   ionoscloud.PtrString(fmt.Sprintf("%s-lan", args.MinecraftResource.GetName())),
			Public: ionoscloud.PtrBool(true),
		},
	}
	lan, _, err := i.client.LanApi.DatacentersLansPost(ctx, *createdDatacenter.GetId()).Lan(lanRequest).Execute()
	if err != nil {
		return nil, err
	}
	lanID, _ := strconv.ParseInt(*lan.GetId(), 10, 0)
	userData, err := i.tmpl.GetTemplate(args.MinecraftResource, "", minctlTemplate.GetTemplateCloudConfigName(args.MinecraftResource.IsProxyServer()))
	if err != nil {
		return nil, err
	}

	images, _, err := i.client.ImageApi.ImagesGet(ctx).Execute()
	if err != nil {
		return nil, err
	}
	var ubuntuImage *string
	for _, image := range *images.Items {
		if *image.GetProperties().GetName() == "Ubuntu-20.04-LTS-server-cloud-init.qcow2" &&
			*image.GetProperties().GetLocation() == args.MinecraftResource.GetRegion() {
			ubuntuImage = image.GetId()
		}
	}
	plan := strings.Split(args.MinecraftResource.GetSize(), "-")
	cpu, _ := strconv.ParseInt(plan[0], 10, 0)
	ram, _ := strconv.ParseInt(plan[1], 10, 0)
	cpuFamiliy := plan[2]
	request := ionoscloud.Server{
		Properties: &ionoscloud.ServerProperties{
			Name:      ionoscloud.PtrString(args.MinecraftResource.GetName()),
			Cores:     ionoscloud.PtrInt32(int32(cpu)),
			Ram:       ionoscloud.PtrInt32(int32(ram)),
			CpuFamily: ionoscloud.PtrString(cpuFamiliy),
		},
		Entities: &ionoscloud.ServerEntities{
			Nics: &ionoscloud.Nics{
				Items: &[]ionoscloud.Nic{
					{
						Properties: &ionoscloud.NicProperties{
							Name:           ionoscloud.PtrString(fmt.Sprintf("%s-nic", args.MinecraftResource.GetName())),
							FirewallActive: ionoscloud.PtrBool(false),
							Dhcp:           ionoscloud.PtrBool(true),
							Lan:            ionoscloud.PtrInt32(int32(lanID)),
						},
					},
				},
			},
			Volumes: &ionoscloud.AttachedVolumes{
				Items: &[]ionoscloud.Volume{
					{
						Properties: &ionoscloud.VolumeProperties{
							Name:     ionoscloud.PtrString(fmt.Sprintf("%s-vol", args.MinecraftResource.GetName())),
							Image:    ubuntuImage,
							SshKeys:  &[]string{string(pubKeyFile)},
							Size:     ionoscloud.PtrFloat32(30),
							Type:     ionoscloud.PtrString("HDD"),
							UserData: ionoscloud.PtrString(base64.StdEncoding.EncodeToString([]byte(userData))),
						},
					},
				},
			},
		},
	}
	server, _, err := i.client.ServerApi.DatacentersServersPost(ctx, *createdDatacenter.GetId()).Server(request).Execute()
	if err != nil {
		return nil, err
	}
	stillCreating := true
	for stillCreating {
		server, resp, err := i.client.ServerApi.DatacentersServersFindById(ctx, *createdDatacenter.GetId(), *server.Id).Execute()
		if err != nil {
			if resp.StatusCode != 404 {
				return nil, err
			}
			time.Sleep(2 * time.Second)
		} else {
			if *server.Metadata.GetState() == "AVAILABLE" {
				stillCreating = false
			} else {
				time.Sleep(2 * time.Second)
			}
		}
	}
	server, _, err = i.client.ServerApi.DatacentersServersFindById(ctx, *createdDatacenter.GetId(), *server.Id).Execute()
	if err != nil {
		return nil, err
	}
	ips := *(*server.Entities.Nics.Items)[0].Properties.Ips
	return &automation.RessourceResults{
		ID:       *createdDatacenter.GetId(),
		Name:     *server.Properties.GetName(),
		Region:   *createdDatacenter.GetProperties().GetLocation(),
		PublicIP: ips[0],
		Tags:     *server.Metadata.Etag,
	}, err
}

func (i *IONOS) DeleteServer(id string, _ automation.ServerArgs) error {
	ctx := context.Background()
	_, _, err := i.client.DataCenterApi.DatacentersDelete(ctx, id).Execute()
	if err != nil {
		return err
	}
	return nil
}

func (i *IONOS) ListServer() ([]automation.RessourceResults, error) {
	ctx := context.Background()
	var result []automation.RessourceResults
	datacenters, _, err := i.client.DataCenterApi.DatacentersGet(ctx).Execute()
	if err != nil {
		return nil, err
	}
	for _, datacenter := range *datacenters.Items {
		if *datacenter.GetProperties().GetDescription() == common.InstanceTag {
			servers, _, err := i.client.ServerApi.DatacentersServersGet(ctx, *datacenter.GetId()).Execute()
			if err != nil {
				return nil, err
			}
			for _, serverItem := range *servers.Items {
				server, _, err := i.client.ServerApi.DatacentersServersFindById(ctx, *datacenter.GetId(), *serverItem.Id).Execute()
				if err != nil {
					return nil, err
				}
				ips := *(*server.Entities.Nics.Items)[0].Properties.Ips
				result = append(result, automation.RessourceResults{
					ID:       *datacenter.GetId(),
					Name:     *server.Properties.GetName(),
					Region:   *datacenter.GetProperties().GetLocation(),
					PublicIP: ips[0],
					Tags:     *server.Metadata.Etag,
				})
			}

		}
	}
	return result, nil
}

func (i *IONOS) UpdateServer(id string, args automation.ServerArgs) error {
	server, err := i.GetServer(id, args)
	if err != nil {
		return err
	}
	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSH(), server.PublicIP, "ubuntu")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	zap.S().Infow("minectl server updated", "name", server.Name)
	return nil
}

func (i *IONOS) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	server, err := i.GetServer(id, args)
	if err != nil {
		return err
	}

	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSH(), server.PublicIP, "ubuntu")
	err = remoteCommand.TransferFile(plugin, filepath.Join(destination, filepath.Base(plugin)))
	if err != nil {
		return err
	}
	_, err = remoteCommand.ExecuteCommand("systemctl restart minecraft.service")
	if err != nil {
		return err
	}
	zap.S().Infow("Minecraft plugin uploaded", "plugin", plugin, "instance", server)
	return nil
}

func (i *IONOS) GetServer(id string, _ automation.ServerArgs) (*automation.RessourceResults, error) {
	ctx := context.Background()
	datacenter, _, err := i.client.DataCenterApi.DatacentersFindById(ctx, id).Execute()
	if err != nil {
		return nil, err
	}
	servers, _, err := i.client.ServerApi.DatacentersServersGet(ctx, *datacenter.GetId()).Execute()
	if err != nil {
		return nil, err
	}
	for _, serverItem := range *servers.Items {
		zap.S().Infow("Ionos found server", "server", serverItem)
		server, _, err := i.client.ServerApi.DatacentersServersFindById(ctx, *datacenter.GetId(), *serverItem.Id).Execute()
		if err != nil {
			return nil, err
		}
		ips := *(*server.Entities.Nics.Items)[0].Properties.Ips
		if len(ips) > 0 {
			return &automation.RessourceResults{
				ID:       *datacenter.GetId(),
				Name:     *server.Properties.GetName(),
				Region:   *datacenter.GetProperties().GetLocation(),
				PublicIP: ips[0],
				Tags:     *server.Metadata.Etag,
			}, nil
		}
	}
	return nil, errors.New("no Minecraft server found in Datacenter")
}
