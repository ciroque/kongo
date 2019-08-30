package main

import (
	"fmt"
	"github.com/hbagdi/go-kong/kong"
	"kongo/kongClientWrapper"
	"os"
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

	args :=os.Args[1:]

	//baseUrl := "http://qa-swarmcluster1.sea1.marchex.com:8001"
	baseUrl := "http://localhost:8001"
	kongo, _ := kongClientWrapper.NewKongo(&baseUrl)
	v, _ := kongo.GetVersion()
	fmt.Println(*v)

	command := args[0]
	switch command {
		case "truncate":
			truncateKong(kongo)
		case "list":
			listAllThings(kongo)
		case "generate-test-data":
			generateTestThings(kongo)
		default:
			fmt.Println(fmt.Sprintf("Unknown command '%s'", command))
	}
}

func listAllThings(kongo *kongClientWrapper.Kongo) error {
	upstreams, err := kongo.ListUpstreams()
	if err != nil {
		return fmt.Errorf("error listing Upstreams: %v", err)
	}
	for _, upstream := range upstreams {
		output := fmt.Sprintf("Upstream{ Name: %s, ID: %s}", *upstream.Name, *upstream.ID)
		fmt.Println(output)

		targets, err := kongo.ListTargets(*upstream.ID)
		if err != nil {
			return fmt.Errorf("error listing Targets for Upstream '%s': %v", *upstream.Name, err)
		}
		for _, target := range targets {
			output := fmt.Sprintf("Target{ Target: %s, ID: %s }", *target.Target, *target.ID)
			fmt.Println(output)
		}
	}

	services, err := kongo.ListServices()
	if err != nil {
		return fmt.Errorf("error listing Services: %v", err)
	}

	for _, service := range services {
		output := fmt.Sprintf("Service{ Name: %s, ID: %s. Path: %s, Port: %v }", *service.Name, *service.ID, *service.Path, *service.Port)
		fmt.Println(output)
	}

	routes, err := kongo.ListRoutes()
	if err != nil {
		return fmt.Errorf("error listing Routes: %v", err)
	}

	for _, route := range routes {
		output := fmt.Sprintf("Route{ Name: %s, ID: %s, ServiceName: %v, StripPath: %v}", *route.Name, *route.ID, *route.Service.Name, *route.StripPath)
		fmt.Println(output)
	}

	return nil
}

func truncateKong(kongo *kongClientWrapper.Kongo) {
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

func generateTestThings(kongo *kongClientWrapper.Kongo) error {
	k8sService := kongClientWrapper.K8sService{
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
