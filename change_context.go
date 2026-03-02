package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	var showCurrentContext = flag.Bool("c", false, "Show Current Context")
	var showListContext = flag.Bool("s", false, "Show Context List")
	var setContext = flag.String("n", "", "Set new context")

	flag.Parse()

	var kubeConfig string

	if home := homedir.HomeDir(); home != "" {
		kubeConfig = filepath.Join(home, ".kube", "config")
	}

	if *showCurrentContext {
		ctx, _ := getCurrentContext(kubeConfig)
		fmt.Printf("Current Context: %s\n", ctx)
	}

	if *showListContext {
		ListContext(kubeConfig)
	}

	if *setContext != "" {
		newContext := fmt.Sprint(*setContext)
		fmt.Println(newContext)
		setNewContext(kubeConfig, newContext)
	}

}

func getCurrentContext(kubeConfigPath string) (ctx string, err error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.RawConfig()
	if err != nil {
		return "", err
	}

	return strings.Trim(config.CurrentContext, "\n"), nil
}

func ListContext(kubeconfigPath string) (ctx string, err error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.RawConfig()
	if err != nil {
		return "", err
	}

	for name := range config.Contexts {

		if name == config.CurrentContext {
			fmt.Printf("--> %s <--\n", name)
		} else {
			fmt.Printf("    %s\n", name)
		}
	}

	return "", nil
}

func setNewContext(kubeconfigPath string, newContext string) (err error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.RawConfig()
	if err != nil {
		return fmt.Errorf("error getting RawConfig: %w", err)
	}

	if config.Contexts[newContext] == nil {
		return fmt.Errorf("context %s doesn't exists", newContext)
	}

	config.CurrentContext = newContext
	err = clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), config, true)
	if err != nil {
		return fmt.Errorf("error ModifyConfig: %w", err)
	}

	return nil

}
