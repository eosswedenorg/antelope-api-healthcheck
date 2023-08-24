package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeV1(t *testing.T) {
	api, err := Make("v1", ApiArguments{})
	assert.NoError(t, err)
	assert.IsType(t, AntelopeV1{}, api)
}

func TestMakeV2(t *testing.T) {
	api, err := Make("v2", ApiArguments{})
	assert.NoError(t, err)
	assert.IsType(t, AntelopeV2{}, api)
}

func TestMakeAtomic(t *testing.T) {
	api, err := Make("atomic", ApiArguments{})
	assert.NoError(t, err)
	assert.IsType(t, AtomicAsset{}, api)
}

func TestMakeDebug(t *testing.T) {
	api, err := Make("debug", ApiArguments{})
	assert.NoError(t, err)
	assert.IsType(t, DebugApi{}, api)
}

func TestMakeInvalid(t *testing.T) {
	api, err := Make("invalid", ApiArguments{})
	assert.Error(t, err)
	assert.Nil(t, api)
}
