// Package inference defines the local ML model inference port.
// Products using ONNX or other local inference runtimes implement this port
// and register their adapter in the DI container.
package inference

import "context"

// ModelType identifies which loaded model session to use for inference.
// Products derive constants from tech-spec §3.5 Model Registry type_enum column.
// Example:
//
//	const (
//	    ModelTypePPE      ModelType = "ppe"
//	    ModelTypeIdentity ModelType = "identity"
//	)
type ModelType string

// Input is the raw payload sent to an inference session.
//
// Vision models (detection/embedding) populate Data + ContentType.
// Tabular classifiers (XGBoost/sklearn → ONNX) populate Features instead —
// see docs/ML_INFERENCE_PATTERN.md "Tabular Classifiers".
type Input struct {
	// Data is the raw bytes of the frame, image, or clip to infer on.
	Data []byte
	// ContentType is the MIME type of Data (e.g. "image/jpeg", "video/mp4").
	ContentType string
	// Features is the named feature vector for classification-tabular models.
	// Empty for vision models. The preprocessor maps these into the model's
	// input column order (from tech-spec §3.5).
	Features map[string]float64
}

// BoundingBox is a single detected region returned by YOLO-style models.
type BoundingBox struct {
	X, Y, Width, Height float32
	Confidence          float32
	ClassID             int
	Label               string
}

// Result is the structured output from a single inference call.
// Not all fields are populated for every model type:
//   - BoundingBoxes: populated by detection/segmentation models (YOLO-bbox output)
//   - Embedding: populated by embedding/recognition models (ArcFace float32[N])
//   - Metadata: arbitrary key-value for model-specific extras (distances, sizes, etc.)
type Result struct {
	ModelType     ModelType
	BoundingBoxes []BoundingBox
	Embedding     []float32
	Metadata      map[string]any
}

// Engine is the inference port. Implementations load model sessions at startup
// and dispatch Infer calls to the correct session based on ModelType.
// The Engine must be safe for concurrent use across goroutines.
type Engine interface {
	Infer(ctx context.Context, modelType ModelType, input Input) (*Result, error)
}
