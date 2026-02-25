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

func tableAlicloudEcsRegion(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_ecs_region",
		Description: "Elastic Compute Region",
		List: &plugin.ListConfig{
			Hydrate: listEcsRegions,
			Tags:    map[string]string{"service": "ecs", "action": "DescribeRegions"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegionId"),
				Description: ColumnDescriptionRegion,
			},
			{
				Name:        "local_name",
				Type:        proto.ColumnType_STRING,
				Description: "The local name of the region.",
			},
			{
				Name:        "region_endpoint",
				Type:        proto.ColumnType_STRING,
				Description: "The endpoint of the region.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates whether the cluster is sold out.",
			},

			// steampipe standard columns
			{
				Name:        "akas",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRegionAkas,
				Transform:   transform.FromValue(),
				Description: ColumnDescriptionAkas,
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegionId"),
				Description: ColumnDescriptionTitle,
			},

			// alicloud common columns
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

func listEcsRegions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := GetDefaultRegion(d.Connection)

	// Create service connection
	client, err := ECSRegionService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ecs.listEcsRegions", "connection_error", err)
		return nil, err
	}
	request := &ecs.DescribeRegionsRequest{
		AcceptLanguage: tea.String("en-US"),
	}

	response, err := client.DescribeRegions(request)
	if err != nil {
		logQueryError(ctx, d, h, "alicloud_ecs.listEcsRegions", err, "request", request)
		return nil, err
	}
	for _, i := range response.Body.Regions.Region {
		d.StreamListItem(ctx, *i)
		// This will return zero if context has been cancelled (i.e due to manual cancellation) or
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

func getRegionAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRegionAkas")
	data := h.Item.(ecs.DescribeRegionsResponseBodyRegionsRegion)

	// Get project details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*alicloudCommonColumnData)
	accountID := commonColumnData.AccountID

	return []string{"acs:ecs::" + accountID + ":region/" + tea.StringValue(data.RegionId)}, nil
}
