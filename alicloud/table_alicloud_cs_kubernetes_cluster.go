package alicloud

import (
	"context"
	"encoding/json"

	cs "github.com/alibabacloud-go/cs-20151215/v7/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAlicloudCsKubernetesCluster(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "alicloud_cs_kubernetes_cluster",
		Description: "Alicloud Container Service Kubernetes Cluster",
		List: &plugin.ListConfig{
			Hydrate: listCsKubernetesClusters,
			Tags:    map[string]string{"service": "cs", "action": "DescribeClustersV1"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("cluster_id"),
			Hydrate:    getCsKubernetesCluster,
			Tags:       map[string]string{"service": "cs", "action": "DescribeClusterDetail"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCsKubernetesClusterLog,
				Tags: map[string]string{"service": "cs", "action": "DescribeClusterLogs"},
			},
			{
				Func: getCsKubernetesClusterNamespace,
				Tags: map[string]string{"service": "cs", "action": "DescribeClusterNamespaces"},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("name"),
			},
			{
				Name:        "cluster_id",
				Description: "The ID of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("cluster_id"),
			},
			{
				Name:        "arn",
				Description: "The Alibaba Cloud Resource Name (ARN) of the cluster.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCsKubernetesClusterARN,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "state",
				Description: "The status of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("state"),
			},
			{
				Name:        "size",
				Description: "The number of nodes in the cluster.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("size"),
			},
			{
				Name:        "created_at",
				Description: "The time when the cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("created"),
			},
			{
				Name:      "capabilities",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("capabilities"),
			},
			{
				Name:        "cluster_healthy",
				Description: "The health status of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("cluster_healthy"),
			},
			{
				Name:      "cluster_spec",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("cluster_spec"),
			},
			{
				Name:        "cluster_type",
				Description: "The type of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("cluster_type"),
			},
			{
				Name:        "current_version",
				Description: "The version of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("current_version"),
			},
			{
				Name:        "data_disk_category",
				Description: "The type of data disks.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("data_disk_category"),
			},
			{
				Name:        "data_disk_size",
				Description: "The size of a data disk.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("data_disk_size"),
			},
			{
				Name:        "deletion_protection",
				Description: "Indicates whether deletion protection is enabled for the cluster.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("deletion_protection"),
			},
			{
				Name:        "docker_version",
				Description: "The version of Docker.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("docker_version"),
			},
			{
				Name:      "enabled_migration",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("enabled_migration"),
			},
			{
				Name:        "external_loadbalancer_id",
				Description: "The ID of the Server Load Balancer (SLB) instance deployed in the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("external_loadbalancer_id"),
			},
			{
				Name:      "gw_bridge",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("gw_bridge"),
			},
			{
				Name:        "init_version",
				Description: "The initial version of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("init_version"),
			},
			{
				Name:        "instance_type",
				Description: "The Elastic Compute Service (ECS) instance type of cluster nodes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("instance_type"),
			},
			{
				Name:      "maintenance_info",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("maintenance_info"),
			},
			{
				Name:      "need_update_agent",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("need_update_agent"),
			},
			{
				Name:        "network_mode",
				Description: "The network type of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("network_mode"),
			},
			{
				Name:      "next_version",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("next_version"),
			},
			{
				Name:        "node_status",
				Description: "The status of cluster nodes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("node_status"),
			},
			{
				Name:      "outputs",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("outputs"),
			},
			{
				Name:      "parameters",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("parameters"),
			},
			{
				Name:        "port",
				Description: "Container port in Kubernetes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("port"),
			},
			{
				Name:        "private_zone",
				Description: "Indicates whether PrivateZone is enabled for the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("private_zone"),
			},
			{
				Name:        "profile",
				Description: "The identifier of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("profile"),
			},
			{
				Name:        "resource_group_id",
				Description: "The ID of the resource group to which the cluster belongs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("resource_group_id"),
			},
			{
				Name:      "service_discovery_types",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("service_discovery_types"),
			},
			{
				Name:        "subnet_cidr",
				Description: "The CIDR block of pods in the cluster.",
				Type:        proto.ColumnType_CIDR,
				Transform:   transform.FromField("subnet_cidr"),
			},
			{
				Name:      "swarm_mode",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("swarm_mode"),
			},
			{
				Name:        "updated",
				Description: "The time when the cluster was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("updated"),
			},
			{
				Name:      "upgrade_components",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("upgrade_components"),
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC used by the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("vpc_id"),
			},
			{
				Name:        "vswitch_id",
				Description: "The IDs of VSwitches.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("vswitch_id"),
			},
			{
				Name:        "vswitch_cidr",
				Description: "The CIDR block of VSwitches.",
				Type:        proto.ColumnType_CIDR,
				Transform:   transform.FromField("vswitch_cidr"),
			},
			{
				Name:        "worker_ram_role_name",
				Description: "The name of the RAM role for worker nodes in the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("worker_ram_role_name"),
			},
			{
				Name:        "zone_id",
				Description: "The ID of the zone where the cluster is deployed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("zone_id"),
			},
			{
				Name:        "cluster_log",
				Description: "The logs of a cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCsKubernetesClusterLog,
				Transform:   transform.FromValue(),
			},
			{
				Name:      "maintenance_window",
				Type:      proto.ColumnType_JSON,
				Transform: transform.FromField("maintenance_window"),
			},
			{
				Name:        "master_url",
				Description: "The endpoints that are open for connections to the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("master_url"),
			},
			{
				Name:        "meta_data",
				Description: "The metadata of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("meta_data"),
			},
			{
				Name:      "cluster_namespace",
				Type:      proto.ColumnType_JSON,
				Hydrate:   getCsKubernetesClusterNamespace,
				Transform: transform.FromValue(),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags attached with the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("tags"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("name"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("tags").Transform(csKubernetesClusterAkaTagsToMap),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCsKubernetesClusterARN,
				Transform:   transform.FromValue().Transform(transform.EnsureStringArray),
			},

			// Alicloud standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("region_id"),
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

func listCsKubernetesClusters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := GetDefaultRegion(d.Connection)

	// Create service connection
	client, err := ContainerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCsKubernetesClusters", "connection_error", err)
		return nil, err
	}
	request := &cs.DescribeClustersV1Request{
		PageSize:   tea.Int64(50),
		PageNumber: tea.Int64(1),
		RegionId:   tea.String(region),
	}

	count := 0
	for {
		d.WaitForListRateLimit(ctx)
		response, err := client.DescribeClustersV1(request)
		if err != nil {
			logQueryError(ctx, d, h, "listCsKubernetesClusters", err, "request", request)
			return nil, err
		}

		for _, cluster := range response.Body.Clusters {
			// Convert v2 struct (with lowercase JSON tags) to map[string]interface{}
			clusterBytes, err := json.Marshal(cluster)
			if err != nil {
				plugin.Logger(ctx).Error("listCsKubernetesClusters", "json.marshal", err)
				return nil, err
			}
			var clusterMap map[string]interface{}
			if err := json.Unmarshal(clusterBytes, &clusterMap); err != nil {
				plugin.Logger(ctx).Error("listCsKubernetesClusters", "json.unmarshal", err)
				return nil, err
			}
			d.StreamListItem(ctx, clusterMap)
			// This will return zero if context has been cancelled (i.e due to manual cancellation) or
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
			count++
		}

		if response.Body.PageInfo == nil || response.Body.PageInfo.TotalCount == nil {
			break
		}
		if count >= int(*response.Body.PageInfo.TotalCount) {
			break
		}
		request.SetPageNumber(*request.PageNumber + 1)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCsKubernetesCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCsKubernetesCluster")

	// Create service connection
	client, err := ContainerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getCsKubernetesCluster", "connection_error", err)
		return nil, err
	}

	var id string
	if h.Item != nil {
		clusterData := h.Item.(map[string]interface{})
		id = clusterData["cluster_id"].(string)
	} else {
		id = d.EqualsQuals["cluster_id"].GetStringValue()
	}

	response, err := client.DescribeClusterDetail(&id)
	if err != nil {
		if serverErr, ok := err.(*tea.SDKError); ok {
			logQueryError(ctx, d, h, "getCsKubernetesCluster", serverErr)
			return nil, serverErr
		}
		return nil, err
	}

	// Convert response body (with lowercase JSON tags) to map[string]interface{}
	clusterBytes, err := json.Marshal(response.Body)
	if err != nil {
		plugin.Logger(ctx).Error("getCsKubernetesCluster", "json_marshal", err)
		return nil, err
	}
	var cluster map[string]interface{}
	if err := json.Unmarshal(clusterBytes, &cluster); err != nil {
		plugin.Logger(ctx).Error("getCsKubernetesCluster", "json_unmarshal", err)
		return nil, err
	}

	return cluster, nil
}

func getCsKubernetesClusterLog(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCsKubernetesClusterLog")

	// Create service connection
	client, err := ContainerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getCsKubernetesClusterLog", "connection_error", err)
		return nil, err
	}

	id := h.Item.(map[string]interface{})["cluster_id"].(string)

	response, err := client.DescribeClusterLogs(&id)
	if err != nil {
		if serverErr, ok := err.(*tea.SDKError); ok {
			logQueryError(ctx, d, h, "getCsKubernetesClusterLog", serverErr)
			return nil, serverErr
		}
		return nil, err
	}

	if len(response.Body) > 0 {
		return response.Body, nil
	}

	return nil, nil
}

func getCsKubernetesClusterNamespace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCsKubernetesClusterNamespace")

	id := h.Item.(map[string]interface{})["cluster_id"].(string)

	client, err := ContainerService(ctx, d)
	if err != nil {
		return nil, nil
	}

	response, err := client.DescribeUserClusterNamespaces(&id)
	if err != nil {
		return nil, nil
	}

	return response.Body, nil
}

func getCsKubernetesClusterARN(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCsKubernetesClusterARN")

	data := h.Item.(map[string]interface{})

	// Get project details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*alicloudCommonColumnData)
	accountID := commonColumnData.AccountID

	arn := "arn:acs:cs:" + data["region_id"].(string) + ":" + accountID + ":cluster/" + data["cluster_id"].(string)

	return arn, nil
}

//// TRANSFORM FUNCTIONS

func csKubernetesClusterAkaTagsToMap(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	tags := d.Value.([]interface{})

	turbotTagsMap := map[string]string{}
	for _, i := range tags {
		tagDetails := i.(map[string]interface{})
		turbotTagsMap[tagDetails["key"].(string)] = tagDetails["value"].(string)
	}

	return turbotTagsMap, nil
}
