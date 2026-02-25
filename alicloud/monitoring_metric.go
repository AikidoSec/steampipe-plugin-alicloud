package alicloud

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	cms "github.com/alibabacloud-go/cms-20190101/v10/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/sethvargo/go-retry"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// append the common cloud monitoring metric columns onto the column list
func cmMetricColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonCMMetricColumns()...)
}

func commonCMMetricColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "metric_name",
			Description: "The name of the metric.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "namespace",
			Description: "The metric namespace.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "average",
			Description: "The average of the metric values that correspond to the data point.",
			Type:        proto.ColumnType_DOUBLE,
			Default:     0,
		},
		{
			Name:        "maximum",
			Description: "The maximum metric value for the data point.",
			Type:        proto.ColumnType_DOUBLE,
			Default:     0,
		},
		{
			Name:        "minimum",
			Description: "The minimum metric value for the data point.",
			Type:        proto.ColumnType_DOUBLE,
			Default:     0,
		},
		{
			Name:        "timestamp",
			Description: "The timestamp used for the data point.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "account_id",
			Description: ColumnDescriptionAccount,
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCommonColumns,
			Transform:   transform.FromField("AccountID"),
		},
	}
}

type CMMetricRow struct {
	// The (single) metric Dimension name
	DimensionName string

	// The value for the (single) metric Dimension
	DimensionValue string

	// The namespace of the metric
	Namespace string

	// The name of the metric
	MetricName string

	// The average of the metric values that correspond to the data point.
	Average float64

	// The percentile statistic for the data point.
	// ExtendedStatistics map[string]*float64 `type:"map"`

	// The maximum metric value for the data point.
	Maximum float64

	// The minimum metric value for the data point.
	Minimum float64

	// The timestamp used for the data point.
	Timestamp string
}

func getCMStartDateForGranularity(granularity string) string {
	str := "2006-01-02T15:04:05Z"
	switch strings.ToUpper(granularity) {
	case "DAILY":
		// 30 days
		return time.Now().AddDate(0, 0, -30).Format(str)
	case "HOURLY":
		// 30 days
		return time.Now().AddDate(0, 0, -30).Format(str)
	}
	// else 5 days
	return time.Now().AddDate(0, 0, -5).Format(str)
}

func getCMPeriodForGranularity(granularity string) string {
	switch strings.ToUpper(granularity) {
	case "DAILY":
		// 24 hours
		return "86400"
	case "HOURLY":
		// 1 hour
		return "3600"
	}
	// else 5 minutes
	return "300"
}

func getCustomError(errorMessage string) error {
	return tea.NewSDKError(map[string]interface{}{
		"message":    errorMessage,
		"statusCode": 500,
	})
}

func listCMMetricStatistics(ctx context.Context, d *plugin.QueryData, granularity string, namespace string, metricName string, dimensionName string, dimensionValue string) (*cms.DescribeMetricListResponse, error) {
	// Create service connection
	client, err := CmsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCMMetricStatistics", "connection_error", err)
		return nil, err
	}

	request := &cms.DescribeMetricListRequest{
		Namespace:  &namespace,
		MetricName: &metricName,
		Dimensions: tea.String("[{\"" + dimensionName + "\": \"" + dimensionValue + "\"}]"),
		StartTime:  tea.String(getCMStartDateForGranularity(granularity)),
		EndTime:    tea.String(time.Now().Format("2006-01-02T15:04:05Z")),
		Period:     tea.String(getCMPeriodForGranularity(granularity)),
	}

	var stats *cms.DescribeMetricListResponse

	b := retry.NewFibonacci(100 * time.Millisecond)

	err = retry.Do(ctx, retry.WithMaxRetries(5, b), func(ctx context.Context) error {
		var err error
		stats, err = client.DescribeMetricList(request)
		if err != nil || stats.Body.Datapoints == nil || *stats.Body.Datapoints == "" {
			// Common server error retry
			if serverErr, ok := err.(*tea.SDKError); ok {
				if serverErr.Code != nil && *serverErr.Code == "Throttling" {
					return retry.RetryableError(err)
				}
				return err
			}
			/**
			* At some point of the time we are getting the error as success response(%!v(PANIC=String method: runtime error: invalid memory address or nil pointer dereference)") which is not expected.
			* If we will retry the api call then we will able to get the data.
			**/
			if *stats.Body.Datapoints == "" && !*stats.Body.Success {
				err = getCustomError(fmt.Sprint(stats))
				return retry.RetryableError(err)
			}

		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	// As some point of the time we are getting the error in response not in the error part.
	// Response in stats variable: "%!v(PANIC=String method: runtime error: invalid memory address or nil pointer dereference)"
	if stats.Body.Datapoints == nil || *stats.Body.Datapoints == "" {
		return nil, nil
	}

	err = json.Unmarshal([]byte(*stats.Body.Datapoints), &results)
	if err != nil {
		return nil, err
	}
	for _, pointValue := range results {
		d.StreamListItem(ctx, &CMMetricRow{
			DimensionName:  dimensionName,
			DimensionValue: pointValue[dimensionName].(string),
			Namespace:      namespace,
			MetricName:     metricName,
			Average:        pointValue["Average"].(float64),
			Maximum:        pointValue["Maximum"].(float64),
			Minimum:        pointValue["Minimum"].(float64),
			Timestamp:      formatTime(pointValue["timestamp"].(float64)),
		})
	}

	if stats.Body.NextToken != nil && *stats.Body.NextToken != "" {
		request.NextToken = stats.Body.NextToken
	}

	return nil, nil
}

func formatTime(timestamp float64) string {
	timeInSec := math.Floor(timestamp / 1000)
	unixTimestamp := time.Unix(int64(timeInSec), 0)
	timestampRFC3339Format := unixTimestamp.Format(time.RFC3339)
	return timestampRFC3339Format
}
