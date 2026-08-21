// ML_RUNTIME: onnxruntime_go
//
// Package onnxruntime implements the inference.Engine port using
// github.com/yalue/onnxruntime_go and libonnxruntime.so.
//
// # Setup
//
// The Docker image must install libonnxruntime at the version matching the
// onnxruntime_go binding:
//
//	RUN apt-get install -y libonnxruntime-dev=<version>
//
// Model files are downloaded from object storage during `docker build` and
// baked into /models/ inside the image.
//
// # Usage
//
// Products register ModelType constants and preprocessors/postprocessors
// matching their tech-spec §3.5 Model Registry. See processor.go for the
// pattern each model type must follow.
package onnxruntime

import (
	"context"
	"fmt"
	"path/filepath"

	inference "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/inference"
	// ort "github.com/yalue/onnxruntime_go"  // uncomment after go get
)

// Engine holds a pool of OrtSession objects, one per registered ModelType.
// All sessions are loaded at startup (fail-fast if any model file is missing).
type Engine struct {
	// sessions map[inference.ModelType]*ort.DynamicAdvancedSession
	// Replace the line above with the real type after go get:
	//   cd backend && go get github.com/yalue/onnxruntime_go
	sessions map[inference.ModelType]any // placeholder — replace with *ort.DynamicAdvancedSession
}

// NewEngine loads all model sessions from modelDir at startup.
// modelDir is typically "/models" inside the container.
// Products add one loadSession call per row in tech-spec §3.5 Model Registry.
//
// Example (cosmos product):
//
//	sessions[inference.ModelTypePPE], err = loadSession(filepath.Join(modelDir, "ppe.onnx"))
func NewEngine(modelDir string) (*Engine, error) {
	sessions := make(map[inference.ModelType]any)
	// ← Add one block per model from §3.5 Model Registry:
	//
	// var err error
	// sessions[inference.ModelType<TypeEnum>], err = loadSession(filepath.Join(modelDir, "<file_name>"))
	// if err != nil {
	//     return nil, fmt.Errorf("load <type_enum> model: %w", err)
	// }
	_ = filepath.Join // remove this line when the first loadSession call is added below
	return &Engine{sessions: sessions}, nil
}

func loadSession(path string) (any, error) {
	// implement: ort.NewDynamicAdvancedSession(path, inputNames, outputNames, nil)
	return nil, fmt.Errorf("loadSession not implemented for %s", path)
}

// Infer dispatches to the correct session for modelType, runs preprocessing,
// executes inference, and returns postprocessed results.
func (e *Engine) Infer(ctx context.Context, modelType inference.ModelType, input inference.Input) (*inference.Result, error) {
	_, ok := e.sessions[modelType]
	if !ok {
		return nil, fmt.Errorf("inference: unknown model type %q", modelType)
	}
	preprocessed, err := preprocess(modelType, input)
	if err != nil {
		return nil, fmt.Errorf("inference: preprocess %q: %w", modelType, err)
	}
	_ = preprocessed
	// implement: run OrtSession with preprocessed tensor, then call postprocess
	return nil, fmt.Errorf("inference: Infer not yet implemented for %q", modelType)
}

func preprocess(modelType inference.ModelType, input inference.Input) ([]float32, error) {
	switch modelType {
	// ← add one case per model (copy from processor.go pattern):
	// case inference.ModelType<TypeEnum>:
	//     return preprocess<TypeEnum>(input)
	}
	return nil, fmt.Errorf("inference: no preprocessor registered for %q", modelType)
}

// compile-time check: Engine must implement inference.Engine.
var _ inference.Engine = (*Engine)(nil)

func postprocess(modelType inference.ModelType, raw []float32) (*inference.Result, error) {
	switch modelType {
	// ← add one case per model: return postprocess<TypeEnum>(raw)
	}
	return nil, fmt.Errorf("inference: no postprocessor registered for %q", modelType)
}
