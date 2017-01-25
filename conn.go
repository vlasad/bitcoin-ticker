package main

import (
	"time"

	"golang.org/x/net/websocket"
)

type Conn struct {
	OnMessage   func([]byte, *Conn)
	OnError     func(error)
	OnConnected func(*Conn)
	Reconnect   bool
	ws          *websocket.Conn
	url         string
	protocol    string
	closed      bool
}

type Message []byte

func (c *Conn) Dial(url, protocol string) {
	c.closed = true
	c.url = url
	c.protocol = protocol
	var err error
	c.ws, err = websocket.Dial(url, protocol, "http://localhost/")
	if err != nil {
		if c.OnError != nil {
			c.OnError(err)
		}
		c.close()
		return
	}

	c.closed = false
	if c.OnConnected != nil {
		go c.OnConnected(c)
	}

	go func() {
		defer c.close()

		for {
			var msg = make([]byte, 8192)
			var n int
			if n, err = c.ws.Read(msg); err != nil {
				if c.OnError != nil {
					c.OnError(err)
				}
				return
			}
			if c.OnMessage != nil {
				go c.OnMessage(msg[:n], c)
			}
		}
	}()
}

func (c *Conn) Send(msg Message) {
	if c.closed {
		return
	}
	if _, err := c.ws.Write(msg); err != nil {
		if c.OnError != nil {
			c.OnError(err)
		}
		c.close()
	}
}

func (c *Conn) close() {
	if !c.closed {
		c.ws.Close()
		c.closed = true
	}
	if c.Reconnect {
		for {
			c.Dial(c.url, c.protocol)
			if !c.closed {
				break
			}
			time.Sleep(time.Second * 1)
		}
	}
}
