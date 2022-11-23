package utils

import (
	"reflect"
	"testing"

	log "github.com/inconshreveable/log15"
)

func Test_ParseLogFormatter(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want log.Format
	}{
		{"Default", "", log.TerminalFormat()},
		{"LogFmt", "logfmt", log.LogfmtFormat()},
		{"Json", "json", log.JsonFormat()},
		{"JsonPretty", "json-pretty", log.JsonFormat()},
		{"Unknown", "unknown", log.TerminalFormat()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseLogFormatter(tt.arg); reflect.ValueOf(got).Pointer() != reflect.ValueOf(tt.want).Pointer() {
				t.Errorf("parseLogFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}
