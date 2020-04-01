package load

import (
	"errors"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Load(data interface{}) (a *api.Config, err error) {
	switch data.(type) {
	case string:
		return clientcmd.LoadFromFile(data.(string))
	case []byte:
		return clientcmd.Load(data.([]byte))
	}
	return nil, errors.New("wrong data type. must be a string or binary")
}
