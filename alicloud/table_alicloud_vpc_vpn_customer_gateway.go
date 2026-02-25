package alicloud

import (
	"context"

	"github.com/alibabacloud-go/tea/tea"
	vpc "github.com/alibabacloud-go/vpc-20160428/v7/client"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAlicloudVpcVpnCustomerGateway(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_vpc_vpn_customer_gateway",
		Description: "Alicloud VPC VPN Customer Gateway.",
		List: &plugin.ListConfig{
			Hydrate: listVpcCustomerGateways,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeCustomerGateways"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("customer_gateway_id"),
			Hydrate:    getVpcCustomerGateway,
			Tags:       map[string]string{"service": "vpc", "action": "DescribeCustomerGateways"},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the customer gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "customer_gateway_id",
				Description: "The ID of the customer gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "asn",
				Description: "Specifies the ASN of the customer gateway.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "create_time",
				Description: "The time when the customer gateway was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreateTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "description",
				Description: "The description of the customer gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_address",
				Description: "The IP address of the customer gateway.",
				Type:        proto.ColumnType_IPADDR,
			},

			// Steampipe standard columns
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcCustomerGatewayAka,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(vpcCustomerGatewayTitle),
			},

			// Alicloud standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVpnCustomerGatewayRegion,
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

func listVpcCustomerGateways(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := VpcService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_vpc_vpn_customer_gateway.listVpcCustomerGateways", "connection_error", err)
		return nil, err
	}
	request := &vpc.DescribeCustomerGatewaysRequest{
		PageSize:   tea.Int32(50),
		PageNumber: tea.Int32(1),
		RegionId:   tea.String(d.EqualsQualString(matrixKeyRegion)),
	}

	count := 0
	for {
		d.WaitForListRateLimit(ctx)
		response, err := client.DescribeCustomerGateways(request)
		if err != nil {
			logQueryError(ctx, d, h, "alicloud_vpc_vpn_customer_gateway.listVpcCustomerGateways", err, "request", request)
			return nil, err
		}
		for _, i := range response.Body.CustomerGateways.CustomerGateway {
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

func getVpcCustomerGateway(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := VpcService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_vpc_vpn_customer_gateway.getVpcCustomerGateway", "connection_error", err)
		return nil, err
	}
	id := d.EqualsQuals["customer_gateway_id"].GetStringValue()

	request := &vpc.DescribeCustomerGatewaysRequest{
		CustomerGatewayId: &id,
	}

	response, err := client.DescribeCustomerGateways(request)
	if err != nil {
		logQueryError(ctx, d, h, "alicloud_vpc_vpn_customer_gateway.getVpcCustomerGateway", err, "request", request)
		return nil, err
	}

	if len(response.Body.CustomerGateways.CustomerGateway) > 0 {
		return *response.Body.CustomerGateways.CustomerGateway[0], nil
	}

	return nil, nil
}

func getVpcCustomerGatewayAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcCustomerGatewayAka")
	data := h.Item.(vpc.DescribeCustomerGatewaysResponseBodyCustomerGatewaysCustomerGateway)
	region := d.EqualsQualString(matrixKeyRegion)

	// Get project details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*alicloudCommonColumnData)
	accountID := commonColumnData.AccountID

	akas := []string{"acs:vpc:" + region + ":" + accountID + ":customergateway/" + tea.StringValue(data.CustomerGatewayId)}

	return akas, nil
}

func getVpnCustomerGatewayRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpnCustomerGatewayRegion")
	region := d.EqualsQualString(matrixKeyRegion)

	return region, nil
}

//// TRANSFORM FUNCTIONS

func vpcCustomerGatewayTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(vpc.DescribeCustomerGatewaysResponseBodyCustomerGatewaysCustomerGateway)

	// Build resource title
	title := tea.StringValue(data.CustomerGatewayId)

	if len(tea.StringValue(data.Name)) > 0 {
		title = tea.StringValue(data.Name)
	}

	return title, nil
}
