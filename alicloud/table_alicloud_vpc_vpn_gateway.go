package alicloud

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/alibabacloud-go/tea/tea"
	vpc "github.com/alibabacloud-go/vpc-20160428/v7/client"
)

//// TABLE DEFINITION

func tableAlicloudVpcVpnGateway(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_vpc_vpn_gateway",
		Description: "Alicloud VPC VPN Gateway.",
		List: &plugin.ListConfig{
			Hydrate: listVpcVpnGateways,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeVpnGateways"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("vpn_gateway_id"),
			Hydrate:    getVpcVpnGateway,
			Tags:       map[string]string{"service": "vpc", "action": "DescribeVpnGateways"},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the VPN gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpn_gateway_id",
				Description: "The ID of the VPN gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the VPN gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auto_propagate",
				Description: "Indicates whether auto propagate is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "billing_method",
				Description: "The billing method of the VPN gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ChargeType"),
			},
			{
				Name:        "business_status",
				Description: "The business state of the VPN gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time when the VPN gateway was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreateTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "description",
				Description: "The description of the VPN gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enable_bgp",
				Description: "Indicates whether bgp is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "end_time",
				Description: "The creation time of the VPC.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EndTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "internet_ip",
				Description: "The public IP address of the VPN gateway.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "ipsec_vpn",
				Description: "Indicates whether the IPsec-VPN feature is enabled.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "spec",
				Description: "The maximum bandwidth of the VPN gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssl_max_connections",
				Description: "The maximum number of concurrent SSL-VPN connections.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "ssl_vpn",
				Description: "Indicates whether the SSL-VPN feature is enabled.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tag",
				Description: "The tag of the VPN gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vswitch_id",
				Description: "The ID of the VSwitch to which the VPN gateway belongs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VSwitchId"),
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC for which the VPN gateway is created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VpcId"),
			},
			{
				Name:        "reservation_data",
				Description: "A set of reservation details.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached with the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags.Tag"),
			},

			// steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags.Tag").Transform(modifyGenericSourceTags),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcVpnGatewayAka,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(vpcVpnGatewayTitle),
			},

			// alicloud common columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVpnGatewayRegion,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "account_id",
				Description: ColumnDescriptionAccount,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCommonColumns,
				Transform:   transform.FromField("AccountID"),
			},
		},
	}
}

//// LIST FUNCTION

func listVpcVpnGateways(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := VpcService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_vpc_vpn_gateway.listVpcVpnGateways", "connection_error", err)
		return nil, err
	}
	request := &vpc.DescribeVpnGatewaysRequest{
		PageSize:   tea.Int32(50),
		PageNumber: tea.Int32(1),
		RegionId:   tea.String(d.EqualsQualString(matrixKeyRegion)),
	}

	count := 0
	for {
		d.WaitForListRateLimit(ctx)
		response, err := client.DescribeVpnGateways(request)
		if err != nil {
			logQueryError(ctx, d, h, "alicloud_vpc_vpn_gateway.listVpcVpnGateways", err, "request", request)
			return nil, err
		}
		for _, i := range response.Body.VpnGateways.VpnGateway {
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

func getVpcVpnGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := VpcService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_vpc_vpn_gateway.getVpcVpnGateway", "connection_error", err)
		return nil, err
	}
	id := d.EqualsQuals["vpn_gateway_id"].GetStringValue()

	request := &vpc.DescribeVpnGatewaysRequest{
		VpnGatewayId: &id,
	}

	response, err := client.DescribeVpnGateways(request)
	if err != nil {
		logQueryError(ctx, d, h, "alicloud_vpc_vpn_gateway.getVpcVpnGateway", err, "request", request)
		return nil, err
	}

	if len(response.Body.VpnGateways.VpnGateway) > 0 {
		return *response.Body.VpnGateways.VpnGateway[0], nil
	}
	return nil, nil
}

func getVpcVpnGatewayAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcVpnGatewayAka")
	data := h.Item.(vpc.DescribeVpnGatewaysResponseBodyVpnGatewaysVpnGateway)
	region := d.EqualsQualString(matrixKeyRegion)

	// Get project details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*alicloudCommonColumnData)
	accountID := commonColumnData.AccountID

	akas := []string{"acs:vpc:" + region + ":" + accountID + ":vpngateway/" + tea.StringValue(data.VpnGatewayId)}

	return akas, nil
}

func getVpnGatewayRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpnGatewayRegion")
	region := d.EqualsQualString(matrixKeyRegion)

	return region, nil
}

//// TRANSFORM FUNCTIONS

func vpcVpnGatewayTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(vpc.DescribeVpnGatewaysResponseBodyVpnGatewaysVpnGateway)

	// Build resource title
	title := tea.StringValue(data.VpnGatewayId)

	if len(tea.StringValue(data.Name)) > 0 {
		title = tea.StringValue(data.Name)
	}

	return title, nil
}
