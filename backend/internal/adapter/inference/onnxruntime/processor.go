package onnxruntime

// This file shows the preprocessor/postprocessor pattern every model type must
// implement. Copy this file, rename it <type_enum>_processor.go, fill in the
// spec values from tech-spec §3.5 Model Registry, and add the matching cases
// to the preprocess/postprocess switch statements in engine.go.
//
// Preprocessor contract:
//   - Decode raw bytes from Input.Data according to Input.ContentType
//   - Resize to the model's input_shape (H×W)
//   - Normalize pixel values using the normalization spec (mean/std)
//   - Reorder from HWC → CHW layout
//   - Flatten to []float32 matching the model's expected tensor shape
//
// Postprocessor contract:
//   - Decode raw float32 tensor from OrtSession output
//   - For YOLO-bbox models: decode [x,y,w,h,conf,cls,...] → []BoundingBox + apply NMS
//   - For embedding models: return float32 slice directly as Result.Embedding
//   - Populate Result.Metadata with any model-specific extras

import "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/inference"

// preprocessExample is a template for new model preprocessors.
// Rename this function to preprocess<TypeEnum> and fill in the spec values.
//
// tech-spec §3.5 fields to use:
//
//	input_shape:   e.g. [1,3,640,640]  → resize to 640×640, 3 channels
//	normalization: e.g. mean=[0.485,0.456,0.406] std=[0.229,0.224,0.225]
func preprocessExample(input inference.Input) ([]float32, error) {
	// implement: resize to <H>×<W>, normalize channels, convert HWC→CHW
	return nil, nil
}

// postprocessExample is a template for new model postprocessors.
// Rename this function to postprocess<TypeEnum> and fill in the spec values.
//
// tech-spec §3.5 fields to use:
//
//	output_format: e.g. "YOLO-bbox [x,y,w,h,conf,cls,…]" or "ArcFace float32[512]"
func postprocessExample(raw []float32) (*inference.Result, error) {
	// implement: decode output tensors → BoundingBoxes / Embedding per output_format
	return nil, nil
}
