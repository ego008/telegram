package config

import (
	"errors"
	"path/filepath"

	log "github.com/kirillDanshin/dlog"
	"github.com/spf13/viper"
)

var (
	Config *viper.Viper

	ErrInvalidPath = errors.New("invalid path to config file")
)

// Open just open configuration file for parsing some data in other functions
func Open(path string) (*viper.Viper, error) {
	log.Ln("Opening config on path:", path)

	dir, file := filepath.Split(path)
	if file == "" {
		return nil, ErrInvalidPath
	}

	fileExt := filepath.Ext(file)[1:]
	fileName := file[:(len(file)-len(fileExt))-1]

	log.Ln("dir:", dir)
	log.Ln("file:", file)
	log.Ln("fileName:", fileName)
	log.Ln("fileExt:", fileExt)

	cfg := viper.New()

	cfg.AddConfigPath(dir)
	cfg.SetConfigName(fileName)
	cfg.SetConfigType(fileExt)

	log.Ln("Reading", file)
	err := cfg.ReadInConfig()
	return cfg, err
}
