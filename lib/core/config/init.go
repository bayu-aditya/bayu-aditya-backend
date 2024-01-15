package config

import (
	"fmt"
	"github.com/bayu-aditya/bayu-aditya-backend/lib/core/util"
	"gopkg.in/yaml.v3"
	"os"
)

func New(path string) (*Config, error) {
	var (
		prefix = "config"
		cfg    Config
	)

	cfgBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, util.ErrWrap(prefix, err, fmt.Sprintf("reading config from %s", path))
	}

	if err = yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		return nil, util.ErrWrap(prefix, err, "unmarshall yaml")
	}

	return &cfg, nil
}
