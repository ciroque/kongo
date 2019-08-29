package kongClient

import (
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
	t.Skip("Not a thing yet")
}
