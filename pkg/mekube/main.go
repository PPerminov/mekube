package mekube

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	//"flag"
	"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/tools/clientcmd/api"
	"os/user"
	"path/filepath"
	//"github.com/davecgh/go-spew/spew"
)

func resolvePath(path string) (p string) {

	if path[0] == []byte("~")[0] {
		usr, err := user.Current()
		if err != nil {
			return ""
		}
		p, err = filepath.Abs(usr.HomeDir + path[1:])
		if err != nil {
			return ""
		}
	} else {
		p, err := filepath.Abs(path)
		if err != nil {
			//spew.Dump(p)
			return ""
		}
		return p
	}
	return p
}

func readFromInput() string {
	reader := bufio.NewReader(os.Stdin) //create new reader, assuming bufio imported
	var storageString string
	storageString, _ = reader.ReadString('\n')
	return storageString
}

func MErgeKUBErnetesconfigfiles() error {
	//var configs []*api.Config
	var defaultConfig = resolvePath("~/.kube/config")

	file := flag.String("file", "./kubeconfig", "File to insert into config")
	flag.Parse()
	newFileFullPath := resolvePath(*file)
	newFileConfig, err := clientcmd.LoadFromFile(newFileFullPath)
	if err != nil {
		return err
	}

	var newCluster string
	for k := range newFileConfig.Clusters {
		newCluster = k
		break
	}
	clusters := map[string]bool{}
	a, _ := clientcmd.LoadFromFile(defaultConfig)
	for k := range a.Clusters {
		clusters[k] = true
	}
	newClusterName := newCluster
	for {
		if _, ok := clusters[newClusterName]; !ok {
			break
		}
		fmt.Println("BadClusterName")
		newClusterName = readFromInput()
	}

	a.Clusters[newClusterName] = newFileConfig.Clusters[newCluster]
	a.AuthInfos[newClusterName] = newFileConfig.AuthInfos[newCluster]
	a.Contexts[newClusterName] = newFileConfig.Contexts[newCluster]
	a.Extensions[newClusterName] = newFileConfig.Extensions[newCluster]
	return clientcmd.WriteToFile(*a, defaultConfig)
}
