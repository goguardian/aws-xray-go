package segment

import (
	"fmt"
	"testing"
)

func TestSetDaemonHostAndPort(t *testing.T) {
	emitter := NewEmitter()

	if daemonHost == "" && daemonPort == "" && emitter.daemonAddress != "127.0.0.1:2000" {
		t.Errorf("daemonAddress should equal 127.0.0.1:2000")
	}

	if emitter.daemonAddress != fmt.Sprintf("%s:%s", daemonHost, daemonPort) {
		t.Errorf("daemonAddress should equal %s:%s", daemonHost, daemonPort)
	}

	host := "host"
	port := "port"
	emitter.SetDaemonHostAndPort(host, port)

	if emitter.daemonAddress != fmt.Sprintf("%s:%s", host, port) {
		t.Errorf("daemonAddress should equal %s:%s", host, port)
	}

	emitter.SetDaemonHostAndPort("127.0.0.1", "2000")

	seg := New("segment", nil)

	if err := emitter.Send(seg); err != nil {
		t.Error("Emitter send should not error")
	}
}

func TestSend(t *testing.T) {
	emitter := NewEmitter()

	seg := New("segment", nil)

	// For race condition testing
	for i := 0; i < 1000; i++ {
		go func() {
			emitter.SetDaemonHostAndPort("127.0.0.1", "1000")
			emitter.Send(seg)
		}()
	}

	for i := 0; i < 1000; i++ {
		seg.AddNewSubsegment("subsegment")
	}

	if err := emitter.Send(seg); err == nil {
		t.Error("Emitter should error for size")
	}
}
