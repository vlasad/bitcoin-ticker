package main

import (
	"testing"
	"time"
)

func TestConn_Dial(t *testing.T) {
	type args struct {
		url      string
		protocol string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"ws-normal",
			args{"ws://echo.websocket.org", ""},
		},
		{
			"ws-tls",
			args{"wss://echo.websocket.org", ""},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := &Conn{
				OnError: func(err error) {
					t.Errorf("Conn.Dial() error = %v", err)
				},
			}
			c.Dial(tt.args.url, tt.args.protocol)
		})
	}
}

func TestConn_Send(t *testing.T) {
	type fields struct {
		OnConnected func(*Conn)
		OnMessage   func([]byte, *Conn)
		OnError     func(error)
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"send-1",
			fields{
				OnConnected: func(con *Conn) {
					m := Message("Hello")
					con.Send(m)
				},
				OnMessage: func(msg []byte, con *Conn) {
					if string(msg) != "Hello" {
						t.Errorf("OnMessage() expected = 'Hello', got = '%s'", msg)
					}
				},
				OnError: func(err error) {
					t.Errorf("Error: %v", err)
				},
			},
			args{
				"ws://echo.websocket.org",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := &Conn{
				OnConnected: tt.fields.OnConnected,
				OnMessage:   tt.fields.OnMessage,
				OnError:     tt.fields.OnError,
			}
			c.Dial(tt.args.url, "")

			// Wait for response
			time.Sleep(time.Second * 2)
		})
	}
}

func TestConn_Close(t *testing.T) {
	type args struct {
		url      string
		protocol string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"ws-normal",
			args{"ws://echo.websocket.org", ""},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := &Conn{}
			c.Dial(tt.args.url, tt.args.protocol)

			if c.closed {
				t.Error("Conn must be open")
			}
			c.close()
			if !c.closed {
				t.Error("Conn must be closed")
			}
		})
	}
}
