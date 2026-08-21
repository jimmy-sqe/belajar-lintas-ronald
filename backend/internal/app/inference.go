package app

import (
	"fmt"
	"strings"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/inference"

	// boilerplate:axis=inference option=noop START
	noopinference "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/inference/noop"
	// boilerplate:axis=inference option=noop END

	// boilerplate:axis=inference option=onnxruntime START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/inference/onnxruntime"
	// boilerplate:axis=inference option=onnxruntime END
)

// newInference builds the inference.Engine for the selected backend.
// Mirrors newStorage: per-option marker-wrapped, unknown value fails fast.
func newInference(cfg *config.Config) (inference.Engine, error) {
	var valid []string
	// boilerplate:axis=inference option=noop START
	valid = append(valid, "noop")
	// boilerplate:axis=inference option=noop END
	// boilerplate:axis=inference option=onnxruntime START
	valid = append(valid, "onnxruntime")
	// boilerplate:axis=inference option=onnxruntime END

	switch cfg.InferenceBackend {
	// boilerplate:axis=inference option=noop START
	case "noop":
		return noopinference.New(), nil
	// boilerplate:axis=inference option=noop END

	// boilerplate:axis=inference option=onnxruntime START
	case "onnxruntime":
		eng, err := onnxruntime.NewEngine(cfg.InferenceModelDir)
		if err != nil {
			return nil, fmt.Errorf("inference: onnxruntime: %w", err)
		}
		return eng, nil
	// boilerplate:axis=inference option=onnxruntime END

	default:
		return nil, fmt.Errorf("inference: INFERENCE_BACKEND must be one of: [%s]; got %q",
			strings.Join(valid, ", "), cfg.InferenceBackend)
	}
}
