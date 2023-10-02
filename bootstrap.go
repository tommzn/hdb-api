package main

import (
	"context"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"
	core "github.com/tommzn/hdb-core"
	datasource "github.com/tommzn/hdb-message-client"
)

// Bootstrap creates a new HTTP server.
func bootstrap(conf config.Config, ctx context.Context) (*core.Minion, error) {

	secretsManager := newSecretsManager()
	if conf == nil {
		conf = loadConfig()
	}
	logger := newLogger(conf, secretsManager, ctx)
	handlerList := requestHandlers(conf, logger)
	server := newServer(conf, logger, handlerList)
	return core.NewMinion(server), nil
}

// loadConfig from config file.
func loadConfig() config.Config {

	configSource, err := config.NewS3ConfigSourceFromEnv()
	if err != nil {
		exitOnError(err)
	}

	conf, err := configSource.Load()
	if err != nil {
		exitOnError(err)
	}
	return conf
}

// newSecretsManager retruns a new secrets manager from passed config.
func newSecretsManager() secrets.SecretsManager {
	secretsManager := secrets.NewDockerecretsManager("/run/secrets/token")
	secrets.ExportToEnvironment([]string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"}, secretsManager)
	return secretsManager
}

// newLogger creates a new logger from  passed config.
func newLogger(conf config.Config, secretsMenager secrets.SecretsManager, ctx context.Context) log.Logger {
	logger := log.NewLoggerFromConfig(conf, secretsMenager)
	logger = log.WithNameSpace(logger, "hdb-api")
	return log.WithK8sContext(logger)
}

// RequestHandlers returns all active request hadnlers.
func requestHandlers(conf config.Config, logger log.Logger) []RequestHandler {
	return []RequestHandler{
		&IndoorClimateRequestHandler{
			logger:      logger,
			climateData: make(map[string]ClimateData),
			locations:   locationsFromConfig(conf, "hdb.locations"),
			datasource:  datasource.New(conf, logger),
		},
		&HealthRequestHandler{},
	}
}

// locationsFromConfig will try to extract a device id to location map from passed config.
func locationsFromConfig(conf config.Config, configKey string) map[string]string {

	locations := make(map[string]string)
	configList := conf.GetAsSliceOfMaps(configKey)
	for _, locationCfg := range configList {
		if deviceId, ok := locationCfg["id"]; ok {
			if location, ok1 := locationCfg["location"]; ok1 {
				locations[deviceId] = location
			}

		}
	}
	return locations
}
