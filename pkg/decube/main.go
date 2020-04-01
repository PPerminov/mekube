package decube

import (
	"github.com/PPerminov/mekube/pkg/load"
	"github.com/sirupsen/logrus"
)

func Run(file string) {
	a, err := load.Load(file)
	if err != nil {
		logrus.Fatal(err)
	}

}
