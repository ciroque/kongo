package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hbagdi/go-kong/kong"
	jsoniter "github.com/json-iterator/go"
	"kongo/client"
	"log"
	"os"
	"strings"
)

type Arguments struct {
	KongUri     *string
	Command     *string
	Namespace   *string
	ServiceName *string
}

func (a Arguments) String() string {
	json, _ := jsoniter.Marshal(a)
	return string(json)
}

var arguments Arguments

type Command struct {
	function    func(kongo *client.Kongo, args Arguments) error
	description string
}

func init() {
	arguments.KongUri = flag.String("kongUri", "http://localhost:8001", "Url for the Kong admin API")
	arguments.Command = flag.String("command", "usage", "Describes the usage of kongo")
	arguments.Namespace = flag.String("namespace", "", "The target namespace")
	arguments.ServiceName = flag.String("service", "", "The target service name")
}

func main() {

	flag.Parse()

	fmt.Println("arguments: ", arguments)

	kongo, _ := client.NewKongo(arguments.KongUri)

	commands := getCommands()
	command, found := commands[*arguments.Command]
	if !found {
		command = commands["usage"]
	}
	err := command.function(kongo, arguments)
	if err != nil {
		log.Fatal("Not so fast: ", err)
	}
}

func getCommands() map[string]Command {
	commands := make(map[string]Command)

	commands["clear-entries"] = Command{clearEntries, "Removes all entries identified by the given namespace and name"}
	commands["register-test-resources"] = Command{registerTestResources, "Generates test entities in Kong"}
	commands["deregister-test-resources"] = Command{deregisterTestResources, "Removes test resources from Kong"}
	commands["list"] = Command{listAllThings, "Lists all entities within Kong"}
	commands["truncate"] = Command{truncateKong, "Deletes all entities from Kong (USE WITH CAUTION)"}
	commands["usage"] = Command{printUsage, "Shows the usage of the tool and available commands"}

	return commands
}

func clearEntries(kongo *client.Kongo, args Arguments) error {
	if *arguments.Namespace == "" || *arguments.ServiceName == "" {
		return fmt.Errorf("clear-entries expects the namespace and name, these were not provided. %v", args)
	}

	baseName := fmt.Sprintf("%s.%s", *arguments.Namespace, *arguments.ServiceName)
	return kongo.DeregisterK8sService(baseName)
}

func deregisterTestResources(kongo *client.Kongo, args Arguments) error {
	k8sService := client.K8sService{
		Addresses: []*string{kong.String("localhost")},
		Name:      "kongo.test-service-one",
		Path:      "/testing-1-2-3",
		Port:      80,
	}

	err := kongo.DeregisterK8sService(k8sService.Name)
	if err != nil {
		return fmt.Errorf("None delete the things: %v", err)
	}
	return nil
}

func registerTestResources(kongo *client.Kongo, args Arguments) error {
	k8sService := client.K8sService{
		Addresses: []*string{kong.String("localhost")},
		Name:      "kongo.test-service-one",
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

func listAllThings(kongo *client.Kongo, args Arguments) error {
	upstreams, err := kongo.ListUpstreams()
	if err != nil {
		return fmt.Errorf("error listing Upstreams: %v", err)
	}
	for _, upstream := range upstreams {
		fmt.Println(jsonize(upstream))

		targets, err := kongo.ListTargets(*upstream.ID)
		if err != nil {
			return fmt.Errorf("error listing Targets for Upstream '%s': %v", *upstream.Name, err)
		}
		for _, target := range targets {
			fmt.Println(jsonize(target))
		}
	}

	services, err := kongo.ListServices()
	if err != nil {
		return fmt.Errorf("error listing Services: %v", err)
	}

	for _, service := range services {
		fmt.Println(jsonize(service))
	}

	routes, err := kongo.ListRoutes()
	if err != nil {
		return fmt.Errorf("error listing Routes: %v", err)
	}

	for _, route := range routes {
		fmt.Println(jsonize(route))
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

func printUsage(kongo *client.Kongo, args Arguments) error {
	commands := getCommands()
	fmt.Println("kongo usage:")
	fmt.Println("./kongo [command], where command is one of:")
	for k, v := range commands {
		fmt.Printf("\t- '%s': %s\n", k, v.description)
	}

	return nil
}

func truncateKong(kongo *client.Kongo, args Arguments) error {
	if askForConfirmation("THIS WILL DELETE ALL KONG ENTRIES, ARE YOU SURE?") {
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

	return nil
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [YES/no]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.TrimSpace(response)

		if response == "YES" {
			return true
		} else if response == "n" || response == "no" || response == "NO" {
			return false
		}
	}
}

func jsonize(v interface{}) string {
	json, _ := jsoniter.Marshal(v)
	return string(json)
}
