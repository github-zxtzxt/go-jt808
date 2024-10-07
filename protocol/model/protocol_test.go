package model

import (
	"github.com/cuteLittleDevil/go-jt808/shared/consts"
	"testing"
)

func TestReplyProtocol(t *testing.T) {
	type Handler interface {
		Protocol() uint16
		ReplyProtocol() uint16
	}
	tests := []struct {
		name              string
		args              Handler
		wantProtocol      uint16
		wantReplyProtocol consts.PlatformReplyRequest
	}{
		{
			name:              "T0X0001 终端-通用应答",
			args:              &T0x0001{},
			wantProtocol:      uint16(consts.T0001GeneralRespond),
			wantReplyProtocol: consts.P8001GeneralRespond,
		},
		{
			name:              "P0X8001 平台-通用应答",
			args:              &P0x8001{},
			wantProtocol:      uint16(consts.P8001GeneralRespond),
			wantReplyProtocol: 0,
		},
		{
			name:              "P0X8100 终端-注册消息应答",
			args:              &P0x8100{},
			wantProtocol:      uint16(consts.P8100RegisterRespond),
			wantReplyProtocol: 0,
		},
		{
			name:              "T0X0002 终端-心跳",
			args:              &T0x0002{},
			wantProtocol:      uint16(consts.T0002HeartBeat),
			wantReplyProtocol: consts.P8001GeneralRespond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.Protocol() != tt.wantProtocol {
				t.Errorf("Protocol() = %v, want %v", tt.args.Protocol(), tt.wantProtocol)
			}
			if tt.args.ReplyProtocol() != uint16(tt.wantReplyProtocol) {
				t.Errorf("ReplyProtocol() = %v, want %v", tt.args.ReplyProtocol(), uint16(tt.wantReplyProtocol))
			}
		})
	}
}
