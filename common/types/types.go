package types

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
)

// ================================================================================================

// CloudOperatorConfig top level configuration
type CloudOperatorConfig struct {

	// LogLevel Level to log
	LogLevel LogLevel `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`

	// API Prisma API Endpoint eg. https://prisma.tld
	API string `json:"api" yaml:"api"`

	// DisableProtectConfig by default Prisma Config is protected from deletion. Setting this to true
	// overries the protection.
	DisableProtectConfig bool

	// OrgTenant is the Prisma Account ID. It should be a large number (for example 806775361903163392)
	OrgTenant string

	// OrgCloudAccount is the cloud account name
	OrgCloudAccount string

	// Ops operations
	Ops []Op `json:"ops" yaml:"ops"`

	// Filter the filter
	Filter Filter `json:"filter" yaml:"filter"`
}

// SetFromEnv sets attributes and types from env variables as defined in
// constants file. If attribute is not of the expected type an error will be
// returned. If child entities exist and are initialized (not nil) then a call
// to the child entities SetFromEnv() will be executed. Any errors will be aggregated
// and returned.
func (t *CloudOperatorConfig) SetFromEnv() error {

	var errors *multierror.Error

	logLevelString := os.Getenv(LogLevelEnv)
	api := os.Getenv(APIEnv)
	orgTenant := os.Getenv(OrgTenantEnv)
	orgCloudAccount := os.Getenv(OrgCloudAccountEnv)
	opsString := os.Getenv(OpsEnv)

	if logLevelString != "" {
		logLevel, err := LogLevelFromString(logLevelString)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			t.LogLevel = logLevel
		}
	}

	if api != "" {
		t.API = api
	}

	disableProtectConfig, err := GetEnvBool(DisableProtectConfigEnv)
	if err != nil {
		errors = multierror.Append(errors, err)
	} else {
		t.DisableProtectConfig = disableProtectConfig
	}

	if orgTenant != "" {
		t.OrgTenant = orgTenant
	}

	if orgCloudAccount != "" {
		t.OrgCloudAccount = orgCloudAccount
	}

	if opsString != "" {
		for _, v := range strings.Split(opsString, ",") {
			v = strings.TrimSpace(v)
			if v != "" {
				op, err := OpFromString(v)
				if err == nil {
					t.AddOps(op)
				} else {
					errors = multierror.Append(errors, err)
				}
			}
		}
	}

	err = t.Filter.SetFromEnv()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	return errors.ErrorOrNil()
}

// SetLogLevel sets type and returns self
func (t *CloudOperatorConfig) SetLogLevel(v LogLevel) *CloudOperatorConfig {
	t.LogLevel = v
	return t
}

// GetLogLevel returns type and returns self
func (t *CloudOperatorConfig) GetLogLevel() LogLevel {
	if t.LogLevel == LogLevelInvalid {
		return LogLevelInfo
	}
	return t.LogLevel
}

// GetNamespace returns attribute
func (t *CloudOperatorConfig) GetNamespace() (string, error) {

	var errors *multierror.Error

	orgTenant, err := t.GetOrgTenant()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	orgCloudAccount, err := t.GetOrgCloudAccount()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	err = errors.ErrorOrNil()
	if err != nil {
		return "", err
	}

	return "/" + orgTenant + "/" + orgCloudAccount, nil
}

// SetAPI sets attribute and returns self
func (t *CloudOperatorConfig) SetAPI(v string) *CloudOperatorConfig {
	t.API = v
	return t
}

// GetAPI returns attribute or error
func (t *CloudOperatorConfig) GetAPI() (string, error) {
	var err error
	if t.API == "" {
		err = fmt.Errorf("attribute API (env var %s) is required", APIEnv)
	}
	return t.API, err
}

// SetOrgTenant sets attribute and returns self
func (t *CloudOperatorConfig) SetOrgTenant(v string) *CloudOperatorConfig {
	t.OrgTenant = v
	return t
}

// GetOrgTenant returns attribute or error
func (t *CloudOperatorConfig) GetOrgTenant() (string, error) {
	var err error
	if t.OrgTenant == "" {
		err = fmt.Errorf("attribute OrgTenant (env var %s) is required", OrgTenantEnv)
	}
	return t.OrgTenant, err
}

// SetOrgCloudAccount sets attribute and returns self
func (t *CloudOperatorConfig) SetOrgCloudAccount(v string) *CloudOperatorConfig {
	t.OrgCloudAccount = v
	return t
}

// GetOrgCloudAccount returns attribute or error
func (t *CloudOperatorConfig) GetOrgCloudAccount() (string, error) {
	var err error
	if t.OrgCloudAccount == "" {
		err = fmt.Errorf("attribute OrgCloudAccount (env var %s) is required", OrgCloudAccountEnv)
	}
	return t.OrgCloudAccount, err
}

// AddOps adds type(s) and returns self
func (t *CloudOperatorConfig) AddOps(v ...Op) *CloudOperatorConfig {
	t.Ops = append(t.Ops, v...)
	return t
}

// HasOp returns true if type is present in parent
func (t *CloudOperatorConfig) HasOp(op Op) bool {
	for _, v := range t.Ops {
		if v == op {
			return true
		}
	}
	return false
}

// ================================================================================================

// Filter is a match filter
type Filter struct {
	KubeMatchTags  map[string]string `json:"kubeMatchTags" yaml:"kubeMatchTags"`
	KubeMatchNames []string          `json:"kubeMatchNames" yaml:"kubeMatchNames"`
	KubeMatchAny   bool              `json:"kubeMatchAny" yaml:"kubeMatchAny"`
}

// NewFilter returns new intance of entity
func NewFilter() *Filter {
	return &Filter{}
}

// SetFromEnv sets attributes and types from env variables as defined in
// constants file. If attribute is not of the expected type an error will be
// returned. If child entities exist and are initialized (not nil) then a call
// to the child entities SetFromEnv() will be executed. Any errors will be aggregated
// and returned.
func (t *Filter) SetFromEnv() error {

	var errors *multierror.Error

	kubeMatchTagsString := os.Getenv(KubeMatchTagsEnv)
	kubeMatchNamesString := os.Getenv(KubeMatchNamesEnv)

	// Tag format is pairs are delinated by a comma and key/values are delinated by colon
	// Example: key1:value1,key2:value2,keyN:valueN
	if kubeMatchTagsString != "" {

		if t.KubeMatchTags == nil {
			t.KubeMatchTags = make(map[string]string)
		}

		for _, keyValuePair := range strings.Split(kubeMatchTagsString, ",") {
			keyValuePair = strings.TrimSpace(keyValuePair)
			keyValuePairSplit := strings.Split(keyValuePair, ":")
			if len(keyValuePairSplit) != 2 {
				errors = multierror.Append(errors, fmt.Errorf("keyValuePair %s has invalid syntax. Expected format is key:value", keyValuePair))
			} else {
				key := keyValuePairSplit[0]
				value := keyValuePairSplit[1]
				t.KubeMatchTags[strings.TrimSpace(key)] = strings.TrimSpace(value)
			}
		}
	}

	if kubeMatchNamesString != "" {
		for _, v := range strings.Split(kubeMatchNamesString, ",") {
			v = strings.TrimSpace(v)
			if v != "" {
				t.AddKubeMatchNames(v)
			}
		}
	}

	kubeMatchAny, err := GetEnvBool(KubeMatchAnyEnv)
	if err != nil {
		errors = multierror.Append(errors, err)
	} else {
		t.KubeMatchAny = kubeMatchAny
	}

	return errors.ErrorOrNil()
}

// SetKubeMatchTags sets entity and returns self
func (t *Filter) SetKubeMatchTags(v map[string]string) *Filter {
	t.KubeMatchTags = v
	return t
}

// AddKubeMatchTag adds attributes (Key/Value tag) and returns self
func (t *Filter) AddKubeMatchTag(key, value string) *Filter {
	if t.KubeMatchTags == nil {
		t.KubeMatchTags = make(map[string]string)
	}
	t.KubeMatchTags[key] = value
	return t
}

// HasKubeMatchTag returns true if attribute (Key/Value tag) is present
func (t *Filter) HasKubeMatchTag(key, value string) bool {

	if t.KubeMatchTags == nil {
		return false
	}

	matchValue := t.KubeMatchTags[key]
	if matchValue == value {
		return true
	}

	return false
}

// SetKubeMatchNames sets entity and returns self
func (t *Filter) SetKubeMatchNames(v []string) *Filter {
	t.KubeMatchNames = v
	return t
}

// AddKubeMatchNames adds attributes(s) and returns self
func (t *Filter) AddKubeMatchNames(v ...string) *Filter {
	t.KubeMatchNames = append(t.KubeMatchNames, v...)
	return t
}

// HasKubeMatchName returns true if kubeMatchName is present
func (t *Filter) HasKubeMatchName(v string) bool {
	for _, match := range t.KubeMatchNames {
		if match == v {
			return true
		}
	}
	return false
}

// SetKubeMatchAny sets attribute and return self
func (t *Filter) SetKubeMatchAny(v bool) *Filter {
	t.KubeMatchAny = v
	return t
}

// ================================================================================================

// Report Aggregated Report
type Report struct {
	CloudProvider string             `json:"cloudProvider" yaml:"cloudProvider"`
	RunTime       int64              `json:"runTime" yaml:"runTime"`
	Notes         string             `json:"notes" yaml:"notes"`
	TotalCount    int                `json:"totalCount" yaml:"totalCount"`
	ErrorCount    int                `json:"errorCount" yaml:"errorCount"`
	Namespace     *NamespaceReports  `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	DHCP          *DHCPReport        `json:"dhcp,omitempty" yaml:"dhcp,omitempty"`
	Auth          *AuthReport        `json:"auth,omitempty" yaml:"auth,omitempty"`
	Kubernetes    *KubernetesReports `json:"kubernetes,omitempty" yaml:"kubernetes,omitempty"`
}

// NewReport returns new intance of entity
func NewReport(cloudProvider string) *Report {
	return &Report{
		CloudProvider: cloudProvider,
		Notes:         "Reports will only be shown if operation is enabled",
		RunTime:       time.Now().Unix(),
	}
}

// SetNamespace set entity and return self
func (t *Report) SetNamespace(v *NamespaceReports) *Report {
	t.Namespace = v
	return t
}

// SetDHCP set entity and return self
func (t *Report) SetDHCP(v *DHCPReport) *Report {
	t.DHCP = v
	return t
}

// SetAuth set entity and return self
func (t *Report) SetAuth(v *AuthReport) *Report {
	t.Auth = v
	return t
}

// SetKubernetes set entity and return self
func (t *Report) SetKubernetes(v *KubernetesReports) *Report {
	t.Kubernetes = v
	return t
}

// Build adds entity(s) and returns self
func (t *Report) Build() *Report {

	if t.Namespace != nil {
		t.Namespace.Build()
		t.TotalCount = t.TotalCount + t.Namespace.TotalCount
		t.ErrorCount = t.ErrorCount + t.Namespace.ErrorCount
	}

	if t.DHCP != nil {
		t.TotalCount++
		if t.DHCP.Error != nil {
			t.ErrorCount++
		}
	}

	if t.Auth != nil {
		t.TotalCount++
		if t.Auth.Error != nil {
			t.ErrorCount++
		}
	}

	if t.Kubernetes != nil {
		t.Kubernetes.Build()
		t.TotalCount = t.TotalCount + t.Kubernetes.TotalCount
		t.ErrorCount = t.ErrorCount + t.Kubernetes.ErrorCount
	}

	return t
}

// Errors returns aggregated errors or nil
func (t *Report) Errors() error {

	var errors *multierror.Error

	if t.Namespace != nil {
		err := t.Namespace.Errors()
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	if t.DHCP != nil {
		if t.DHCP.Error != nil {
			errors = multierror.Append(errors, t.DHCP.Error)
		}
	}

	if t.Auth != nil {
		if t.Auth.Error != nil {
			errors = multierror.Append(errors, t.Auth.Error)
		}
	}

	if t.Kubernetes != nil {
		err := t.Kubernetes.Errors()
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	return errors.ErrorOrNil()
}

// ================================================================================================

// NamespaceReports Namespace Reports
type NamespaceReports struct {
	TotalCount int                `json:"totalCount" yaml:"totalCount"`
	ErrorCount int                `json:"errorCount" yaml:"errorCount"`
	Namespaces []*NamespaceReport `json:"namespaces,omitempty" yaml:"namespaces,omitempty"`
}

// NewNamespaceReports returns new entity instance
func NewNamespaceReports() *NamespaceReports {
	return &NamespaceReports{}
}

// AddNamespaces adds entity and returns self
func (t *NamespaceReports) AddNamespaces(v ...*NamespaceReport) *NamespaceReports {
	t.Namespaces = append(t.Namespaces, v...)
	return t
}

// Build adds entity(s) and returns self
func (t *NamespaceReports) Build() *NamespaceReports {

	for _, report := range t.Namespaces {
		t.TotalCount++
		if report.Error != nil {
			t.ErrorCount++
		}
	}

	return t
}

// Errors returns aggregated errors or nil
func (t *NamespaceReports) Errors() error {

	var errors *multierror.Error

	for _, report := range t.Namespaces {
		if report.Error != nil {
			errors = multierror.Append(errors, report.Error)
		}
	}

	return errors.ErrorOrNil()
}

// ==============================================

// NamespaceReport Namespace Report
type NamespaceReport struct {
	Name      string             `json:"name" yaml:"name"`
	Type      CloudEntityType    `json:"type" yaml:"type"`
	Status    OperationStatus    `json:"status" yaml:"status"`
	Operation NamespaceOperation `json:"operation" yaml:"operation"`
	Error     error              `json:"error,omitempty" yaml:"error,omitempty"`
}

// NewNamespaceReport returns new entity instance
func NewNamespaceReport(v string) *NamespaceReport {
	return &NamespaceReport{
		Name: v,
	}
}

// SetType sets type and returns self
func (t *NamespaceReport) SetType(v CloudEntityType) *NamespaceReport {
	t.Type = v
	return t
}

// SetOperation sets type and returns self
func (t *NamespaceReport) SetOperation(v NamespaceOperation) *NamespaceReport {
	t.Operation = v
	return t
}

// SetStatus sets type and returns self
func (t *NamespaceReport) SetStatus(v OperationStatus) *NamespaceReport {
	t.Status = v
	return t
}

// SetError sets entity and returns self
func (t *NamespaceReport) SetError(v error) *NamespaceReport {
	t.Error = v
	return t
}

// ================================================================================================

// KubernetesReports is a report for each Kubernete's clusters
type KubernetesReports struct {
	TotalCount int                 `json:"totalCount" yaml:"totalCount"`
	ErrorCount int                 `json:"errorCount" yaml:"errorCount"`
	Reports    []*KubernetesReport `json:"reports" yaml:"reports"`
}

// AddReports adds entity and returns self
func (t *KubernetesReports) AddReports(v ...*KubernetesReport) *KubernetesReports {
	t.Reports = append(t.Reports, v...)
	return t
}

// NewKubernetesReports returns new entity instance
func NewKubernetesReports() *KubernetesReports {
	return &KubernetesReports{}
}

// Build adds entity(s) and returns self
func (t *KubernetesReports) Build() *KubernetesReports {

	for _, report := range t.Reports {
		t.TotalCount++
		if report.Error != nil {
			t.ErrorCount++
		}
	}

	return t
}

// Errors returns aggregated errors or nil
func (t *KubernetesReports) Errors() error {

	var errors *multierror.Error

	for _, report := range t.Reports {
		if report.Error != nil {
			errors = multierror.Append(errors, report.Error)
		}
	}

	return errors.ErrorOrNil()
}

// ==============================================

// KubernetesReport is a report for each Kubernetes cluster
type KubernetesReport struct {
	Name   string   `json:"name" yaml:"name"`
	Status OpStatus `json:"status" yaml:"status"`
	Error  error    `json:"error,omitempty" yaml:"error,omitempty"`
}

// NewKubernetesReport returns new instance
func NewKubernetesReport(v string) *KubernetesReport {
	return &KubernetesReport{
		Name: v,
	}
}

// SetStatus sets type and returns self
func (t *KubernetesReport) SetStatus(v OpStatus) *KubernetesReport {
	t.Status = v
	return t
}

// SetError sets entity and returns self
func (t *KubernetesReport) SetError(v error) *KubernetesReport {
	t.Error = v
	return t
}

// ================================================================================================

// DHCPReport DHCP Report
type DHCPReport struct {
	Status OpStatus `json:"status" yaml:"status"`
	Error  error    `json:"error,omitempty" yaml:"error,omitempty"`
}

// NewDHCPReport returns new instance
func NewDHCPReport() *DHCPReport {
	return &DHCPReport{}
}

// SetStatus sets type and returns self
func (t *DHCPReport) SetStatus(status OpStatus) *DHCPReport {
	t.Status = status
	return t
}

// SetError sets entity and returns self
func (t *DHCPReport) SetError(err error) *DHCPReport {
	t.Error = err
	return t
}

// ================================================================================================

// AuthReport Auth Report
type AuthReport struct {
	Status OpStatus `json:"status" yaml:"status"`
	Error  error    `json:"error,omitempty" yaml:"error,omitempty"`
}

// NewAuthReport returns new instance
func NewAuthReport() *AuthReport {
	return &AuthReport{}
}

// SetStatus sets type and returns self
func (t *AuthReport) SetStatus(v OpStatus) *AuthReport {
	t.Status = v
	return t
}

// SetError sets entity and returns self
func (t *AuthReport) SetError(v error) *AuthReport {
	t.Error = v
	return t
}

// ================================================================================================

// NewExampleFilterMatchNames returns new example Filter
func NewExampleFilterMatchNames() *Filter {
	return NewFilter().AddKubeMatchNames("cluster1", "cluster2", "clusterN")
}

// NewExampleFilterMatchTags returns new example Filter
func NewExampleFilterMatchTags() *Filter {
	return NewFilter().
		AddKubeMatchTag("key1", "value1").
		AddKubeMatchTag("key2", "value2").
		AddKubeMatchTag("key3", "value3")
}

// NewExampleFilterMatchAny returns new example Filter
func NewExampleFilterMatchAny() *Filter {
	return NewFilter().SetKubeMatchAny(true)
}

// ================================================================================================
