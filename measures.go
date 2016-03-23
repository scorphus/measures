// Copyright 2015 Measures authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package measures

import (
	"encoding/json"
	"net"
	"time"
)

type Dimensions map[string]interface{}

type Client interface {
	Connect() error
	Disconnect() error
	Write(b []byte) (int, error)
}

type client struct {
	address string
	conn    net.Conn
}

func NewClient(address string) Client {
	c := client{}
	if address != "" {
		c.address = address
	}
	return &c
}

func (c *client) Connect() (err error) {
	c.conn, err = net.Dial("udp", c.address)
	return err
}

func (c *client) Disconnect() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *client) Write(b []byte) (n int, err error) {
	if c.conn == nil {
		err = c.Connect()
	}
	if err != nil {
		return 0, err
	}
	n, err = c.conn.Write(b)
	if err == nil {
		return n, nil
	}
	err = c.Connect()
	if err != nil {
		return 0, err
	}
	return c.conn.Write(b)
}

type Measures struct {
	client     Client
	clientName string
}

func New(client, address string) *Measures {
	m := Measures{clientName: client}
	if address != "" {
		m.client = NewClient(address)
	}
	return &m
}

func (m *Measures) send(d Dimensions) (n int, err error) {
	b, err := json.Marshal(d)
	if err != nil {
		return 0, err
	}
	return m.client.Write(b)
}

func (m *Measures) CleanUp() {
	m.client.Disconnect()
}

func (m *Measures) Count(metric string, counter int, dimensions Dimensions) error {
	d := make(Dimensions, len(dimensions)+3)
	d["client"] = m.clientName
	d["count"] = counter
	d["metric"] = metric
	for k, v := range dimensions {
		d[k] = v
	}
	_, err := m.send(d)
	if err != nil {
		return err
	}
	return nil
}

func (m *Measures) SetClient(client Client) error {
	m.client = client
	return nil
}

func (m *Measures) Time(metric string, startTime time.Time, dimensions Dimensions) error {
	d := make(Dimensions, len(dimensions)+3)
	d["client"] = m.clientName
	d["metric"] = metric
	d["time"] = time.Since(startTime).Seconds()
	for k, v := range dimensions {
		d[k] = v
	}
	_, err := m.send(d)
	if err != nil {
		return err
	}
	return nil
}
