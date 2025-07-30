package config

type ServerConfig struct {
	Address string `key:"address" default:":1337"`
}

type TraefikConfig struct {
	Endpoint   string `key:"endpoint" validate:"required,http_url"`
	RuleFilter string `key:"ruleFilter"`
	Timeout    int64  `key:"timeout" default:"5"`
}

type GeneratorConfig struct {
	Entrypoints   []string `key:"entrypoints" validate:"required"`
	TargetServers []string `key:"targets" validate:"required"`
	Middlewares   []string `key:"middlewares"`
}

type Config struct {
	Server    ServerConfig
	Traefik   TraefikConfig
	Generator GeneratorConfig

	Environment     string `key:"env" default:"production" validate:"oneof=production development"`
	RefreshInterval int64  `key:"refresh" default:"30"`
}
