package alicloud

import (
	"context"

	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type roleInfo = struct {
	RoleId                   string
	RoleName                 string
	Arn                      string
	Description              string
	AssumeRolePolicyDocument string
	CreateDate               string
	UpdateDate               string
	MaxSessionDuration       int64
}

//// TABLE DEFINITION

func tableAlicloudRAMRole(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_ram_role",
		Description: "Resource Access Management roles who can login via the console or access keys.",
		List: &plugin.ListConfig{
			Hydrate: listRAMRoles,
			Tags:    map[string]string{"service": "ram", "action": "ListRoles"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getRAMRole,
			Tags:       map[string]string{"service": "ram", "action": "GetRole"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getRAMRolePolicies,
				Tags: map[string]string{"service": "ram", "action": "ListPoliciesForRole"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the RAM role.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleName"),
			},
			{
				Name:        "arn",
				Description: "The Alibaba Cloud Resource Name (ARN) of the RAM role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role_id",
				Description: "The ID of the RAM role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description of the RAM role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_session_duration",
				Description: "The maximum session duration of the RAM role.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "create_date",
				Description: "The time when the RAM role was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_date",
				Description: "The time when the RAM role was modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "assume_role_policy_document",
				Description: "The content of the policy that specifies one or more entities entrusted to assume the RAM role.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRAMRole,
				Transform:   transform.FromField("AssumeRolePolicyDocument").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "assume_role_policy_document_std",
				Description: "The standard content of the policy that specifies one or more entities entrusted to assume the RAM role.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRAMRole,
				Transform:   transform.FromField("AssumeRolePolicyDocument").Transform(policyToCanonical),
			},
			{
				Name:        "attached_policy",
				Description: "A list of policies attached to a RAM role.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRAMRolePolicies,
				Transform:   transform.FromField("Policies.Policy"),
			},

			// steampipe standard columns
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(ensureStringArray),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleName"),
			},

			// alicloud standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("global"),
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

func listRAMRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := RAMService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ram_role.listRAMRoles", "connection_error", err)
		return nil, err
	}

	request := &ram.ListRolesRequest{}

	for {
		response, err := client.ListRoles(request)
		if err != nil {
			logQueryError(ctx, d, h, "alicloud_ram_role.listRAMRoles", err, "request", request)
			return nil, err
		}
		for _, i := range response.Body.Roles.Role {
			plugin.Logger(ctx).Warn("listRAMRoles", "item", *i)
			d.StreamListItem(ctx, roleInfo{
				tea.StringValue(i.RoleId),
				tea.StringValue(i.RoleName),
				tea.StringValue(i.Arn),
				tea.StringValue(i.Description),
				"",
				tea.StringValue(i.CreateDate),
				tea.StringValue(i.UpdateDate),
				tea.Int64Value(i.MaxSessionDuration),
			})
			// This will return zero if context has been cancelled (i.e due to manual cancellation) or
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if !tea.BoolValue(response.Body.IsTruncated) {
			break
		}
		request.Marker = response.Body.Marker
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRAMRole(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRAMRole")

	// Create service connection
	client, err := RAMService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ram_role.getRAMRole", "connection_error", err)
		return nil, err
	}

	var name string
	if h.Item != nil {
		i := h.Item.(roleInfo)
		name = i.RoleName
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	request := &ram.GetRoleRequest{
		RoleName: &name,
	}

	response, err := client.GetRole(request)
	if err != nil {
		logQueryError(ctx, d, h, "alicloud_ram_role.getRAMRole", err, "request", request)
		return nil, err
	}

	data := response.Body.Role
	return roleInfo{
		tea.StringValue(data.RoleId),
		tea.StringValue(data.RoleName),
		tea.StringValue(data.Arn),
		tea.StringValue(data.Description),
		tea.StringValue(data.AssumeRolePolicyDocument),
		tea.StringValue(data.CreateDate),
		tea.StringValue(data.UpdateDate),
		tea.Int64Value(data.MaxSessionDuration),
	}, nil
}

func getRAMRolePolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRAMRolePolicies")
	data := h.Item.(roleInfo)

	// Create service connection
	client, err := RAMService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ram_group.getRAMRolePolicies", "connection_error", err)
		return nil, err
	}

	request := &ram.ListPoliciesForRoleRequest{
		RoleName: &data.RoleName,
	}

	response, err := client.ListPoliciesForRole(request)
	if serverErr, ok := err.(*tea.SDKError); ok {
		logQueryError(ctx, d, h, "alicloud_ram_group.getRAMRolePolicies", serverErr, "request", request)
		return nil, serverErr
	}

	return response, nil
}
