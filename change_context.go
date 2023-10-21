package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/akamensky/argparse"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	// Create new parser object
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	var showCurrentContext *bool = parser.Flag("c", "showCurrentContext", &argparse.Options{Required: false, Help: "Show Current Context"})
	var showListContext *bool = parser.Flag("s", "showListContext", &argparse.Options{Required: false, Help: "Show Context List"})
	var setContext *string = parser.String("n", "newContext", &argparse.Options{Required: false, Help: "Set new context"})
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}

	var kubeconfig string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	if *showCurrentContext {
		ctx, _ := getCurrentContext(kubeconfig)
		fmt.Printf(ctx)
	}

	if *showListContext {
		ListContext(kubeconfig)
	}

	if *setContext != "" {
		newContext := fmt.Sprint(*setContext)
		fmt.Println(newContext)
		setNewContext(kubeconfig, newContext)
	}

}

func getCurrentContext(kubeconfigPath string) (ctx string, err error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
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

	//var contexts []string
	for name, _ := range config.Contexts {

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
