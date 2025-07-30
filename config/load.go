package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

func init() {
	config.AddDriver(yaml.Driver)

	config.WithOptions(config.ParseEnv, config.ParseTime, config.ParseDefault)

	config.WithOptions(func(opt *config.Options) {
		opt.DecoderConfig.TagName = "key"
	})
}

func (c *Config) Load() {
	err := config.LoadFiles("config.yml")

	err = config.Decode(c)
	if err != nil {
		log.Fatal(err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(c)
	if err != nil {
		log.Fatal(err)
	}
}
