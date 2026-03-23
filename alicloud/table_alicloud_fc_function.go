package alicloud

import (
	"context"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	fc "github.com/alibabacloud-go/fc-20230330/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAlicloudFcFunction(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_fc_function",
		Description: "Alicloud FC Function",
		List: &plugin.ListConfig{
			Hydrate: listFunctions,
			Tags:    map[string]string{"service": "fc", "action": "ListFunctions"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getFunction,
			Tags:       map[string]string{"service": "fc", "action": "GetFunction"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getFunction,
				Tags: map[string]string{"service": "fc", "action": "GetFunction"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FunctionName"),
			},
			{
				Name:        "function_id",
				Description: "The ID of the function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Alibaba Cloud Resource Name (ARN) of the FC Function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FunctionArn"),
			},
			{
				Name:        "code_checksum",
				Description: "The CRC-64 value of the function code package.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "code_size",
				Description: "The size of the function code package, in bytes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "cpu",
				Description: "The CPU specification of the function, in vCPUs. The value must be a multiple of 0.05 vCPU. The minimum value is 0.05 and the maximum value is 16. The ratio of CPU to memory size (in GB) must be between 1:1 and 1:4.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "created_time",
				Description: "The time when the function was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "custom_container_config",
				Description: "The configuration for the custom container runtime. After you configure this parameter, the function can use a custom container image to run. You must specify either `code` or `customContainerConfig`.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "custom_runtime_config",
				Description: "The custom runtime configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "custom_dns",
				Description: "The custom DNS configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CustomDNS"),
			},
			{
				Name:        "description",
				Description: "The description of the function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "disk_size",
				Description: "The disk size of the function, in MB. Valid values are 512 MB and 10240 MB.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "environment_variables",
				Description: "The environment variables of the function. The configured environment variables can be accessed in the runtime environment.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "gpu_config",
				Description: "The GPU configuration of the function.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "handler",
				Description: "The entry point for the function execution. The format varies based on the runtime.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_concurrency",
				Description: "The maximum concurrency of the instance.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "instance_lifecycle_config",
				Description: "The configuration of the instance lifecycle hook.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "internet_access",
				Description: "Specifies whether the function can access the Internet. Default value: true.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "last_modified_time",
				Description: "The time when the function was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_update_status",
				Description: "The status of the last update operation on the function. When a function is created, the value is `Successful`. Valid values: `Successful`, `Failed`, and `InProgress`.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getFunction,
			},
			{
				Name:        "last_update_status_reason",
				Description: "The reason for the status of the last update operation on the function",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getFunction,
			},
			{
				Name:        "last_update_status_reason_code",
				Description: "The status code for the reason of the last update operation on the function.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getFunction,
			},
			{
				Name:        "layers",
				Description: "The list of layers.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "log_config",
				Description: "The log configuration. The logs that are generated by the function are written to the configured Logstore.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "memory_size",
				Description: "The memory size of the function, in MB. The value must be a multiple of 64 MB. The minimum value is 128 MB and the maximum value is 32 GB. The ratio of CPU to memory size (in GB) must be between 1:1 and 1:4.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "nas_config",
				Description: "The NAS configuration. After you configure this parameter, the function can access the specified NAS resources.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "oss_mount_config",
				Description: "The OSS mount configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "role",
				Description: "The RAM role that you grant to Function Compute. After you set this parameter, Function Compute assumes the role to generate temporary access credentials. The function can use the temporary access credentials of the role to access specified Alibaba Cloud services, such as OSS and Tablestore.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "runtime",
				Description: "The runtime environment of the function. The following runtimes are supported: nodejs12, nodejs14, nodejs16, nodejs18, nodejs20, go1, python3, python3.9, python3.10, python3.12, java8, java11, php7.2, dotnetcore3.1, custom, custom.debian10, custom.debian11, custom.debian12, and custom-container.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The current state of the function.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getFunction,
			},
			{
				Name:        "state_reason",
				Description: "The reason why the function is in the current state.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getFunction,
			},
			{
				Name:        "state_reason_code",
				Description: "The status code for the reason why the function is in the current state.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getFunction,
			},
			{
				Name:        "timeout",
				Description: "The timeout period for the function to run, in seconds. The minimum value is 1 second and the maximum value is 86,400 seconds. The default value is 3 seconds. If the function runs longer than this period, the execution is stopped.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "tracing_config",
				Description: "The Tracing Analysis configuration. Integrating Function Compute with Tracing Analysis lets you record the time that requests consume in Function Compute, view the cold start time of functions, and record the time consumed within functions.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vpc_config",
				Description: "The VPC configuration. After you configure this parameter, the function can access the specified VPC resources.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "disable_ondemand",
				Description: "Specifies whether to disable the creation of on-demand instances. If this feature is enabled, on-demand instances are not created. Only provisioned instances can be used.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "invocation_restriction",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "session_affinity",
				Description: "The affinity policy for Function Compute invocation requests. To implement request affinity for the MCP SSE protocol, set this parameter to `MCP_SSE`. To use cookie-based affinity, set this parameter to `GENERATED_COOKIE`. To use header-based affinity, set this parameter to `HEADER_FIELD`. If you do not set this parameter or set it to `NONE`, no affinity is used, and requests are routed based on the default scheduling policy of Function Compute.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enable_long_living",
				Description: "When you set the `sessionAffinity` type, you must configure the related affinity settings. For `MCP_SSE` affinity, specify the `MCPSSESessionAffinityConfig` settings. For cookie-based affinity, specify the `CookieSessionAffinityConfig` settings. For header field-based affinity, specify the `HeaderFieldSessionAffinityConfig` settings.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "resource_group_id",
				Description: "The ID of the resource group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_isolation_mode",
				Description: "The isolation mode for the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "session_affinity_config",
				Description: "When you set the `sessionAffinity` type, you must configure the related affinity settings. For `MCP_SSE` affinity, specify the `MCPSSESessionAffinityConfig` settings. For cookie-based affinity, specify the `CookieSessionAffinityConfig` settings. For header field-based affinity, specify the `HeaderFieldSessionAffinityConfig` settings.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "idle_timeout",
				Description: "The amount of time that an instance can remain idle before it is released.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "disable_inject_credentials",
				Description: "Specifies whether to prevent the injection of the Security Token Service (STS) token. Valid values are `None`, `Env`, `Request`, and `All`. `None` means the token is injected. `Env` means the token is not injected into environment variables. `Request` means the token is not injected into the request, including the context and header. `All` means the token is not injected.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "polar_fs_config",
				Description: "The PolarFS configuration. After you configure this parameter, the function can access the specified PolarFS resources.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached with the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFunctionTags,
				Transform:   transform.FromField("TagResources").Transform(modifyGenericSourceTags),
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFunctionTags,
				Transform:   transform.FromField("TagResources").Transform(genericTagsToMap),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FunctionArn").Transform(transform.EnsureStringArray),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FunctionName"),
			},

			// Alicloud standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getFunctionRegion,
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

func listFunctions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := FCService(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_fc_function.listFunctions", "connection_error", err)
		return nil, err
	}

	request := &fc.ListFunctionsRequest{
		Limit: tea.Int32(100),
	}

	// If the request no of items is less than the paging max limit
	// update limit to the requested no of results.
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		maxResults := int64(tea.Int32Value(request.Limit))
		if *limit < maxResults {
			request.Limit = tea.Int32(int32(*limit))
		}
	}

	pageLeft := true
	for pageLeft {
		d.WaitForListRateLimit(ctx)
		response, err := client.ListFunctions(request)
		if err != nil {
			logQueryError(ctx, d, h, "alicloud_fc_function.listFunctions", err, "request", request)
			return nil, err
		}
		for _, fn := range response.Body.Functions {
			plugin.Logger(ctx).Warn("alicloud_fc_function.listFunctions", "item", fn)
			d.StreamListItem(ctx, fn)
			// This will return zero if context has been cancelled (i.e due to manual cancellation) or
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if tea.StringValue(response.Body.NextToken) != "" {
			request.NextToken = response.Body.NextToken
		} else {
			pageLeft = false
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getFunction(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFunction")

	client, err := FCService(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_gc_function.getFunction", "connection_error", err)
		return nil, err
	}

	var name *string
	if h.Item != nil {
		fn := h.Item.(*fc.Function)
		name = fn.FunctionName
	} else {
		name = tea.String(d.EqualsQuals["function_name"].GetStringValue())
	}

	response, err := client.GetFunction(name, &fc.GetFunctionRequest{})
	if serverErr, ok := err.(*tea.SDKError); ok {
		logQueryError(ctx, d, h, "alicloud_gc_function.getFunction", serverErr, "function", *name)
		return nil, serverErr
	}

	return response.Body, nil
}

func getFunctionTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFunctionTags")
	data := h.Item.(*fc.Function)

	// Create service connection
	client, err := FCService(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_gc_function.getFunctionTags", "connection_error", err)
		return nil, err
	}

	params := &openapi.Params{ // API Name
		Action:      tea.String("ListTagResources"),
		Version:     tea.String("2023-03-30"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("FC"),
		Pathname:    tea.String("/2023-03-30/tags-v2"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}

	queries := map[string]interface{}{}
	queries["ResourceType"] = tea.String("function")
	queries["ResourceId"] = tea.String("[\"" + *data.FunctionArn + "\"]")
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	resp, err := client.CallApi(params, request, &util.RuntimeOptions{})
	if err != nil {
		if serverErr, ok := err.(*tea.SDKError); ok {
			logQueryError(ctx, d, h, "alicloud_gc_function.getFunctionTags", serverErr, "request", request)
			return nil, serverErr
		}
		return nil, err
	}

	return resp["body"], nil
}

func getFunctionRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	if region != "" {
		return region, nil
	}

	return GetDefaultRegion(d.Connection), nil
}
