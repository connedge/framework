package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"strings"
	"sync"
	"time"
)

var mu = &sync.Mutex{}

type Option struct {
	FilePath   string
	ConfigType string
	ConfigName string
}

type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetInt64(key string) int64
	GetFloat(key string) float64
	GetBool(key string) bool
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetStringArray(key string) []string
	GetIntArray(key string) []int
	Get(key string) interface{}
	GetStringMap(key string) map[string]interface{}
	GetCustomConfigMap(config interface{}) interface{}
}

type config struct {
	configReader *viper.Viper
	filePath     string
	configType   string
	configName   string
}

func New(options ...func(*Option)) Config {

	opts := &Option{}
	for _, o := range options {
		o(opts)
	}

	configReader := viper.New()
	configReader.AddConfigPath(opts.FilePath)
	configReader.SetConfigName(opts.ConfigName)
	configReader.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if opts.ConfigType != "" {
		configReader.SetConfigType(opts.ConfigType)
	} else {
		configReader.SetConfigType("json")
	}

	if err := configReader.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("config read error: %s", err.Error()))
	}

	configReader.WatchConfig()
	configReader.OnConfigChange(func(in fsnotify.Event) {
		mu.Lock()
		fmt.Println("config watch: updated")
		mu.Unlock()
	})

	return &config{
		configReader: configReader,
		filePath:     opts.FilePath,
		configType:   opts.ConfigType,
		configName:   opts.ConfigName,
	}
}

func WithFilePath(path, defaultPath string) func(*Option) {
	return func(o *Option) {
		if o.FilePath == "" {
			o.FilePath = defaultPath
			return
		}
		o.FilePath = path
	}
}

func WithConfigType(configType, defaultConfigType string) func(*Option) {
	return func(o *Option) {
		if o.ConfigType == "" {
			o.ConfigType = defaultConfigType
			return
		}
		o.ConfigType = configType
	}
}

func WithConfigName(cfgName string, defaultName string) func(*Option) {
	return func(o *Option) {

		if o.ConfigName == "" {
			o.ConfigName = defaultName
			return
		}

		o.ConfigName = cfgName
	}
}

func (c *config) GetInt(key string) int {
	return c.configReader.GetInt(key)
}

func (c *config) GetString(key string) string {
	return c.configReader.GetString(key)
}

func (c *config) GetInt64(key string) int64 {
	return c.configReader.GetInt64(key)
}

func (c *config) GetStringArray(key string) []string {
	return c.configReader.GetStringSlice(key)
}

func (c *config) GetIntArray(key string) []int {
	return c.configReader.GetIntSlice(key)
}

func (c *config) GetBool(key string) bool {
	return c.configReader.GetBool(key)
}

func (c *config) GetFloat(key string) float64 {
	return c.configReader.GetFloat64(key)
}

func (c *config) GetTime(key string) time.Time {
	return c.configReader.GetTime(key)
}

func (c *config) GetDuration(key string) time.Duration {
	return c.configReader.GetDuration(key)
}

func (c *config) Get(key string) interface{} {
	return c.configReader.Get(key)
}

func (c *config) GetStringMap(key string) map[string]interface{} {
	return c.configReader.GetStringMap(key)
}

func (c *config) GetCustomConfigMap(config interface{}) interface{} {
	_ = c.configReader.Unmarshal(config)
	return config
}
