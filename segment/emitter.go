package segment

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/goguardian/aws-xray-go/utils"
	"net"
	"sync"
)

const (
	maxBodySize = 64 * 1024 // 64KB
)

var (
	protocolHeader    = []byte(`{"format": "json", "version": 1}`)
	protocolDelimiter = []byte("\n")
	daemonHost        = utils.GetenvOrDefault("XRAY_DAEMON_HOST", "127.0.0.1")
	daemonPort        = utils.GetenvOrDefault("XRAY_DAEMON_PORT", "2000")
)

// Emitter represents an emitter.  The purpose of this struct is to provide a
// cleaner interface within the "segment" package, allowing for an Emitter
// instance (e.g. "emitter") to execute Send with "emitter.Send(seg)"
type Emitter struct {
	daemonAddress string
	udpConn       net.Conn

	sync.RWMutex
}

// NewEmitter creates a new emitter with default daemon address settings.
func NewEmitter() *Emitter {
	return &Emitter{daemonAddress: fmt.Sprintf("%s:%s", daemonHost, daemonPort)}
}

func (e *Emitter) getConnection() (net.Conn, error) {
	e.Lock()
	defer e.Unlock()

	if e.udpConn != nil {
		return e.udpConn, nil
	}

	var err error
	e.udpConn, err = net.Dial("udp", e.daemonAddress)

	return e.udpConn, err
}

// SetDaemonHostAndPort updates the daemon address.
func (e *Emitter) SetDaemonHostAndPort(host string, port string) {
	e.Lock()
	defer e.Unlock()

	e.daemonAddress = fmt.Sprintf("%s:%s", host, port)

	if e.udpConn != nil {
		e.udpConn.Close()
		e.udpConn = nil
	}
}

// Send sends the segment packet to the daemon.
func (e *Emitter) Send(segment *Segment) error {
	body, err := segment.Bytes()
	if err != nil {
		return fmt.Errorf("error encoding segment: %s", err.Error())
	}

	if len(body) > maxBodySize {
		return errors.New("segment too large. >64KB")
	}

	conn, err := e.getConnection()
	if err != nil {
		return fmt.Errorf("error dialing UDP: %s", err.Error())
	}

	var buf bytes.Buffer
	buf.Write(protocolHeader)
	buf.Write(protocolDelimiter)
	buf.Write(body)

	_, err = conn.Write(buf.Bytes())
	return err
}
