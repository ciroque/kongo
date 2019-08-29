package kongClient

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/hbagdi/go-kong/kong"
	"net/http"
	"strings"
)

type Kongo struct {
	Kong        *kong.Client
	context     context.Context
	listOptions kong.ListOpt
	tags        []*string
}

func NewKongo(baseUrl *string) (*Kongo, error) {
	headers := []string{"Content-Type: application/json", "Accept: application/json"}

	kongo := new(Kongo)

	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = true

	defaultTransport := http.DefaultTransport.(*http.Transport)

	defaultTransport.TLSClientConfig = &tlsConfig

	httpClient := http.DefaultClient
	httpClient.Transport = &KongoRoundTripper{
		headers:      headers,
		roundTripper: defaultTransport,
	}
	kongClient, err := kong.NewClient(baseUrl, httpClient)
	if err != nil {
		fmt.Errorf("Something went horribly horribly wrong: %s", err)
	}

	kongo.Kong = kongClient

	return kongo, nil
}

type RouteDef struct {
	Name string
	Paths []*string
	Service *kong.Service
	StripPath bool
}

func (kongo *Kongo) CreateRoute(routeDef *RouteDef) (*kong.Route, error) {
	kongRoute := kong.Route{
		CreatedAt:               nil,
		Hosts:                   nil,
		Headers:                 nil,
		ID:                      nil,
		Name:                    kong.String(routeDef.Name),
		Methods:                 nil,
		Paths:                   routeDef.Paths,
		PreserveHost:            nil,
		Protocols:               nil,
		RegexPriority:           nil,
		Service:                 routeDef.Service,
		StripPath:               kong.Bool(routeDef.StripPath),
		UpdatedAt:               nil,
		SNIs:                    nil,
		Sources:                 nil,
		Destinations:            nil,
		Tags:                    nil,
		HTTPSRedirectStatusCode: nil,
	}
	return kongo.Kong.Routes.Create(kongo.context, &kongRoute)
}

type ServiceDef struct {
	Name     string
	Host     string
	Path     string
	Port     int
	Protocol string `default:"GET"`
}

func (kongo *Kongo) CreateService(serviceDef *ServiceDef) (*kong.Service, error) {
	kongService := kong.Service{
		ClientCertificate: nil,
		CreatedAt:         nil,
		Host:              kong.String(serviceDef.Host),
		ID:                nil,
		Name:              kong.String(serviceDef.Name),
		Path:              kong.String(serviceDef.Path),
		Port:              kong.Int(serviceDef.Port),
		Protocol:          nil,
		ReadTimeout:       nil,
		Retries:           nil,
		WriteTimeout:      nil,
		Tags:              kongo.tags,
	}
	return kongo.Kong.Services.Create(kongo.context, &kongService)
}

type TargetDef struct {
	Target string
	Upstream *kong.Upstream
	Weight int
}

func (kongo *Kongo) CreateTarget(targetDef *TargetDef) (*kong.Target, error) {
	kongTarget := kong.Target{
		Target:    kong.String(targetDef.Target),
		Upstream:  targetDef.Upstream,
		Weight:    kong.Int(targetDef.Weight),
		Tags:      kongo.tags,
	}
	return kongo.Kong.Targets.Create(kongo.context, targetDef.Upstream.Name, &kongTarget)
}

type UpstreamDef struct {
	Name string
	// TODO: Add Healthchecks configuration
}

func (kongo *Kongo) CreateUpstream(upstreamDef *UpstreamDef) (*kong.Upstream, error) {
	kongUpstream := kong.Upstream{
		ID:                 nil,
		Name:               kong.String(upstreamDef.Name),
		Algorithm:          nil,
		Slots:              nil,
		Healthchecks:       nil,
		CreatedAt:          nil,
		HashOn:             nil,
		HashFallback:       nil,
		HashOnHeader:       nil,
		HashFallbackHeader: nil,
		HashOnCookie:       nil,
		HashOnCookiePath:   nil,
		Tags:               kongo.tags,
	}
	return kongo.Kong.Upstreams.Create(kongo.context, &kongUpstream)
}

func (kongo *Kongo) DeleteRoute(idOrName string) (*kong.Route, error) {
	return nil, kongo.Kong.Routes.Delete(kongo.context, kong.String(idOrName))
}

func (kongo *Kongo) DeleteService(idOrName string) (*kong.Service, error) {
	return nil, kongo.Kong.Services.Delete(kongo.context, kong.String(idOrName))
}

func (kongo *Kongo) DeleteTarget(targetDef *TargetDef) (*kong.Target, error) {
	return nil, kongo.Kong.Targets.Delete(kongo.context, targetDef.Upstream.Name, kong.String(targetDef.Target))
}

func (kongo *Kongo) DeleteUpstream(idOrName string) (*kong.Upstream, error) {
	err := kongo.Kong.Upstreams.Delete(kongo.context, kong.String(idOrName))
	return nil, err
}

func (kongo *Kongo) GetVersion() (*string, error) {
	root, err := kongo.Kong.Root(nil)
	return kong.String(root["version"].(string)), err
}

func (kongo *Kongo) ListRoutes() ([]*kong.Route, error) {
	services, _, err := kongo.Kong.Routes.List(kongo.context, &kongo.listOptions)
	return services, err
}

func (kongo *Kongo) ListServices() ([]*kong.Service, error) {
	services, _, err := kongo.Kong.Services.List(kongo.context, &kongo.listOptions)
	return services, err
}

func (kongo *Kongo) ListTargets(upstreamId string) ([]*kong.Target, error) {
	targets, _, err := kongo.Kong.Targets.List(kongo.context, kong.String(upstreamId), &kongo.listOptions)
	return targets, err
}

func (kongo *Kongo) ListUpstreams() ([]*kong.Upstream, error) {
	upstreams, _, err := kongo.Kong.Upstreams.List(kongo.context, &kongo.listOptions)
	return upstreams, err
}

type K8sService struct {
	Addresses []*string
	Name string
	Path string
	Port int
}

type RegisteredK8sService struct {
	Service *kong.Service
	Targets []*kong.Target
	Route *kong.Route
	Upstream *kong.Upstream
}

func (kongo *Kongo) RegisterK8sService(k8sService *K8sService) (*RegisteredK8sService, error) {
	// 1 - Create Upstream
	upstreamName := strings.Join([]string{k8sService.Name, "upstream"}, "-")
	upstreamDef := UpstreamDef{Name: upstreamName}
	kongUpstream, err := kongo.CreateUpstream(&upstreamDef)
	if err != nil {
		return nil, err
	}

	// retval
	var registeredK8sService RegisteredK8sService
	registeredK8sService.Upstream = kongUpstream

	// 2 - Create Target(s)
	targets := []*kong.Target{}
	for _, target := range k8sService.Addresses {
		targetDef := TargetDef{
			Target:   *target,
			Upstream: kongUpstream,
			Weight:   0,
		}
		kongTarget, err := kongo.CreateTarget(&targetDef)
		if err != nil {
			return &registeredK8sService, err
		}
		targets = append(targets, kongTarget)
	}

	registeredK8sService.Targets = targets

	// 3 - Create Service
	serviceName := strings.Join([]string{k8sService.Name, "service"}, "-")
	serviceDef := ServiceDef{
		Name:     serviceName,
		Host:     upstreamName,
		Path:     k8sService.Path,
		Port:     k8sService.Port,
	}
	kongService, err := kongo.CreateService(&serviceDef)
	if err != nil {
		return &registeredK8sService, err
	}

	registeredK8sService.Service = kongService

	// 4 - Create Route
	routeName := strings.Join([]string{k8sService.Name, "service"}, "-")
	routeDef := RouteDef{
		Name:      routeName,
		Paths:     nil,
		Service:   kongService,
		StripPath: false,
	}
	kongRoute, err := kongo.CreateRoute(&routeDef)
	if err != nil {
		return &registeredK8sService, err
	}

	registeredK8sService.Route = kongRoute

	return &registeredK8sService, nil
}
