package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	dynamic "github.com/traefik/traefik/v3/pkg/config/dynamic"
)

func strToBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func (a *App) fetchTraefikRouters() ([]routerRepresentation, error) {
	url := fmt.Sprintf("%s/api/http/routers?status=enabled", a.Config.Traefik.Endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if a.Config.Traefik.BasicAuth != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Basic %s", strToBase64(a.Config.Traefik.BasicAuth)))
	}

	res, err := a.http.Do(req)
	if err != nil {
		a.ready = false
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		a.ready = false
		return nil, err
	}

	data := []routerRepresentation{}
	json.Unmarshal(body, &data)

	return data, nil
}

func (a *App) buildTraefikConfiguration(src []routerRepresentation) dynamic.Configuration {
	routers := map[string]*dynamic.Router{}
	services := map[string]*dynamic.Service{}

	servers := make([]dynamic.Server, len(a.Config.Generator.TargetServers))
	for i, server := range a.Config.Generator.TargetServers {
		servers[i].URL = server
	}

	for _, router := range src {
		if router.Status != "enabled" {
			continue
		}
		if a.Config.Traefik.RuleFilter != "" && !strings.Contains(router.Rule, a.Config.Traefik.RuleFilter) {
			continue
		}

		name := strings.Split(router.Name, "@")[0]
		routers[name] = &dynamic.Router{
			Rule:        router.Rule,
			Service:     name,
			EntryPoints: a.Config.Generator.Entrypoints,
			Middlewares: a.Config.Generator.Middlewares,
			TLS:         &dynamic.RouterTLSConfig{},
		}
		services[name] = &dynamic.Service{
			LoadBalancer: &dynamic.ServersLoadBalancer{
				PassHostHeader: &a.Config.Generator.PassHostHeader,
				Servers:        servers,
			},
		}
	}

	configuration := dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:  routers,
			Services: services,
		},
	}

	return configuration
}
