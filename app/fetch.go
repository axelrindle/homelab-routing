package app

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	dynamic "github.com/traefik/traefik/v3/pkg/config/dynamic"
)

func (a *App) fetchTraefikRouters() ([]routerRepresentation, error) {
	url := fmt.Sprintf("%s/api/http/routers?status=enabled", a.Config.Traefik.Endpoint)

	res, err := a.http.Get(url)
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
		}
		services[name] = &dynamic.Service{
			LoadBalancer: &dynamic.ServersLoadBalancer{
				PassHostHeader: TruePointer(),
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
