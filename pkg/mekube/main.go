package mekube

import (
	"bufio"
	"fmt"
	"github.com/PPerminov/mekube/pkg/load"
	"github.com/PPerminov/random"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func init() {
	if os.Getenv("DEBUG") == "1" {
		spew.Dump()
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func resolvePath(path string) (p string) {
	var err error
	if path[0] == []byte("~")[0] {
		usr, err := user.Current()
		if err != nil {
			logrus.Debugln(err)
			return ""
		}
		p, err = filepath.Abs(usr.HomeDir + path[1:])
		if err != nil {
			logrus.Debugln(err)
			return ""
		}
	} else {
		p, err = filepath.Abs(path)
		if err != nil {
			logrus.Debugln(err)
			return ""
		}
	}
	return p
}

func readFromInput() string {
	var storageString string
	reader := bufio.NewReader(os.Stdin)
	storageString, err := reader.ReadString('\n')
	if err != nil {
		logrus.Debugln(err)
	}
	storageString = strings.TrimSpace(storageString)
	return storageString
}

func MErgeKUBErnetesconfigfiles(file string) {
	var defaultConfig = resolvePath("~/.kube/config")
	logrus.Debugln("Default Config Path:", defaultConfig)

	logrus.Debugln("New Config File Path:", file)
	newFileFullPath := resolvePath(file)
	logrus.Debugln("New Config File Full Path:", newFileFullPath)
	newFileConfig, err := load.Load(newFileFullPath)
	logrus.Debugln("Load From New File:", newFileConfig)
	logrus.Debugln("Load From New File Error:", err)
	logrus.Debugln("New File Auth:", newFileConfig.AuthInfos)
	logrus.Debugln("New File Clusters:", newFileConfig.Clusters)
	logrus.Debugln("New File Contexts:", newFileConfig.Contexts)
	logrus.Debugln("New File Extensions:", newFileConfig.Extensions)
	if err != nil {
		return
	}
	var newClusterList []string
	for k := range newFileConfig.Contexts {
		newClusterList = append(newClusterList, k)
		break
	}

	logrus.Debugln("Cluster's names to import:", newClusterList)

	a, err := load.Load(defaultConfig)

	if err != nil {
		a = &api.Config{}
		a.APIVersion = "v1"
		a.CurrentContext = ""
		a.Kind = "Config"

		a.Preferences = api.Preferences{}
		a.Extensions = map[string]runtime.Object{}
		a.Contexts = map[string]*api.Context{}
		a.AuthInfos = map[string]*api.AuthInfo{}
		a.Clusters = map[string]*api.Cluster{}
	}
	logrus.Debugln("Load From Default File:", a)
	logrus.Debugln("Load From Default File Error:", err)
	logrus.Debugln("Old File Auth:", a.AuthInfos)
	logrus.Debugln("Old File Clusters:", a.Clusters)
	logrus.Debugln("Old File Contexts:", a.Contexts)
	logrus.Debugln("Old File Extensions:", a.Extensions)
	contexts := map[string]bool{}
	for k := range a.Contexts {
		contexts[k] = true
	}
	//clusters := map[string]bool{}
	//for k := range a.Clusters {
	//	clusters[k] = true
	//}
	//auth := map[string]bool{}
	//for k := range a.AuthInfos {
	//	auth[k] = true
	//}

	logrus.Debugln("Existed Clusters:", contexts)
	CCounter := 0

	for _, CurrentContextName := range newClusterList {
		CurrentContext := newFileConfig.Contexts[CurrentContextName]
		CurrentContextAuthInfoName := CurrentContext.AuthInfo
		CurrentContextClusterName := CurrentContext.Cluster
		CurrentContextAuthInfo := newFileConfig.AuthInfos[CurrentContextAuthInfoName]
		CurrentContextCluster := newFileConfig.Clusters[CurrentContextClusterName]
	Contexts:
		for {
			if _, ok := contexts[CurrentContextName]; !ok {
				break Contexts
			}
			fmt.Printf("Bad Context Name %s. Input new one: ", CurrentContextName)
			CurrentContextName = readFromInput()
		}

		CurrentContextAuthInfoName = fmt.Sprintf("%s-%s-%s", CurrentContextAuthInfoName, "auth", random.RandomString(5, 21))
		CurrentContextClusterName = fmt.Sprintf("%s-%s-%s", CurrentContextClusterName, "cluster", random.RandomString(5, 21))

		logrus.Debugln("Context name to import:", CurrentContextName)
		a.AuthInfos[CurrentContextAuthInfoName] = CurrentContextAuthInfo
		a.Clusters[CurrentContextClusterName] = CurrentContextCluster
		CurrentContext.Cluster = CurrentContextClusterName
		CurrentContext.AuthInfo = CurrentContextAuthInfoName
		a.Contexts[CurrentContextName] = CurrentContext
		if a.CurrentContext == "" {
			a.CurrentContext = CurrentContextName
		}
		CCounter += 1
	}

	logrus.Debugln("Import result", clientcmd.WriteToFile(*a, defaultConfig))
	logrus.Infof("Imported %d contexts", CCounter)
}
