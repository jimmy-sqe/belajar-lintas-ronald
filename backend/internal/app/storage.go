package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/storage"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"

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

	// boilerplate:axis=storage option=noop START
	noopstorage "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/storage/noop"
	// boilerplate:axis=storage option=noop END
)

func newStorage(ctx context.Context, cfg *config.Config) (storage.Storage, func(context.Context) error, error) {
	var valid []string
	// boilerplate:axis=storage option=gcs START
	valid = append(valid, "gcs")
	// boilerplate:axis=storage option=gcs END
	// boilerplate:axis=storage option=s3 START
	valid = append(valid, "s3")
	// boilerplate:axis=storage option=s3 END
	// boilerplate:axis=storage option=minio START
	valid = append(valid, "minio")
	// boilerplate:axis=storage option=minio END
	// boilerplate:axis=storage option=local START
	valid = append(valid, "local")
	// boilerplate:axis=storage option=local END
	// boilerplate:axis=storage option=noop START
	valid = append(valid, "noop")
	// boilerplate:axis=storage option=noop END

	switch cfg.StorageBackend {
	// boilerplate:axis=storage option=gcs START
	case "gcs":
		c, err := gcs.New(ctx, cfg.GCS)
		if err != nil {
			return nil, nil, fmt.Errorf("storage: gcs: %w", err)
		}
		return c, func(_ context.Context) error { return c.Close() }, nil
	// boilerplate:axis=storage option=gcs END

	// boilerplate:axis=storage option=s3 START
	case "s3":
		c, err := s3.New(ctx, cfg.S3)
		if err != nil {
			return nil, nil, fmt.Errorf("storage: s3: %w", err)
		}
		return c, nil, nil
	// boilerplate:axis=storage option=s3 END

	// boilerplate:axis=storage option=minio START
	case "minio":
		c, err := minio.New(ctx, cfg.Minio)
		if err != nil {
			return nil, nil, fmt.Errorf("storage: minio: %w", err)
		}
		return c, nil, nil
	// boilerplate:axis=storage option=minio END

	// boilerplate:axis=storage option=local START
	case "local":
		c, err := localfs.New(cfg.Local)
		if err != nil {
			return nil, nil, fmt.Errorf("storage: local: %w", err)
		}
		return c, nil, nil
	// boilerplate:axis=storage option=local END

	// boilerplate:axis=storage option=noop START
	case "noop":
		return noopstorage.New(), nil, nil
	// boilerplate:axis=storage option=noop END

	default:
		return nil, nil, fmt.Errorf("storage: STORAGE_BACKEND must be one of: [%s]; got %q",
			strings.Join(valid, ", "), cfg.StorageBackend)
	}
}
