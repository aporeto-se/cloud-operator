package operator

import (
	"context"
	"fmt"
	"sync"

	builder "github.com/aporeto-se/enforcerd-kube-builder"
	prisma_api "github.com/aporeto-se/prisma-sdk-go-v2/api"
	prisma_types "github.com/aporeto-se/prisma-sdk-go-v2/types"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"

	"github.com/aporeto-se/cloud-operator/common/processors"
	"github.com/aporeto-se/cloud-operator/common/reportwrapper"
	"github.com/aporeto-se/cloud-operator/common/tag"
	"github.com/aporeto-se/cloud-operator/gcp/operator/cache"
	"github.com/aporeto-se/cloud-operator/gcp/types"

	lib_types "github.com/aporeto-se/cloud-operator/common/types"
)

// Client the AWS Client implementation
type Client struct {
	*cache.Cache
	cloudAccountPrismaClient *prisma_api.Client
	cloudOperatorConfig      *types.CloudOperatorConfig
	api                      string
	accountID                string
	protectConfig            bool
	orgTenant                string
	orgCloudAccount          string
	namespace                string
}

// NewClient returns new Client or error
func NewClient(ctx context.Context, config *Config) (*Client, error) {

	var errors *multierror.Error

	if config.CloudOperatorConfig == nil {
		errors = multierror.Append(errors, fmt.Errorf("entity CloudOperatorConfig is required"))
	}

	if config.PrismaClient == nil {
		errors = multierror.Append(errors, fmt.Errorf("entity PrismaClient is required"))
	}

	err := errors.ErrorOrNil()
	if err != nil {
		zap.L().Debug("returning Build with error(s)")
		return nil, err
	}

	project, err := config.CloudOperatorConfig.GetGCloudProject()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	zone, err := config.CloudOperatorConfig.GetGCloudZone()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	api, err := config.CloudOperatorConfig.GetAPI()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	orgTenant, err := config.CloudOperatorConfig.GetOrgTenant()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	orgCloudAccount, err := config.CloudOperatorConfig.GetOrgCloudAccount()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	protectConfig := true
	if config.CloudOperatorConfig.DisableProtectConfig {
		protectConfig = false
	}

	err = errors.ErrorOrNil()
	if err != nil {
		zap.L().Debug("returning Build with error(s)")
		return nil, err
	}

	// At this point prismaClient is NOT nil
	accountID, err := config.PrismaClient.AccountID(ctx)
	if err != nil {
		zap.L().Debug("returning Build with error(s)")
		return nil, err
	}

	cache, err := cache.NewConfig().SetProject(project).SetZone(zone).Build(ctx)
	if err != nil {
		zap.L().Debug("returning Build with error(s)")
		return nil, err
	}

	namespace := "/" + orgTenant + "/" + orgCloudAccount

	if err != nil {
		zap.L().Debug("returning Build with error(s)")
		return nil, err
	}

	return &Client{
		cloudAccountPrismaClient: config.PrismaClient,
		cloudOperatorConfig:      config.CloudOperatorConfig,
		accountID:                accountID,
		api:                      api,
		protectConfig:            protectConfig,
		orgTenant:                orgTenant,
		orgCloudAccount:          orgCloudAccount,
		namespace:                namespace,
		Cache:                    cache,
	}, nil

}

// Run run and returns report. Any error(s) will be wrapped in report.
func (t *Client) Run(ctx context.Context, filter *lib_types.Filter) *lib_types.Report {

	zap.L().Debug("entering Run")

	tagMatcher, _ := tag.NewMatcher(&t.cloudOperatorConfig.Filter, filter)

	report := lib_types.NewReport(cloudProvider)

	// DHCP
	if t.cloudOperatorConfig.HasOp(lib_types.OpDHCP) {
		report.SetDHCP(t.dhcpReport(ctx))
	} else {
		zap.L().Debug("DHCP operation is disabled")
	}

	// The Namespace and Auth options are present for both Compute and Kubernetes. If either Compute or Kubernetes
	// has the operation set we will run it. The subfunction will determine if it is to be ran for Compute, Kubernetes
	// or both.

	// Namespace
	if t.cloudOperatorConfig.HasOp(lib_types.OpNamespaceRogueDelete) ||
		t.cloudOperatorConfig.HasOp(lib_types.OpNamespaceComputeCreate) ||
		t.cloudOperatorConfig.HasOp(lib_types.OpNamespaceComputeDelete) ||
		t.cloudOperatorConfig.HasOp(lib_types.OpNamespaceKubeCreate) ||
		t.cloudOperatorConfig.HasOp(lib_types.OpNamespaceKubeDelete) {

		zap.L().Debug("Namespace operation is enabled")

		nsprocessor, _ := processors.NewNamespaceProcessor(&t.cloudOperatorConfig.CloudOperatorConfig, t.cloudAccountPrismaClient)

		for _, account := range t.ServiceAccounts {
			if account.ComputeInstancesLen() > 0 {
				nsprocessor.AddCompute(account.NamespaceName)
				zap.L().Debug(fmt.Sprintf("compute namespace %s added to add list", account.NamespaceName))
			} else {
				zap.L().Debug(fmt.Sprintf("Service Account %s has NO compute instances", account.Email))
			}
		}

		for _, cluster := range t.Clusters {
			nsprocessor.AddKube(cluster.Name)
			zap.L().Debug(fmt.Sprintf("kubernetes namespace %s added to add list", cluster.Name))
		}

		report.SetNamespace(nsprocessor.Process(ctx))

	} else {
		zap.L().Debug("Namespace operation is disabled")
	}

	// DHCP is neither a Compute or Kubernetes op
	if t.cloudOperatorConfig.HasOp(lib_types.OpComputeAuth) || t.cloudOperatorConfig.HasOp(lib_types.OpKubeAuth) {
		report.SetAuth(t.authReport(ctx))
	} else {
		zap.L().Debug("Auth operation is disabled")
	}

	// Kubernetes is a Kubernetes only op. We check to see if any Kubernetes Ops are present and if so we will
	// run it. The subfunction will determine which options to execute.
	runKube := false
	for _, op := range t.cloudOperatorConfig.Ops {

		switch op {

		// case types.OpDHCP:
		// case types.OpComputeNamespace:
		// case types.OpComputeAuth:
		// case types.OpKubeNamespace:

		case lib_types.OpKubeAuth:
			runKube = true

		case lib_types.OpKubeAPINet:
			runKube = true

		case lib_types.OpKubeDNSNet:
			runKube = true

		case lib_types.OpKubeNodesNet:
			runKube = true

		case lib_types.OpKubeEnforcer:
			runKube = true

		}

	}

	if runKube {
		report.SetKubernetes(t.kubernetesReports(ctx, tagMatcher))
	} else {
		zap.L().Debug("Kubernetes operations are disabled")
	}

	zap.L().Debug("returning Run")
	return report.Build()
}

func (t *Client) dhcpReport(ctx context.Context) *lib_types.DHCPReport {

	zap.L().Debug("entering dhcpReport")

	report := lib_types.NewDHCPReport().
		SetStatus(lib_types.OpStatusFailed)

	prismaConfig := prisma_types.NewPrismaConfig(dhcpImportLabel)

	prismaConfig.AddExternalnetwork(
		prisma_types.NewExternalnetwork("GCP DHCP").
			SetDescription("auto-generated by Cloud Operator").
			SetProtected(t.protectConfig).
			SetPropagate(true).
			AddEntry("169.254.169.254"))

	egressRule := prisma_types.NewRule().
		SetTrafficActionAllow().
		AddUDPProtocolPort(67).
		AddUDPProtocolPort(68).
		AddObject(
			"@org:cloudaccount="+t.orgCloudAccount,
			"@org:tenant="+t.orgTenant,
			"externalnetwork:name=GCP DHCP")

	prismaConfig.AddNetworkrulesetpolicy(
		prisma_types.NewNetworkrulesetpolicy("GCP DHCP Linux").
			SetDescription("auto-generated by Cloud Operator").
			AddOutgoingRule(egressRule).
			AddSubject(
				"@org:cloudaccount="+t.orgCloudAccount,
				"@org:tenant="+t.orgTenant,
				"@os:host=linux",
			).
			SetProtected(t.protectConfig).
			SetPropagate(true))

	prismaConfig.AddNetworkrulesetpolicy(
		prisma_types.NewNetworkrulesetpolicy("GCP DHCP Windows").
			SetDescription("auto-generated by Cloud Operator").
			AddOutgoingRule(egressRule).
			AddSubject(
				"@org:cloudaccount="+t.orgCloudAccount,
				"@org:tenant="+t.orgTenant,
				"@os:host=windows",
			).
			SetProtected(t.protectConfig).
			SetPropagate(true))

	zap.L().Debug(fmt.Sprintf("Importing Prisma API config for %s", dhcpImportLabel))
	err := t.cloudAccountPrismaClient.ImportPrismaConfig(ctx, prismaConfig)

	if err != nil {
		zap.L().Debug("returning dhcpReport with error(s)")
		return report.SetError(err)
	}

	zap.L().Debug("returning dhcpReport")
	return report.SetStatus(lib_types.OpStatusCompleted)
}

func (t *Client) authReport(ctx context.Context) *lib_types.AuthReport {

	report := lib_types.NewAuthReport().
		SetStatus(lib_types.OpStatusFailed)

	prismaConfig := prisma_types.NewPrismaConfig(authImportLabel)

	// Auth Policies for Kubernetes

	// OpDHCP OpComputeNamespace OpKubeNamespace OpKubeAuth OpKubeAPINet OpKubeDNSNet OpKubeNodesNet OpKubeEnforcer

	if t.cloudOperatorConfig.HasOp(lib_types.OpComputeAuth) {

		for _, account := range t.ServiceAccounts {

			add := false
			if account.ComputeInstancesLen() > 0 {
				// If the service account has any compute instances we add it
				add = true
			}

			if !add {
				// If the service account has NO Kubernetes instances we add it
				// but if it does have Kubernetes instances (and no compute) then
				// we do NOT add it.
				if account.ClustersLen() <= 0 {
					add = true
				}
			}

			if add {
				prismaConfig.AddApiauthorizationpolicy(
					prisma_types.NewAPIAuthorizationPolicy("instances:"+account.Email).
						SetDescription("auto-generated cloud operator policy").
						SetProtected(t.protectConfig).
						SetAuthorizedNamespace(t.namespace+"/"+account.NamespaceName).
						AddAuthorizedIdentity("@auth:role=enforcer").
						AddSubject(realm, accountLabel+t.accountID, role+account.Email))

				zap.L().Debug(fmt.Sprintf("Service Account %s added to Auth Policy for instance namespace %s", account.Email, account.NamespaceName))

			} else {
				zap.L().Debug(fmt.Sprintf("Service Account %s NOT added to Auth Policy", account.Email))
			}

		}

	} else {
		zap.L().Debug("Auth operation for Compute disabled")
	}

	if t.cloudOperatorConfig.HasOp(lib_types.OpKubeAuth) {

		for _, cluster := range t.Clusters {
			for _, account := range cluster.ServiceAccounts {
				prismaConfig.AddApiauthorizationpolicy(
					prisma_types.NewAPIAuthorizationPolicy(cluster.Name+":"+account.Email).
						SetDescription("auto-generated cloud operator policy").
						SetProtected(t.protectConfig).
						SetAuthorizedNamespace(t.namespace+"/"+cluster.Name).
						AddAuthorizedIdentity("@auth:role=enforcer").
						AddSubject(realm, accountLabel+t.accountID, role+account.Email))

				zap.L().Debug(fmt.Sprintf("Service Account %s added to Auth Policy for cluster namespace %s", account.Email, cluster.Name))

			}
		}

	} else {
		zap.L().Debug("Auth operation for Kubernetes disabled")
	}

	zap.L().Debug(fmt.Sprintf("Importing Prisma API config for %s", authImportLabel))
	err := t.cloudAccountPrismaClient.ImportPrismaConfig(ctx, prismaConfig)

	if err != nil {
		zap.L().Debug("returning dhcpReport with error(s)")
		return report.SetError(err)
	}

	return report.SetStatus(lib_types.OpStatusCompleted)

}

func (t *Client) kubernetesReports(ctx context.Context, tagMatcher *tag.Matcher) *lib_types.KubernetesReports {

	zap.L().Debug("entering kubernetesReports")

	var wg sync.WaitGroup

	wrapper := reportwrapper.NewWrapper()

	for _, _cluster := range t.Clusters {

		cluster := _cluster

		if tagMatcher.MatchKubeCluster(cluster.Name, cluster.ResourceLabels) {
			zap.L().Debug(fmt.Sprintf("Cluster %s is a match", cluster.Name))

			wg.Add(1)
			go func() {
				defer wg.Done()
				wrapper.AddKube(t.kubernetesReport(ctx, cluster))
			}()

		} else {
			zap.L().Debug(fmt.Sprintf("Cluster %s is NOT a match", cluster.Name))
		}

	}

	wg.Wait()

	zap.L().Debug("returning kubernetesReports")
	return wrapper.Build()
}

func (t *Client) kubernetesReport(ctx context.Context, cluster *cache.Cluster) *lib_types.KubernetesReport {

	zap.L().Debug("entering kubernetesReport")

	report := lib_types.NewKubernetesReport(cluster.Name).
		SetStatus(lib_types.OpStatusFailed)

	if cluster.Status != "RUNNING" {
		zap.L().Debug("returning kubernetesReport (not active)")
		return report.SetStatus(lib_types.OpStatusNotReady)
	}

	endpoint := cluster.Endpoint
	if endpoint == "" {
		zap.L().Debug("returning kubernetesReport with error(s)")
		return report.SetError(fmt.Errorf("unable to determine cluster endpoint"))
	}

	endpoint = "https://" + endpoint

	zap.L().Debug(fmt.Sprintf("endpoint=%s", endpoint))

	kubernetesClientset, err := getKubernetesClientset(cluster)
	if err != nil {
		zap.L().Debug("returning kubernetesReport with error(s)")
		return report.SetError(err)
	}

	prismaClient, err := t.cloudAccountPrismaClient.NewClient(ctx, cluster.Name)
	if err != nil {
		zap.L().Debug("returning kubernetesReport with wrapped error(s)")
		return report.SetError(err)
	}

	kubeprocessor, _ := processors.NewKubeProcessor(cluster.Name, &t.cloudOperatorConfig.CloudOperatorConfig, prismaClient)
	err = kubeprocessor.
		AddCidrBlocks(cluster.ClusterIpv4Cidr).
		SetKubernetesDaemonsetBuilder(builder.NewGke(prismaClient.GetNamespacePath(), t.api)).
		SetEndpoint(endpoint).
		SetKubernetesClientset(kubernetesClientset).
		Process(ctx)

	if err != nil {
		zap.L().Debug("returning kubernetesReport with wrapped error(s)")
		return report.SetError(err)
	}

	zap.L().Debug("returning kubernetesReport")

	return report.SetStatus(lib_types.OpStatusCompleted)

}
