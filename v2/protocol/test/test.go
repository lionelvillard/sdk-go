// Package test provides re-usable functions for binding tests.
package test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/protocol"
)

// SendReceive does s.Send(in), then it receives the message in r.Receive() and executes outAssert
// Halt test on error.
func SendReceive(t *testing.T, ctx context.Context, in binding.Message, s protocol.Sender, r protocol.Receiver, outAssert func(binding.Message)) {
	t.Helper()
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		out, recvErr := r.Receive(ctx)
		require.NoError(t, recvErr)
		outAssert(out)
		finishErr := out.Finish(nil)
		require.NoError(t, finishErr)
	}()

	go func() {
		defer wg.Done()
		finished := false
		in = binding.WithFinish(in, func(err error) {
			require.NoError(t, err)
			finished = true
		})
		err := s.Send(ctx, in)
		require.NoError(t, err)
		require.True(t, finished)
	}()

	wg.Wait()
}
