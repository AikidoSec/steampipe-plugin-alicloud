package alicloud

import (
	"context"

	ecs "github.com/alibabacloud-go/ecs-20140526/v7/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAlicloudEcsEni(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_ecs_network_interface",
		Description: "Alicloud ECS Network Interface.",
		List: &plugin.ListConfig{
			Hydrate: listEcsEni,
			Tags:    map[string]string{"service": "ecs", "action": "DescribeNetworkInterfaces"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("network_interface_id"),
			Hydrate:    getEcsEni,
			Tags:       map[string]string{"service": "ecs", "action": "DescribeNetworkInterfaces"},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the ENI.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NetworkInterfaceName"),
			},
			{
				Name:        "network_interface_id",
				Description: "An unique identifier for the ENI.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the ENI. Valid values: 'Primary' and 'Secondary'",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the ENI.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The ID of the account that owns the ENI.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_managed",
				Description: "Indicates whether the user is an Alibaba Cloud service or a distributor.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "creation_time",
				Description: "The time when the ENI was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description of the ENI.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_id",
				Description: "The ID of the instance to which the ENI is bound.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mac_address",
				Description: "The MAC address of the ENI.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_ip_address",
				Description: "The private IP address of the ENI.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "queue_number",
				Description: "The number of queues supported by the ENI.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "resource_group_id",
				Description: "The ID of the resource group to which the ENI belongs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_id",
				Description: "The ID of the distributor to which the ENI belongs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceID"),
			},
			{
				Name:        "vswitch_id",
				Description: "The ID of the VSwitch to which the ENI is connected.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VSwitchId"),
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC to which the ENI belongs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "zone_id",
				Description: "The zone ID of the ENI.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "associated_public_ip_address",
				Description: "The public IP address of the instance.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("AssociatedPublicIp.PublicIpAddress"),
			},
			{
				Name:        "associated_public_ip_allocation_id",
				Description: "The allocation ID of the EIP.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssociatedPublicIp.AllocationId"),
			},
			{
				Name:        "attachment",
				Description: "Attachments of the ENI",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ipv6_sets",
				Description: "The IPv6 addresses assigned to the ENI.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Ipv6Sets.Ipv6Set"),
			},
			{
				Name:        "private_ip_sets",
				Description: "The private IP addresses of the ENI.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PrivateIpSets.PrivateIpSet"),
			},
			{
				Name:        "security_group_ids",
				Description: "The IDs of the security groups to which the ENI belongs.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SecurityGroupIds.SecurityGroupId"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached with the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags.Tag").Transform(modifyGenericSourceTags),
			},

			// steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags.Tag").Transform(genericTagsToMap),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(ecsEniAka),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(ecsEniTitle),
			},

			// alibaba standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ZoneId").Transform(zoneToRegion),
			},
			{
				Name:        "account_id",
				Description: "The alicloud Account ID in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OwnerId"),
			},
		},
	}
}

//// LIST FUNCTION

func listEcsEni(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := ECSService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ecs_network_interface.listEcsEni", "connection_error", err)
		return nil, err
	}
	request := &ecs.DescribeNetworkInterfacesRequest{
		MaxResults: tea.Int32(100),
		RegionId:   tea.String(d.EqualsQualString(matrixKeyRegion)),
	}

	// If the request no of items is less than the paging max limit
	// update limit to the requested no of results.
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.MaxResults) {
			request.MaxResults = tea.Int32(int32(*limit))
		}
	}

	pageLeft := true
	for pageLeft {
		d.WaitForListRateLimit(ctx)
		response, err := client.DescribeNetworkInterfaces(request)
		if err != nil {
			logQueryError(ctx, d, h, "alicloud_ecs_network_interface.listEcsEni", err, "request", request)
			return nil, err
		}
		for _, eni := range response.Body.NetworkInterfaceSets.NetworkInterfaceSet {
			d.StreamListItem(ctx, *eni)
			// This will return zero if context has been cancelled (i.e due to manual cancellation) or
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if tea.StringValue(response.Body.NextToken) != "" {
			request.NextToken = response.Body.NextToken
		} else {
			pageLeft = false
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEcsEni(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEcsEni")

	// Create service connection
	client, err := ECSService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ecs_network_interface.getEcsEni", "connection_error", err)
		return nil, err
	}
	id := d.EqualsQuals["network_interface_id"].GetStringValue()

	request := &ecs.DescribeNetworkInterfacesRequest{
		NetworkInterfaceId: []*string{&id},
	}

	response, err := client.DescribeNetworkInterfaces(request)
	if serverErr, ok := err.(*tea.SDKError); ok {
		logQueryError(ctx, d, h, "alicloud_ecs_network_interface.getEcsEni", serverErr, "request", request)
		return nil, serverErr
	}

	if len(response.Body.NetworkInterfaceSets.NetworkInterfaceSet) > 0 {
		return *response.Body.NetworkInterfaceSets.NetworkInterfaceSet[0], nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func ecsEniAka(_ context.Context, d *transform.TransformData) (interface{}, error) {
	eni := d.HydrateItem.(ecs.DescribeNetworkInterfacesResponseBodyNetworkInterfaceSetsNetworkInterfaceSet)
	akas := []string{"acs:ecs:" + tea.StringValue(eni.ZoneId) + ":" + tea.StringValue(eni.OwnerId) + ":eni/" + tea.StringValue(eni.NetworkInterfaceId)}

	return akas, nil
}

func ecsEniTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	eni := d.HydrateItem.(ecs.DescribeNetworkInterfacesResponseBodyNetworkInterfaceSetsNetworkInterfaceSet)

	// Build resource title
	title := tea.StringValue(eni.NetworkInterfaceId)

	if len(tea.StringValue(eni.NetworkInterfaceName)) > 0 {
		title = tea.StringValue(eni.NetworkInterfaceName)
	}

	return title, nil
}
