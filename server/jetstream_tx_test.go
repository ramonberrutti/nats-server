package server

import (
	"testing"

	"github.com/nats-io/nats.go"
)

func TestJetStreamTransaction(t *testing.T) {
	s := RunBasicJetStreamServer(t)
	defer s.Shutdown()

	nc, js := jsClientConnect(t, s)
	defer nc.Close()

	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "S",
		Subjects: []string{"foo.>"},
	})
	require_NoError(t, err)

	pub, err := js.Publish("foo.bar", []byte("ok"))
	require_NoError(t, err)

	m := nats.NewMsg("foo.bar")
	m.Data = []byte("no ok")
	m.Header.Add(JSTransactionId, "tx-1")

	pub2, err := js.PublishMsg(m, nats.ExpectLastSequencePerSubject(pub.Sequence))
	require_NoError(t, err)

	m2 := nats.NewMsg("foo.bar")
	m2.Data = []byte("no ok 2")
	m2.Header.Add(JSTransactionId, "tx-1")
	pub3, err := js.PublishMsg(m2, nats.ExpectLastSequencePerSubject(pub2.Sequence))
	require_NoError(t, err)

	m3 := nats.NewMsg("foo.bar")
	m3.Data = []byte("no ok 3")
	m3.Header.Add(JSTransactionId, "tx-1")
	m3.Header.Add(JSTransactionCommit, "true")
	js.PublishMsg(m3, nats.ExpectLastSequencePerSubject(pub3.Sequence))

	t.Fail()
}
