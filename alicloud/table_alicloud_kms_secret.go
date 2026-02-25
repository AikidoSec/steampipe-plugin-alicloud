package alicloud

import (
	"context"
	"strings"
	"time"

	"github.com/sethvargo/go-retry"

	kms "github.com/alibabacloud-go/kms-20160120/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAlicloudKmsSecret(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_kms_secret",
		Description: "Alicloud KMS Secret",
		List: &plugin.ListConfig{
			Hydrate: listKmsSecret,
			Tags:    map[string]string{"service": "kms", "action": "ListSecrets"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "region"}),
			Hydrate:    getKmsSecret,
			Tags:       map[string]string{"service": "kms", "action": "DescribeSecret"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: listKmsSecretVersionIds,
				Tags: map[string]string{"service": "kms", "action": "ListSecretVersionIds"},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the secret.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SecretName"),
			},
			{
				Name:        "arn",
				Description: "The Alibaba Cloud Resource Name (ARN).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKmsSecret,
			},
			{
				Name:        "secret_type",
				Description: "The type of the secret.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "automatic_rotation",
				Description: "Specifies whether automatic key rotation is enabled.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKmsSecret,
			},
			{
				Name:        "create_time",
				Description: "The time when the KMS Secret was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description of the secret.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKmsSecret,
			},
			{
				Name:        "encryption_key_id",
				Description: "The ID of the KMS customer master key (CMK) that is used to encrypt the secret value.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKmsSecret,
			},
			{
				Name:        "last_rotation_date",
				Description: "Date of last rotation of Secret.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getKmsSecret,
			},
			{
				Name:        "next_rotation_date",
				Description: "The date of next rotation of Secret.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getKmsSecret,
			},
			{
				Name:        "planned_delete_time",
				Description: "The time when the KMS Secret is planned to delete.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "rotation_interval",
				Description: "The rotation perion of Secret.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKmsSecret,
			},
			{
				Name:        "update_time",
				Description: "The time when the KMS Secret was modifies.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "extended_config",
				Description: "The extended configuration of Secret.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKmsSecret,
			},
			{
				Name:        "version_ids",
				Description: "The list of secret versions.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listKmsSecretVersionIds,
				Transform:   transform.FromField("VersionId"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKmsSecret,
				Transform:   transform.FromField("Tags.Tag").Transform(modifyGenericSourceTags),
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKmsSecret,
				Transform:   transform.FromField("Tags.Tag").Transform(genericTagsToMap),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getKmsSecret,
				Transform:   transform.FromField("Arn").Transform(ensureStringArray),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SecretName"),
			},

			// Alicloud standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKmsSecret,
				Transform:   transform.From(fetchRegionFromArn),
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

func listKmsSecret(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := KMSService(ctx, d)
	if err != nil {
		logQueryError(ctx, d, h, "alicloud_kms_secret.listKmsSecret", err)
		return nil, err
	}

	request := &kms.ListSecretsRequest{
		PageSize:   tea.Int32(50),
		PageNumber: tea.Int32(1),
	}

	count := 0
	for {
		d.WaitForListRateLimit(ctx)
		response, err := client.ListSecrets(request)
		if err != nil {
			logQueryError(ctx, d, h, "alicloud_kms_secret.listKmsSecret", err, "request", request)
			return nil, err
		}
		for _, i := range response.Body.SecretList.Secret {
			d.StreamListItem(ctx, &kms.ListSecretsResponseBodySecretListSecret{
				CreateTime:        i.CreateTime,
				PlannedDeleteTime: i.PlannedDeleteTime,
				SecretName:        i.SecretName,
				UpdateTime:        i.UpdateTime,
				SecretType:        i.SecretType,
				Tags: &kms.ListSecretsResponseBodySecretListSecretTags{
					Tag: i.Tags.Tag,
				},
			})
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
		request.PageNumber = tea.Int32(tea.Int32Value(response.Body.PageNumber) + 1)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getKmsSecret(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKmsSecret")

	// Create service connection
	client, err := KMSService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_kms_secret.getKmsSecret", "connection_error", err)
		return nil, err
	}

	var name *string
	var response *kms.DescribeSecretResponse
	if h.Item != nil {
		data := h.Item.(*kms.DescribeSecretResponse)
		name = data.Body.SecretName
	} else {
		name = tea.String(d.EqualsQuals["name"].GetStringValue())
	}

	request := &kms.DescribeSecretRequest{
		SecretName: name,
		FetchTags:  tea.String("true"),
	}

	b := retry.NewFibonacci(100 * time.Millisecond)

	err = retry.Do(ctx, retry.WithMaxRetries(5, b), func(ctx context.Context) error {
		var err error
		response, err = client.DescribeSecret(request)
		if err != nil {
			if serverErr, ok := err.(*tea.SDKError); ok {
				if *serverErr.Code == "Throttling" {
					return retry.RetryableError(err)
				}
				logQueryError(ctx, d, h, "alicloud_kms_key.getKmsSecret", err, "request", request)
				return err
			}
		}
		return nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_kms_secret.getKmsSecret", "query_retry_error", err, "request", request)
		return nil, err
	}

	return response, nil
}

func listKmsSecretVersionIds(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listKmsSecretVersionIds")

	// Create service connection
	client, err := KMSService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_kms_secret.getKmsSecret", "connection_error", err)
		return nil, err
	}
	secretData := h.Item.(*kms.DescribeSecretResponse)
	var response *kms.ListSecretVersionIdsResponse

	request := &kms.ListSecretVersionIdsRequest{
		SecretName:        secretData.Body.SecretName,
		IncludeDeprecated: tea.String("true"),
	}

	b := retry.NewFibonacci(100 * time.Millisecond)

	err = retry.Do(ctx, retry.WithMaxRetries(5, b), func(ctx context.Context) error {
		var err error
		response, err = client.ListSecretVersionIds(request)
		if err != nil {
			if serverErr, ok := err.(*tea.SDKError); ok {
				if *serverErr.Code == "Throttling" {
					return retry.RetryableError(err)
				}
				logQueryError(ctx, d, h, "alicloud_kms_key.listKmsSecretVersionIds", err, "request", request)
				return err
			}
		}
		return nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_kms_key.listKmsSecretVersionIds", "retry_query_error", err, "request", request)
		return nil, err
	}

	if len(response.Body.VersionIds.VersionId) > 0 {
		return response.Body.VersionIds, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func fetchRegionFromArn(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(*kms.DescribeSecretResponseBody)

	resourceArn := data.Arn
	region := strings.Split(tea.StringValue(resourceArn), ":")[2]
	return region, nil
}
