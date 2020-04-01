package main

import (
	"fmt"
	"github.com/PPerminov/mekube/pkg/load"
	"golang.org/x/crypto/sha3"
	"k8s.io/client-go/tools/clientcmd/api"
	"strconv"
)

func Run(a api.Config) {

}

type Cluster struct {
	name    string
	hash    string
	cluster *api.Cluster
}

type Auth struct {
	name    string
	hash    string
	cluster *api.Cluster
}

type Auth struct {
	name    string
	hash    string
	cluster *api.Cluster
}

type Config struct {
	ContextsByName map[string]*api.Context
	ContextsByHash map[string]*api.Context
	ClustersByName map[string]*Cluster
	ClustersByHash map[string]*Cluster
	AuthByName     map[string]*api.AuthInfo
	AuthBtHash     map[string]*api.AuthInfo
}

func DisassembleCluster() {
	newConfi := Config{}
	newConfi.ClustersByHash = map[string]*Cluster{}
	newConfi.ClustersByName = map[string]*Cluster{}
	a, _ := load.Load("/Users/pavelperminov/.kube/config")
	Clusters := a.Clusters
	for name, value := range Clusters {

		CA := value.CertificateAuthority
		CAD := string(value.CertificateAuthorityData)
		TlsVerify := strconv.FormatBool(value.InsecureSkipTLSVerify)
		Server := value.Server
		hash := fmt.Sprintf("%x", sha3.Sum512([]byte(fmt.Sprintf("%s%s%s%s", CA, CAD, TlsVerify, Server))))
		c := Cluster{cluster: value, name: name, hash: hash}
		newConfi.ClustersByHash[hash] = &c
		newConfi.ClustersByName[name] = &c
	}
	fmt.Println(newConfi)
}

func main() {

	DisassembleCluster()
}
