package config

import (
	"time"

	"github.com/rotationalio/confire"
	"github.com/rs/zerolog"
	"go.rtnl.ai/vanity/logger"
)

// All environment variables will have this prefix unless otherwise defined in struct
// tags. For example, the conf.LogLevel environment variable will be VANITY_LOG_LEVEL
// because of this prefix and the split_words struct tag in the conf below.
const Prefix = "vanity"

// Configures the vanityd server from the environment.
type Config struct {
	Maintenance   bool                `default:"false" desc:"if true, the server will start in maintenance mode"`
	Domain        string              `default:"" desc:"specify the domain of the vanity URLs, otherwise the request host will be used"`
	DefaultBranch string              `split_words:"true" default:"main" desc:"specify the default branch to use for repository references"`
	ConfigMap     string              `split_words:"true" required:"true" desc:"location of yaml file with module configurations"`
	LogLevel      logger.LevelDecoder `split_words:"true" default:"info" desc:"specify the verbosity of logging (trace, debug, info, warn, error, fatal panic)"`
	ConsoleLog    bool                `split_words:"true" default:"false" desc:"if true logs colorized human readable output instead of json"`
	BindAddr      string              `split_words:"true" default:":3264" desc:"the ip address and port to bind the server on"`
	ReadTimeout   time.Duration       `split_words:"true" default:"20s" desc:"amount of time allowed to read request headers before server decides the request is too slow"`
	WriteTimeout  time.Duration       `split_words:"true" default:"20s" desc:"maximum amount of time before timing out a write to a response"`
	IdleTimeout   time.Duration       `split_words:"true" default:"10m" desc:"maximum amount of time to wait for the next request while keep alives are enabled"`
	processed     bool
}

func New() (conf Config, err error) {
	if err = confire.Process(Prefix, &conf); err != nil {
		return Config{}, err
	}

	if err = conf.Validate(); err != nil {
		return Config{}, err
	}

	conf.processed = true
	return conf, nil
}

// Returns true if the config has not been correctly processed from the environment.
func (c Config) IsZero() bool {
	return !c.processed
}

// Custom validations are added here, particularly validations that require one or more
// fields to be processed before the validation occurs.
func (c Config) Validate() (err error) {
	return err
}

func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}
