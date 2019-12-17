package main

import (
	"fmt"
	"github.com/PPerminov/mekube/pkg/mekube"
)

func main() {
	err := mekube.MErgeKUBErnetesconfigfiles()
	if err != nil {
		fmt.Println(err)
	}
}