package config

import "github.com/traefik/traefik/v3/pkg/config/dynamic"

type ServerConfig struct {
	Address string `key:"address" default:":1337"`
}

type TraefikConfig struct {
	Endpoint   string   `key:"endpoint" validate:"required,http_url"`
	BasicAuth  string   `key:"basicAuth"`
	RuleFilter []string `key:"ruleFilter"`
	Timeout    int64    `key:"timeout" default:"5"`
}

type GeneratorConfig struct {
	Entrypoints    []string                 `key:"entrypoints" validate:"required"`
	TargetServers  []string                 `key:"targets" validate:"required"`
	Middlewares    []string                 `key:"middlewares"`
	PassHostHeader bool                     `key:"passHostHeader" validate:"boolean" default:"false"`
	TLS            *dynamic.RouterTLSConfig `key:"tls,omitempty"`
}

type Config struct {
	Server    ServerConfig
	Traefik   TraefikConfig
	Generator GeneratorConfig

	Environment     string `key:"env" default:"production" validate:"oneof=production development"`
	RefreshInterval int64  `key:"refresh" default:"30"`
}
