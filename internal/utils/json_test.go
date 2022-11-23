package utils

import "testing"

func TestJsonGetInt64(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  int64
	}{
		{"String", "test", 0},
		{"Int", 1234, 0},
		{"Float", float64(1234), 1234},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JsonGetInt64(tt.input); got != tt.want {
				t.Errorf("JsonGetInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}
