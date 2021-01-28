// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libp2p

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/ethersphere/bee/pkg/p2p"
	"github.com/libp2p/go-libp2p-core/network"
)

var (
	closeDeadline  = 30 * time.Second
	errExpectedEof = errors.New("read: expected eof")
)
var _ p2p.Stream = (*stream)(nil)

type stream struct {
	network.Stream
	headers map[string][]byte
}

func NewStream(s network.Stream) p2p.Stream {
	return &stream{Stream: s}
}

func newStream(s network.Stream) *stream {
	return &stream{Stream: s}
}
func (s *stream) Headers() p2p.Headers {
	return s.headers
}

func (s *stream) FullClose() error {
	fmt.Println("close write000")

	if err := s.CloseWrite(); err != nil {
		_ = s.Reset()
		return err
	}
	fmt.Println("close write")

	// So we don't wait forever
	_ = s.SetDeadline(time.Now().Add(closeDeadline))
	fmt.Println("close write11")

	// We *have* to observe the EOF. Otherwise, we leak the stream.
	// Now, technically, we should do this *before*
	// returning from SendMessage as the message
	// hasn't really been sent yet until we see the
	// EOF but we don't actually *know* what
	// protocol the other side is speaking.
	n, err := s.Read([]byte{0})
	if n > 0 || err == nil {
		fmt.Println("read not nil")
		_ = s.Reset()
		return errExpectedEof
	}
	if err != io.EOF {
		fmt.Println("read not nil", err)

		_ = s.Reset()
		return err
	}
	return nil

}
