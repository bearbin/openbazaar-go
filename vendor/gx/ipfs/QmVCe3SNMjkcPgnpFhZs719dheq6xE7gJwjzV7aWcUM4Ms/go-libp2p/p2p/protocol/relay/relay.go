package relay

import (
	"fmt"
	"io"
	"time"

	logging "gx/ipfs/QmNQynaz7qfriSUJkiEZUrm2Wen1u3Kj9goZzWtrPyu7XR/go-log"
	peer "gx/ipfs/QmRBqJF7hb8ZSpRcMwUt8hNhydWcxGEhtk81HKq6oUwKvs/go-libp2p-peer"
	host "gx/ipfs/QmVCe3SNMjkcPgnpFhZs719dheq6xE7gJwjzV7aWcUM4Ms/go-libp2p/p2p/host"
	inet "gx/ipfs/QmVCe3SNMjkcPgnpFhZs719dheq6xE7gJwjzV7aWcUM4Ms/go-libp2p/p2p/net"
	protocol "gx/ipfs/QmVCe3SNMjkcPgnpFhZs719dheq6xE7gJwjzV7aWcUM4Ms/go-libp2p/p2p/protocol"
	mh "gx/ipfs/QmYf7ng2hG5XBtJA3tN34DQ2GUN5HNksEw1rLDkmr6vGku/go-multihash"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
)

var log = logging.Logger("github.com/ipfs/go-libp2p/p2p/protocol/relay")

// ID is the protocol.ID of the Relay Service.
const ID protocol.ID = "/ipfs/relay/line/0.1.0"

// Relay is a structure that implements ProtocolRelay.
// It is a simple relay service which forwards traffic
// between two directly connected peers.
//
// the protocol is very simple:
//
//   /ipfs/relay\n
//   <multihash src id>
//   <multihash dst id>
//   <data stream>
//
type RelayService struct {
	host    host.Host
	handler inet.StreamHandler // for streams sent to us locally.
}

func NewRelayService(h host.Host, sh inet.StreamHandler) *RelayService {
	s := &RelayService{
		host:    h,
		handler: sh,
	}
	h.SetStreamHandler(ID, s.requestHandler)
	return s
}

// requestHandler is the function called by clients
func (rs *RelayService) requestHandler(s inet.Stream) {
	if err := rs.handleStream(s); err != nil {
		log.Debugf("RelayService error:", err)
	}
}

// handleStream is our own handler, which returns an error for simplicity.
func (rs *RelayService) handleStream(s inet.Stream) error {
	defer s.Close()

	// read the header (src and dst peer.IDs)
	src, dst, err := ReadHeader(s)
	if err != nil {
		return fmt.Errorf("stream with bad header: %s", err)
	}

	local := rs.host.ID()

	switch {
	case src == local:
		return fmt.Errorf("relaying from self")
	case dst == local: // it's for us! yaaay.
		log.Debugf("%s consuming stream from %s", local, src)
		return rs.consumeStream(s)
	default: // src and dst are not local. relay it.
		log.Debugf("%s relaying stream %s <--> %s", local, src, dst)
		return rs.pipeStream(src, dst, s)
	}
}

// consumeStream connects streams directed to the local peer
// to our handler, with the header now stripped (read).
func (rs *RelayService) consumeStream(s inet.Stream) error {
	rs.handler(s) // boom.
	return nil
}

// pipeStream relays over a stream to a remote peer. It's like `cat`
func (rs *RelayService) pipeStream(src, dst peer.ID, s inet.Stream) error {
	// TODO: find a good way to pass contexts into here
	nsctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancel()

	s2, err := rs.openStreamToPeer(nsctx, dst)
	if err != nil {
		return fmt.Errorf("failed to open stream to peer: %s -- %s", dst, err)
	}
	cancel() // cancel here because this function might last a while

	if err := WriteHeader(s2, src, dst); err != nil {
		return err
	}

	// connect the series of tubes.
	done := make(chan retio, 2)
	go func() {
		n, err := io.Copy(s2, s)
		done <- retio{n, err}
	}()
	go func() {
		n, err := io.Copy(s, s2)
		done <- retio{n, err}
	}()

	r1 := <-done
	r2 := <-done
	log.Infof("%s relayed %d/%d bytes between %s and %s", rs.host.ID(), r1.n, r2.n, src, dst)

	if r1.err != nil {
		return r1.err
	}
	return r2.err
}

// openStreamToPeer opens a pipe to a remote endpoint
// for now, can only open streams to directly connected peers.
// maybe we can do some routing later on.
func (rs *RelayService) openStreamToPeer(ctx context.Context, p peer.ID) (inet.Stream, error) {
	return rs.host.NewStream(ctx, ID, p)
}

func ReadHeader(r io.Reader) (src, dst peer.ID, err error) {

	mhr := mh.NewReader(r)

	s, err := mhr.ReadMultihash()
	if err != nil {
		return "", "", err
	}

	d, err := mhr.ReadMultihash()
	if err != nil {
		return "", "", err
	}

	return peer.ID(s), peer.ID(d), nil
}

func WriteHeader(w io.Writer, src, dst peer.ID) error {
	// write header to w.
	mhw := mh.NewWriter(w)
	if err := mhw.WriteMultihash(mh.Multihash(src)); err != nil {
		return fmt.Errorf("failed to write relay header: %s -- %s", dst, err)
	}
	if err := mhw.WriteMultihash(mh.Multihash(dst)); err != nil {
		return fmt.Errorf("failed to write relay header: %s -- %s", dst, err)
	}

	return nil
}

type retio struct {
	n   int64
	err error
}
