package alicloud

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Constants for Standard Column Descriptions
const (
	ColumnDescriptionAkas    = "Array of globally unique identifier strings (also known as) for the resource."
	ColumnDescriptionTags    = "A map of tags for the resource."
	ColumnDescriptionTitle   = "Title of the resource."
	ColumnDescriptionAccount = "The Alicloud Account ID in which the resource is located."
	ColumnDescriptionRegion  = "The Alicloud region in which the resource is located."
)

type resourceTags = struct {
	TagKey   string
	TagValue string
}

func ensureStringArray(_ context.Context, d *transform.TransformData) (interface{}, error) {
	switch v := d.Value.(type) {
	case []string:
		return v, nil
	case string:
		return []string{v}, nil
	default:
		str := fmt.Sprintf("%v", d.Value)
		return []string{string(str)}, nil
	}
}

func csvToStringArray(_ context.Context, d *transform.TransformData) (interface{}, error) {
	s := tea.StringValue(d.Value.(*string))
	if s == "" {
		// Empty string should always be an empty array
		return []string{}, nil
	}
	sep := ","
	if d.Param != nil {
		sep = d.Param.(string)
	}
	return strings.Split(s, sep), nil
}

func getGenericTags(d *transform.TransformData) ([]map[string]interface{}, error) {
	// Strict typing would give a lot of boilerplate, so let's cheat
	b, err := json.Marshal(d.Value)
	if err != nil {
		return nil, err
	}

	var rawTags []map[string]interface{}
	if err := json.Unmarshal(b, &rawTags); err != nil {
		return nil, err
	}
	return rawTags, nil
}

func modifyGenericSourceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	// Strict typing would give a lot of boilerplate, so let's cheat
	tags, err := getGenericTags(d)
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return nil, nil
	}

	// We now have a generic interface, but different naming is used accross alicloud
	// So we'll check for TagKey/Key and TagValue/Value

	var sourceTags []resourceTags

	for _, tag := range tags {
		var tagKey, tagValue string

		if val, ok := tag["Key"].(string); ok {
			tagKey = val
		} else if val, ok := tag["TagKey"].(string); ok {
			tagKey = val
		}

		if val, ok := tag["Value"].(string); ok {
			tagValue = val
		} else if val, ok := tag["TagValue"].(string); ok {
			tagValue = val
		}
		sourceTags = append(sourceTags, resourceTags{tagKey, tagValue})
	}

	return sourceTags, nil
}

func genericTagsToMap(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	genericTags, err := modifyGenericSourceTags(ctx, d)
	if genericTags == nil || err != nil {
		return nil, err
	}

	tags := genericTags.([]resourceTags)

	turbotTagsMap := map[string]string{}
	for _, i := range tags {
		turbotTagsMap[i.TagKey] = i.TagValue
	}

	return turbotTagsMap, nil
}

func zoneToRegion(_ context.Context, d *transform.TransformData) (interface{}, error) {
	region := tea.StringValue(d.Value.(*string))
	return region[:len(region)-1], nil
}

func GetBoolQualValue(quals plugin.KeyColumnQualMap, columnName string) (value *bool, exists bool) {
	exists = false
	if quals[columnName] == nil {
		return nil, exists
	}

	if quals[columnName].Quals == nil {
		return nil, exists
	}

	for _, qual := range quals[columnName].Quals {
		if qual.Value != nil {
			value := qual.Value
			boolValue := value.GetBoolValue()
			switch qual.Operator {
			case "<>":
				return types.Bool(!boolValue), true
			case "=":
				return types.Bool(boolValue), true
			}
			break
		}
	}
	return nil, exists
}

// GetStringQualValue :: Can be used to get equal value
func GetStringQualValue(quals plugin.KeyColumnQualMap, columnName string) (value *string, exists bool) {
	exists = false
	if quals[columnName] == nil {
		return nil, exists
	}

	if quals[columnName].Quals == nil {
		return nil, exists
	}

	for _, qual := range quals[columnName].Quals {
		if qual.Operator != "=" {
			return nil, exists
		}
		if qual.Value != nil {
			value := qual.Value
			// In case of IN caluse the qual value is usally of format vpcid = '[id1 id2]'
			// which can lead to generation of wrong filter
			if value.GetListValue() != nil {
				// Cannot assign array values
				return nil, exists
			} else {
				return types.String(value.GetStringValue()), true
			}
		}
	}
	return nil, exists
}

// GetStringQualValueList :: Can be used to get equal value as a list of strings
// supports only equal operator
func GetStringQualValueList(quals plugin.KeyColumnQualMap, columnName string) (values []string, exists bool) {
	exists = false
	if quals[columnName] == nil {
		return nil, exists
	}

	if quals[columnName].Quals == nil {
		return nil, exists
	}

	for _, qual := range quals[columnName].Quals {
		if qual.Operator != "=" {
			return nil, exists
		}
		if qual.Value != nil {
			value := qual.Value
			if value.GetListValue() != nil {
				for _, q := range value.GetListValue().Values {
					values = append(values, q.GetStringValue())
				}
				return values, true
			} else {
				values = append(values, value.GetStringValue())
				return values, true
			}
		}
	}
	return values, exists
}

type QueryFilterItem struct {
	Key    string
	Values []string
}

// QueryFilters is an array of filters items
type QueryFilters []QueryFilterItem

// To get the stringified value of QueryFilters
func (filters *QueryFilters) String() (string, error) {
	data, err := json.Marshal(filters)
	if err != nil {
		return "", fmt.Errorf("error marshalling filters: %v", err)
	}

	return string(data), nil
}

func logQueryError(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, key string, err error, extra ...any) {
	// Do not pollute the logs with error messages for missing/disable services
	if !shouldIgnoreErrorPluginDefault()(ctx, d, h, err) {
		info := []any{"connection_error", err}
		info = append(info, extra...)
		plugin.Logger(ctx).Error(key, info...)
	}
}
