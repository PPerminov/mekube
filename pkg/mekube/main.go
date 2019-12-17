package mekube

import (
	"errors"
	"flag"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"path/filepath"
	"strings"
)

const ConfigSeparator = ","

func MErgeKUBErnetesconfigfiles() (err error) {
	var configs []*api.Config
	var defaultInside = false
	var defaultConfig, _ = filepath.Abs("~/.kube/config")

	var rawFiles = flag.String("files",
		"",
		"Сomma separated list of configs. First one will be forcibly inserted in next one until the end of the list")
	var mergeToDefault = flag.Bool("todefault",
		false, "If setted to true and there are no default config in files list then all new configs will be finally merged with default one")
	var output = flag.String("output",
		"/dev/stdout",
		"Path to output merged config. Stdout by default")

	flag.Parse()

	files := strings.Split(*rawFiles, ConfigSeparator)

	if files[0] == "" {
		return errors.New("no files presented")
	}

	for index := range files {
		files[index],err = filepath.Abs(files[index])
		if err != nil {
			return
		}
		if files[index] == defaultConfig {
			defaultInside = true
		}
	}

	if *mergeToDefault && !defaultInside {
		files=append(files, defaultConfig)
	}

	for _, file := range files {
		a, err := clientcmd.LoadFromFile(file)
		if err != nil {
			return err
		}
		configs = append(configs, a)
	}

	for len(configs) != 1 {
		config1 := configs[0]
		config2 := configs[1]
		config2.APIVersion = config1.APIVersion
		config2.CurrentContext = config1.CurrentContext
		config2.Kind = config1.Kind
		config2.Preferences = config1.Preferences
		for k, v := range config1.Clusters {
			config2.Clusters[k] = v
		}
		for k, v := range config1.AuthInfos {
			config2.AuthInfos[k] = v
		}
		for k, v := range config1.Contexts {
			config2.Contexts[k] = v
		}
		for k, v := range config1.Extensions {
			config2.Extensions[k] = v
		}
		configs = configs[1:]
		configs[0] = config2
	}

	return clientcmd.WriteToFile(*configs[0], *output)

}
