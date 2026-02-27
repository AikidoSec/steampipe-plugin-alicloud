package alicloud

import (
	"context"
	"fmt"
	"time"

	ims "github.com/alibabacloud-go/ims-20190815/v4/client"
	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sethvargo/go-retry"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type userInfo = struct {
	UserName      string
	UserId        string
	DisplayName   string
	Email         string
	MobilePhone   string
	Comments      string
	CreateDate    string
	UpdateDate    string
	LastLoginDate string
}

//// TABLE DEFINITION

func tableAlicloudRAMUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_ram_user",
		Description: "Resource Access Management users who can login via the console or access keys.",
		List: &plugin.ListConfig{
			Hydrate: listRAMUser,
			Tags:    map[string]string{"service": "ram", "action": "ListUsers"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"EntityNotExist.User", "MissingParameter"}),
			},
			Hydrate: getRAMUser,
			Tags:    map[string]string{"service": "ram", "action": "GetUser"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getRAMUserMfaDevices,
				Tags: map[string]string{"service": "ram", "action": "ListVirtualMFADevices"},
			},
			{
				Func: getRAMUserPolicies,
				Tags: map[string]string{"service": "ram", "action": "ListPoliciesForUser"},
			},
			{
				Func: getRAMUserGroups,
				Tags: map[string]string{"service": "ram", "action": "ListGroupsForUser"},
			},
			{
				Func: getRAMUserPasskeys,
				Tags: map[string]string{"service": "ims", "action": "getRAMUserPasskeys"},
				Depends: []plugin.HydrateFunc{getRAMUserMfaDevices},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{
				Name:        "name",
				Description: "The username of the RAM user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserName"),
			},
			{
				Name:        "arn",
				Description: "The Alibaba Cloud Resource Name (ARN) of the RAM user.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getUserArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "user_id",
				Description: "The unique ID of the RAM user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "The display name of the RAM user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "email",
				Description: "The email address of the RAM user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_login_date",
				Description: "The time when the RAM user last logged on to the console by using the password.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getRAMUser,
			},
			{
				Name:        "mobile_phone",
				Description: "The mobile phone number of the RAM user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "comments",
				Description: "The description of the RAM user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "create_date",
				Description: "The time when the RAM user was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "update_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time when the RAM user was modified.",
			},
			{
				Name:        "mfa_enabled",
				Description: "The MFA status of the user",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getRAMUserPasskeys,
				Transform:   transform.From(userMfaStatus),
			},
			{
				Name:        "mfa_device_serial_number",
				Description: "The serial number of the MFA device.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getRAMUserMfaDevices,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "attached_policy",
				Description: "A list of policies attached to a RAM user.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRAMUserPolicies,
				Transform:   transform.FromField("Policies.Policy"),
			},
			{
				Name:        "cs_user_permissions",
				Description: "User permissions for Container Service Kubernetes clusters.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCsUserPermissions,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "groups",
				Description: "A list of groups attached to the user.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRAMUserGroups,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "virtual_mfa_devices",
				Description: "The list of MFA devices.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRAMUserMfaDevices,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "passkeys",
				Description: "The list of passkeys for the user.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRAMUserPasskeys,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getUserArn,
				Transform:   transform.FromValue().Transform(ensureStringArray),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserName"),
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

func listRAMUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	client, err := RAMService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ram_user.listRAMUser", "connection_error", err)
		return nil, err
	}
	request := &ram.ListUsersRequest{}
	for {
		response, err := client.ListUsers(request)
		if err != nil {
			logQueryError(ctx, d, h, "alicloud_ram_user.listRAMUser", err, "request", request)
			return nil, err
		}
		for _, i := range response.Body.Users.User {
			plugin.Logger(ctx).Warn("listRAMUser", "item", i)
			d.StreamListItem(ctx, userInfo{
				tea.StringValue(i.UserName),
				tea.StringValue(i.UserId),
				tea.StringValue(i.DisplayName),
				tea.StringValue(i.Email),
				tea.StringValue(i.MobilePhone),
				tea.StringValue(i.Comments),
				tea.StringValue(i.CreateDate),
				tea.StringValue(i.UpdateDate),
				"",
			})
			// This will return zero if context has been cancelled (i.e due to manual cancellation) or
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if response.Body.IsTruncated == nil || !*response.Body.IsTruncated {
			break
		}
		request.Marker = response.Body.Marker
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRAMUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRAMUser")

	// Create service connection
	client, err := RAMService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ram_user.getRAMUser", "connection_error", err)
		return nil, err
	}

	var name string
	if h.Item != nil {
		i := h.Item.(userInfo)
		name = i.UserName
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	request := &ram.GetUserRequest{
		UserName: &name,
	}
	var response *ram.GetUserResponse

	b := retry.NewFibonacci(100 * time.Millisecond)

	err = retry.Do(ctx, retry.WithMaxRetries(5, b), func(ctx context.Context) error {
		var err error
		response, err = client.GetUser(request)
		if err != nil {
			if serverErr, ok := err.(*tea.SDKError); ok {
				if serverErr.Code != nil && *serverErr.Code == "Throttling" {
					return retry.RetryableError(err)
				}
				logQueryError(ctx, d, h, "alicloud_ram_user.getRAMUser", err, "request", request)
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	data := response.Body.User
	return userInfo{
		tea.StringValue(data.UserName),
		tea.StringValue(data.UserId),
		tea.StringValue(data.DisplayName),
		tea.StringValue(data.Email),
		tea.StringValue(data.MobilePhone),
		tea.StringValue(data.Comments),
		tea.StringValue(data.CreateDate),
		tea.StringValue(data.UpdateDate),
		tea.StringValue(data.LastLoginDate),
	}, nil
}

func getRAMUserGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRAMUserGroups")
	data := h.Item.(userInfo)

	// Create service connection
	client, err := RAMService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ram_group.getRAMUserGroups", "connection_error", err)
		return nil, err
	}

	request := &ram.ListGroupsForUserRequest{
		UserName: &data.UserName,
	}

	response, err := client.ListGroupsForUser(request)
	if err != nil {
		if serverErr, ok := err.(*tea.SDKError); ok {
			logQueryError(ctx, d, h, "alicloud_ram_group.getRAMUserGroups", serverErr, "request", request)
			return nil, serverErr
		}
		return nil, err
	}

	return response.Body.Groups.Group, nil
}

func getRAMUserPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRAMUserPolicies")
	data := h.Item.(userInfo)

	// Create service connection
	client, err := RAMService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ram_group.getRAMUserPolicies", "connection_error", err)
		return nil, err
	}

	request := &ram.ListPoliciesForUserRequest{
		UserName: &data.UserName,
	}
	var response *ram.ListPoliciesForUserResponse

	b := retry.NewFibonacci(100 * time.Millisecond)

	err = retry.Do(ctx, retry.WithMaxRetries(5, b), func(ctx context.Context) error {
		var err error
		response, err = client.ListPoliciesForUser(request)
		if err != nil {
			if serverErr, ok := err.(*tea.SDKError); ok {
				if serverErr.Code != nil && *serverErr.Code == "Throttling" {
					return retry.RetryableError(err)
				}
				logQueryError(ctx, d, h, "alicloud_ram_group.getRAMUserPolicies", serverErr, "request", request)
				return nil
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return *response.Body, nil
}

func getRAMUserMfaDevices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRAMUserMfaDevices")
	data := h.Item.(userInfo)

	// Create service connection
	client, err := RAMService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ram_group.getRAMUserMfaDevices", "connection_error", err)
		return nil, err
	}

	response, err := client.ListVirtualMFADevices()
	if err != nil {
		if serverErr, ok := err.(*tea.SDKError); ok {
			logQueryError(ctx, d, h, "alicloud_ram_group.getRAMUserMfaDevices", serverErr)
			return nil, serverErr
		}
		return nil, err
	}

	var items []*ram.ListVirtualMFADevicesResponseBodyVirtualMFADevicesVirtualMFADevice
	for _, i := range response.Body.VirtualMFADevices.VirtualMFADevice {
		if i.User != nil && tea.StringValue(i.User.UserName) == data.UserName {
			items = append(items, i)
		}
	}

	return items, nil
}

func getRAMUserPasskeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRAMUserPasskeys")
	data := h.Item.(userInfo)

	// Get project details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*alicloudCommonColumnData)
	accountID := commonColumnData.AccountID

	// Create service connection
	client, err := IMSService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_ram_user.getRAMUserPasskeys", "connection_error", err)
		return nil, err
	}

	response, err := client.ListPasskeys(&ims.ListPasskeysRequest{
		UserPrincipalName: tea.String(fmt.Sprintf("%s@%s.onaliyun.com", data.UserName, accountID)),
	})
	if err != nil {
		if serverErr, ok := err.(*tea.SDKError); ok {
			logQueryError(ctx, d, h, "alicloud_ram_user.getRAMUserPasskeys", serverErr)
			return nil, serverErr
		}
		return nil, err
	}

	return response.Body.Passkeys, nil
}

func getUserArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getUserAkas")
	data := h.Item.(userInfo)

	// Get project details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*alicloudCommonColumnData)
	accountID := commonColumnData.AccountID

	return "acs:ram::" + accountID + ":user/" + data.UserName, nil
}

func getCsUserPermissions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCsUserPermissions")

	// Create service connection
	client, err := ContainerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getCsUserPermissions", "connection_error", err)
		return nil, err
	}

	data := h.Item.(userInfo)

	response, err := client.DescribeUserPermission(&data.UserId)
	if err != nil {
		logQueryError(ctx, d, h, "getCsUserPermissions", err)
		return nil, err
	}

	return response.Body, nil
}

//// TRANSFORM FUNCTION

func userMfaStatus(_ context.Context, d *transform.TransformData) (interface{}, error) {
	passkeys := d.HydrateResults["getRAMUserPasskeys"].([]*ims.ListPasskeysResponseBodyPasskeys)
	mfaDevices := d.HydrateResults["getRAMUserMfaDevices"].([]*ram.ListVirtualMFADevicesResponseBodyVirtualMFADevicesVirtualMFADevice)

	if len(passkeys) > 0 || len(mfaDevices) > 0 {
		return true, nil
	}

	return false, nil
}
