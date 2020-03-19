package buffering

import (
	"context"
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/types"
)

func TestWithAcksBeforeFinish(t *testing.T) {
	var testEvent = event.Event{
		DataEncoded: []byte(`"data"`),
		Context: event.EventContextV1{
			DataContentType: event.StringOfApplicationJSON(),
			Source:          types.URIRef{URL: url.URL{Path: "source"}},
			ID:              "id",
			Type:            "type"}.AsV1(),
	}

	finishCalled := false
	finishMessage := binding.WithFinish((*binding.EventMessage)(&testEvent), func(err error) {
		finishCalled = true
	})

	wg := sync.WaitGroup{}

	messageToTest := WithAcksBeforeFinish(finishMessage, 1000)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(m binding.Message) {
			ch := make(chan binding.Message, 1)
			assert.NoError(t, binding.ChanSender(ch).Send(context.Background(), m))
			<-ch
			wg.Done()
		}(messageToTest)
	}

	wg.Wait()
	assert.True(t, finishCalled)
}

func TestCopyAndWithAcksBeforeFinish(t *testing.T) {
	var testEvent = event.Event{
		DataEncoded: []byte(`"data"`),
		Context: event.EventContextV1{
			DataContentType: event.StringOfApplicationJSON(),
			Source:          types.URIRef{URL: url.URL{Path: "source"}},
			ID:              "id",
			Type:            "type"}.AsV1(),
	}

	finishCalled := false
	finishMessage := binding.WithFinish((*binding.EventMessage)(&testEvent), func(err error) {
		finishCalled = true
	})

	copiedMessage, err := BufferMessage(context.Background(), finishMessage)
	assert.NoError(t, err)

	wg := sync.WaitGroup{}

	messageToTest := WithAcksBeforeFinish(copiedMessage, 1000)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(m binding.Message) {
			ch := make(chan binding.Message, 1)
			assert.NoError(t, binding.ChanSender(ch).Send(context.Background(), m))
			<-ch
			wg.Done()
		}(messageToTest)
	}

	wg.Wait()
	assert.True(t, finishCalled)
}