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
	fmt.Println(kongo.GetVersion())

		services, _ := kongo.ListServices()
	printServices(services)

	fmt.Println("CREATING UPSTREAM")
	upstream, _ := kongo.CreateUpstream("orange-service-upstream")
	upstreamName := *upstream.Name
	fmt.Println(upstreamName)

	upstreams, _ := kongo.ListUpstreams()
	printUpstreams(upstreams)

	fmt.Println("DELETING UPSTREAM")
	kongo.DeleteUpstream(upstreamName)

	upstreams, _ = kongo.ListUpstreams()
	printUpstreams(upstreams)


	//var ctx context.Context
	//
	//routes, _, err := kongClient.Routes.List(ctx, &opts)
	//if err != nil {
	//	fmt.Errorf("Bad Juju: %s", err)
	//}
	//fmt.Println("ROUTES: ", len(routes))
	//
	//for idx, route := range routes {
	//	fmt.Printf("Route %i => %s: %s\n", idx, *route.ID, *route.Name)
	//}
	//
	//services, _, err := kongClient.Services.List(ctx, &opts)
	//if err != nil {
	//	fmt.Errorf("Bad Juju: %s", err)
	//}
	//fmt.Println("SERVICES: ", len(services))

	//orangeService := kong.Service{
	//	Host:              kong.String("orange-service-upstream"),
	//	Name:              kong.String("orange-service"),
	//	Port:              kong.Int(8080),
	//	Tags:              kong.StringSlice("app-services", "k8s-kong-federated-ingress"),
	//}
	//service, err := kongClient.Services.Create(ctx, &orangeService)
	//if err != nil {
	//	fmt.Errorf("Bad Juju: %s", err)
	//}
	//fmt.Println("CREATED SERVICE: ", service)
	//
	//orangeRoute := kong.Route{
	//	Hosts:                   kong.StringSlice("10.235.33.199", "10.235.35.197", "10.235.45.91", "10.235.41.231"),
	//	ID:                      kong.String("orange-service-ingress-route"),
	//	Name:                    kong.String("orange-service-ingress-route"),
	//	Methods:                 kong.StringSlice("GET"),
	//	Paths:                   kong.StringSlice("/orange"),
	//	PreserveHost:            kong.Bool(true),
	//	Service:				 service,
	//	StripPath:               kong.Bool(false),
	//}
	//
	//route, err := kongClient.Routes.Create(ctx, &orangeRoute)
	//
	//if err != nil {
	//	fmt.Errorf("BOOO: %s", err)
	//}
	//
	//fmt.Println("CREATED ROUTE: ", route)
}


