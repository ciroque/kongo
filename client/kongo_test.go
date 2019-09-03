package client

import (
	"fmt"
	"github.com/hbagdi/go-kong/kong"
	"testing"
)


func TestUpstreams(t *testing.T) {
	upstreamName := "kongo-test-upstream"
	baseUrl := "http://localhost:8001"
	kongo, _ := NewKongo(&baseUrl)

	kongo.DeleteUpstream(upstreamName)

	upstreams, _ := kongo.ListUpstreams()
	startUpstreamCount := len(upstreams)

	upstreamDef := UpstreamDef{Name: upstreamName}

	upstream, err := kongo.CreateUpstream(&upstreamDef)
	if err != nil {
		t.Fatalf("Creation of Upstream failed: %s", err)
	}

	if upstream == nil {
		t.Fatalf("Created Upstream was nil (!?)")
	}

	upstreams, _ = kongo.ListUpstreams()
	nextUpstreamCount := len(upstreams)

	if nextUpstreamCount-startUpstreamCount != 1 {
		t.Fatalf("There should be one more Upstream than at the start of the test start: %v, ended: %v", startUpstreamCount, nextUpstreamCount)
	}

	_, err = kongo.DeleteUpstream(upstreamName)
	if err != nil {
		t.Fatalf("Deletion of Upstream failed: %s", err)
	}

	upstreams, _ = kongo.ListUpstreams()
	nextUpstreamCount = len(upstreams)
	if nextUpstreamCount != startUpstreamCount {
		t.Fatalf("There should be the same number of Upstreams as at the start of the test start: %v, ended: %v", startUpstreamCount, nextUpstreamCount)
	}
}

func TestServices(t *testing.T) {
	serviceName := "kongo-test-service"
	serviceHost := "kongo-test-service-host"
	baseUrl := "http://localhost:8001"
	kongo, _ := NewKongo(&baseUrl)

	kongo.DeleteService(serviceName)

	services, _ := kongo.ListServices()
	startServiceCount := len(services)

	serviceDef := ServiceDef{
		Name: serviceName,
		Host: serviceHost,
		Path: "/rootPath",
		Port: 8080,
	}

	upstream, err := kongo.CreateService(&serviceDef)
	if err != nil {
		t.Fatalf("Creation of Service failed: %s", err)
	}

	if upstream == nil {
		t.Fatalf("Created Service was nil (!?)")
	}

	services, _ = kongo.ListServices()
	nextServiceCount := len(services)

	if nextServiceCount-startServiceCount != 1 {
		t.Fatalf("There should be one more Services than at the start of the test start: %v, ended: %v", startServiceCount, nextServiceCount)
	}

	_, err = kongo.DeleteService(serviceName)
	if err != nil {
		t.Fatalf("Deletion of Upstream failed: %s", err)
	}

	services, _ = kongo.ListServices()
	nextServiceCount = len(services)
	if nextServiceCount != startServiceCount {
		t.Fatalf("There should be the same number of Services as at the start of the test start: %v, ended: %v", startServiceCount, nextServiceCount)
	}
}

func TestTargets(t *testing.T) {
	baseUrl := "http://localhost:8001"
	kongo, _ := NewKongo(&baseUrl)

	upstreamDef := UpstreamDef{Name: "kongo-test-target-upstream"}

	kongo.DeleteUpstream(upstreamDef.Name)

	upstream, err := kongo.CreateUpstream(&upstreamDef)
	if err != nil {
		t.Fatalf("Error creating Upstream for Target: %v", err)
	}

	targetDef := TargetDef{
		Target:   "kongo-test-target-1:80",
		Upstream: upstream,
		Weight:   10,
	}

	targets, err := kongo.ListTargets(upstreamDef.Name)
	startTargetCount := len(targets)

	target, err := kongo.CreateTarget(&targetDef)
	if err != nil {
		t.Fatalf("Error creating Target: %v", err)
	}

	if target == nil {
		t.Fatal("The Target was not created")
	}

	targets, err = kongo.ListTargets(upstreamDef.Name)
	newTargetCount := len(targets)

	if newTargetCount - startTargetCount != 1 {
		t.Fatalf("Target should have been created")
	}

	_, err = kongo.DeleteTarget(&targetDef)
	if err != nil {
		t.Fatalf("Failed to remove created Target: %v", err)
	}

	targets, err = kongo.ListTargets(upstreamDef.Name)
	newTargetCount = len(targets)
	if newTargetCount != startTargetCount {
		t.Fatalf("Target should have been deleted")
	}

	_, err = kongo.DeleteUpstream(upstreamDef.Name)
	if err != nil {
		t.Fatalf("Failed to remove created Upstream: %v", err)
	}
}

func TestRoutes(t *testing.T) {
	baseUrl := "http://localhost:8001"
	kongo, _ := NewKongo(&baseUrl)

	routeName := "kongo-routes-test-route"

	serviceDef := ServiceDef{
		Name:     "kongo-routes-test-service",
		Host:     "localhost",
		Path:     "/orange",
		Port:     8080,
		Protocol: "HTTP",
	}

	_, err := kongo.DeleteRoute(routeName)
	_, err = kongo.DeleteService(serviceDef.Name)
	if err != nil {
		fmt.Println("Failed to delete previously existing Service for Route: ", err)
	}

	service, err := kongo.CreateService(&serviceDef)
	if err != nil {
		t.Fatalf("Failed to create Service for Route: %v", err)
	}

	routeDef := RouteDef{
		Name:      routeName,
		Paths:     kong.StringSlice("/orange", "/orange-whip"),
		Service:   service,
		StripPath: false,
	}

	routes, err := kongo.ListRoutes()
	startRouteCount := len(routes)

	_, err = kongo.CreateRoute(&routeDef)
	if err != nil {
		t.Fatalf("Failed to create Route: %v", err)
	}

	routes, err = kongo.ListRoutes()
	nextRouteCount := len(routes)

	if nextRouteCount - startRouteCount != 1 {
		t.Fatalf("A Route should have been created.")
	}

	_, err = kongo.DeleteRoute(routeDef.Name)
	if err != nil {
		t.Fatalf("Failed to delete Route: %v", err)
	}

	routes, err = kongo.ListRoutes()
	nextRouteCount = len(routes)

	if nextRouteCount != startRouteCount {
		t.Fatalf("A route should have been deleted.")
	}
}
