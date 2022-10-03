package api

import (
	"reflect"
	"testing"
    "github.com/stretchr/testify/assert"
	"github.com/eosswedenorg-go/haproxy/agentcheck"
)

func TestNewDebugApi(t *testing.T) {
	type args struct {
		response string
	}
	tests := []struct {
		name string
		args args
		want DebugApi
	}{
		{"Up", args{"up"}, DebugApi{response: agentcheck.NewStatusResponse(agentcheck.Up)}},
		{"Down", args{"down"}, DebugApi{response: agentcheck.NewStatusResponse("down")}},
		{"DownMessage", args{"down#some message"}, DebugApi{response: agentcheck.NewStatusMessageResponse(agentcheck.Down, "some message")}},
		{"Ready", args{"ready"}, DebugApi{response: agentcheck.NewStatusResponse(agentcheck.Ready)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDebugApi(tt.args.response); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDebugApi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDebugApi_LogInfo(t *testing.T) {

    expected := LogParams{"type", "Debug", "response", "up"}

    api := DebugApi{
        response: agentcheck.NewStatusResponse(agentcheck.Up),
    }

    assert.Equal(t, api.LogInfo(), expected)
}

func TestDebugApi_Call(t *testing.T) {

    expected := agentcheck.NewStatusMessageResponse(agentcheck.Stopped, "message")

    api := DebugApi{
        response: expected,
    }

    response, msg := api.Call()

    assert.Equal(t, response, expected)
    assert.Equal(t, msg, "")
}
