package loader

import (
	"catalog_service/internal/conf"
	pkgConf "catalog_service/pkg/config"
	"catalog_service/pkg/val"
	"errors"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type typeCode int

const (
	typeString typeCode = iota
	typeInt
	typeDuration

	viperEnvPrefix = "mb_cat"

	defaultTimeout            = 1 * time.Second
	defaultIdleTimeout        = 1 * time.Minute
	defaultHttpMaxHeaderBytes = 8 * 1024
	defaultConfigKey          = "config"
	defaultConfigFilePath     = "config.yaml"

	defaultLogLevel  = "warn"
	defaultLogFormat = "json"

	defaultMigrationTimeout = 10 * time.Second
	defaultMigrationVersion = 0
	defaultMigrationDir     = ""

	defaultDatabaseSSL        = "prefer"
	defaultMaxOpenConnections = 100
	defaultMaxIdleConnections = 20
	defaultConnMaxLifetime    = 5 * time.Minute
	defaultInitTimeout        = 2 * time.Second

	defaultJwtTokenLifetime        = 1 * time.Hour
	defaultJwtRefreshTokenLifetime = 24 * time.Hour

	defaultRedisvalue       = "6379"
	defaultRedisInitTimeout = 2 * time.Second
)

type viperKey struct {
	name         string
	cmdlineName  string
	defaultValue any
	usage        string
	typeCode     typeCode
}

type loader struct{}

func NewLoader() conf.Loader {
	return &loader{}
}

func (l *loader) Load() (*conf.Config, error) {
	var cfg conf.Config

	viperKeys := genViperKeys()

	setDefaultValues(viperKeys)

	err := setCmdlineArgs(viperKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cmdline args :%w", err)
	}

	loadEnv(viperKeys)

	err = loadConfigFile(viper.GetString(defaultConfigKey))
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	if err = viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from environment: %w", err)
	}

	if err = val.ValidateStruct(&cfg); err != nil {
		return nil, fmt.Errorf("failed to validate loaded config: %w", err)
	}
	return &cfg, nil
}

func genViperKeys() []viperKey {
	return []viperKey{
		{
			name:         defaultConfigKey,
			cmdlineName:  defaultConfigKey,
			defaultValue: defaultConfigFilePath,
			usage:        "Path to configuration file",
			typeCode:     typeString,
		},
		// Http options
		{
			name:         "http.host",
			cmdlineName:  "http-host",
			defaultValue: nil,
			usage:        "Host address(es) to bind the HTTP server to",
			typeCode:     typeString,
		},
		{
			name:         "http.port",
			cmdlineName:  "http-port",
			defaultValue: nil,
			usage:        "Port to run the HTTP server on",
			typeCode:     typeInt,
		},
		{
			name:         "http.max_header_bytes",
			cmdlineName:  "http-max-header-bytes",
			defaultValue: defaultHttpMaxHeaderBytes,
			usage:        "Maximum size of HTTP request headers in bytes",
			typeCode:     typeInt,
		},
		{
			name:         "http.read_timeout",
			cmdlineName:  "http-read-timeout",
			defaultValue: defaultTimeout,
			usage:        "Maximum duration for reading the entire request",
			typeCode:     typeDuration,
		},
		{
			name:         "http.write_timeout",
			cmdlineName:  "http-write-timeout",
			defaultValue: defaultTimeout,
			usage:        "Maximum duration for writing the response",
			typeCode:     typeDuration,
		},
		{
			name:         "http.idle_timeout",
			cmdlineName:  "http-idle-timeout",
			defaultValue: defaultIdleTimeout,
			usage:        "Maximum time to wait for the next request when keep-alives are enabled",
			typeCode:     typeDuration,
		},
		{
			name:         "http.read_header_timeout",
			cmdlineName:  "http-read-header-timeout",
			defaultValue: defaultTimeout,
			usage:        "Maximum time to read request headers",
			typeCode:     typeDuration,
		},
		{
			name:         "http.shutdown_timeout",
			cmdlineName:  "http-shutdown-timeout",
			defaultValue: defaultTimeout,
			usage:        "",
			typeCode:     typeDuration,
		},
		// Logger options
		{
			name:         "log.format",
			cmdlineName:  "log-format",
			defaultValue: defaultLogFormat,
			usage:        "Log format (json, text)",
			typeCode:     typeString,
		},
		{
			name:         "log.level",
			cmdlineName:  "log-level",
			defaultValue: defaultLogLevel,
			usage:        "Log level (debug infoUser warn error)",
			typeCode:     typeString,
		},
		// Migration options
		{
			name:         "migrate.version",
			cmdlineName:  "migrate-version",
			defaultValue: defaultMigrationVersion,
			usage:        "",
			typeCode:     typeInt,
		},
		{
			name:         "migrate.dir",
			cmdlineName:  "migrate-dir",
			defaultValue: defaultMigrationDir,
			usage:        "",
			typeCode:     typeString,
		},
		{
			name:         "migrate.timeout",
			cmdlineName:  "migrate-timeout",
			defaultValue: defaultMigrationTimeout,
			usage:        "",
			typeCode:     typeDuration,
		},
		// Database options
		{
			name:         "db.host",
			cmdlineName:  "db-host",
			defaultValue: nil,
			usage:        "",
			typeCode:     typeString,
		},
		{
			name:         "db.port",
			cmdlineName:  "db-port",
			defaultValue: nil,
			usage:        "",
			typeCode:     typeInt,
		},
		{
			name:         "db.user",
			cmdlineName:  "db-user",
			defaultValue: nil,
			usage:        "",
			typeCode:     typeString,
		},
		{
			name:         "db.pwd",
			cmdlineName:  "db-pwd",
			defaultValue: nil,
			usage:        "",
			typeCode:     typeString,
		},
		{
			name:         "db.database",
			cmdlineName:  "db-database",
			defaultValue: nil,
			usage:        "",
			typeCode:     typeString,
		},
		{
			name:         "db.ssl",
			cmdlineName:  "db-ssl",
			defaultValue: defaultDatabaseSSL,
			usage:        "",
			typeCode:     typeString,
		},
		{
			name:         "db.schema",
			cmdlineName:  "db-schema",
			defaultValue: nil,
			usage:        "",
			typeCode:     typeString,
		},
		{
			name:         "db.max_open",
			cmdlineName:  "db-max-open",
			defaultValue: defaultMaxOpenConnections,
			usage:        "",
			typeCode:     typeInt,
		},
		{
			name:         "db.max_idle",
			cmdlineName:  "db-max-idle",
			defaultValue: defaultMaxIdleConnections,
			usage:        "",
			typeCode:     typeInt,
		},
		{
			name:         "db.max_lifetime",
			cmdlineName:  "db-max-lifetime",
			defaultValue: defaultConnMaxLifetime,
			usage:        "",
			typeCode:     typeDuration,
		},
		{
			name:         "db.init_timeout",
			cmdlineName:  "db-init-timeout",
			defaultValue: defaultInitTimeout,
			usage:        "",
			typeCode:     typeDuration,
		},
		// JWT options
		{
			name:         "jwt.25519key",
			cmdlineName:  "jwt-25519key",
			defaultValue: nil,
			usage:        "Path to Ed25519 private key file for JWT signing",
			typeCode:     typeString,
		},
		{
			name:         "jwt.25519pub",
			cmdlineName:  "jwt-25519pub",
			defaultValue: nil,
			usage:        "Path to Ed25519 public key file for JWT verification",
			typeCode:     typeString,
		},
		{
			name:         "jwt.token_lifetime",
			cmdlineName:  "jwt-token-lifetime",
			defaultValue: defaultJwtTokenLifetime,
			usage:        "",
			typeCode:     typeDuration,
		},
		{
			name:         "jwt.refresh_token_lifetime",
			cmdlineName:  "jwt-refresh-token-lifetime",
			defaultValue: defaultJwtRefreshTokenLifetime,
			usage:        "",
			typeCode:     typeDuration,
		},
		{
			name:         "jwt.issuer",
			cmdlineName:  "jwt-issuer",
			defaultValue: nil,
			usage:        "",
			typeCode:     typeString,
		},
		// Redis options
		{
			name:         "redis.host",
			cmdlineName:  "redis-host",
			defaultValue: nil,
			usage:        "Redis host address",
			typeCode:     typeString,
		},
		{
			name:         "redis.port",
			cmdlineName:  "redis-port",
			defaultValue: defaultRedisvalue,
			usage:        "Redis port",
			typeCode:     typeInt,
		},
		{
			name:         "redis.password",
			cmdlineName:  "redis-password",
			defaultValue: "",
			usage:        "Redis password (if any)",
			typeCode:     typeString,
		},
		{
			name:         "redis.db",
			cmdlineName:  "redis-db",
			defaultValue: 0,
			usage:        "Redis database number",
			typeCode:     typeInt,
		},
		{
			name:         "redis.init_timeout",
			cmdlineName:  "redis-init-timeout",
			defaultValue: defaultRedisInitTimeout,
			usage:        "Timeout for Redis client initialization",
			typeCode:     typeDuration,
		},
	}
}

func setDefaultValues(viperKeys []viperKey) {
	for _, v := range viperKeys {
		if v.defaultValue != nil {
			viper.SetDefault(v.name, v.defaultValue)
		}
	}
}

func setCmdlineArgs(viperKeys []viperKey) error {
	for _, v := range viperKeys {
		switch v.typeCode {
		case typeInt:
			pflag.Int(v.cmdlineName, 0, v.usage)
		case typeString:
			pflag.String(v.cmdlineName, "", v.usage)
		case typeDuration:
			pflag.Duration(v.cmdlineName, 0, v.usage)
		default:
			return fmt.Errorf("%w (%T)", pkgConf.ErrUnexpectedType, v.typeCode)
		}

		if err := viper.BindPFlag(v.name, pflag.Lookup(v.cmdlineName)); err != nil {
			return err
		}
	}

	pflag.Parse()

	return viper.BindPFlags(pflag.CommandLine)
}

func loadEnv(viperKeys []viperKey) {
	viper.SetEnvPrefix(viperEnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for _, v := range viperKeys {
		_ = viper.BindEnv(v.name)
	}

	viper.AutomaticEnv()
}

func loadConfigFile(configPath string) error {
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		var configFileNotFoundErr viper.ConfigFileNotFoundError
		var pathErr *os.PathError

		if errors.As(err, &configFileNotFoundErr) || errors.As(err, &pathErr) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
