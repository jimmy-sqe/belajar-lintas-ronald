package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	// boilerplate:axis=auth option=jwt START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/auth/jwt"
	// boilerplate:axis=auth option=jwt END

	// boilerplate:axis=cache option=redis START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/cache/redis"
	// boilerplate:axis=cache option=redis END

	// boilerplate:axis=observability option=otel START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/observability/otel"
	// boilerplate:axis=observability option=otel END

	// boilerplate:axis=observability option=datadog START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/observability/datadog"
	// boilerplate:axis=observability option=datadog END

	// boilerplate:axis=persistence option=postgres START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/postgres"
	// boilerplate:axis=persistence option=postgres END

	// boilerplate:axis=persistence option=mysql START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/mysql"
	// boilerplate:axis=persistence option=mysql END

	// boilerplate:axis=persistence option=mongodb START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/mongodb"
	// boilerplate:axis=persistence option=mongodb END

	// boilerplate:axis=storage option=gcs START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/storage/gcs"
	// boilerplate:axis=storage option=gcs END

	// boilerplate:axis=storage option=s3 START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/storage/s3"
	// boilerplate:axis=storage option=s3 END

	// boilerplate:axis=storage option=minio START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/storage/minio"
	// boilerplate:axis=storage option=minio END

	// boilerplate:axis=storage option=local START
	localfs "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/storage/local"
	// boilerplate:axis=storage option=local END
)

// Config is the application configuration root.
type Config struct {
	App  AppConfig
	HTTP HTTPConfig
	CORS CORSConfig

	// boilerplate:axis=persistence option=postgres START
	Postgres postgres.Config
	// boilerplate:axis=persistence option=postgres END

	// boilerplate:axis=persistence option=mysql START
	MySQL mysql.Config
	// boilerplate:axis=persistence option=mysql END

	// boilerplate:axis=persistence option=mongodb START
	Mongo mongodb.Config
	// boilerplate:axis=persistence option=mongodb END

	// boilerplate:axis=cache option=redis START
	Redis redis.Config
	// boilerplate:axis=cache option=redis END

	// boilerplate:axis=auth option=jwt START
	JWT jwt.Config
	// boilerplate:axis=auth option=jwt END

	// boilerplate:axis=observability option=otel START
	OTel otel.Config
	// boilerplate:axis=observability option=otel END

	// boilerplate:axis=observability option=datadog START
	Datadog datadog.Config
	// boilerplate:axis=observability option=datadog END

	// ObservabilityBackend selects the telemetry adapter (otel|datadog|noop).
	ObservabilityBackend string

	PersistenceBackend string
	StorageBackend     string
	CacheBackend       string

	// InferenceBackend selects the ML inference adapter (noop|onnxruntime).
	InferenceBackend  string
	InferenceModelDir string

	// boilerplate:axis=storage option=gcs START
	GCS gcs.Config
	// boilerplate:axis=storage option=gcs END

	// boilerplate:axis=storage option=s3 START
	S3 s3.Config
	// boilerplate:axis=storage option=s3 END

	// boilerplate:axis=storage option=minio START
	Minio minio.Config
	// boilerplate:axis=storage option=minio END

	// boilerplate:axis=storage option=local START
	Local localfs.Config
	// boilerplate:axis=storage option=local END

	// boilerplate:axis=rpc option=grpc START
	GRPCHost string
	GRPCPort string
	// boilerplate:axis=rpc option=grpc END
}

// AppConfig contains generic application metadata.
type AppConfig struct {
	Env  string `mapstructure:"APP_ENV"`
	Name string `mapstructure:"APP_NAME"`
}

// HTTPConfig contains HTTP server settings.
type HTTPConfig struct {
	Port int `mapstructure:"HTTP_PORT"`
}

// CORSConfig contains CORS middleware settings. AllowedOrigins is the
// comma-separated env var CORS_ALLOWED_ORIGINS, split into a slice.
// Wildcards are supported per Echo middleware (e.g.
// "https://*.preview.squantumengine.com"). Empty slice means CORS is
// effectively disabled (no origins allowed).
type CORSConfig struct {
	AllowedOrigins []string
}

// Load reads env/env.<APP_ENV> via Viper plus actual env overrides.
func Load() (*Config, error) {
	env := viper.GetString("APP_ENV")
	if env == "" {
		env = "local"
	}

	v := viper.New()
	v.SetConfigName("env." + env)
	v.SetConfigType("env")
	v.AddConfigPath("./env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config: read env file: %w", err)
	}

	cfg := &Config{
		App: AppConfig{
			Env:  env,
			Name: getString(v, "APP_NAME", "backend-belajar-lintas-ronald"),
		},
		HTTP: HTTPConfig{
			Port: getInt(v, "HTTP_PORT", 8000),
		},
		CORS: CORSConfig{
			AllowedOrigins: splitAndTrim(v.GetString("CORS_ALLOWED_ORIGINS")),
		},
		// boilerplate:axis=persistence option=postgres START
		Postgres: postgres.Config{
			Host:     v.GetString("POSTGRES_HOST"),
			Port:     getInt(v, "POSTGRES_PORT", 5432),
			User:     v.GetString("POSTGRES_USER"),
			Password: v.GetString("POSTGRES_PASSWORD"),
			DB:       v.GetString("POSTGRES_DB"),
			SSLMode:  v.GetString("POSTGRES_SSLMODE"),
		},
		// boilerplate:axis=persistence option=postgres END
		// boilerplate:axis=persistence option=mysql START
		MySQL: mysql.Config{
			Host:     v.GetString("MYSQL_HOST"),
			Port:     getInt(v, "MYSQL_PORT", 3306),
			User:     v.GetString("MYSQL_USER"),
			Password: v.GetString("MYSQL_PASSWORD"),
			DB:       v.GetString("MYSQL_DB"),
		},
		// boilerplate:axis=persistence option=mysql END
		// boilerplate:axis=persistence option=mongodb START
		Mongo: mongodb.Config{
			URI: v.GetString("MONGODB_URI"),
			DB:  v.GetString("MONGODB_DB"),
		},
		// boilerplate:axis=persistence option=mongodb END
		// boilerplate:axis=cache option=redis START
		Redis: redis.Config{
			Host:     v.GetString("REDIS_HOST"),
			Port:     getInt(v, "REDIS_PORT", 6379),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
		// boilerplate:axis=cache option=redis END
		// boilerplate:axis=auth option=jwt START
		JWT: jwt.Config{
			Secret:        v.GetString("JWT_SECRET"),
			Issuer:        v.GetString("JWT_ISSUER"),
			Audience:      v.GetString("JWT_AUDIENCE"),
			AccessTTLSec:  getInt(v, "JWT_ACCESS_TTL_SEC", 3600),
			RefreshTTLSec: getInt(v, "JWT_REFRESH_TTL_SEC", 86400),
		},
		// boilerplate:axis=auth option=jwt END
		// boilerplate:axis=observability option=otel START
		OTel: otel.Config{
			Endpoint:    v.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"),
			ServiceName: getString(v, "OTEL_SERVICE_NAME", "backend-belajar-lintas-ronald"),
			Insecure:    v.GetBool("OTEL_EXPORTER_OTLP_INSECURE"),
		},
		// boilerplate:axis=observability option=otel END
		// boilerplate:axis=observability option=datadog START
		Datadog: datadog.Config{
			AgentHost:   v.GetString("DD_AGENT_HOST"),
			AgentPort:   v.GetString("DD_TRACE_AGENT_PORT"),
			ServiceName: getString(v, "DD_SERVICE", "backend-belajar-lintas-ronald"),
			Environment: v.GetString("DD_ENV"),
			Enabled:     v.GetBool("DD_TRACE_ENABLED"),
		},
		// boilerplate:axis=observability option=datadog END

		ObservabilityBackend: getString(v, "OBSERVABILITY_BACKEND", "noop"),

		PersistenceBackend: getString(v, "PERSISTENCE_BACKEND", "postgres"),
		StorageBackend:     v.GetString("STORAGE_BACKEND"),
		CacheBackend:       getString(v, "CACHE_BACKEND", "inmemory"),

		InferenceBackend:  getString(v, "INFERENCE_BACKEND", "noop"),
		InferenceModelDir: getString(v, "INFERENCE_MODEL_DIR", "/models"),

		// boilerplate:axis=storage option=gcs START
		GCS: gcs.Config{
			Bucket: v.GetString("GCS_BUCKET"),
		},
		// boilerplate:axis=storage option=gcs END

		// boilerplate:axis=storage option=s3 START
		S3: s3.Config{
			Region: v.GetString("AWS_REGION"),
			Bucket: v.GetString("S3_BUCKET"),
		},
		// boilerplate:axis=storage option=s3 END

		// boilerplate:axis=storage option=minio START
		Minio: minio.Config{
			Endpoint:  v.GetString("MINIO_ENDPOINT"),
			AccessKey: v.GetString("MINIO_ACCESS_KEY"),
			SecretKey: v.GetString("MINIO_SECRET_KEY"),
			Bucket:    v.GetString("MINIO_BUCKET"),
			UseSSL:    v.GetBool("MINIO_USE_SSL"),
		},
		// boilerplate:axis=storage option=minio END

		// boilerplate:axis=storage option=local START
		Local: localfs.Config{
			BasePath: getString(v, "LOCAL_STORAGE_PATH", "./tmp/storage"),
		},
		// boilerplate:axis=storage option=local END

		// boilerplate:axis=rpc option=grpc START
		GRPCHost: getString(v, "GRPC_HOST", "0.0.0.0"),
		GRPCPort: getString(v, "GRPC_PORT", "50051"),
		// boilerplate:axis=rpc option=grpc END
	}
	return cfg, nil
}

func getString(v *viper.Viper, key, def string) string {
	if s := v.GetString(key); s != "" {
		return s
	}
	return def
}

func getInt(v *viper.Viper, key string, def int) int {
	if i := v.GetInt(key); i != 0 {
		return i
	}
	return def
}

// splitAndTrim splits a comma-separated string and trims whitespace.
// Empty entries are filtered out. Returns nil for empty input.
func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
