package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/akamensky/argparse"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Context struct {
	isCurrent bool
	name      string
}

func main() {

	// Create new parser object
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	// Create string flag
	var next *bool = parser.Flag("n", "nextContext", &argparse.Options{Required: false, Help: "Set to next context"})
	var previus *bool = parser.Flag("p", "previusContext", &argparse.Options{Required: false, Help: "Set to previus context"})
	var showCurrentContext *bool = parser.Flag("c", "showCurrentContext", &argparse.Options{Required: false, Help: "Show Current Context"})
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
		getCurrentContext(kubeconfig)
	}

	if *next {
		goNextContext(kubeconfig)
	}

	if *previus {
		goPreviusContext(kubeconfig)
	}
}

func goNextContext(kubeconfigPath string) (err error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.RawConfig()
	if err != nil {
		return fmt.Errorf("error getting RawConfig: %w", err)
	}
	keys := make([]string, 0, len(config.Contexts))
	for name := range config.Contexts {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	// index := 0
	indexPrev := 0
	for index, name := range keys {
		if name == config.CurrentContext {
			indexPrev = index + 1
			if index == len(keys)-1 {
				indexPrev = 0
			}
			break
		}
	}
	for index, name := range keys {
		if index == indexPrev {
			switchContext(name, kubeconfigPath)
		}
	}
	getCurrentContext(kubeconfigPath)
	return nil
}

func goPreviusContext(kubeconfigPath string) (err error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.RawConfig()
	if err != nil {
		return fmt.Errorf("error getting RawConfig: %w", err)
	}
	keys := make([]string, 0, len(config.Contexts))
	for name := range config.Contexts {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	// index := 0
	indexPrev := 0
	for index, name := range keys {
		if name == config.CurrentContext {
			indexPrev = index - 1
			if index == 0 {
				indexPrev = len(keys) - 1
			}
			break
		}
	}
	for index, name := range keys {
		if index == indexPrev {
			switchContext(name, kubeconfigPath)
			break
		}
	}
	getCurrentContext(kubeconfigPath)
	return nil
}

func getCurrentContext(kubeconfigPath string) (err error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.RawConfig()
	if err != nil {
		return fmt.Errorf("error getting RawConfig: %w", err)
	}
	fmt.Println(config.CurrentContext)

	return nil
}

func switchContext(ctx, kubeconfigPath string) (err error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.RawConfig()
	if err != nil {
		return fmt.Errorf("error getting RawConfig: %w", err)
	}

	if config.Contexts[ctx] == nil {
		return fmt.Errorf("context %s doesn't exists", ctx)
	}

	config.CurrentContext = ctx
	err = clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), config, true)
	if err != nil {
		return fmt.Errorf("error ModifyConfig: %w", err)
	}

	return nil
}
