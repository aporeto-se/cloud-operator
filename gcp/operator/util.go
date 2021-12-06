package operator

import (
	"encoding/base64"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // register GCP auth provider
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/aporeto-se/cloud-operator/gcp/operator/cache"
)

func getKubernetesClientset(cluster *cache.Cluster) (*kubernetes.Clientset, error) {

	zap.L().Debug("entering getKubernetesClientset")

	cert, err := base64.StdEncoding.DecodeString(cluster.MasterAuth.ClusterCaCertificate)
	if err != nil {
		zap.L().Debug("returning getKubernetesClientset with error(s)")
		return nil, err
	}

	ret := api.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters:   map[string]*api.Cluster{},  // Clusters is a map of referencable names to cluster configs
		AuthInfos:  map[string]*api.AuthInfo{}, // AuthInfos is a map of referencable names to user configs
		Contexts:   map[string]*api.Context{},  // Contexts is a map of referencable names to context configs
	}

	ret.Clusters[cluster.Name] = &api.Cluster{
		CertificateAuthorityData: cert,
		Server:                   "https://" + cluster.Endpoint,
	}
	// Just reuse the context name as an auth name.
	ret.Contexts[cluster.Name] = &api.Context{
		Cluster:  cluster.Name,
		AuthInfo: cluster.Name,
	}
	// GCP specific configation; use cloud platform scope.
	ret.AuthInfos[cluster.Name] = &api.AuthInfo{
		AuthProvider: &api.AuthProviderConfig{
			Name: "gcp",
			Config: map[string]string{
				"scopes": "https://www.googleapis.com/auth/cloud-platform",
			},
		},
	}

	kubeConfig, err := clientcmd.NewNonInteractiveClientConfig(ret, cluster.Name, &clientcmd.ConfigOverrides{CurrentContext: cluster.Name}, nil).ClientConfig()
	if err != nil {
		zap.L().Debug("returning getKubernetesClientset with error(s)")
		return nil, err
	}

	kubernetesClientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		zap.L().Debug("returning getKubernetesClientset with error(s)")
		return nil, err
	}

	zap.L().Debug("returning getKubernetesClientset")
	return kubernetesClientset, nil

}
