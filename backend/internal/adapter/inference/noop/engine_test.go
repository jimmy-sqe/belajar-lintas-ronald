package noop

import (
	"context"
	"testing"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/inference"

	"github.com/stretchr/testify/require"
)

func TestNoopEngine_InferReturnsEmptyResult(t *testing.T) {
	eng := New()
	res, err := eng.Infer(context.Background(), inference.ModelType("ppe"), inference.Input{})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, inference.ModelType("ppe"), res.ModelType)
	require.Empty(t, res.BoundingBoxes)
	require.Empty(t, res.Embedding)
}

func TestNoopEngine_ImplementsPort(t *testing.T) {
	var _ inference.Engine = New()
}
