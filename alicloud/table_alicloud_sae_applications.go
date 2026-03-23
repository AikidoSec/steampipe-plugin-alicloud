package alicloud

import (
	"context"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	sae "github.com/alibabacloud-go/sae-20190506/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAlicloudSaeApplication(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_sae_application",
		Description: "Alicloud Serverless App Engine Application",
		List: &plugin.ListConfig{
			Hydrate: listApplications,
			Tags:    map[string]string{"service": "sae", "action": "ListApplications"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getApplication,
			Tags:       map[string]string{"service": "sae", "action": "GetApplication"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getApplication,
				Tags: map[string]string{"service": "sae", "action": "GetApplication"},
			},
			{
				Func: describeApplicationConfig,
				Tags: map[string]string{"service": "sae", "action": "DescribeApplicationConfig"},
			},
			{
				Func:    getSaeAppArn,
				Tags:    map[string]string{"service": "sae", "action": "getSaeAppArn"},
				Depends: []plugin.HydrateFunc{describeApplicationConfig},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the application.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppName"),
			},
			{
				Name:        "namespace_id",
				Description: "The namespace ID.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "arn",
				Description: "The Alibaba Cloud Resource Name (ARN) of the SAE Application.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSaeAppArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "id",
				Description: "The application ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppId"),
			},
			{
				Name:        "readiness",
				Description: "The readiness probe of the application. A container that fails the health check multiple times is shut down and restarted.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "liveness",
				Description: "The liveness probe of the container. A container that fails the health check is shut down and restored.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "config_map_mount_desc",
				Description: "The information about the ConfigMap.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "security_group_id",
				Description: "The ID of the security group.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "batch_wait_time",
				Description: "The interval between batches in a phased release. Unit: seconds.",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "sls_configs",
				Description: "The configuration of collecting logs to SLS.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "jdk",
				Description: "The version of the JDK on which the deployment package depends.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "image_url",
				Description: "The address of the image. This parameter is required if PackageType is set to Image.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "package_url",
				Description: "The URL of the deployment package. If you upload the deployment package using SAE.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "package_type",
				Description: "The type of the application package.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "post_start",
				Description: "The script that is run after the container is started.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "pre_stop",
				Description: "The script that is run before the container is stopped",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "package_version",
				Description: "The version of the deployment package.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "jar_start_args",
				Description: "The startup parameters of the JAR package.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "jar_start_options",
				Description: "The startup options of the JAR package.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "replicas",
				Description: "The number of application instances.",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "update_strategy",
				Description: "The release policy.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "min_ready_instances",
				Description: "The minimum number of ready instances",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "min_ready_instance_ratio",
				Description: "The percentage of the minimum number of ready instances",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "memory",
				Description: "The amount of memory that is required by each instance. Unit: MB.",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "php",
				Description: "The version of PHP that is used for the deployment package.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "php_config",
				Description: "The content of the PHP configuration file.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "php_config_location",
				Description: "The path on which the startup configuration file of the PHP application is mounted.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "termination_grace_period_seconds",
				Description: "The timeout period for a graceful shutdown. Unit: seconds. ",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "command_args",
				Description: "The parameters of the image startup command.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "mount_host",
				Description: "The mount target of the NAS.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "tomcat_config",
				Description: "The Tomcat configuration.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "v_switch_id",
				Description: "The ID of the vSwitch.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "cpu",
				Description: "The CPU quota of the application instance. Unit: millicores.",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "envs",
				Description: "The environment variables of the container.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "mount_desc",
				Description: "The description of the NAS mounting configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "enable_ahas",
				Description: "Specifies whether to access AHAS.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "custom_host_alias",
				Description: "The custom domain mappings that are configured to map the domain names to the IP addresses in the pod.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "web_container",
				Description: "The Tomcat version that is used for the deployment package.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "command",
				Description: "The image startup command.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "war_start_options",
				Description: "The startup options of the WAR package.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "php_arms_config_location",
				Description: "The path on which the PHP application monitors the configuration file.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "nas_id",
				Description: "The ID of the NAS file system.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "oss_ak_id",
				Description: "The AccessKey ID of the OSS bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "oss_ak_secret",
				Description: "The AccessKey secret of the OSS bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "oss_mount_descs",
				Description: "The description of the OSS mounting configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "edas_container_version",
				Description: "The version of the container on which the deployment package depends.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "timezone",
				Description: "The timezone of the application.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "app_description",
				Description: "The description of the application.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "enable_grey_tag_route",
				Description: "Specifies whether to enable the canary release rules.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "mse_application_id",
				Description: "The ID of the application in MSE.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "acr_instance_id",
				Description: "The ID of the Container Registry instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "acr_assume_role_arn",
				Description: "The ARN of the RAM role that is used to pull images across accounts.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "image_pull_secrets",
				Description: "The name of the Secret that is used to pull the image.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "associate_eip",
				Description: "Specifies whether to associate an EIP with the application.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "kafka_configs",
				Description: "The configuration of collecting logs to Kafka.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "programming_language",
				Description: "The programming language used by the application. Valid values: java, php, other.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "pvtz_discovery",
				Description: "The configuration of service registration and discovery.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "swimlane_pvtz_discovery",
				Description: "The configuration for service registration and discovery based on Kubernetes Service.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "micro_registration",
				Description: "The type of the registry. 0: SAE built-in registry. 1: custom registry. 2: MSE registry.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "nas_configs",
				Description: "The NAS mounting configuration.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "python",
				Description: "The version of the Python environment.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "python_modules",
				Description: "The Python modules that need to be installed.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "mse_application_name",
				Description: "The name of the application in the MSE registry.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "app_source",
				Description: "The type of the application. Valid values: micro_service, web, job.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "service_tags",
				Description: "The canary release tags of the application.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "dotnet",
				Description: "The .NET version.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "micro_registration_config",
				Description: "The configuration of the registry. This parameter is applicable when you set MicroRegistration to 2 (MSE Nacos registry).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "enable_idle",
				Description: "Specifies whether to enable the idle mode.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "enable_new_arms",
				Description: "Specifies whether to enable the new ARMS feature.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "enable_cpu_burst",
				Description: "Specifies whether to enable the CPU burst feature.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "oidc_role_name",
				Description: "The name of the RAM role for identity authentication.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "secret_mount_desc",
				Description: "The description of the Secret mounting configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "sidecar_containers_config",
				Description: "The configuration of the sidecar container.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "init_containers_config",
				Description: "The configuration of the init container.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "custom_image_network_type",
				Description: "The network type of the custom image.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "resource_type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "base_app_id",
				Description: "The ID of the baseline application.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "microservice_engine_config",
				Description: "The microservice governance configurations.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "new_sae_version",
				Description: "The version of SAE. Valid values: lite, std, pro.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "cluster_id",
				Description: "The ID of the cluster.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "disk_size",
				Description: "The disk size of the application instance. Unit: GiB.",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "startup_probe",
				Description: "The startup probe.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "gpu_count",
				Description: "The number of GPUs.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "gpu_type",
				Description: "The GPU type.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "enable_prometheus",
				Description: "Specifies whether to enable Prometheus monitoring.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "is_stateful",
				Description: "Specifies whether the application is a stateful application.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "cms_service_id",
				Description: "The ID of the CMS service.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "headless_pvtz_discovery",
				Description: "The configuration of headless service registration and discovery.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "html",
				Description: "The HTML content.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "empty_dir_desc",
				Description: "The shared temporary storage configurations for the main container and sidecar containers.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "deployment_name",
				Description: "The name of the deployment.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "alb_ingress_readiness_gate",
				Description: "The ALB Ingress readiness configuration.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "max_surge_instances",
				Description: "The maximum number of new instances that can be created in a batch.",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "max_surge_instance_ratio",
				Description: "The maximum percentage of new instances that can be created in a batch.",
				Type:        proto.ColumnType_INT,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "agent_version",
				Description: "The version of the agent.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "enable_namespace_agent_version",
				Description: "Specifies whether to enable the namespace agent version.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "sls_log_env_tags",
				Description: "The environment tags of the SLS logs.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "loki_configs",
				Description: "The Loki log configuration.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "idle_hour",
				Description: "The configuration of the idle hour.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     describeApplicationConfig,
			},
			{
				Name:        "labels",
				Description: "The labels of the application.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     describeApplicationConfig,
			},

			{
				Name:        "tags_src",
				Description: "A list of tags attached with the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags").Transform(modifyGenericSourceTags),
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags").Transform(genericTagsToMap),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSaeAppArn,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AppName"),
			},

			// Alicloud standard columns
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

func listApplications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listApplications")
	client, err := SAEService(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_sae_application.listApplications", "connection_error", err)
		return nil, err
	}

	request := &sae.ListApplicationsRequest{
		CurrentPage: tea.Int32(1),
		PageSize:    tea.Int32(100),
	}

	count := 0
	for {
		d.WaitForListRateLimit(ctx)
		response, err := client.ListApplications(request)
		if err != nil {
			logQueryError(ctx, d, h, "alicloud_sae_application.listApplications", err, "request", request)
			return nil, err
		}

		for _, app := range response.Body.Data.Applications {
			plugin.Logger(ctx).Warn("alicloud_sae_application.listApplications", "item", app)
			d.StreamListItem(ctx, app)
			// This will return zero if context has been cancelled (i.e due to manual cancellation) or
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
			count++
		}
		if count >= int(*response.Body.TotalSize) {
			break
		}
		request.SetCurrentPage(tea.Int32Value(response.Body.CurrentPage) + 1)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getApplication(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getApplication")

	client, err := SAEService(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_sae_application.getApplication", "connection_error", err)
		return nil, err
	}

	var id *string
	if h.Item != nil {
		app := h.Item.(*sae.ListApplicationsResponseBodyDataApplications)
		id = app.AppId
	} else {
		id = tea.String(d.EqualsQuals["id"].GetStringValue())
	}

	params := &openapi.Params{ // API Name
		Action:      tea.String("GetApplication"),
		Version:     tea.String("2019-05-06"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		Pathname:    tea.String("/pop/v1/sam/app/getApplication"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}

	queries := map[string]interface{}{}
	queries["AppId"] = id
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	response, err := client.CallApi(params, request, &util.RuntimeOptions{})
	if serverErr, ok := err.(*tea.SDKError); ok {
		logQueryError(ctx, d, h, "alicloud_sae_application.getApplication", serverErr, "app", *id)
		return nil, serverErr
	}
	return response["body"], nil
}

func describeApplicationConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("describeApplicationConfig")

	client, err := SAEService(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("alicloud_sae_application.describeApplicationConfig", "connection_error", err)
		return nil, err
	}

	var id *string
	if h.Item != nil {
		app := h.Item.(*sae.ListApplicationsResponseBodyDataApplications)
		id = app.AppId
	} else {
		id = tea.String(d.EqualsQuals["id"].GetStringValue())
	}

	request := &sae.DescribeApplicationConfigRequest{
		AppId: id,
	}
	response, err := client.DescribeApplicationConfig(request)
	if serverErr, ok := err.(*tea.SDKError); ok {
		logQueryError(ctx, d, h, "alicloud_sae_application.describeApplicationConfig", serverErr, "app", *id)
		return nil, serverErr
	}
	plugin.Logger(ctx).Warn("alicloud_sae_application.describeApplicationConfig", "item", response.Body.Data)

	return response.Body.Data, nil
}

func getSaeAppArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSaeAppArn")
	res, ok := h.HydrateResults["describeApplicationConfig"]
	if !ok || res == nil {
		return nil, fmt.Errorf("could not get required info for ARN construction")
	}
	app := res.(*sae.DescribeApplicationConfigResponseBodyData)

	// Get project details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*alicloudCommonColumnData)
	accountID := commonColumnData.AccountID
	arn := "acs:sae:" + tea.StringValue(app.RegionId) + ":" + accountID + ":application/" + tea.StringValue(app.NamespaceId) + "/" + tea.StringValue(app.AppId)

	return arn, nil
}

func getRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	if region != "" {
		return region, nil
	}

	return GetDefaultRegion(d.Connection), nil
}
