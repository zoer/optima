package main

import (
	"errors"

	"github.com/sirupsen/logrus"
)

var ErrNotFound = errors.New("Not found")

type ConfigRepo interface {
	GetParams(string, string) ([]byte, error)
}

type ConfigRepoLogger struct {
	ConfigRepo
	logger *logrus.Logger
}

// NewConfigRepoLogger allocates ConfigRepoLogger to log the ConfigRepo's calls.
func NewConfigRepoLogger(repo ConfigRepo, logger *logrus.Logger) ConfigRepo {
	return &ConfigRepoLogger{
		ConfigRepo: repo,
		logger:     logger,
	}
}

// GetParams gets config params with given keys.
func (l *ConfigRepoLogger) GetParams(config, param string) ([]byte, error) {
	data, err := l.ConfigRepo.GetParams(config, param)
	if err != nil && err != ErrNotFound {
		l.logger.WithFields(logrus.Fields{
			"config": config,
			"param":  param,
		}).WithError(err).Error("error while executing ConfigRepo.GetParams")
	}
	return data, err
}
