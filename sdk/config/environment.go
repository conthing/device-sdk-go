package config

import (
	"os"
	"strings"

	"github.com/pelletier/go-toml"
)

const (
	envKeyUrl = "edgex_registry"
)

type environment struct {
	env map[string]interface{}
}

func NewEnvironment() *environment {
	osEnv := os.Environ()
	e := &environment{
		env: make(map[string]interface{}, len(osEnv)),
	}
	for _, env := range osEnv {
		kv := strings.Split(env, "=")
		if len(kv) == 2 && len(kv[0]) > 0 && len(kv[1]) > 0 {
			e.env[kv[0]] = kv[1]
		}
	}
	return e
}

func (e *environment) OverrideUseRegistryFromEnvironment(useRegistry string) string {
	registry, registryOk := e.env[envKeyUrl]
	if registryOk && registry != "" {
		useRegistry = registry.(string)
	}
	return useRegistry
}

func (e *environment) OverrideFromEnvironment(tree *toml.Tree) *toml.Tree {
	for k, v := range e.env {
		k = strings.Replace(k, "_", ".", -1)
		switch {
		case tree.Has(k):
			tree.Set(k, v)
		default:

		}
	}
	return tree
}
