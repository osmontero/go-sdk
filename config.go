package go_sdk

import (
	"path"
	"sync"
	"time"

	"github.com/threatwinds/logger"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Pipeline      []Pipeline                        `yaml:"pipeline"`
	DisabledRules []int64                           `yaml:"disabledRules,omitempty"`
	Tenants       []Tenant                          `yaml:"tenants,omitempty"`
	Patterns      map[string]string                 `yaml:"patterns,omitempty"`
	Plugins       map[string]map[string]interface{} `yaml:"plugins,omitempty"`
	Env           Env                               `yaml:"-"`
}

func (c *Config) loadCfg() {
	cFiles := ListFiles(path.Join(getEnv().Workdir, "pipeline"), ".yaml")
	for _, cFile := range cFiles {
		nCfg, e := ReadYAML[Config](cFile)
		if e != nil {
			continue
		}

		c.Pipeline = append(c.Pipeline, nCfg.Pipeline...)

		c.DisabledRules = append(c.DisabledRules, nCfg.DisabledRules...)

		c.Tenants = append(c.Tenants, nCfg.Tenants...)

		for name, pattern := range nCfg.Patterns {
			c.Patterns[name] = pattern
		}

		for name, plugin := range nCfg.Plugins {
			c.Plugins[name] = plugin
		}
	}

	c.Env = getEnv()
}

func updateCfg() {
	cfgMutex.Lock()

	tmpCfg := new(Config)
	tmpCfg.Plugins = make(map[string]map[string]interface{})
	tmpCfg.Patterns = make(map[string]string)
	tmpCfg.loadCfg()

	*cfg = *tmpCfg

	cfgMutex.Unlock()

	cfgFirst = false
}

func GetCfg() *Config {
	cfgOnce.Do(func() {
		cfg = new(Config)

		go func() {
			for {
				updateCfg()
				time.Sleep(60 * time.Second)
			}
		}()
	})

	for cfgFirst {
		time.Sleep(1 * time.Second)
	}

	cfgMutex.RLock()
	defer cfgMutex.RUnlock()

	return cfg
}

func PluginCfg[t any](name string) (*t, *logger.Error) {
	cfg := GetCfg()
	if cfg.Plugins[name] == nil {
		return nil, Logger().ErrorF("plugin %s not found", name)
	}

	tmpYaml, err := yaml.Marshal(cfg.Plugins[name])
	if err != nil {
		return nil, Logger().ErrorF("error reading plugin config: %s", err.Error())
	}

	finalCfg := new(t)

	err = yaml.Unmarshal(tmpYaml, finalCfg)
	if err != nil {
		return nil, Logger().ErrorF("error writing plugin config: %s", err.Error())
	}

	return finalCfg, nil
}

var cfg *Config
var cfgOnce sync.Once
var cfgMutex sync.RWMutex
var cfgFirst bool = true

type Tenant struct {
	Name          string  `yaml:"name"`
	Id            string  `yaml:"id"`
	Assets        []Asset `yaml:"assets"`
	DisabledRules []int64 `yaml:"disabledRules"`
}

type Asset struct {
	Name            string   `yaml:"name"`
	Hostnames       []string `yaml:"hostnames"`
	IPs             []string `yaml:"ips"`
	Confidentiality int32    `yaml:"confidentiality"`
	Availability    int32    `yaml:"availability"`
	Integrity       int32    `yaml:"integrity"`
}