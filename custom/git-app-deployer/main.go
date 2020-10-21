package main

import (
	"fmt"
	"os/exec"

	"github.com/litmuschaos/litmus-go/pkg/clients"
	"github.com/litmuschaos/litmus-go/pkg/status"
)

func main() {

	clients := clients.ClientSets{}
	if err := CreateNamespace("sock-shop"); err != nil {
		fmt.Printf("err: %v", err)
	}

	//kubectl apply -f deploy/sock-shop/
	if err := CreateShockShop("/sock-shop.yaml"); err != nil {
		fmt.Printf("err: %v", err)
	}
	// check status
	if err := status.CheckApplicationStatus("sock-shop", "app=sock-shop", 180, 2, clients); err != nil {
		fmt.Printf("err: %v", err)
	}
}

//CreateNamespace ...
func CreateNamespace(ns string) error {
	//kubectl create ns sock-shop
	if err := exec.Command("kubectl", "create", "ns", "sock-shop").Run(); err != nil {
		return err
	}
	return nil
}

// CreateShockShop ...
func CreateShockShop(path string) error {
	//kubectl create ns sock-shop
	if err := exec.Command("kubectl", "apply", "-f", path, "-n", "sock-shop").Run(); err != nil {
		return err
	}
	return nil
}
