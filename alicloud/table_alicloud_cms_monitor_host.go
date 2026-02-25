package alicloud

import (
	"context"

	cms "github.com/alibabacloud-go/cms-20190101/v10/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAlicloudCmsMonitorHost(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_cms_monitor_host",
		Description: "Alicloud Cloud Monitor Host",
		List: &plugin.ListConfig{
			Hydrate: listCmsMonitorHosts,
			Tags:    map[string]string{"service": "cms", "action": "DescribeMonitoringAgentHosts"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"host_name", "instance_id"}),
			Hydrate:    getCmsMonitorHost,
			Tags:       map[string]string{"service": "cms", "action": "DescribeMonitoringAgentHosts"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCmsMonitoringAgentStatus,
				Tags: map[string]string{"service": "cms", "action": "DescribeMonitoringAgentStatuses"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "host_name",
				Description: "The name of the host.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_id",
				Description: "The ID of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_type_family",
				Description: "The type of the ECS instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "agent_version",
				Description: "The version of the Cloud Monitor agent.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_aliyun_host",
				Description: "Indicates whether the host is provided by Alibaba Cloud.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "eip_id",
				Description: "The ID of the EIP.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "eip_address",
				Description: "The elastic IP address (EIP) of the host.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ali_uid",
				Description: "The ID of the Alibaba Cloud account.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "ip_group",
				Description: "The IP address of the host.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "nat_ip",
				Description: "The IP address of the Network Address Translation (NAT) gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network_type",
				Description: "The type of the network.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "operating_system",
				Description: "The operating system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "serial_number",
				Type:        proto.ColumnType_STRING,
				Description: "The serial number of the host. A host that is not provided by Alibaba Cloud has a serial number instead of an instance ID.",
			},
			{
				Name:        "monitoring_agent_status",
				Description: "The status of the Cloud Monitor agent.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCmsMonitoringAgentStatus,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("HostName"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCmsMonitoringHostAka,
				Transform:   transform.FromValue(),
			},

			// Alicloud standard columns
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

func listCmsMonitorHosts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := CmsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCmsMonitorHosts", "connection_error", err)
		return nil, err
	}
	request := &cms.DescribeMonitoringAgentHostsRequest{
		PageSize:   tea.Int32(50),
		PageNumber: tea.Int32(1),
	}

	count := 0
	for {
		d.WaitForListRateLimit(ctx)
		response, err := client.DescribeMonitoringAgentHosts(request)
		if err != nil {
			logQueryError(ctx, d, h, "listCmsMonitorHosts", err, "request", request)
			return nil, err
		}
		for _, host := range response.Body.Hosts.Host {
			plugin.Logger(ctx).Warn("listCmsMonitorHosts", "item", host)
			d.StreamListItem(ctx, *host)
			// This will return zero if context has been cancelled (i.e due to manual cancellation) or
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
			count++
		}
		if count >= int(*response.Body.Total) {
			break
		}
		request.SetPageNumber(*response.Body.PageNumber + 1)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCmsMonitorHost(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCmsMonitorHost")
	// Create service connection
	client, err := CmsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getCmsMonitorHost", "connection_error", err)
		return nil, err
	}

	hostName := d.EqualsQuals["host_name"].GetStringValue()
	instanceId := d.EqualsQuals["instance_id"].GetStringValue()

	// handle empty hostName or instanceId in get call
	if hostName == "" || instanceId == "" {
		return nil, nil
	}

	request := &cms.DescribeMonitoringAgentHostsRequest{
		HostName:    &hostName,
		InstanceIds: &instanceId,
	}

	response, err := client.DescribeMonitoringAgentHosts(request)
	if err != nil {
		logQueryError(ctx, d, h, "getCmsMonitorHost", err, "request", request)
		return nil, err
	}

	return *response.Body.Hosts.Host[0], nil
}

func getCmsMonitoringAgentStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCmsMonitoringAgentStatus")

	// Create service connection
	client, err := CmsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getCmsMonitoringAgentStatus", "connection_error", err)
		return nil, err
	}

	request := &cms.DescribeMonitoringAgentStatusesRequest{
		InstanceIds: h.Item.(cms.DescribeMonitoringAgentHostsResponseBodyHostsHost).InstanceId,
	}

	response, err := client.DescribeMonitoringAgentStatuses(request)
	if err != nil {
		logQueryError(ctx, d, h, "getCmsMonitoringAgentStatus", err, "request", request)
		return nil, err
	}

	return response.Body.NodeStatusList.NodeStatus, nil
}

func getCmsMonitoringHostAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCmsMonitoringHostAka")

	data := h.Item.(cms.DescribeMonitoringAgentHostsResponseBodyHostsHost)

	// Get project details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*alicloudCommonColumnData)
	accountID := commonColumnData.AccountID

	akas := []string{"arn:acs:cms:" + tea.StringValue(data.Region) + ":" + accountID + ":host/" + tea.StringValue(data.HostName)}

	return akas, nil
}
