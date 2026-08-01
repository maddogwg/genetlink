package genetlink_test

import (
	"errors"
	"iter"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
)

// recordingSocket is a netlink.Socket that records the destination address
// passed to Send.
type recordingSocket struct {
	pid   uint32
	group uint32
	msgs  []netlink.Message
}

func (s *recordingSocket) Close() error { return nil }

func (s *recordingSocket) Send(m netlink.Message, pid uint32, group uint32) error {
	s.pid = pid
	s.group = group
	s.msgs = []netlink.Message{m}
	return nil
}

func (s *recordingSocket) SendMessages(msgs []netlink.Message, pid uint32) error {
	s.pid = pid
	s.group = 0
	s.msgs = msgs
	return nil
}

func (s *recordingSocket) Receive() ([]netlink.Message, error) {
	return nil, errors.New("not implemented")
}

func (s *recordingSocket) ReceiveIter() iter.Seq2[netlink.Message, error] {
	return func(yield func(netlink.Message, error) bool) {
		yield(netlink.Message{}, errors.New("not implemented"))
	}
}

func TestConnPID(t *testing.T) {
	const want = uint32(42)

	c := genetlink.NewConn(netlink.NewConn(&recordingSocket{}, want))
	defer c.Close()

	if diff := cmp.Diff(want, c.PID()); diff != "" {
		t.Fatalf("unexpected PID (-want +got):\n%s", diff)
	}
}

func TestConnSendDestination(t *testing.T) {
	tests := []struct {
		name      string
		send      func(c *genetlink.Conn) error
		wantPID   uint32
		wantGroup uint32
	}{
		{
			name: "Send",
			send: func(c *genetlink.Conn) error {
				_, err := c.Send(genetlink.Message{}, 1, netlink.Request)
				return err
			},
		},
		{
			name: "SendTo",
			send: func(c *genetlink.Conn) error {
				_, err := c.SendTo(genetlink.Message{}, 1, netlink.Request, 42)
				return err
			},
			wantPID: 42,
		},
		{
			name: "Multicast",
			send: func(c *genetlink.Conn) error {
				_, err := c.Multicast(genetlink.Message{}, 1, netlink.Request, 0x5)
				return err
			},
			wantGroup: 0x5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sock := &recordingSocket{}
			c := genetlink.NewConn(netlink.NewConn(sock, 1))
			defer c.Close()

			if err := tt.send(c); err != nil {
				t.Fatalf("failed to send: %v", err)
			}

			if diff := cmp.Diff(tt.wantPID, sock.pid); diff != "" {
				t.Fatalf("unexpected destination pid (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantGroup, sock.group); diff != "" {
				t.Fatalf("unexpected destination group (-want +got):\n%s", diff)
			}
			if len(sock.msgs) == 0 {
				t.Fatal("no messages recorded by socket")
			}
		})
	}
}
