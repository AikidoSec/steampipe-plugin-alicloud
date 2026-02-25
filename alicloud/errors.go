package alicloud

import (
	"context"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// isNotFoundError:: function which returns an ErrorPredicateWithContext for Alicloud API calls
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if err == nil {
			return false
		}

		alicloudConfig := GetConfig(d.Connection)

		// If the get or list hydrate functions have an overriding IgnoreConfig
		// defined using the isNotFoundError function, then it should
		// also check for errors in the "ignore_error_codes" config argument
		allErrors := append(notFoundErrors, alicloudConfig.IgnoreErrorCodes...)

		// Ignore dynamic endpoint resolution failures ("no such host") for services not available in all regions
		if strings.Contains(err.Error(), "no such host") {
			plugin.Logger(ctx).Debug("ignoring no such host error for unreachable region", "error", err)
			return true
		}

		// V2 SDK error handling using SDK error codes
		if sdkErr, ok := err.(*tea.SDKError); ok {
			errCode := tea.StringValue(sdkErr.Code)
			for _, pattern := range allErrors {
				if strings.Contains(errCode, pattern) {
					return true
				}
			}
		}

		for _, pattern := range allErrors {
			if strings.Contains(err.Error(), pattern) {
				return true
			}
		}
		return false
	}
}

// shouldIgnoreErrorPluginDefault:: Plugin level default function to ignore a set errors for hydrate functions based on "ignore_error_codes" config argument
func shouldIgnoreErrorPluginDefault() plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if err == nil {
			return false
		}

		alicloudConfig := GetConfig(d.Connection)

		// Ignore dynamic endpoint resolution failures ("no such host") for services not available in all regions
		if strings.Contains(err.Error(), "no such host") {
			plugin.Logger(ctx).Debug("ignoring no such host error for unreachable region", "error", err)
			return true
		}

		// V2 SDK error handling using SDK error codes
		if sdkErr, ok := err.(*tea.SDKError); ok {
			errCode := tea.StringValue(sdkErr.Code)
			for _, pattern := range alicloudConfig.IgnoreErrorCodes {
				if strings.Contains(errCode, pattern) {
					return true
				}
			}
		}

		for _, pattern := range alicloudConfig.IgnoreErrorCodes {
			if strings.Contains(err.Error(), pattern) {
				return true
			}
		}
		return false
	}
}
