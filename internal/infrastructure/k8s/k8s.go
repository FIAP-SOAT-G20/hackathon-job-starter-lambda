package k8s

import (
	"context"
	"fmt"
	"os"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func ConnectToK8s(ctx context.Context, logger *logger.Logger, cfg *config.LambdaConfig) (*kubernetes.Clientset, error) {
	logger.InfoContext(ctx, "Connecting to k8s")

	var config *rest.Config
	var err error

	// Check if we're running inside a Kubernetes cluster
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
		// We're inside a cluster, use in-cluster config
		logger.InfoContext(ctx, "Using in-cluster Kubernetes configuration")
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to load in-cluster config: %v", err)
		}
	} else {
		// We're outside the cluster, use kubeconfig
		var kubeconfig string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = home + "/.kube/config"
		}

		logger.InfoContext(ctx, "Using kubeconfig", "kubeconfig", kubeconfig)
		if cfg.K8S.ContextName == "" {
			config, err = clientcmd.BuildConfigFromFlags(cfg.K8S.MasterUrl, kubeconfig)
		} else {
			config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
				&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
				&clientcmd.ConfigOverrides{
					CurrentContext: cfg.K8S.ContextName,
				},
			).ClientConfig()
		}

		if err != nil {
			return nil, fmt.Errorf("failed to load kubeconfig: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate Kubernetes Client: %v", err)
	}
	return clientset, nil

}
