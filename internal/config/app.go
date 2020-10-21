package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// These variables exports build related metadata. They set during
// compile time.
var (
	GitHash string
	Version string
	BuildAt string
)

// App describes common application settings.
//
// The library looks up for variables with prefixes as defined by
// namespace(). So the service could look up for names like
// "APP_HOST", "APP_PORT" and so on.
var App app

type app struct {
	Host           string
	Port           int      `default:"8000"`
	AllowedDomains []string `split_words:"true" default:"*"`
	TokenSecret    string   `default:"you-must-replace-this-string-in-production"`
	TokenDuration  int   `default:"10080"` // duration in minutes

	RequestTimeout  time.Duration `default:"1s"` // request serving timeout
	ShutdownTimeout time.Duration `default:"5s"` // timeout wait shutdown
	ReadTimeout     time.Duration `default:"5s"` // server read timeout
	WriteTimeout    time.Duration `default:"5s"` // server write timeout
	IdleTimeout     time.Duration `default:"5s"` // server idle timeout
}

// ListenAt ...
func (a *app) ListenAt() string {
	return fmt.Sprintf("%s:%d", App.Host, App.Port)
}

func init() {
	err := envconfig.Process("APP", &App)
	log.Log("app", App)
	if err != nil {
		panic(err)
	}
}
