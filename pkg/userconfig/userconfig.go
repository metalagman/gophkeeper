package userconfig

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type UserConfig struct {
	appName string
	cfgType string
	cfgDir  string
}

func New(appName, cfgType string) (*UserConfig, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	cfgDir = filepath.Join(cfgDir, appName)

	if err := ensureDir(cfgDir); err != nil {
		return nil, err
	}

	return &UserConfig{
		appName: appName,
		cfgType: cfgType,
		cfgDir:  cfgDir,
	}, nil
}

// Viper creates, loads and returns viper instance for named config file
func (t *UserConfig) Viper(name string) *viper.Viper {
	v := viper.New()
	v.SetConfigName(name)
	v.AddConfigPath(t.cfgDir)
	v.SetConfigType(t.cfgType)
	// init empty config file
	if err := v.ReadInConfig(); err != nil {
		_ = v.WriteConfigAs(filepath.Join(t.cfgDir, name+"."+t.cfgType))
	}
	return v
}

// ensureDir at path exists
func ensureDir(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	return nil
}
