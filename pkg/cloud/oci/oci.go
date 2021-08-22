package oci

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/minectl/pkg/automation"
	common2 "github.com/minectl/pkg/common"
	minctlTemplate "github.com/minectl/pkg/template"
	"github.com/minectl/pkg/update"
	"github.com/oracle/oci-go-sdk/v46/common"
	"github.com/oracle/oci-go-sdk/v46/core"
	"github.com/oracle/oci-go-sdk/v46/example/helpers"
	"github.com/oracle/oci-go-sdk/v46/identity"
	"go.uber.org/zap"
)

type OCI struct {
	compute  core.ComputeClient
	identity identity.IdentityClient
	network  core.VirtualNetworkClient
	tmpl     *minctlTemplate.Template
}

func NewOCI() (*OCI, error) {
	c, err := core.NewComputeClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		return nil, err
	}
	tmpl, err := minctlTemplate.NewTemplateCloudConfig()
	if err != nil {
		return nil, err
	}
	i, err := identity.NewIdentityClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		return nil, err
	}
	n, err := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		return nil, err
	}
	return &OCI{
		compute:  c,
		identity: i,
		network:  n,
		tmpl:     tmpl,
	}, nil
}

func getTags(edition string) map[string]string {
	return map[string]string{
		common2.InstanceTag: "true",
		edition:             "true",
	}
}
func getTagKeys(tags map[string]string) []string {
	var keys []string
	for key := range tags {
		keys = append(keys, key)
	}
	return keys
}

func (o *OCI) CreateServer(args automation.ServerArgs) (*automation.RessourceResults, error) {
	ctx := context.Background()

	tenancyOCID, err := common.DefaultConfigProvider().TenancyOCID()
	if err != nil {
		return nil, err
	}

	compartmentRequest := identity.CreateCompartmentRequest{
		CreateCompartmentDetails: identity.CreateCompartmentDetails{
			CompartmentId: common.String(tenancyOCID),
			Name:          common.String(args.MinecraftResource.GetName()),
			Description:   common.String(fmt.Sprintf("Compartment for %s", args.MinecraftResource.GetName())),
			FreeformTags:  getTags(args.MinecraftResource.GetEdition()),
		},
	}
	compartment, err := o.identity.CreateCompartment(ctx, compartmentRequest)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Oracle compartment created", "compartment", compartment)

	request := identity.ListAvailabilityDomainsRequest{
		CompartmentId: compartment.Id,
	}

	availabilityDomains, err := o.identity.ListAvailabilityDomains(context.Background(), request)
	if err != nil {
		return nil, err
	}
	if len(availabilityDomains.Items) == 0 {
		zap.S().Error("no Availability Domains found")
		return nil, errors.New("no Availability Domains found")
	}
	availabilityDomain := availabilityDomains.Items[0]
	zap.S().Infow("Oracle get Availability Domain", "availabilityDomain", availabilityDomain)

	imagesRequest := core.ListImagesRequest{
		CompartmentId:          compartment.Id,
		OperatingSystem:        common.String("Canonical Ubuntu"),
		OperatingSystemVersion: common.String("20.04"),
		SortBy:                 core.ListImagesSortByTimecreated,
		SortOrder:              core.ListImagesSortOrderDesc,
		Shape:                  common.String(args.MinecraftResource.GetSize()),
		RequestMetadata:        helpers.GetRequestMetadataWithDefaultRetryPolicy(),
	}

	images, err := o.compute.ListImages(ctx, imagesRequest)
	if err != nil {
		return nil, err
	}
	if len(images.Items) == 0 {
		return nil, errors.New("no imags found")
	}
	image := images.Items[0]
	zap.S().Infow("Oracle get Image ", "image", image)

	vcnRequest := core.CreateVcnRequest{
		CreateVcnDetails: core.CreateVcnDetails{
			CidrBlock:     common.String("10.0.0.0/16"),
			CompartmentId: compartment.Id,
			DisplayName:   common.String(fmt.Sprintf("%s-vcn", args.MinecraftResource.GetName())),
			DnsLabel:      common.String(fmt.Sprintf("vcn%s", common2.InstanceTag)),
			FreeformTags:  getTags(args.MinecraftResource.GetEdition()),
		},
		RequestMetadata: helpers.GetRequestMetadataWithDefaultRetryPolicy(),
	}
	vcn, err := o.network.CreateVcn(ctx, vcnRequest)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Oracle VCN created", "vcn", vcn)

	var ingressSecurityRules []core.IngressSecurityRule
	minecraftIngressSecurityRule := core.IngressSecurityRule{
		//Options are supported only for ICMP ("1"), TCP ("6"), UDP ("17"), and ICMPv6 ("58").
		Description: common.String("Minecraft Server Port"),
		Protocol:    common.String("6"),
		IsStateless: common.Bool(true),
		Source:      common.String("0.0.0.0/0"),
		TcpOptions: &core.TcpOptions{
			DestinationPortRange: &core.PortRange{
				Max: common.Int(args.MinecraftResource.GetPort()),
				Min: common.Int(args.MinecraftResource.GetPort()),
			},
		},
	}
	if args.MinecraftResource.GetEdition() == "bedrock" {
		//Options are supported only for ICMP ("1"), TCP ("6"), UDP ("17"), and ICMPv6 ("58").
		minecraftIngressSecurityRule.Protocol = common.String("17")
		minecraftIngressSecurityRule.TcpOptions = nil
		minecraftIngressSecurityRule.UdpOptions = &core.UdpOptions{
			DestinationPortRange: &core.PortRange{
				Max: common.Int(args.MinecraftResource.GetPort()),
				Min: common.Int(args.MinecraftResource.GetPort()),
			},
		}
		zap.S().Info("Oracle bedrock server protocol set")
	} else {
		//RCON
		if args.MinecraftResource.HasRCON() {
			ingressSecurityRules = append(ingressSecurityRules, core.IngressSecurityRule{
				Description: common.String("RCON Port"),
				Protocol:    common.String("6"),
				IsStateless: common.Bool(true),
				Source:      common.String("0.0.0.0/0"),
				TcpOptions: &core.TcpOptions{
					DestinationPortRange: &core.PortRange{
						Max: common.Int(args.MinecraftResource.GetRCONPort()),
						Min: common.Int(args.MinecraftResource.GetRCONPort()),
					},
				},
			})
			zap.S().Info("Oracle RCON ingress security ruleset created")
		}
	}
	ingressSecurityRules = append(ingressSecurityRules, minecraftIngressSecurityRule)
	if args.MinecraftResource.HasMonitoring() {
		ingressSecurityRules = append(ingressSecurityRules, core.IngressSecurityRule{
			Description: common.String("Default Prometheus Port"),
			Protocol:    common.String("6"),
			IsStateless: common.Bool(true),
			Source:      common.String("0.0.0.0/0"),
			TcpOptions: &core.TcpOptions{
				DestinationPortRange: &core.PortRange{
					Max: common.Int(9090),
					Min: common.Int(9090),
				},
			},
		})
		zap.S().Info("Oracle Prometheus ingress security ruleset created")
	}

	securityListRequest := core.CreateSecurityListRequest{
		CreateSecurityListDetails: core.CreateSecurityListDetails{
			VcnId:         vcn.Id,
			CompartmentId: compartment.Id,
			DisplayName:   common.String(fmt.Sprintf("%s-sl", args.MinecraftResource.GetName())),
			FreeformTags:  getTags(args.MinecraftResource.GetEdition()),
			EgressSecurityRules: []core.EgressSecurityRule{
				{
					Protocol:    common.String("all"),
					Destination: common.String("0.0.0.0/0"),
				},
			},
			IngressSecurityRules: ingressSecurityRules,
		},
	}
	securityList, err := o.network.CreateSecurityList(ctx, securityListRequest)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Oracle Security List created", "securityList", securityList)

	subnetRequest := core.CreateSubnetRequest{
		CreateSubnetDetails: core.CreateSubnetDetails{
			CidrBlock:     common.String("10.0.0.0/24"),
			CompartmentId: compartment.Id,
			VcnId:         vcn.Id,
			FreeformTags:  getTags(args.MinecraftResource.GetEdition()),
			SecurityListIds: []string{
				*vcn.DefaultSecurityListId,
				*securityList.Id,
			},
			ProhibitPublicIpOnVnic: common.Bool(false),
			RouteTableId:           vcn.DefaultRouteTableId,
			DhcpOptionsId:          vcn.DefaultDhcpOptionsId,
			DisplayName:            common.String(fmt.Sprintf("%s-subnet", args.MinecraftResource.GetName())),
			DnsLabel:               common.String(fmt.Sprintf("subnet%s", common2.InstanceTag)),
		},
	}
	subnet, err := o.network.CreateSubnet(ctx, subnetRequest)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Oracle subnet created", "subnet", subnet)

	internetGatewayRequest := core.CreateInternetGatewayRequest{
		CreateInternetGatewayDetails: core.CreateInternetGatewayDetails{
			IsEnabled:     common.Bool(true),
			CompartmentId: compartment.Id,
			VcnId:         vcn.Id,
			DisplayName:   common.String(fmt.Sprintf("%s-gw", args.MinecraftResource.GetName())),
			FreeformTags:  getTags(args.MinecraftResource.GetEdition()),
		},
	}
	internetGateway, err := o.network.CreateInternetGateway(ctx, internetGatewayRequest)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Oracle Internet Gateway created", "internetGateway", internetGateway)

	routeTableResponse, err := o.network.GetRouteTable(ctx, core.GetRouteTableRequest{
		RtId: vcn.DefaultRouteTableId,
	})
	if err != nil {
		return nil, err
	}

	routeTableResponse.RouteRules = append(routeTableResponse.RouteRules, core.RouteRule{
		NetworkEntityId: internetGateway.Id,
		Destination:     common.String("0.0.0.0/0"),
		DestinationType: core.RouteRuleDestinationTypeCidrBlock,
	})

	updateRouteTable, err := o.network.UpdateRouteTable(ctx, core.UpdateRouteTableRequest{
		RtId: routeTableResponse.Id,
		UpdateRouteTableDetails: core.UpdateRouteTableDetails{
			RouteRules: routeTableResponse.RouteRules,
		},
	})
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Oracle Route Table updated", "updateRouteTable", updateRouteTable)

	userData, err := o.tmpl.GetTemplate(args.MinecraftResource, "", minctlTemplate.GetTemplateCloudConfigName(args.MinecraftResource.IsProxyServer()))
	if err != nil {
		return nil, err
	}
	pubKeyFile, err := ioutil.ReadFile(fmt.Sprintf("%s.pub", args.MinecraftResource.GetSSH()))
	if err != nil {
		return nil, err
	}

	launchInstanceRequest := core.LaunchInstanceRequest{
		LaunchInstanceDetails: core.LaunchInstanceDetails{
			AvailabilityDomain: availabilityDomain.Name,
			CompartmentId:      compartment.Id,
			Shape:              common.String(args.MinecraftResource.GetSize()),
			SourceDetails: core.InstanceSourceViaImageDetails{
				ImageId: image.Id,
			},
			DisplayName:  common.String(args.MinecraftResource.GetName()),
			FreeformTags: getTags(args.MinecraftResource.GetEdition()),
			CreateVnicDetails: &core.CreateVnicDetails{
				AssignPublicIp: common.Bool(true),
				SubnetId:       subnet.Id,
				FreeformTags:   getTags(args.MinecraftResource.GetEdition()),
			},
			Metadata: map[string]string{
				"user_data":           base64.StdEncoding.EncodeToString([]byte(userData)),
				"ssh_authorized_keys": string(pubKeyFile),
			},
		},
	}
	launchInstance, err := o.compute.LaunchInstance(ctx, launchInstanceRequest)
	if err != nil {
		return nil, err
	}

	zap.S().Infow("Oracle launching instance", "launchInstance", launchInstance)

	instanceRequest := core.GetInstanceRequest{
		InstanceId: launchInstance.Instance.Id,
		RequestMetadata: helpers.GetRequestMetadataWithCustomizedRetryPolicy(func(r common.OCIOperationResponse) bool {
			if converted, ok := r.Response.(core.GetInstanceResponse); ok {
				return converted.LifecycleState != core.InstanceLifecycleStateRunning
			}
			return true
		}),
	}

	instance, err := o.compute.GetInstance(ctx, instanceRequest)
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Oracle instance launched", "instance", instance)

	vnicAttachments, err := o.compute.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
		CompartmentId: compartment.Id,
		InstanceId:    instance.Id,
	})
	if err != nil {
		return nil, err
	}
	for _, vnicAttachment := range vnicAttachments.Items {
		vnicRequest := core.GetVnicRequest{
			VnicId: vnicAttachment.VnicId,
		}
		vnic, err := o.network.GetVnic(ctx, vnicRequest)
		if err != nil {
			return nil, err
		}
		zap.S().Infow("Oracle get vnic", "vnic", vnic)
		if *vnic.IsPrimary {
			return &automation.RessourceResults{
				ID:       *instance.Id,
				Name:     *instance.DisplayName,
				Region:   *instance.Region,
				PublicIP: *vnic.PublicIp,
				Tags:     strings.Join(getTagKeys(instance.FreeformTags), ","),
			}, err
		}
	}
	return nil, errors.New("no instance created")
}

func (o *OCI) DeleteServer(id string, args automation.ServerArgs) error {
	ctx := context.Background()

	terminateInstance, err := o.compute.TerminateInstance(ctx, core.TerminateInstanceRequest{
		InstanceId:         common.String(id),
		PreserveBootVolume: common.Bool(false),
	})
	if err != nil {
		return err
	}
	zap.S().Infow("Oracle delete instance", "terminateInstance", terminateInstance)

	tenancyOCID, err := common.DefaultConfigProvider().TenancyOCID()
	if err != nil {
		return err
	}

	listCompartmentsRequest := identity.ListCompartmentsRequest{
		CompartmentId: common.String(tenancyOCID),
		Name:          common.String(args.MinecraftResource.GetName()),
		SortBy:        identity.ListCompartmentsSortByTimecreated,
		SortOrder:     identity.ListCompartmentsSortOrderDesc,
	}
	listCompartments, err := o.identity.ListCompartments(ctx, listCompartmentsRequest)
	if err != nil {
		return err
	}
	if len(listCompartments.Items) == 0 {
		zap.S().Error("compartment not found")
		return errors.New("compartment not found")
	}
	compartment := listCompartments.Items[0]

	listVcns, err := o.network.ListVcns(ctx, core.ListVcnsRequest{
		CompartmentId: compartment.Id,
		DisplayName:   common.String(fmt.Sprintf("%s-vcn", args.MinecraftResource.GetName())),
	})
	if err != nil {
		return err
	}
	if len(listVcns.Items) == 0 {
		return errors.New("vnc not found")
	}
	vcn := listVcns.Items[0]
	zap.S().Infow("Oracle get vnc ", "vcn", vcn)

	listSubnets, err := o.network.ListSubnets(ctx, core.ListSubnetsRequest{
		VcnId:         vcn.Id,
		CompartmentId: compartment.Id,
		DisplayName:   common.String(fmt.Sprintf("%s-subnet", args.MinecraftResource.GetName())),
	})
	if err != nil {
		return err
	}

	if len(listSubnets.Items) == 0 {
		return errors.New("subnet not found")
	}
	subnet := listSubnets.Items[0]
	deleteSubnet, err := o.network.DeleteSubnet(ctx, core.DeleteSubnetRequest{
		SubnetId:        subnet.Id,
		RequestMetadata: helpers.GetRequestMetadataWithDefaultRetryPolicy(),
	})
	if err != nil {
		return err
	}
	zap.S().Infow("Oracle subnet deleted", "deleteSubnet", deleteSubnet)

	securityLists, err := o.network.ListSecurityLists(ctx, core.ListSecurityListsRequest{
		VcnId:         vcn.Id,
		CompartmentId: compartment.Id,
		DisplayName:   common.String(fmt.Sprintf("%s-sl", args.MinecraftResource.GetName())),
	})
	if err != nil {
		return err
	}
	if len(securityLists.Items) == 0 {
		return errors.New("security list not found")
	}
	securityList := securityLists.Items[0]
	securityListResponse, err := o.network.DeleteSecurityList(ctx, core.DeleteSecurityListRequest{
		SecurityListId: securityList.Id,
	})
	if err != nil {
		return err
	}
	zap.S().Infow("Oracle security list deleted", "securityListResponse", securityListResponse)

	routeTable, err := o.network.GetRouteTable(ctx, core.GetRouteTableRequest{
		RtId: vcn.DefaultRouteTableId,
	})
	if err != nil {
		return err
	}

	routeTableResponse, err := o.network.UpdateRouteTable(ctx, core.UpdateRouteTableRequest{
		RtId: routeTable.Id,
		UpdateRouteTableDetails: core.UpdateRouteTableDetails{
			RouteRules: []core.RouteRule{},
		},
	})
	if err != nil {
		return err
	}
	zap.S().Infow("Oracle route table updated", "routeTableResponse", routeTableResponse)

	listInternetGateways, err := o.network.ListInternetGateways(ctx, core.ListInternetGatewaysRequest{
		VcnId:         vcn.Id,
		CompartmentId: compartment.Id,
		DisplayName:   common.String(fmt.Sprintf("%s-gw", args.MinecraftResource.GetName())),
	})
	if err != nil {
		return err
	}
	if len(listInternetGateways.Items) == 0 {
		return errors.New("internet gateways not found")
	}
	internetGateway := listInternetGateways.Items[0]
	gateway, err := o.network.DeleteInternetGateway(ctx, core.DeleteInternetGatewayRequest{
		IgId: internetGateway.Id,
	})
	if err != nil {
		return err
	}
	zap.S().Infow("Oracle internet gateway deleted", "gateway", gateway)

	deleteVcn, err := o.network.DeleteVcn(ctx, core.DeleteVcnRequest{
		VcnId: vcn.Id,
	})
	if err != nil {
		return err
	}
	zap.S().Infow("Oracle vcn deleted", "deleteVcn", deleteVcn)

	zap.S().Infow("Oracle compartment will be deleted", "compartment", compartment)
	deleteCompartmentResponse, err := o.identity.DeleteCompartment(ctx, identity.DeleteCompartmentRequest{
		CompartmentId: compartment.Id,
	})
	if err != nil {
		return err
	}
	zap.S().Infow("Oracle compartment deleted", "compartment", deleteCompartmentResponse)
	return nil
}

func (o *OCI) ListServer() ([]automation.RessourceResults, error) {
	ctx := context.Background()
	tenancyOCID, err := common.DefaultConfigProvider().TenancyOCID()
	if err != nil {
		return nil, err
	}
	zap.S().Infow("Oracle get tenancyOCID", "tenancyOCID", tenancyOCID)
	listCompartments, err := o.identity.ListCompartments(ctx, identity.ListCompartmentsRequest{
		CompartmentId: common.String(tenancyOCID),
	})

	if err != nil {
		return nil, err
	}
	var result []automation.RessourceResults
	for _, compartment := range listCompartments.Items {
		for key := range compartment.FreeformTags {
			if key == common2.InstanceTag && compartment.LifecycleState == identity.CompartmentLifecycleStateActive {
				listInstances, err := o.compute.ListInstances(ctx, core.ListInstancesRequest{
					CompartmentId: compartment.Id,
				})

				if err != nil {
					return nil, err
				}
				for _, instance := range listInstances.Items {
					attachments, err := o.compute.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
						CompartmentId: compartment.Id,
						InstanceId:    instance.Id,
					})

					if err != nil {
						return nil, err
					}
					for _, vnicAttachment := range attachments.Items {

						vnicRequest := core.GetVnicRequest{
							VnicId: vnicAttachment.VnicId,
						}
						vnic, err := o.network.GetVnic(ctx, vnicRequest)
						zap.S().Infow("Oracle vnic", "vnic", vnic)
						if err != nil {
							return nil, err
						}
						zap.S().Infow("Oracle get vnic", "vnic", vnic)
						if *vnic.IsPrimary {
							result = append(result, automation.RessourceResults{
								ID:       *instance.Id,
								Name:     *instance.DisplayName,
								Region:   *instance.Region,
								PublicIP: *vnic.PublicIp,
								Tags:     strings.Join(getTagKeys(instance.FreeformTags), ","),
							})
						}
					}
				}
			}
		}
	}
	return result, nil
}

func (o *OCI) UpdateServer(id string, args automation.ServerArgs) error {
	server, err := o.GetServer(id, args)
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

func (o *OCI) UploadPlugin(id string, args automation.ServerArgs, plugin, destination string) error {
	server, err := o.GetServer(id, args)
	if err != nil {
		return err
	}
	remoteCommand := update.NewRemoteServer(args.MinecraftResource.GetSSH(), server.PublicIP, "ubuntu")

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
	zap.S().Infow("Minecraft plugin uploaded", "plugin", plugin, "server", server.Name)
	return nil
}

func (o *OCI) GetServer(id string, args automation.ServerArgs) (*automation.RessourceResults, error) {
	ctx := context.Background()
	instance, err := o.compute.GetInstance(ctx, core.GetInstanceRequest{
		InstanceId: common.String(id),
	})
	if err != nil {
		return nil, err
	}
	vnicAttachments, err := o.compute.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
		CompartmentId: instance.CompartmentId,
		InstanceId:    instance.Id,
	})
	if err != nil {
		return nil, err
	}
	for _, vnicAttachment := range vnicAttachments.Items {
		vnicRequest := core.GetVnicRequest{
			VnicId: vnicAttachment.VnicId,
		}
		vnic, err := o.network.GetVnic(ctx, vnicRequest)
		if err != nil {
			return nil, err
		}
		zap.S().Infow("Oracle get vnic", "vnic", vnic)
		if *vnic.IsPrimary {
			return &automation.RessourceResults{
				ID:       *instance.Id,
				Name:     *instance.DisplayName,
				Region:   *instance.Region,
				PublicIP: *vnic.PublicIp,
				Tags:     strings.Join(getTagKeys(instance.FreeformTags), ","),
			}, err
		}
	}
	return nil, errors.New("no instance found")
}
