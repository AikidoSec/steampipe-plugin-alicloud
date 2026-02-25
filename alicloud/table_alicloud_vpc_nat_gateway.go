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

func tableAlicloudVpcNatGateway(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_vpc_nat_gateway",
		Description: "Aliclod VPC NAT Gateway",
		List: &plugin.ListConfig{
			Hydrate: listVpcNatGateways,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeNatGateways"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("nat_gateway_id"),
			Hydrate:    getVpcNatGateway,
			Tags:       map[string]string{"service": "vpc", "action": "DescribeNatGateways"},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the NAT gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "nat_gateway_id",
				Description: "The ID of the NAT gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "nat_type",
				Description: "The type of the NAT gateway. Valid values: 'Normal' and 'Enhanced'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The state of the NAT gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "deletion_protection",
				Description: "Indicates whether deletion protection is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "auto_pay",
				Description: "Indicates whether auto pay is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "billing_method",
				Description: "The billing method of the NAT gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("InstanceChargeType"),
			},
			{
				Name:        "business_status",
				Description: "The status of the NAT gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time when the NAT gateway was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description of the NAT gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ecs_metric_enabled",
				Description: "Indicates whether the traffic monitoring feature is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "expired_ime",
				Description: "The time when the NAT gateway expires.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "internet_charge_type",
				Description: "The billing method of the NAT gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_group_id",
				Description: "The ID of the resource group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "spec",
				Description: "The size of the NAT gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the virtual private cloud (VPC) to which the NAT gateway belongs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "forward_table_ids",
				Description: "The ID of the Destination Network Address Translation (DNAT) table.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ForwardTableIds.ForwardTableId"),
			},
			{
				Name:        "ip_lists",
				Description: "The elastic IP address (EIP) that is associated with the NAT gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("IpLists.IpList"),
			},
			{
				Name:        "nat_gateway_private_info",
				Description: "The information of the virtual private cloud (VPC) to which the enhanced NAT gateway belongs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "snat_table_ids",
				Description: "The ID of the SNAT table for the NAT gateway.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SnatTableIds.SnatTableId"),
			},

			// steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(vpcNatGatewayTitle),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcNatGatewayAka,
				Transform:   transform.FromValue(),
			},

			// alicloud common columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegionId"),
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

func listVpcNatGateways(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := VpcService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_vpc_nat_gateway.listVpcNatGateways", "connection_error", err)
		return nil, err
	}

	request := &vpc.DescribeNatGatewaysRequest{
		PageSize:   tea.Int32(50),
		PageNumber: tea.Int32(1),
		RegionId:   tea.String(d.EqualsQualString(matrixKeyRegion)),
	}

	count := 0
	for {
		d.WaitForListRateLimit(ctx)
		response, err := client.DescribeNatGateways(request)
		if err != nil {
			logQueryError(ctx, d, h, "alicloud_vpc_nat_gateway.listVpcNatGateways", err, "request", request)
			return nil, err
		}
		for _, i := range response.Body.NatGateways.NatGateway {
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

func getVpcNatGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcNatGateway")

	// Create service connection
	client, err := VpcService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_vpc_nat_gateway.getVpcNatGateway", "connection_error", err)
		return nil, err
	}
	id := d.EqualsQuals["nat_gateway_id"].GetStringValue()

	request := &vpc.DescribeNatGatewaysRequest{
		NatGatewayId: &id,
	}

	response, err := client.DescribeNatGateways(request)
	if err != nil {
		logQueryError(ctx, d, h, "alicloud_vpc_nat_gateway.getVpcNatGateway", err, "request", request)
		return nil, err
	}

	if len(response.Body.NatGateways.NatGateway) > 0 {
		return *response.Body.NatGateways.NatGateway[0], nil
	}

	return nil, nil
}

func getVpcNatGatewayAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcNatGatewayAka")
	ngw := h.Item.(vpc.DescribeNatGatewaysResponseBodyNatGatewaysNatGateway)

	// Get project details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*alicloudCommonColumnData)
	accountID := commonColumnData.AccountID

	akas := []string{"acs:vpc:" + tea.StringValue(ngw.RegionId) + ":" + accountID + ":natgateway/" + tea.StringValue(ngw.NatGatewayId)}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func vpcNatGatewayTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(vpc.DescribeNatGatewaysResponseBodyNatGatewaysNatGateway)

	// Build resource title
	title := tea.StringValue(data.NatGatewayId)

	if len(tea.StringValue(data.Name)) > 0 {
		title = *data.Name
	}

	return title, nil
}
