// Copyright 2015 Measures authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package measures_test

import (
	"net"
	"testing"
	"time"

	"github.com/scorphus/measures"
	. "gopkg.in/check.v1"
)

type stubClient struct {
	output string
}

func newStubClient() measures.Client {
	return &stubClient{}
}

func (c *stubClient) Connect() error    { return nil }
func (c *stubClient) Disconnect() error { return nil }

func (c *stubClient) Write(b []byte) (n int, err error) {
	s := string(b)
	c.output += s
	return len(s), nil
}

type S struct {
	conn       *net.UDPConn
	client     measures.Client
	stubClient stubClient
	measures   *measures.Measures
}

var _ = Suite(&S{})

func Test(t *testing.T) { TestingT(t) }

func (s *S) listenUDP() error {
	serverAddr, err := net.ResolveUDPAddr("udp", ":3593")
	if err != nil {
		return err
	}
	s.conn, err = net.ListenUDP("udp", serverAddr)
	if err != nil {
		return err
	}
	return nil
}

func (s *S) readFromUDP() string {
	if s.conn == nil {
		return ""
	}
	buf := make([]byte, 1024)
	n, _, err := s.conn.ReadFromUDP(buf)
	if err != nil {
		return ""
	}
	return string(buf[0:n])
}

func (s *S) SetUpTest(c *C) {
	err := s.listenUDP()
	c.Check(err, IsNil)
	s.client = measures.NewClient("0.0.0.0:3593")
	s.measures = measures.New("tests", "")
	s.stubClient = stubClient{}
	s.measures.SetClient(&s.stubClient)
}

func (s *S) TearDownTest(c *C) {
	if s.client != nil {
		s.client.Disconnect()
	}
	if s.measures != nil {
		s.measures.CleanUp()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *S) TestClient(c *C) {
	err := s.client.Connect()
	c.Check(err, IsNil)
	message := "Measure it!"
	n, err := s.client.Write([]byte(message))
	c.Check(err, IsNil)
	c.Check(n, Equals, len(message))
	c.Check(s.readFromUDP(), Equals, message)
}

func (s *S) TestClientAutoConnect(c *C) {
	message := "You shall auto-connect!"
	n, err := s.client.Write([]byte(message))
	c.Check(err, IsNil)
	c.Check(n, Equals, len(message))
	c.Check(s.readFromUDP(), Equals, message)
}

func (s *S) TestClientReconnect(c *C) {
	message := "You shall auto-connect!!!"
	n, err := s.client.Write([]byte(message))
	c.Check(err, IsNil)
	c.Check(n, Equals, len(message))
	c.Check(s.readFromUDP(), Equals, message)
	c.Check(err, IsNil)
	err = s.client.Disconnect()
	message = "You shall reconnect!"
	n, err = s.client.Write([]byte(message))
	c.Check(err, IsNil)
	c.Check(n, Equals, len(message))
	c.Check(s.readFromUDP(), Equals, message)
}

func (s *S) TestClientAcceptMultipleWrites(c *C) {
	message := "Accept it!"
	n, err := s.client.Write([]byte(message))
	c.Check(err, IsNil)
	c.Check(n, Equals, len(message))
	c.Check(s.readFromUDP(), Equals, message)
	message = "Accept it again!"
	n, err = s.client.Write([]byte(message))
	c.Check(err, IsNil)
	c.Check(n, Equals, len(message))
	c.Check(s.readFromUDP(), Equals, message)
	message = "And again!"
	n, err = s.client.Write([]byte(message))
	c.Check(err, IsNil)
	c.Check(n, Equals, len(message))
	c.Check(s.readFromUDP(), Equals, message)
}

func (s *S) TestCount(c *C) {
	sizes := measures.Dimensions{"XL": 20, "L": 10, "M": 5}
	err := s.measures.Count("sizes", 3, sizes)
	c.Check(err, IsNil)
	c.Check(s.stubClient.output, Equals, `{"L":10,"M":5,"XL":20,"client":"tests","count":3,"metric":"sizes"}`)
}

func (s *S) TestCountDeep(c *C) {
	book := measures.Dimensions{
		"title":      measures.Dimensions{"title": "Blood Meridian", "subtitle": "The Evening Redness in the west"},
		"author":     measures.Dimensions{"first": "Cormac", "last": "McCarthy"},
		"characters": []string{"The Kid", "The Judge"},
	}
	err := s.measures.Count("books", 1, book)
	c.Check(err, IsNil)
	c.Check(s.stubClient.output, Equals, `{"author":{"first":"Cormac","last":"McCarthy"},"characters":["The Kid","The Judge"],"client":"tests","count":1,"metric":"books","title":{"subtitle":"The Evening Redness in the west","title":"Blood Meridian"}}`)
}

func (s *S) TestCountDeeper(c *C) {
	author := measures.Dimensions{
		"Cormac McCarthy": measures.Dimensions{
			"best_books": []measures.Dimensions{{
				"book1": measures.Dimensions{
					"title": "Blood Meridian",
					"pages": 337,
					"reads": 1,
				},
				"book2": measures.Dimensions{
					"title": "No Country for Old Men",
					"pages": 309,
					"reads": nil,
				},
			}},
		},
	}
	err := s.measures.Count("authors", 1, author)
	c.Check(err, IsNil)
	c.Check(s.stubClient.output, Equals, `{"Cormac McCarthy":{"best_books":[{"book1":{"pages":337,"reads":1,"title":"Blood Meridian"},"book2":{"pages":309,"reads":null,"title":"No Country for Old Men"}}]},"client":"tests","count":1,"metric":"authors"}`)
}

func (s *S) TestCountMulti(c *C) {
	sizes := measures.Dimensions{"XL": 20, "L": 10, "M": 5}
	err := s.measures.Count("sizes", 3, sizes)
	c.Check(err, IsNil)
	book := measures.Dimensions{"title": "Blood Meridian", "pages": 337, "stars": 4.19}
	err = s.measures.Count("books", 1, book)
	c.Check(err, IsNil)
	author := measures.Dimensions{"name": "Cormac McCarthy", "magnum_opus": "Blood Meridian", "prize": "Pulitzer 2007"}
	err = s.measures.Count("authors", 1, author)
	c.Check(err, IsNil)
	c.Check(s.stubClient.output, Equals, `{"L":10,"M":5,"XL":20,"client":"tests","count":3,"metric":"sizes"}{"client":"tests","count":1,"metric":"books","pages":337,"stars":4.19,"title":"Blood Meridian"}{"client":"tests","count":1,"magnum_opus":"Blood Meridian","metric":"authors","name":"Cormac McCarthy","prize":"Pulitzer 2007"}`)
}

func (s *S) TestCountWithNoClientErrsNicely(c *C) {
	m := measures.New("tests", "")
	sizes := measures.Dimensions{"XL": 20, "L": 10, "M": 5}
	err := m.Count("sizes", 3, sizes)
	c.Assert(err.Error(), Equals, "no client set")
}

func (s *S) TestTime(c *C) {
	done := make(chan bool)
	go func() {
		sizes := measures.Dimensions{"XL": 20, "L": 10, "M": 5}
		defer s.measures.Time("sizes", time.Now(), sizes)
		done <- true
	}()
	<-done
	c.Check(s.stubClient.output, Matches, `\{"L":10,"M":5,"XL":20,"client":"tests","metric":"sizes","time":[0-9\.e\-]+\}`)
}

func (s *S) TestTimeDeep(c *C) {
	done := make(chan bool)
	go func() {
		book := measures.Dimensions{
			"title":      measures.Dimensions{"title": "Blood Meridian", "subtitle": "The Evening Redness in the west"},
			"author":     measures.Dimensions{"first": "Cormac", "last": "McCarthy"},
			"characters": []string{"The Kid", "The Judge"},
		}
		defer s.measures.Time("books", time.Now(), book)
		done <- true
	}()
	<-done
	c.Check(s.stubClient.output, Matches, `\{"author":\{"first":"Cormac","last":"McCarthy"\},"characters":\["The Kid","The Judge"\],"client":"tests","metric":"books","time":[0-9\.e\-]+,"title":\{"subtitle":"The Evening Redness in the west","title":"Blood Meridian"\}\}`)
}

func (s *S) TestTimeDeeper(c *C) {
	done := make(chan bool)
	go func() {
		author := measures.Dimensions{
			"Cormac McCarthy": measures.Dimensions{
				"best_books": []measures.Dimensions{{
					"book1": measures.Dimensions{
						"title": "Blood Meridian",
						"pages": 337,
						"reads": 1,
					},
					"book2": measures.Dimensions{
						"title": "No Country for Old Men",
						"pages": 309,
						"reads": nil,
					},
				}},
			},
		}
		defer s.measures.Time("authors", time.Now(), author)
		done <- true
	}()
	<-done
	c.Check(s.stubClient.output, Matches, `\{"Cormac McCarthy":\{"best_books":\[\{"book1":\{"pages":337,"reads":1,"title":"Blood Meridian"\},"book2":\{"pages":309,"reads":null,"title":"No Country for Old Men"\}\}\]\},"client":"tests","metric":"authors","time":[0-9\.e\-]+\}`)
}

func (s *S) TestTimeAfterDefer(c *C) {
	done := make(chan bool)
	go func() {
		sizes := make(measures.Dimensions, 3)
		defer s.measures.Time("sizes", time.Now(), sizes)
		sizes["XL"] = 20
		sizes["L"] = 10
		sizes["M"] = 5
		done <- true
	}()
	<-done
	c.Check(s.stubClient.output, Matches, `\{"L":10,"M":5,"XL":20,"client":"tests","metric":"sizes","time":[0-9\.e\-]+\}`)
}
