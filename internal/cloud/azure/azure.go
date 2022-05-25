package azure

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-03-01/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-10-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/minectl/internal/automation"
	"github.com/minectl/internal/common"
	minctlTemplate "github.com/minectl/internal/template"
	"github.com/minectl/internal/update"
)

type Azure struct {
	subscriptionID string
	authorizer     autorest.Authorizer
	tmpl           *minctlTemplate.Template
}

func NewAzure(authFile string) (*Azure, error) {
	authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return nil, err
	}
	authInfo, err := readJSON(authFile)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read JSON")
	}
	tmpl, err := minctlTemplate.NewTemplateCloudConfig()
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Azure set cloud-config template", "name", tmpl.Template.Name())
	return &Azure{
		subscriptionID: (*authInfo)["subscriptionId"].(string),
		authorizer:     authorizer,
		tmpl:           tmpl,
	}, nil
}

func readJSON(path string) (*map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read file")
	}
	contents := make(map[string]interface{})
	_ = json.Unmarshal(data, &contents)
	return &contents, nil
}

func getTags(edition string) map[string]*string {
	return map[string]*string{
		common.InstanceTag: to.StringPtr("true"),
		edition:            to.StringPtr("true"),
	}
}

func getTagKeys(tags map[string]*string) []string {
	var keys []string
	for key := range tags {
		keys = append(keys, key)
	}
	return keys
}

func (a *Azure) CreateServer(args automation.ServerArgs) (*automation.ResourceResults, error) {
	groupsClient := resources.NewGroupsClient(a.subscriptionID)
	groupsClient.Authorizer = a.authorizer

	group, err := groupsClient.CreateOrUpdate(
		context.Background(),
		fmt.Sprintf("%s-rg", args.MinecraftResource.GetName()),
		resources.Group{
			Response: autorest.Response{},
			Location: to.StringPtr(args.MinecraftResource.GetRegion()),
			Tags:     getTags(args.MinecraftResource.GetEdition()),
		})
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Azure resource group created", "name", group.Name)

	virtualNetworksClient := network.NewVirtualNetworksClient(a.subscriptionID)
	virtualNetworksClient.Authorizer = a.authorizer

	virtualNetworksCreateOrUpdateFuture, err := virtualNetworksClient.CreateOrUpdate(
		context.Background(),
		to.String(group.Name),
		fmt.Sprintf("%s-vnet", args.MinecraftResource.GetName()),
		network.VirtualNetwork{
			Name:     to.StringPtr(fmt.Sprintf("%s-vnet", args.MinecraftResource.GetName())),
			Location: to.StringPtr(args.MinecraftResource.GetRegion()),
			Tags:     getTags(args.MinecraftResource.GetEdition()),
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{"10.0.0.0/8"},
				},
			},
		})
	if err != nil {
		return nil, err
	}
	err = virtualNetworksCreateOrUpdateFuture.WaitForCompletionRef(context.Background(), virtualNetworksClient.Client)
	if err != nil {
		return nil, err
	}
	vnet, err := virtualNetworksCreateOrUpdateFuture.Result(virtualNetworksClient)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Azure virtual network created", "name", vnet.Name)

	subnetsClient := network.NewSubnetsClient(a.subscriptionID)
	subnetsClient.Authorizer = a.authorizer
	subnetsCreateOrUpdateFuture, err := subnetsClient.CreateOrUpdate(
		context.Background(),
		to.String(group.Name),
		to.String(vnet.Name),
		fmt.Sprintf("%s-snet", args.MinecraftResource.GetName()),
		network.Subnet{
			SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr("10.0.0.0/16"),
			},
		})
	if err != nil {
		return nil, err
	}

	err = subnetsCreateOrUpdateFuture.WaitForCompletionRef(context.Background(), subnetsClient.Client)
	if err != nil {
		return nil, err
	}
	subnet, err := subnetsCreateOrUpdateFuture.Result(subnetsClient)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Azure subnetwork created", "name", subnet.Name)

	ipClient := network.NewPublicIPAddressesClient(a.subscriptionID)
	ipClient.Authorizer = a.authorizer
	publicIPAddressesCreateOrUpdateFuture, err := ipClient.CreateOrUpdate(
		context.Background(),
		to.String(group.Name),
		fmt.Sprintf("%s-ip", args.MinecraftResource.GetName()),
		network.PublicIPAddress{
			Name:     to.StringPtr(fmt.Sprintf("%s-ip", args.MinecraftResource.GetName())),
			Location: to.StringPtr(args.MinecraftResource.GetRegion()),
			PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   network.IPVersionIPv4,
				PublicIPAllocationMethod: network.IPAllocationMethodStatic,
			},
			Tags: getTags(args.MinecraftResource.GetEdition()),
		},
	)
	if err != nil {
		return nil, err
	}

	err = publicIPAddressesCreateOrUpdateFuture.WaitForCompletionRef(context.Background(), ipClient.Client)
	if err != nil {
		return nil, err
	}
	ip, err := publicIPAddressesCreateOrUpdateFuture.Result(ipClient)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Azure public ip created", "name", ip.Name)

	nicClient := network.NewInterfacesClient(a.subscriptionID)
	nicClient.Authorizer = a.authorizer
	interfacesCreateOrUpdateFuture, err := nicClient.CreateOrUpdate(
		context.Background(),
		to.String(group.Name),
		fmt.Sprintf("%s-nic", args.MinecraftResource.GetName()),
		network.Interface{
			Name:     to.StringPtr(fmt.Sprintf("%s-nic", args.MinecraftResource.GetName())),
			Location: group.Location,
			InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
				IPConfigurations: &[]network.InterfaceIPConfiguration{
					{
						Name: to.StringPtr("ipConfig1"),
						InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
							Subnet:                    &subnet,
							PrivateIPAllocationMethod: network.IPAllocationMethodDynamic,
							PublicIPAddress:           &ip,
						},
					},
				},
			},
			Tags: getTags(args.MinecraftResource.GetEdition()),
		})
	if err != nil {
		return nil, err
	}
	err = interfacesCreateOrUpdateFuture.WaitForCompletionRef(context.Background(), nicClient.Client)
	if err != nil {
		return nil, err
	}
	nic, err := interfacesCreateOrUpdateFuture.Result(nicClient)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Azure network interface controller created", "name", nic.Name)

	var mount string
	var diskID *string
	if args.MinecraftResource.GetVolumeSize() > 0 {
		disksClient := compute.NewDisksClient(a.subscriptionID)
		disksClient.Authorizer = a.authorizer
		disksCreateOrUpdateFuture, err := disksClient.CreateOrUpdate(
			context.Background(),
			to.String(group.Name),
			fmt.Sprintf("%s-vol", args.MinecraftResource.GetName()),
			compute.Disk{
				Location: group.Location,
				DiskProperties: &compute.DiskProperties{
					CreationData: &compute.CreationData{
						CreateOption: compute.DiskCreateOptionEmpty,
					},
					DiskSizeGB: to.Int32Ptr(int32(args.MinecraftResource.GetVolumeSize())),
				},
			})
		if err != nil {
			return nil, err
		}

		err = disksCreateOrUpdateFuture.WaitForCompletionRef(context.Background(), disksClient.Client)
		if err != nil {
			return nil, err
		}
		disk, err := disksCreateOrUpdateFuture.Result(disksClient)
		if err != nil {
			return nil, err
		}
		diskID = disk.ID
		mount = "sda"
		zap.S().Infow("Azure managed disk created", "name", disk.Name)
	}
	vmClient := compute.NewVirtualMachinesClient(a.subscriptionID)
	vmClient.Authorizer = a.authorizer

	pubKeyFile, err := os.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSHKeyFolder()))
	if err != nil {
		return nil, err
	}
	userData, err := a.tmpl.GetTemplate(args.MinecraftResource, &minctlTemplate.CreateUpdateTemplateArgs{Mount: mount, Name: minctlTemplate.GetTemplateCloudConfigName(args.MinecraftResource.IsProxyServer())})
	if err != nil {
		return nil, err
	}

	priority := compute.VirtualMachinePriorityTypesRegular
	var evictionPolicy compute.VirtualMachineEvictionPolicyTypes
	if args.MinecraftResource.IsSpot() {
		priority = compute.VirtualMachinePriorityTypesSpot
		evictionPolicy = compute.VirtualMachineEvictionPolicyTypesDeallocate
	}
	image := &compute.ImageReference{
		Publisher: to.StringPtr("Canonical"),
		Offer:     to.StringPtr("0001-com-ubuntu-minimal-jammy-daily"),
		Sku:       to.StringPtr("minimal-22_04-daily-lts-gen2"),
		Version:   to.StringPtr("latest"),
	}
	if args.MinecraftResource.IsArm() {
		image.Offer = to.StringPtr("0001-com-ubuntu-server-arm-preview-focal")
		image.Sku = to.StringPtr("20_04-lts")
	}

	vmOptions := compute.VirtualMachine{
		Location: group.Location,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			Priority:       priority,
			EvictionPolicy: evictionPolicy,
			HardwareProfile: &compute.HardwareProfile{
				VMSize: compute.VirtualMachineSizeTypes(args.MinecraftResource.GetSize()),
			},
			StorageProfile: &compute.StorageProfile{
				ImageReference: image,
			},
			OsProfile: &compute.OSProfile{
				ComputerName:  to.StringPtr(args.MinecraftResource.GetName()),
				AdminUsername: to.StringPtr("ubuntu"),
				CustomData:    to.StringPtr(base64.StdEncoding.EncodeToString([]byte(userData))),
				LinuxConfiguration: &compute.LinuxConfiguration{
					SSH: &compute.SSHConfiguration{
						PublicKeys: &[]compute.SSHPublicKey{
							{
								Path:    to.StringPtr("/home/ubuntu/.ssh/authorized_keys"),
								KeyData: to.StringPtr(string(pubKeyFile)),
							},
						},
					},
				},
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: nic.ID,
						NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
							Primary: to.BoolPtr(true),
						},
					},
				},
			},
		},
		Tags: getTags(args.MinecraftResource.GetEdition()),
	}

	if args.MinecraftResource.GetVolumeSize() > 0 {
		vmOptions.StorageProfile.DataDisks = &[]compute.DataDisk{{
			CreateOption: compute.DiskCreateOptionTypesAttach,
			Lun:          to.Int32Ptr(0),
			ManagedDisk: &compute.ManagedDiskParameters{
				ID: diskID,
			},
		}}
	}

	virtualMachinesCreateOrUpdateFuture, err := vmClient.CreateOrUpdate(
		context.Background(),
		to.String(group.Name),
		args.MinecraftResource.GetName(),
		vmOptions)
	if err != nil {
		return nil, err
	}
	err = virtualMachinesCreateOrUpdateFuture.WaitForCompletionRef(context.Background(), vmClient.Client)
	if err != nil {
		return nil, err
	}
	instance, err := virtualMachinesCreateOrUpdateFuture.Result(vmClient)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Azure virtual machine created", "name", instance.Name)

	virtualMachinesStartFuture, err := vmClient.Start(context.Background(), to.String(group.Name), to.String(instance.Name))
	if err != nil {
		return nil, err
	}

	err = virtualMachinesStartFuture.WaitForCompletionRef(context.Background(), vmClient.Client)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Azure virtual machine started", "name", instance.Name, "ip", ip.IPAddress, "id", instance.Name)

	return &automation.ResourceResults{
		ID:       to.String(instance.Name),
		Name:     to.String(instance.Name),
		Region:   to.String(group.Location),
		PublicIP: to.String(ip.IPAddress),
		Tags:     strings.Join(getTagKeys(instance.Tags), ","),
	}, err
}

func (a *Azure) DeleteServer(id string, args automation.ServerArgs) error {
	resourceGroupName := fmt.Sprintf("%s-rg", args.MinecraftResource.GetName())
	zap.S().Infow("Azure delete resource group", "name", resourceGroupName)
	groupsClient := resources.NewGroupsClient(a.subscriptionID)
	groupsClient.Authorizer = a.authorizer
	groupsDeleteFuture, err := groupsClient.Delete(context.Background(), resourceGroupName)
	if err != nil {
		return err
	}
	err = groupsDeleteFuture.WaitForCompletionRef(context.Background(), groupsClient.Client)
	if err != nil {
		return err
	}
	zap.S().Infow("Azure resource group deleted", "name", resourceGroupName)
	return nil
}

func (a *Azure) ListServer() ([]automation.ResourceResults, error) {
	vmClient := compute.NewVirtualMachinesClient(a.subscriptionID)
	vmClient.Authorizer = a.authorizer
	virtualMachineListResultPage, err := vmClient.ListAll(
		context.Background(), "false")
	if err != nil {
		return nil, err
	}
	var result []automation.ResourceResults
	for _, instance := range virtualMachineListResultPage.Values() {
		for key := range instance.Tags {
			if key == common.InstanceTag {
				publicIPAddressesClient := network.NewPublicIPAddressesClient(a.subscriptionID)
				publicIPAddressesClient.Authorizer = a.authorizer
				ip, err := publicIPAddressesClient.Get(
					context.Background(),
					fmt.Sprintf("%s-rg", to.String(instance.Name)),
					fmt.Sprintf("%s-ip", to.String(instance.Name)),
					"")
				if err != nil {
					return nil, err
				}
				result = append(result, automation.ResourceResults{
					ID:       to.String(instance.Name),
					Name:     to.String(instance.Name),
					Region:   to.String(instance.Location),
					PublicIP: to.String(ip.IPAddress),
					Tags:     strings.Join(getTagKeys(instance.Tags), ","),
				})
			}
		}
	}
	if len(result) > 0 {
		zap.S().Infow("Azure list all minectl vms", "list", result)
	} else {
		zap.S().Infow("No minectl vms found")
	}
	return result, nil
}

func (a *Azure) UpdateServer(id string, args automation.ServerArgs) error {
	server, err := a.GetServer(id, args)
	if err != nil {
		return err
	}
	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), server.PublicIP, "ubuntu")
	err = remoteCommand.UpdateServer(args.MinecraftResource)
	if err != nil {
		return err
	}
	zap.S().Infow("minectl server updated", "name", server.Name)
	return nil
}

func (a *Azure) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	server, err := a.GetServer(id, args)
	if err != nil {
		return err
	}
	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSHKeyFolder(), server.PublicIP, "ubuntu")

	// as we are not allowed to login via root user, we need to add sudo to the command
	source := filepath.Join("/tmp", filepath.Base(plugin))
	err = remoteCommand.TransferFile(plugin, source, args.MinecraftResource.GetSSHPort())
	if err != nil {
		return err
	}
	_, err = remoteCommand.ExecuteCommand(fmt.Sprintf("sudo mv %s %s\nsudo systemctl restart minecraft.service", source, destination), args.MinecraftResource.GetSSHPort())
	if err != nil {
		return err
	}
	zap.S().Infow("Minecraft plugin uploaded", "plugin", plugin, "server", server.Name)
	return nil
}

func (a *Azure) GetServer(id string, args automation.ServerArgs) (*automation.ResourceResults, error) {
	vmClient := compute.NewVirtualMachinesClient(a.subscriptionID)
	vmClient.Authorizer = a.authorizer

	instance, err := vmClient.Get(
		context.Background(),
		fmt.Sprintf("%s-rg", args.MinecraftResource.GetName()),
		id,
		compute.InstanceViewTypesInstanceView,
	)
	if err != nil {
		return nil, err
	}

	ipClient := network.NewPublicIPAddressesClient(a.subscriptionID)
	ipClient.Authorizer = a.authorizer
	publicIPAddress, err := ipClient.Get(
		context.Background(),
		fmt.Sprintf("%s-rg", args.MinecraftResource.GetName()),
		fmt.Sprintf("%s-ip", args.MinecraftResource.GetName()),
		"")
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Get Azure minctl server", "name", instance.Name, "ip", publicIPAddress.IPAddress)
	return &automation.ResourceResults{
		ID:       to.String(instance.Name),
		Name:     to.String(instance.Name),
		Region:   args.MinecraftResource.GetRegion(),
		PublicIP: to.String(publicIPAddress.IPAddress),
		Tags:     strings.Join(getTagKeys(instance.Tags), ","),
	}, err
}
