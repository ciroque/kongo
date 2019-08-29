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

func (kongo *Kongo) CreateService(serviceParams *ServiceDef) (*kong.Service, error) {
	serviceDef := kong.Service{
		ClientCertificate: nil,
		CreatedAt:         nil,
		Host:              kong.String(serviceParams.Host),
		ID:                nil,
		Name:              kong.String(serviceParams.Name),
		Path:              kong.String(serviceParams.Path),
		Port:              kong.Int(serviceParams.Port),
		Protocol:          nil,
		ReadTimeout:       nil,
		Retries:           nil,
		WriteTimeout:      nil,
		Tags:              kongo.tags,
	}
	service, err := kongo.Kong.Services.Create(kongo.context, &serviceDef)
	return service, err
}

type UpstreamDef struct {
	Name string
	// TODO: Add Healthchecks configuration
}

func (kongo *Kongo) CreateUpstream(upstreamParams *UpstreamDef) (*kong.Upstream, error) {
	upstreamDef := kong.Upstream{
		ID:                 nil,
		Name:               kong.String(upstreamParams.Name),
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
	upstream, err := kongo.Kong.Upstreams.Create(kongo.context, &upstreamDef)
	return upstream, err
}

func (kongo *Kongo) DeleteService(id string) (*kong.Upstream, error) {
	err := kongo.Kong.Services.Delete(kongo.context, kong.String(id))
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

func (kongo *Kongo) ListUpstreams() ([]*kong.Upstream, error) {
	upstreams, _, err := kongo.Kong.Upstreams.List(kongo.context, &kongo.listOptions)
	return upstreams, err
}
