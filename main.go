package main

import (
	"fmt"
	"github.com/hbagdi/go-kong/kong"
	"kongo/kongClient"
)

func printServices(services []*kong.Service) {
	for idx, service := range services {
		fmt.Println("Service #", idx, ": ", *service.Name)
	}
}

func printUpstreams(upstreams []*kong.Upstream) {
	for idx, upstream := range upstreams {
		fmt.Println("Upstream #", idx, ": ", *upstream.Name)
	}
}

func main() {
	baseUrl := "http://localhost:8001"
	kongo, _ := kongClient.NewKongo(&baseUrl)
	v, _ := kongo.GetVersion()
	fmt.Println(*v)


	truncateKong(kongo)
	generateTestThings(kongo)
}

func truncateKong(kongo *kongClient.Kongo) {
	err := kongo.DeleteAllTargets()
	if err != nil {
		fmt.Println("Error deleting all Targets: ", err)
	}

	err = kongo.DeleteAllUpstreams()
	if err != nil {
		fmt.Println("Error deleting all Upstreams: ", err)
	}

	err = kongo.DeleteAllRoutes()
	if err != nil {
		fmt.Println("Error deleting all Routes: ", err)
	}

	err = kongo.DeleteAllServices()
	if err != nil {
		fmt.Println("Error deleting all Streams: ", err)
	}
}

func generateTestThings(kongo *kongClient.Kongo) error {
	k8sService := kongClient.K8sService{
		Addresses: []*string{kong.String("localhost")},
		Name:      "steve-test-service-one",
		Path:      "/testing-1-2-3",
		Port:      80,
	}

	registered, err := kongo.RegisterK8sService(&k8sService)
	if err != nil {
		_, err = fmt.Println("None create the things: ", err)
		return err
	}

	_, err = fmt.Println(registered)
	return err
}
