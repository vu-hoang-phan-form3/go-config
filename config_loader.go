package main

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ConfigLoader struct {
	funcMap     template.FuncMap
	viper       *viper.Viper
	vaultClient *api.Client
}

func NewConfigLoader(opts ...OptFunc) *ConfigLoader {
	cl := &ConfigLoader{
		viper: viper.New(),
	}
	cl.funcMap = template.FuncMap{
		"env": cl.env,
	}
	for _, opt := range opts {
		opt(cl)
	}
	return cl
}

type OptFunc func(*ConfigLoader)

func WithVaultClient(vaultClient *api.Client) OptFunc {
	return func(cl *ConfigLoader) {
		cl.vaultClient = vaultClient
		cl.funcMap["vault"] = cl.vault
	}
}

func (cl *ConfigLoader) Viper() *viper.Viper {
	return cl.viper
}

func (cl *ConfigLoader) LoadConfigFiles(fileNames ...string) error {
	for _, fileName := range fileNames {
		configType, err := getConfigType(fileName)
		if err != nil {
			return err
		}

		baseName := filepath.Base(fileName)
		tmpl, err := template.New(baseName).Funcs(cl.funcMap).ParseFiles(fileName)
		if err != nil {
			return fmt.Errorf("failed to parse config file '%s': %w", fileName, err)
		}

		buf := &bytes.Buffer{}
		if err := tmpl.Execute(buf, nil); err != nil {
			return fmt.Errorf("failed to render config file '%s': %w", fileName, err)
		}

		cl.viper.SetConfigType(configType)
		if err := cl.viper.MergeConfig(buf); err != nil {
			return fmt.Errorf("failed to merge config file '%s': %w", fileName, err)
		}
	}

	return nil
}

func (cl *ConfigLoader) Unmarshal(v interface{}) error {
	if err := cl.viper.Unmarshal(v); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return nil
}

func (cl *ConfigLoader) env(envName string, defaultVal ...string) string {
	val, isSet := os.LookupEnv(envName)
	if !isSet && len(defaultVal) > 0 {
		val = defaultVal[0]
	}
	return val
}

func (cl *ConfigLoader) vault(path string, key string, defaultVal ...string) (string, error) {
	secrets, err := cl.vaultClient.Logical().Read(path)
	if err != nil {
		return "", fmt.Errorf("failed to read secrets from vault path '%s': %w", path, err)
	}

	if secrets != nil {
		val, ok := secrets.Data[key]
		if ok {
			return fmt.Sprintf("%v", val), nil
		}
	}

	if len(defaultVal) > 0 {
		return defaultVal[0], nil
	}
	return "", fmt.Errorf("vault: key '%s' does not exist in '%s' and no default value has been provided", key, path)
}

func getConfigType(fileName string) (string, error) {
	ext := filepath.Ext(fileName)
	ext = strings.TrimPrefix(ext, ".")
	for _, supported := range viper.SupportedExts {
		if ext == supported {
			return ext, nil
		}
	}
	return "", viper.UnsupportedConfigError(ext)
}
