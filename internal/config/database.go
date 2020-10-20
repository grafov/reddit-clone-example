package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// DB describes Postgres related settings.
var DB db

type db struct {
	Host      string `required:"true"`
	Port      int    `required:"true"`
	Name      string `required:"true"`
	Namespace string `default:"public"`
	Username  string `required:"true"`
	Password  string `required:"true"`
	SSLMode   string `default:"disable"`
	Debug     bool   `default:"false"`

	MigrateDown int
}

// ConnectionString merges db config data into connection string.
func (d *db) ConnectionString() string {
	s := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password='%s' sslmode=%s search_path=%s",
		d.Host, d.Port, d.Name, d.Username, d.Password, d.SSLMode, d.Namespace)
	if d.Debug {
		s = s + " log_error_verbosity=verbose"
	}
	return s
}

func init() {
	err := envconfig.Process("DB", &DB)
	log.Log("db", DB)
	if err != nil {
		panic(err)
	}
}
