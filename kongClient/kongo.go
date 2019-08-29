package kongClient

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/hbagdi/go-kong/kong"
	"net/http"
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

func (kongo *Kongo) DeleteService(id string) (*kong.Service, error) {
	err := kongo.Kong.Services.Delete(kongo.context, kong.String(id))
	return nil, err
}

func (kongo *Kongo) DeleteTarget(targetDef *TargetDef) (*kong.Target, error) {
	err := kongo.Kong.Targets.Delete(kongo.context, targetDef.Upstream.Name, kong.String(targetDef.Target))
	return nil, err
}

func (kongo *Kongo) DeleteUpstream(id string) (*kong.Upstream, error) {
	err := kongo.Kong.Upstreams.Delete(kongo.context, kong.String(id))
	return nil, err
}

func (kongo *Kongo) GetVersion() (*string, error) {
	root, err := kongo.Kong.Root(nil)
	return kong.String(root["version"].(string)), err
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
