package main

import (
	"fmt"
	"github.com/hbagdi/go-kong/kong"
	"kongo/client"
	"os"
)

type Command struct {
	function func(kongo *client.Kongo) error
	description string
}

func getCommands() map[string]Command {
	commands := make(map[string]Command)

	commands["register-test-resources"] = Command{registerTestResources, "Generates test entities in Kong"}
	commands["deregister-test-resources"] = Command{deregisterTestResources, "Removes test resources from Kong"}
	commands["list"] = Command{listAllThings, "Lists all entities within Kong"}
	commands["truncate"] = Command{truncateKong, "Deletes all entities from Kong (USE WITH CAUTION)"}
	commands["usage"] = Command{printUsage, "Shows the usage of the tool and available commands"}

	return commands
}

func main() {
	args :=os.Args[1:]

	baseUrl := "http://localhost:8001"
	kongo, _ := client.NewKongo(&baseUrl)

	commands := getCommands()
	command, found := commands[args[0]]
	if !found {
		command = commands["usage"]
	}
	command.function(kongo)
}

func deregisterTestResources(kongo *client.Kongo) error {
	k8sService := client.K8sService{
		Addresses: []*string{kong.String("localhost")},
		Name:      "steve-test-service-one",
		Path:      "/testing-1-2-3",
		Port:      80,
	}

	err := kongo.DeregisterK8sService(&k8sService)
	if err != nil {
		return fmt.Errorf("None delete the things: %v", err)
	}
	return nil
}

func registerTestResources(kongo *client.Kongo) error {
	k8sService := client.K8sService{
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

func listAllThings(kongo *client.Kongo) error {
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

	getServiceName := func (route *kong.Route) string {
		name := "<undefined>"
		if route.Service.Name != nil {
			name = *route.Service.Name
		}
		return name
	}

	for _, route := range routes {

		output := fmt.Sprintf("Route{ Name: %s, ID: %s, ServiceName: %v, StripPath: %v}", *route.Name, *route.ID, getServiceName(route), *route.StripPath)
		fmt.Println(output)
	}

	return nil
}

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

func printUsage(kongo *client.Kongo) error {
	commands := getCommands()
	fmt.Println("kongo usage:")
	fmt.Println("./kongo [command], where command is one of:")
	for k, v := range commands {
		fmt.Printf("\t- '%s': %s\n", k, v.description)
	}

	return nil
}

func truncateKong(kongo *client.Kongo) error {
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

	return nil
}
