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

func tableAlicloudEcskeyPair(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_ecs_key_pair",
		Description: "An SSH key pair is a secure and convenient authentication method provided by Alibaba Cloud for instance logon. An SSH key pair consists of a public key and a private key. You can use SSH key pairs to log on to only Linux instances.",
		List: &plugin.ListConfig{
			Hydrate: listEcsKeypair,
			Tags:    map[string]string{"service": "ecs", "action": "DescribeKeyPairs"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getEcsKeypair,
			Tags:       map[string]string{"service": "ecs", "action": "DescribeKeyPairs"},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the key pair.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyPairName"),
			},
			{
				Name:        "key_pair_finger_print",
				Description: "The fingerprint of the key pair.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time when the key pair was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "resource_group_id",
				Description: "The ID of the resource group to which the key pair belongs.",
				Type:        proto.ColumnType_STRING,
			},

			{
				Name:        "tags_src",
				Description: "A list of tags attached with the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags.Tag").Transform(modifyGenericSourceTags),
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags.Tag").Transform(genericTagsToMap),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyPairName"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getEcsKeypairAka,
				Transform:   transform.FromValue(),
			},
			// Alibaba standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKeyPairRegion,
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

func listEcsKeypair(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := ECSService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ecs_keypair.listEcsKeypair", "connection_error", err)
		return nil, err
	}
	request := &ecs.DescribeKeyPairsRequest{
		PageSize:   tea.Int32(50),
		PageNumber: tea.Int32(1),
		RegionId:   tea.String(d.EqualsQualString(matrixKeyRegion)),
	}
	count := 0
	for {
		d.WaitForListRateLimit(ctx)
		response, err := client.DescribeKeyPairs(request)
		if err != nil {
			logQueryError(ctx, d, h, "alicloud_ecs_keypair.listEcsKeypair", err, "request", request)
			return nil, err
		}
		for _, keypair := range response.Body.KeyPairs.KeyPair {
			plugin.Logger(ctx).Warn("listEcsKeypair", "item", *keypair)
			d.StreamListItem(ctx, *keypair)
			// This will return zero if context has been cancelled (i.e due to manual cancellation) or
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
			count++
		}
		if count >= int(*response.Body.TotalCount) {
			break
		}
		request.PageNumber = tea.Int32((*response.Body.PageNumber) + int32(1))
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getEcsKeypair(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEcsSnapshot")

	// Create service connection
	client, err := ECSService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ecs_keypair.getEcsKeypair", "connection_error", err)
		return nil, err
	}

	var name string
	if h.Item != nil {
		keypair := h.Item.(ecs.DescribeKeyPairsResponseBodyKeyPairsKeyPair)
		name = *keypair.KeyPairName
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	request := &ecs.DescribeKeyPairsRequest{
		KeyPairName: &name,
	}

	response, err := client.DescribeKeyPairs(request)
	if serverErr, ok := err.(*tea.SDKError); ok {
		logQueryError(ctx, d, h, "alicloud_ecs_keypair.getEcsKeypair", serverErr, "request", request)
		return nil, serverErr
	}

	if len(response.Body.KeyPairs.KeyPair) > 0 {
		return *response.Body.KeyPairs.KeyPair[0], nil
	}

	return nil, nil
}

func getEcsKeypairAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEcsKeypairAka")
	data := h.Item.(ecs.DescribeKeyPairsResponseBodyKeyPairsKeyPair)
	region := d.EqualsQualString(matrixKeyRegion)

	// Get account details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*alicloudCommonColumnData)
	accountID := commonColumnData.AccountID

	akas := []string{"acs:ecs:" + region + ":" + accountID + ":keypair/" + tea.StringValue(data.KeyPairName)}

	return akas, nil
}

func getKeyPairRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	return region, nil
}
