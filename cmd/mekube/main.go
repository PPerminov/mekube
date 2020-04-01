package main

import (
	"flag"
	"github.com/PPerminov/mekube/pkg/decube"
	"github.com/PPerminov/mekube/pkg/mekube"
)

func main() {
	file := flag.String("file", "./kubeconfig", "File to insert into config")
	Delete := flag.Bool("delete", false, "Set this flag if you want to delete some contexts from default config")
	flag.Parse()
	if *Delete {
		decube.Run()
	} else {
		mekube.MErgeKUBErnetesconfigfiles(*file)
	}
}
