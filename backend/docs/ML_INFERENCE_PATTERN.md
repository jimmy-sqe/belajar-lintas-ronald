# ML Inference Pattern

This document is the authoritative guide for implementing local ML model inference in
`backend-belajar-lintas-ronald`-based services. It defines the **InferenceEngine port**, the
preprocessor/postprocessor contract, and how to wire a new model type end-to-end.

**Used by:** `sdlc-be-go:implement-feature` reads this file (via `backend/docs/`) to
generate correct inference scaffolding when tech-spec §3.5 is present.

---

## Overview

Local inference (ONNX, libtorch, etc.) follows the same Ports & Adapters pattern used
for persistence and cache:

```
internal/domain/inference/        ← Port (interface + types)
├── engine.go                     ← Engine interface, ModelType, Input, Result
└── mock/engine_mock.go           ← Test double

internal/adapter/inference/<runtime>/  ← Adapter (implements Engine)
├── engine.go                     ← Runtime-specific session pool + Infer dispatch
└── <type_enum>_processor.go      ← Per-model preprocessor + postprocessor (one file each)
```

Domain services that need inference declare a dependency on `inference.Engine` (the
interface), never on a concrete adapter. This keeps domain logic testable without
loading real model files.

---

## Task Taxonomy

The inference engine supports three broad task families. The tech-spec §3.5 `output_format`
column determines which family a model belongs to:

| Family | `output_format` examples | Typical architectures |
|---|---|---|
| `vision-detection` | `YOLO-bbox [x,y,w,h,...]` | YOLOv8, RT-DETR |
| `vision-embedding` | `ArcFace float32[N]` | ArcFace, FaceNet |
| `classification-tabular` | `softmax float32[N]`, `sigmoid float32[1]` | XGBoost→ONNX, sklearn→ONNX |

Each family drives a different preprocessor/postprocessor pattern:
- **vision-detection / vision-embedding** — decode `Input.Data` bytes, resize, normalize, flatten.
- **classification-tabular** — read `Input.Features` directly; no image decode step.

---

## inference.Engine Port

```go
// Engine is the inference port. Thread-safe for concurrent use.
type Engine interface {
    Infer(ctx context.Context, modelType ModelType, input Input) (*Result, error)
}
```

`ModelType` is a product-defined string enum. Values come from tech-spec §3.5 Model
Registry `type_enum` column:

```go
// Defined by the product, not the boilerplate.
// Example (cosmos):
const (
    ModelTypePPE      ModelType = "ppe"
    ModelTypeIdentity ModelType = "identity"
    ModelTypePresence ModelType = "presence"
)
```

`Input` carries the raw frame/image/clip bytes plus MIME type. `Result` is a union
struct — not every field is populated for every model:

| output_format (§3.5)       | Populated Result fields      |
|---|---|
| `YOLO-bbox [x,y,w,h,...]`  | `BoundingBoxes []BoundingBox` |
| `ArcFace float32[N]`        | `Embedding []float32`         |
| Custom / mixed              | `Metadata map[string]any`     |

---

## Tabular Classifiers

Tabular models (XGBoost, scikit-learn pipelines) receive a feature vector instead of image
bytes. The `Input` type carries a `Features` field for this purpose:

```go
type Input struct {
    // vision models: populate Data + ContentType
    Data        []byte
    ContentType string   // "image/jpeg", "image/png", "video/mp4"

    // classification-tabular models: populate Features instead of Data
    Features    map[string]float64
}
```

The preprocessor for a `classification-tabular` model reads `input.Features` and encodes
it into a flat `[]float32` tensor matching the model's input column order from tech-spec §3.5:

```go
func preprocessTabularRisk(input inference.Input) ([]float32, error) {
    cols := []string{"age", "income", "credit_score"}  // column order from tech-spec §3.5
    out := make([]float32, len(cols))
    for i, col := range cols {
        v, ok := input.Features[col]
        if !ok {
            return nil, fmt.Errorf("missing feature: %s", col)
        }
        out[i] = float32(v)
    }
    return out, nil
}
```

The postprocessor reads the raw `float32` tensor output. For binary classifiers, index 1
is the positive-class probability; for multi-class, the argmax is the predicted class index.

---

## Adding a New Model Type (Step-by-Step)

When tech-spec §3.5 adds a new row to the Model Registry, follow these steps:

### 1. Add the ModelType constant

In `internal/domain/inference/engine.go`:

```go
const (
    // existing constants ...
    ModelTypeCoalSizing ModelType = "coal-sizing"  // ← new
)
```

### 2. Create a processor file

Copy `internal/adapter/inference/onnxruntime/processor.go` to
`internal/adapter/inference/onnxruntime/coal_sizing_processor.go`.

Rename the example functions and fill in the spec values from §3.5:

```go
// implement: preprocess input for coal-sizing model
// Input shape:   [1, 3, 640, 640]
// Normalization: mean=[0.485,0.456,0.406] std=[0.229,0.224,0.225]
func preprocessCoalSizing(input inference.Input) ([]float32, error) {
    // 1. Decode input.Data according to input.ContentType (JPEG/PNG/MP4 frame)
    // 2. Resize to 640×640
    // 3. Normalize each channel: pixel = (pixel/255 - mean[c]) / std[c]
    // 4. Reorder from HWC → CHW
    // 5. Flatten to []float32 of length 1*3*640*640
    return nil, nil
}

// implement: postprocess raw ONNX output for coal-sizing model
// Output format: YOLO-bbox [x,y,w,h,conf,cls,...]
func postprocessCoalSizing(raw []float32) (*inference.Result, error) {
    // 1. Parse [x,y,w,h,confidence,class_id,...] from raw tensor rows
    // 2. Filter by confidence threshold (read from system_config or pass as param)
    // 3. Apply Non-Maximum Suppression (NMS)
    // 4. Map class IDs to label strings
    // 5. Return Result{BoundingBoxes: [...]}
    return nil, nil
}
```

### 3. Register the model in engine.go

In `NewEngine`, add one session load block:

```go
sessions[inference.ModelTypeCoalSizing], err = loadSession(filepath.Join(modelDir, "coal-sizing.onnx"))
if err != nil {
    return nil, fmt.Errorf("load coal-sizing model: %w", err)
}
```

In the `preprocess` switch statement:

```go
// in preprocess()
case inference.ModelTypeCoalSizing:
    return preprocessCoalSizing(input)
```

In the `postprocess` switch statement:

```go
// in postprocess()
case inference.ModelTypeCoalSizing:
    return postprocessCoalSizing(raw)
```

---

## Pipeline Composition

Some products run multiple models in sequence (e.g., detect presence → crop face → compute
embedding). This is a **pipeline**: the output of one `Infer` call becomes the input of the
next.

Pipeline logic belongs in a dedicated file at:

```
internal/domain/<feature>/pipeline.go
```

`pipeline.go` lives in the domain layer because it encodes business rules (which models run
in which order, under what conditions). It depends only on `inference.Engine` (the port),
never on a concrete adapter.

```go
// pipeline.go — multi-model inference pipeline
func (s *PresenceService) RunPipeline(ctx context.Context, frame []byte) (*PipelineResult, error) {
    // Step 1: presence detection
    presenceResult, err := s.engine.Infer(ctx, inference.ModelTypePresence,
        inference.Input{Data: frame, ContentType: "image/jpeg"})
    if err != nil {
        return nil, fmt.Errorf("presence infer: %w", err)
    }
    if len(presenceResult.BoundingBoxes) == 0 {
        return &PipelineResult{PersonDetected: false}, nil  // short-circuit
    }

    // Step 2: identity embedding — only if person detected.
    // Crop the original frame to the first detection's box (product helper).
    box := presenceResult.BoundingBoxes[0]
    faceCrop, err := cropRegion(frame, box.X, box.Y, box.Width, box.Height)
    if err != nil {
        return nil, fmt.Errorf("crop face region: %w", err)
    }
    identityResult, err := s.engine.Infer(ctx, inference.ModelTypeIdentity,
        inference.Input{Data: faceCrop, ContentType: "image/jpeg"})
    if err != nil {
        return nil, fmt.Errorf("identity infer: %w", err)
    }
    return &PipelineResult{PersonDetected: true, Embedding: identityResult.Embedding}, nil
}
```

Key rules for `pipeline.go`:
- Call `inference.Engine` (the port), never a concrete adapter.
- Handle each step's error before passing output downstream.
- Short-circuit early when upstream output is empty (no detections → skip later steps).
- Each model in the pipeline still needs its own processor file and session registration
  (see **Adding a New Model Type** above).

---

## Preprocessor Contract

Every `preprocess<TypeEnum>` function must:

1. **Decode** `input.Data` based on `input.ContentType`:
   - `image/jpeg`, `image/png` → standard image decode
   - `video/mp4` → extract the relevant frame (product-specific logic)

2. **Resize** to the model's `input_shape` H×W (from tech-spec §3.5).

3. **Normalize** each channel using `mean` and `std` from tech-spec §3.5:
   ```
   pixel_normalized = (pixel_value / 255.0 - mean[channel]) / std[channel]
   ```

4. **Reorder** from HWC (height × width × channels) → CHW (channels × height × width).

5. **Flatten** to `[]float32` matching the model's full input tensor shape
   (e.g., `[1, 3, 640, 640]` → 1×3×640×640 = 1,228,800 float32 values).

---

## Postprocessor Contract

Every `postprocess<TypeEnum>` function must:

**For YOLO-bbox output_format:**
1. Parse raw float32 tensor rows as `[x_center, y_center, width, height, confidence, class_scores...]`
2. Filter detections below the model's confidence threshold (read from `system_config`).
3. Apply Non-Maximum Suppression (NMS) to remove duplicate detections.
4. Map class IDs to human-readable labels.
5. Return `Result{BoundingBoxes: [...]}`

**For ArcFace/embedding output_format:**
1. Read the embedding directly from the raw tensor (e.g., 512 float32 values).
2. Optionally L2-normalize the embedding vector.
3. Return `Result{Embedding: raw[:]}`

**For custom output_format:**
- Document the tensor layout as a comment at the top of the function.
- Populate `Result.Metadata` with model-specific key-value pairs.

---

## Incremental Model Wiring

Models are wired in `NewEngine` incrementally, driven by the batch that runs their linked
FR. Four states exist in `adapter/inference/onnxruntime/engine.go`:

| State | What it looks like in `engine.go` | Transition |
|---|---|---|
| `first-time active` | No block yet; model's FR is in the very first batch | Lintas generates an active `loadSession` directly → `wired` |
| `wired` | Active `loadSession` call | Re-run skips (idempotent) |
| `deferred` | Block prefixed with `// deferred <FR>:` | Batch containing FR runs → `activate` |
| `activate` | Same as `deferred`, but FR now in active batch | Lintas uncomments the block → `wired` |

`first-time active` and `activate` both converge to the `wired` end-state — they differ
only in the generation path: the first generates a fresh active block, the second uncomments
an existing deferred one. A model that was never deferred (its FR was active from the start)
appears directly as `wired` and has no `// deferred` history.

When you see a commented block like:

```go
// deferred FR-13: uncomment when identity model batch runs
// sessions[inference.ModelTypeIdentity], err = loadSession(filepath.Join(modelDir, "identity.onnx"))
// if err != nil { return nil, fmt.Errorf("load identity model: %w", err) }
```

→ **Do NOT uncomment manually.** Run `/start-feature` again when the identity model
  batch is scheduled — Lintas detects the `// deferred FR-13:` prefix and activates it.

To advance a model from `deferred` → `wired`:
1. Ensure the model's ONNX file is available in `/models/` (delivered by DS team).
2. Add or confirm a BE task for the model's FR exists in `tasks.md`.
3. Run `/start-feature` — Lintas detects the deferred block and uncomments it.
4. Verify the service starts cleanly with the new model loaded.

---

## Idempotence & Runtime Mismatch

`implement-feature` identifies adapter files by the `// ML_RUNTIME: <runtime>` header
comment on the first line of each generated file under `adapter/inference/`. If the
header of an existing file does not match the current `ML_RUNTIME` from tech-spec §3.5,
the skill stops and asks the developer to delete `adapter/inference/` manually before
re-running. This prevents silent corruption when switching runtimes.

---

## Performance Notes (from §9 / §14)

- Model sessions are loaded **once at startup** (`startup-eager`). Do not load on first
  request — this causes an unacceptable latency spike for the first caller.
- `Infer` is called concurrently. The `sessions` map in `Engine` is **read-only** after
  `NewEngine` returns, so concurrent map reads are safe in Go. However, the underlying
  `ort.DynamicAdvancedSession` object itself may not be goroutine-safe depending on the
  ORT C runtime version — serialize `Infer` calls per model type with a mutex, or allocate
  one session per goroutine (session pool), and verify under load test before go-live.
- Memory footprint: a single YOLOv8n ONNX is ~6 MB loaded; YOLOv8s–m is ~22–85 MB.
  Confirm total footprint fits container memory limits before sizing pods.

---

## ONNX Conversion Prerequisites

All models must be in ONNX format before they can be loaded by the runtime. The DS team
produces `.onnx` files; BE should verify format compatibility before wiring.

### YOLO → ONNX (vision-detection)

Use the Ultralytics `yolo export` CLI:

```bash
pip install ultralytics
yolo export model=yolov8n.pt format=onnx opset=17 simplify=True
# Output: yolov8n.onnx
```

Validate the output before committing it to model storage:
```bash
python -c "import onnx; m = onnx.load('yolov8n.onnx'); onnx.checker.check_model(m); print('OK')"
```

### scikit-learn / XGBoost → ONNX (classification-tabular)

Use `skl2onnx` to convert fitted scikit-learn pipelines or estimators:

```bash
pip install skl2onnx
```

```python
from skl2onnx import convert_sklearn
from skl2onnx.common.data_types import FloatTensorType

# clf is a fitted sklearn Pipeline or standalone estimator
initial_type = [("float_input", FloatTensorType([None, n_features]))]
onx = convert_sklearn(clf, initial_types=initial_type)
with open("model.onnx", "wb") as f:
    f.write(onx.SerializeToString())
```

After export, verify input/output node names — the processor code must reference them exactly:
```python
import onnxruntime as ort
sess = ort.InferenceSession("model.onnx")
print("inputs:", [i.name for i in sess.get_inputs()])
print("outputs:", [o.name for o in sess.get_outputs()])
```

---

## Dockerfile Dependency (Manual Step)

The `libonnxruntime.so` shared library must be installed in the Docker image at the
version matching the `yalue/onnxruntime_go` Go binding. This is **not** generated by
`implement-feature` — add it to the Dockerfile manually:

```dockerfile
RUN apt-get install -y libonnxruntime-dev=<version>
```

Check the current required version at:
https://github.com/yalue/onnxruntime_go#requirements

---

## Reference Files in This Boilerplate

| File | Purpose |
|---|---|
| `internal/domain/inference/engine.go` | Port definition — copy as-is, add ModelType constants |
| `internal/domain/inference/mock/engine_mock.go` | Test double — copy verbatim |
| `internal/adapter/inference/onnxruntime/engine.go` | OrtEngine scaffold — customize NewEngine + switch cases |
| `internal/adapter/inference/onnxruntime/processor.go` | Processor pattern — copy per model, rename functions |
