package kafka

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmailProducer_SendAsync_EnqueuesMessage(t *testing.T) {
	p := &EmailProducer{
		writer: nil,                          
		queue:  make(chan queuedMessage, 10), 
	}

	msg := EmailMessage{
		Email:   "test@example.com",
		Subject: "Hello",
		Code:    "1234",
	}

	p.SendAsync(msg)

	assert.Equal(t, 1, len(p.queue))

	qm := <-p.queue
	assert.Equal(t, msg.Email, qm.msg.Email)
	assert.Equal(t, msg.Subject, qm.msg.Subject)
	assert.Equal(t, msg.Code, qm.msg.Code)
	assert.NotNil(t, qm.ctx) 
}

func TestEmailProducer_SendAsync_DoesNotBlockWhenQueueFull(t *testing.T) {
	p := &EmailProducer{
		writer: nil,
		queue:  make(chan queuedMessage, 1), 
	}

	p.SendAsync(EmailMessage{Email: "a@test.com"})
	assert.Equal(t, 1, len(p.queue))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		p.SendAsync(EmailMessage{Email: "b@test.com"})
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("SendAsync заблокировался при полной очереди")
	}

	assert.Equal(t, 1, len(p.queue))
}
