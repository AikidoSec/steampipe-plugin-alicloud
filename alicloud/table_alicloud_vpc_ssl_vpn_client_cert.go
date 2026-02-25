package alicloud

import (
	"context"

	"github.com/alibabacloud-go/tea/tea"
	vpc "github.com/alibabacloud-go/vpc-20160428/v7/client"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type vpnSslClientCertInfo = struct {
	Name               string
	SslVpnClientCertId string
	SslVpnServerId     string
	Status             string
	CreateTime         int64
	EndTime            int64
	CaCert             string
	ClientCert         string
	ClientKey          string
	ClientConfig       string
	Region             string
}

//// TABLE DEFINITION

func tableAlicloudVpcSslVpnClientCert(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_vpc_ssl_vpn_client_cert",
		Description: "SSL Client is responsible for managing client certificates. The client needs to first complete certificate verification in order to connect to the SSL Server.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("ssl_vpn_client_cert_id"),
			Hydrate:    getVpcSslVpnClientCert,
			Tags:       map[string]string{"service": "vpc", "action": "DescribeSslVpnClientCert"},
		},
		List: &plugin.ListConfig{
			Hydrate: listVpcSslVpnClientCerts,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeSslVpnClientCerts"},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the SSL client certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssl_vpn_client_cert_id",
				Description: "The ID of the SSL client certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssl_vpn_server_id",
				Description: "The ID of the SSL-VPN server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the client certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_time",
				Description: "The time when the SSL client certificate was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CreateTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "end_time",
				Description: "The time when the SSL client certificate expires.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EndTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "ca_cert",
				Description: "The CA certificate.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVpcSslVpnClientCert,
			},
			{
				Name:        "client_cert",
				Description: "The client certificate.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVpcSslVpnClientCert,
			},
			{
				Name:        "client_key",
				Description: "The client key.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVpcSslVpnClientCert,
			},
			{
				Name:        "client_config",
				Description: "The client configuration.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVpcSslVpnClientCert,
			},

			// steampipe standard columns
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVpcSslVpnClientCertCertAka,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(sslVpnClientCertTitle),
			},

			// alicloud standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
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

func listVpcSslVpnClientCerts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := VpcService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_vpc_vpn_ssl_client.listVpcSslVpnClientCerts", "connection_error", err)
		return nil, err
	}
	request := &vpc.DescribeSslVpnClientCertsRequest{
		PageSize:   tea.Int32(50),
		PageNumber: tea.Int32(1),
		RegionId:   tea.String(d.EqualsQualString(matrixKeyRegion)),
	}

	count := 0
	for {
		d.WaitForListRateLimit(ctx)
		response, err := client.DescribeSslVpnClientCerts(request)
		if err != nil {
			logQueryError(ctx, d, h, "alicloud_vpc_vpn_ssl_client.listVpcSslVpnClientCerts", err, "request", request)
			return nil, err
		}
		for _, i := range response.Body.SslVpnClientCertKeys.SslVpnClientCertKey {
			d.StreamListItem(ctx, vpnSslClientCertInfo{
				tea.StringValue(i.Name),
				tea.StringValue(i.SslVpnClientCertId),
				tea.StringValue(i.SslVpnServerId),
				tea.StringValue(i.Status),
				tea.Int64Value(i.CreateTime),
				tea.Int64Value(i.EndTime),
				"",
				"",
				"",
				"",
				tea.StringValue(i.RegionId),
			})
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

func getVpcSslVpnClientCert(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcSslVpnClientCert")

	// Create service connection
	client, err := VpcService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_vpc_vpn_ssl_client.getVpcSslVpnClientCert", "connection_error", err)
		return nil, err
	}

	var id string
	if h.Item != nil {
		data := h.Item.(vpnSslClientCertInfo)
		id = data.SslVpnClientCertId
	} else {
		id = d.EqualsQuals["ssl_vpn_client_cert_id"].GetStringValue()
	}

	request := &vpc.DescribeSslVpnClientCertRequest{
		SslVpnClientCertId: &id,
	}

	data, err := client.DescribeSslVpnClientCert(request)
	if err != nil {
		logQueryError(ctx, d, h, "alicloud_vpc_vpn_ssl_client.getVpcSslVpnClientCert", err, "request", request)
		return nil, err
	}

	return vpnSslClientCertInfo{
		tea.StringValue(data.Body.Name),
		tea.StringValue(data.Body.SslVpnClientCertId),
		tea.StringValue(data.Body.SslVpnServerId),
		tea.StringValue(data.Body.Status),
		tea.Int64Value(data.Body.CreateTime),
		tea.Int64Value(data.Body.EndTime),
		tea.StringValue(data.Body.CaCert),
		tea.StringValue(data.Body.ClientCert),
		tea.StringValue(data.Body.ClientKey),
		tea.StringValue(data.Body.ClientConfig),
		tea.StringValue(data.Body.RegionId),
	}, nil
}

func getVpcSslVpnClientCertCertAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVpcSslVpnClientCertCertAka")

	data := h.Item.(vpnSslClientCertInfo)

	// Get project details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*alicloudCommonColumnData)
	accountID := commonColumnData.AccountID

	akas := []string{"acs:vpc:" + data.Region + ":" + accountID + ":sslclientcert/" + data.SslVpnClientCertId}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func sslVpnClientCertTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(vpnSslClientCertInfo)

	// Build resource title
	title := data.SslVpnClientCertId

	if len(data.Name) > 0 {
		title = data.Name
	}

	return title, nil
}
