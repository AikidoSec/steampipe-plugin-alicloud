package alicloud

import (
	"context"
	"strconv"

	"github.com/alibabacloud-go/tea/tea"
	vpc "github.com/alibabacloud-go/vpc-20160428/v7/client"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAlicloudVpcVSwitch(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_vpc_vswitch",
		Description: "VSwitches to divide the VPC network into one or more subnets.",
		List: &plugin.ListConfig{
			Hydrate: listVSwitch,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeVSwitches"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getVSwitchAttributes,
				Tags: map[string]string{"service": "vpc", "action": "DescribeVSwitchAttributes"},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			// Top columns
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VSwitchName"),
				Description: "The name of the VPC.",
			},
			{
				Name:        "vswitch_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VSwitchId"),
				Description: "The unique ID of the VPC.",
			},
			{
				Name:        "vpc_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the VPC to which the VSwitch belongs.",
			},
			// Other columns
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "The status of the VPC. Pending: The VPC is being configured. Available: The VPC is available.",
			},
			{
				Name:        "cidr_block",
				Type:        proto.ColumnType_CIDR,
				Description: "The IPv4 CIDR block of the VPC.",
			},
			{
				Name:        "ipv6_cidr_block",
				Type:        proto.ColumnType_CIDR,
				Transform:   transform.FromField("Ipv6CidrBlock"),
				Description: "The IPv6 CIDR block of the VPC.",
			},
			{
				Name:        "zone_id",
				Type:        proto.ColumnType_STRING,
				Description: "The zone to which the VSwitch belongs.",
			},
			{
				Name:        "available_ip_address_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of available IP addresses in the VSwitch.",
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of the VPC.",
			},
			{
				Name:        "creation_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The creation time of the VPC.",
			},
			{
				Name:        "is_default",
				Type:        proto.ColumnType_BOOL,
				Description: "True if the VPC is the default VPC in the region.",
			},
			{
				Name:        "resource_group_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the resource group to which the VPC belongs.",
			},
			{
				Name:        "network_acl_id",
				Type:        proto.ColumnType_STRING,
				Description: "A list of IDs of NAT Gateways.",
			},
			{
				Name:        "owner_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the owner of the VPC.",
			},
			{
				Name:        "share_type",
				Type:        proto.ColumnType_STRING,
				Description: "",
			},
			{
				Name:        "route_table",
				Type:        proto.ColumnType_JSON,
				Description: "Details of the route table.",
			},
			{
				Name:        "cloud_resources",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVSwitchAttributes,
				Transform:   transform.FromField("CloudResourceSetType"),
				Description: "The list of resources in the VSwitch.",
			},
			{
				Name:        "tags_src",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags.Tag"),
				Description: ColumnDescriptionTags,
			},

			//  steampipe common columns
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags.Tag").Transform(modifyGenericSourceTags),
				Description: ColumnDescriptionTags,
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(vswitchTitle),
				Description: ColumnDescriptionTitle,
			},
			{
				Name:        "akas",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(vswitchAkas),
				Description: ColumnDescriptionAkas,
			},

			//  alicloud common columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ZoneId").Transform(zoneToRegion),
			},
			{
				Name:        "account_id",
				Description: ColumnDescriptionAccount,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("OwnerId"),
			},
		},
	}
}

//// LIST FUNCTION

func listVSwitch(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := VpcService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listVSwitch", "connection_error", err)
		return nil, err
	}
	request := &vpc.DescribeVSwitchesRequest{
		PageSize:   tea.Int32(50),
		PageNumber: tea.Int32(1),
		RegionId:   tea.String(d.EqualsQualString(matrixKeyRegion)),
	}

	quals := d.EqualsQuals
	if quals["is_default"] != nil {
		request.IsDefault = tea.Bool(quals["is_default"].GetBoolValue())
	}
	if quals["vswitch_id"] != nil {
		request.VSwitchId = tea.String(quals["vswitch_id"].GetStringValue())
	}

	count := 0
	for {
		d.WaitForListRateLimit(ctx)
		response, err := client.DescribeVSwitches(request)
		if err != nil {
			logQueryError(ctx, d, h, "listVSwitch", err, "request", request)
			return nil, err
		}
		for _, i := range response.Body.VSwitches.VSwitch {
			plugin.Logger(ctx).Warn("listVSwitch", "tags", i.Tags, "item", i)
			d.StreamListItem(ctx, *i)
			// This will return zero if context has been cancelled (i.e due to manual cancellation) or
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
			count++
		}
		if count >= int(tea.Int32Value(response.Body.TotalCount)) {
			break
		}
		request.PageNumber = tea.Int32(tea.Int32Value(response.Body.PageNumber) + 1)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getVSwitchAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := VpcService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getVSwitchAttributes", "connection_error", err)
		return nil, err
	}
	i := h.Item.(vpc.DescribeVSwitchesResponseBodyVSwitchesVSwitch)
	request := &vpc.DescribeVSwitchAttributesRequest{
		VSwitchId: i.VSwitchId,
	}
	response, err := client.DescribeVSwitchAttributes(request)
	if err != nil {
		logQueryError(ctx, d, h, "getVSwitchAttributes", err, "request", request)
		return nil, err
	}
	return *response.Body, nil
}

//// TRANSFORM FUNCTIONS

func vswitchAkas(_ context.Context, d *transform.TransformData) (interface{}, error) {
	i := d.HydrateItem.(vpc.DescribeVSwitchesResponseBodyVSwitchesVSwitch)
	return []string{"acs:vswitch:" + tea.StringValue(i.ZoneId) + ":" + strconv.FormatInt(tea.Int64Value(i.OwnerId), 10) + ":vswitch/" + tea.StringValue(i.VSwitchId)}, nil
}

func vswitchTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	i := d.HydrateItem.(vpc.DescribeVSwitchesResponseBodyVSwitchesVSwitch)

	// Build resource title
	title := tea.StringValue(i.VSwitchId)
	if len(tea.StringValue(i.VSwitchName)) > 0 {
		title = tea.StringValue(i.VSwitchName)
	}

	return title, nil
}
