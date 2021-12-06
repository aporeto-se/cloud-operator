package types

const (
	// PrismaPrependEnv is appended to the env var of each Prisma env var
	PrismaPrependEnv = "PRISMA_"

	// ConfigVersionEnv enviroment variable
	ConfigVersionEnv = PrismaPrependEnv + "CONFIG_VERSION"

	// LogLevelEnv enviroment variable
	LogLevelEnv = PrismaPrependEnv + "LOG_LEVEL"

	// APIEnv enviroment variable
	APIEnv = PrismaPrependEnv + "API"

	// DisableProtectConfigEnv enviroment variable
	DisableProtectConfigEnv = PrismaPrependEnv + "DISABLE_PROTECT_CONFIG"

	// OrgTenantEnv enviroment variable
	OrgTenantEnv = PrismaPrependEnv + "ORG_TENANT"

	// OrgCloudAccountEnv enviroment variable
	OrgCloudAccountEnv = PrismaPrependEnv + "CLOUD_ACCOUNT"

	// OpsEnv enviroment variable
	OpsEnv = PrismaPrependEnv + "OPS"

	// KubeMatchTagsEnv enviroment variable
	KubeMatchTagsEnv = PrismaPrependEnv + "KUBE_MATCH_TAGS"

	// KubeMatchNamesEnv enviroment variable
	KubeMatchNamesEnv = PrismaPrependEnv + "KUBE_MATCH_NAMES"

	// KubeMatchAnyEnv enviroment variable
	KubeMatchAnyEnv = PrismaPrependEnv + "KUBE_MATCH_ANY"
)
