
package utils

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestJsonGetInt64WithString(t *testing.T) {

    n := JsonGetInt64("test")
    assert.Equal(t, int64(0), n)
}

func TestJsonGetInt64WithInt(t *testing.T) {

    n := JsonGetInt64(1234)
    assert.Equal(t, int64(0), n)
}

func TestJsonGetInt64WithFloat64(t *testing.T) {

    n := JsonGetInt64(float64(1234))
    assert.Equal(t, int64(1234), n)
}
