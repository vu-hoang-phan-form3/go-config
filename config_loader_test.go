package main

import (
	"encoding/json"
	"testing"
)

type Config struct {
	Server   ServerCfg
	Log      LogCfg
	Database DatabaseCfg
}

type ServerCfg struct {
	Host string
	Port string
}

type LogCfg struct {
	Level  string
	Format string
}

type DatabaseCfg struct {
	Postgres PostgresCfg
	MongoDB  MongoDBCfg
}

type PostgresCfg struct {
	ConnectionString string
}

type MongoDBCfg struct {
	URL      string
	Username string
	Password string
}

func TestConfig(t *testing.T) {
	cl := NewConfigLoader()
	err := cl.LoadConfigFiles(
		"tests/default.yaml",
		"tests/overrides/1/override.yaml",
		"tests/overrides/2/override.yaml",
	)
	if err != nil {
		t.Fatal(err)
	}
	cfg := &Config{}
	err = cl.Unmarshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", cfg)

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s", string(b))
}
