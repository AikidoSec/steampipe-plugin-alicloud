package alicloud

import (
	"context"
	"fmt"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	"github.com/alibabacloud-go/tea/dara"

	actiontrail "github.com/alibabacloud-go/actiontrail-20200706/v3/client"
	alidns "github.com/alibabacloud-go/alidns-20150109/v5/client"
	cas "github.com/alibabacloud-go/cas-20200407/v4/client"
	cms "github.com/alibabacloud-go/cms-20190101/v10/client"
	cs "github.com/alibabacloud-go/cs-20151215/v7/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v7/client"
	ess "github.com/alibabacloud-go/ess-20220222/v2/client"
	ims "github.com/alibabacloud-go/ims-20190815/v4/client"
	kms "github.com/alibabacloud-go/kms-20160120/v3/client"
	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	rds "github.com/alibabacloud-go/rds-20140815/v16/client"
	sas "github.com/alibabacloud-go/sas-20181203/v8/client"
	slb "github.com/alibabacloud-go/slb-20140515/v4/client"
	sts "github.com/alibabacloud-go/sts-20150401/v2/client"
	vpc "github.com/alibabacloud-go/vpc-20160428/v7/client"

	credential "github.com/aliyun/credentials-go/credentials"
	credentialProviders "github.com/aliyun/credentials-go/credentials/providers"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	ossCred "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	sls "github.com/aliyun/aliyun-log-go-sdk"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Credential configuration
type CredentialConfig struct {
	Cred          credential.Credential
	DefaultRegion string
}

// newOpenAPIConfig creates an OpenAPI config for the given region using the credential
func newOpenAPIConfig(cred credential.Credential, region string) *openapi.Config {
	return &openapi.Config{
		Credential: cred,
		RegionId:   dara.String(region),
	}
}

// AliDNSService returns the service connection for Alicloud DNS service
func AliDNSService(ctx context.Context, d *plugin.QueryData) (*alidns.Client, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	if region == "" {
		return nil, fmt.Errorf("region must be passed AliDNSService")
	}

	serviceCacheKey := fmt.Sprintf("alidns-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*alidns.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := alidns.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// AutoscalingService returns the service connection for Alicloud Autoscaling service
func AutoscalingService(ctx context.Context, d *plugin.QueryData) (*ess.Client, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	if region == "" {
		return nil, fmt.Errorf("region must be passed AutoscalingService")
	}
	serviceCacheKey := fmt.Sprintf("ess-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*ess.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := ess.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// CasService returns the service connection for Alicloud SSL service
func CasService(ctx context.Context, d *plugin.QueryData, region string) (*cas.Client, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed CasService")
	}
	serviceCacheKey := fmt.Sprintf("cas-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cas.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := cas.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// CmsService returns the service connection for Alicloud CMS service
func CmsService(ctx context.Context, d *plugin.QueryData) (*cms.Client, error) {
	region := GetDefaultRegion(d.Connection)

	if region == "" {
		return nil, fmt.Errorf("region must be passed CmsService")
	}
	serviceCacheKey := fmt.Sprintf("cms-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cms.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := cms.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ECSService returns the service connection for Alicloud ECS service
func ECSService(ctx context.Context, d *plugin.QueryData) (*ecs.Client, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	if region == "" {
		return nil, fmt.Errorf("region must be passed ECSService")
	}
	serviceCacheKey := fmt.Sprintf("ecs-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*ecs.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := ecs.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ECSRegionService returns the service connection for Alicloud ECS Region service
func ECSRegionService(ctx context.Context, d *plugin.QueryData, region string) (*ecs.Client, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed ECSRegionService")
	}
	serviceCacheKey := fmt.Sprintf("ecsregion-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*ecs.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := ecs.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// KMSService returns the service connection for Alicloud KMS service
func KMSService(ctx context.Context, d *plugin.QueryData) (*kms.Client, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	if region == "" {
		return nil, fmt.Errorf("region must be passed KMSService")
	}
	serviceCacheKey := fmt.Sprintf("kms-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*kms.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := kms.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// RAMService returns the service connection for Alicloud RAM service
func RAMService(ctx context.Context, d *plugin.QueryData) (*ram.Client, error) {
	region := GetDefaultRegion(d.Connection)

	serviceCacheKey := fmt.Sprintf("ram-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*ram.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := ram.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

func IMSService(ctx context.Context, d *plugin.QueryData) (*ims.Client, error) {
	region := GetDefaultRegion(d.Connection)

	serviceCacheKey := fmt.Sprintf("ims-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*ims.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := ims.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// SLBService returns the service connection for Alicloud Server Load Balancer service
func SLBService(ctx context.Context, d *plugin.QueryData) (*slb.Client, error) {
	region := GetDefaultRegion(d.Connection)

	serviceCacheKey := fmt.Sprintf("slb-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*slb.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := slb.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// StsService returns the service connection for Alicloud STS service
func StsService(ctx context.Context, d *plugin.QueryData) (*sts.Client, error) {
	region := GetDefaultRegion(d.Connection)
	serviceCacheKey := fmt.Sprintf("sts-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*sts.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := sts.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// VpcService returns the service connection for Alicloud VPC service
func VpcService(ctx context.Context, d *plugin.QueryData) (*vpc.Client, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	// Fallback to the default region in connection config
	if region == "" {
		region = GetDefaultRegion(d.Connection)
	}

	if region == "" {
		return nil, fmt.Errorf("region could not be determined for VpcService")
	}

	serviceCacheKey := fmt.Sprintf("vpc-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*vpc.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := vpc.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// OssService returns the service connection for Alicloud OSS service
func OssService(ctx context.Context, d *plugin.QueryData, region string) (*oss.Client, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be provided to initialize the OSS service")
	}

	serviceCacheKey := fmt.Sprintf("oss-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*oss.Client), nil
	}

	endpoint := "oss-" + region + ".aliyuncs.com"

	ossCfg := oss.NewConfig()
	ossCfg.WithEndpoint(endpoint)
	ossCfg.WithRegion(region)
	ossCfg.WithProxyFromEnvironment(true)

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve cached credentials: %v", err)
	}

	cfg := credCfg.(*CredentialConfig)

	// Extract credentials from the v2 credential interface
	accessKeyID, err := cfg.Cred.GetAccessKeyId()
	if err != nil {
		return nil, fmt.Errorf("failed to get access key id: %v", err)
	}
	accessKeySecret, err := cfg.Cred.GetAccessKeySecret()
	if err != nil {
		return nil, fmt.Errorf("failed to get access key secret: %v", err)
	}
	securityToken, err := cfg.Cred.GetSecurityToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get security token: %v", err)
	}

	ossCfg.CredentialsProvider = ossCred.NewStaticCredentialsProvider(
		dara.StringValue(accessKeyID),
		dara.StringValue(accessKeySecret),
		dara.StringValue(securityToken),
	)

	svc := oss.NewClient(ossCfg)
	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ActionTrailService returns the service connection for Alicloud ActionTrail service
func ActionTrailService(ctx context.Context, d *plugin.QueryData) (*actiontrail.Client, error) {
	region := d.EqualsQualString(matrixKeyRegion)

	if region == "" {
		return nil, fmt.Errorf("region must be passed ActionTrailService")
	}
	serviceCacheKey := fmt.Sprintf("actiontrail-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*actiontrail.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := actiontrail.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ContainerService returns the service connection for Alicloud Container service
func ContainerService(ctx context.Context, d *plugin.QueryData) (*cs.Client, error) {
	region := GetDefaultRegion(d.Connection)

	if region == "" {
		return nil, fmt.Errorf("region must be passed ContainerService")
	}
	serviceCacheKey := fmt.Sprintf("cs-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cs.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := cs.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// SecurityCenterService returns the service connection for Alicloud Security Center service
func SecurityCenterService(ctx context.Context, d *plugin.QueryData, region string) (*sas.Client, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed SecurityCenterService")
	}

	serviceCacheKey := fmt.Sprintf("sas-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*sas.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := sas.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// RDSService returns the service connection for Alicloud RDS service
func RDSService(ctx context.Context, d *plugin.QueryData, region string) (*rds.Client, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed RDSService")
	}
	serviceCacheKey := fmt.Sprintf("rds-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*rds.Client), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	cfg := credCfg.(*CredentialConfig)

	svc, err := rds.NewClient(newOpenAPIConfig(cfg.Cred, region))
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// SLSService returns the client interface for Alicloud Log Service (SLS)
func SLSService(ctx context.Context, d *plugin.QueryData, region string) (sls.ClientInterface, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be provided to initialize the SLS service")
	}

	serviceCacheKey := fmt.Sprintf("sls-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(sls.ClientInterface), nil
	}

	credCfg, err := getCredentialSessionCached(ctx, d, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve cached credentials: %v", err)
	}
	cfg := credCfg.(*CredentialConfig)

	// Extract credentials from the v2 credential interface
	accessKeyId, err := cfg.Cred.GetAccessKeyId()
	if err != nil {
		return nil, fmt.Errorf("failed to get access key id: %v", err)
	}
	accessKeySecret, err := cfg.Cred.GetAccessKeySecret()
	if err != nil {
		return nil, fmt.Errorf("failed to get access key secret: %v", err)
	}
	securityToken, err := cfg.Cred.GetSecurityToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get security token: %v", err)
	}

	staticProvider := sls.NewStaticCredentialsProvider(
		dara.StringValue(accessKeyId),
		dara.StringValue(accessKeySecret),
		dara.StringValue(securityToken),
	)
	endpoint := region + ".log.aliyuncs.com"
	client := sls.CreateNormalInterfaceV2(endpoint, staticProvider)

	d.ConnectionManager.Cache.Set(serviceCacheKey, client)
	return client, nil
}

// GetDefaultRegion returns the default region used
func GetDefaultRegion(connection *plugin.Connection) string {
	alicloudConfig := GetConfig(connection)

	var regions []string
	var region string

	if alicloudConfig.Regions != nil {
		regions = alicloudConfig.Regions
	}

	if len(regions) > 0 {
		region = regions[0]
		if len(getInvalidRegions([]string{region})) > 0 {
			panic("\n\nConnection config have invalid region: " + region + ". Edit your connection configuration file and then restart Steampipe")
		}
		return region
	}

	if region == "" {
		region = os.Getenv("ALIBABACLOUD_REGION_ID")
		if region == "" {
			region = os.Getenv("ALICLOUD_REGION_ID")
			if region == "" {
				region = os.Getenv("ALICLOUD_REGION")
			}
		}
	}

	if region == "" {
		region = "cn-hangzhou"
	}

	return region
}

// https://github.com/aliyun/aliyun-cli/blob/master/README.md#supported-environment-variables
func getEnvForProfile(_ context.Context, d *plugin.QueryData) (profile string) {
	alicloudConfig := GetConfig(d.Connection)
	if alicloudConfig.Profile != nil {
		profile = *alicloudConfig.Profile
	} else {
		var ok bool
		if profile, ok = os.LookupEnv("ALIBABACLOUD_PROFILE"); !ok {
			if profile, ok = os.LookupEnv("ALIBABA_CLOUD_PROFILE"); !ok {
				if profile, ok = os.LookupEnv("ALICLOUD_PROFILE"); !ok {
					return ""
				}
			}
		}
	}
	return profile
}

func getEnv(_ context.Context, d *plugin.QueryData) (secretKey string, accessKey string, sessionToken string, err error) {
	alicloudConfig := GetConfig(d.Connection)

	if alicloudConfig.AccessKey != nil {
		accessKey = *alicloudConfig.AccessKey
	} else {
		var ok bool
		if accessKey, ok = os.LookupEnv("ALIBABACLOUD_ACCESS_KEY_ID"); !ok {
			if accessKey, ok = os.LookupEnv("ALICLOUD_ACCESS_KEY_ID"); !ok {
				if accessKey, ok = os.LookupEnv("ALICLOUD_ACCESS_KEY"); !ok {
					panic("\n'access_key' or 'profile' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe.")
				}
			}
		}
	}

	if alicloudConfig.SecretKey != nil {
		secretKey = *alicloudConfig.SecretKey
	} else {
		var ok bool
		if secretKey, ok = os.LookupEnv("ALIBABACLOUD_ACCESS_KEY_SECRET"); !ok {
			if secretKey, ok = os.LookupEnv("ALICLOUD_ACCESS_KEY_SECRET"); !ok {
				if secretKey, ok = os.LookupEnv("ALICLOUD_SECRET_KEY"); !ok {
					panic("\n'secret_key' or 'profile' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe.")
				}
			}
		}
	}

	if alicloudConfig.SessionToken != nil {
		sessionToken = *alicloudConfig.SessionToken
	} else {
		var ok bool
		if sessionToken, ok = os.LookupEnv("ALIBABACLOUD_ACCESS_KEY_SECRET"); !ok {
			if sessionToken, ok = os.LookupEnv("ALICLOUD_ACCESS_KEY_SECRET"); !ok {
				sessionToken, _ = os.LookupEnv("ALICLOUD_SECRET_KEY")
			}
		}
	}

	return accessKey, secretKey, sessionToken, nil
}

// Get credential from the profile configuration for Alicloud CLI
func getProfileConfigurations(_ context.Context, d *plugin.QueryData) (*CredentialConfig, error) {
	alicloudConfig := GetConfig(d.Connection)
	profile := alicloudConfig.Profile

	cfg, err := getCredentialConfigByProfile(*profile, d)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func getCredentialConfigByProfile(profile string, d *plugin.QueryData) (*CredentialConfig, error) {
	defaultRegion := GetDefaultRegion(d.Connection)

	provider, err := credentialProviders.NewCLIProfileCredentialsProviderBuilder().
		WithProfileName(profile).
		Build()
	if err != nil {
		return nil, err
	}

	cred := credential.FromCredentialsProvider("cli_profile", provider)

	return &CredentialConfig{cred, defaultRegion}, nil
}

var getCredentialSessionCached = plugin.HydrateFunc(getCredentialSessionUncached).Memoize()

func getCredentialSessionUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	config := GetConfig(d.Connection)
	defaultRegion := GetDefaultRegion(d.Connection)

	// Profile based client
	if config.Profile != nil {
		return getProfileConfigurations(ctx, d)
	}

	profileEnv := getEnvForProfile(ctx, d)
	if profileEnv != "" {
		return getCredentialConfigByProfile(profileEnv, d)
	}

	// Access key and Secret Key from environment variable
	accessKey, secretKey, sessionToken, err := getEnv(ctx, d)
	if err != nil {
		return nil, err
	}
	if sessionToken != "" && accessKey != "" && secretKey != "" {
		credConfig := &credential.Config{
			Type:            dara.String("sts"),
			AccessKeyId:     dara.String(accessKey),
			AccessKeySecret: dara.String(secretKey),
			SecurityToken:   dara.String(sessionToken),
		}
		cred, err := credential.NewCredential(credConfig)
		if err != nil {
			return nil, err
		}
		return &CredentialConfig{cred, defaultRegion}, nil
	}
	if accessKey != "" && secretKey != "" {
		credConfig := &credential.Config{
			Type:            dara.String("access_key"),
			AccessKeyId:     dara.String(accessKey),
			AccessKeySecret: dara.String(secretKey),
		}
		cred, err := credential.NewCredential(credConfig)
		if err != nil {
			return nil, err
		}
		return &CredentialConfig{cred, defaultRegion}, nil
	}

	return nil, nil
}
