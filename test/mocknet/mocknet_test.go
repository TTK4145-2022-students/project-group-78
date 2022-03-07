package mocknet_test

import (
	"Network-go/network/bcast"
	"testing"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/test/mocknet"
	"github.com/stretchr/testify/assert"
)

func TestMocknet(t *testing.T) {
	mocknet := mocknet.New(1234, 5678)
	defer mocknet.Close()

	N := 100

	t1 := make(chan string, N)
	r1, r2 := make(chan string, N), make(chan string, N)

	go bcast.Transmitter(1234, t1)
	go bcast.Receiver(1234, r1)
	go bcast.Receiver(5678, r2)

	t.Run("Without packet loss", func(t *testing.T) {
		c1, c2 := transmit(N, t1, r1, r2)
		assert.Equal(t, N, c1)
		assert.Equal(t, N, c2)
	})

	mocknet.LossPercentage = 50
	t.Run("With packet loss", func(t *testing.T) {
		c1, c2 := transmit(N, t1, r1, r2)
		assert.True(t, 15 <= c1 && c1 <= 35)
		assert.True(t, 15 <= c2 && c2 <= 35)
	})

	mocknet.LossPercentage = 100
	t.Run("With 100 % packet loss", func(t *testing.T) {
		c1, c2 := transmit(N, t1, r1, r2)
		assert.Equal(t, 0, c1)
		assert.Equal(t, 0, c2)
	})

	mocknet.LossPercentage = 0
	mocknet.Disconnect(5678)
	t.Run("With one disconnect", func(t *testing.T) {
		c1, c2 := transmit(N, t1, r1, r2)
		assert.Equal(t, N, c1)
		assert.Equal(t, 0, c2)
	})

	mocknet.Connect(5678)
	mocknet.Disconnect(1234)
	t.Run("With one other disconnect", func(t *testing.T) {
		c1, c2 := transmit(N, t1, r1, r2)
		assert.Equal(t, 0, c1)
		assert.Equal(t, 0, c2)
	})
}

func transmit(N int, t1 chan string, r1 chan string, r2 chan string) (c1 int, c2 int) {
	for i := 0; i < N; i++ {
		t1 <- "hello"
	}
	for {
		select {
		case <-r1:
			c1++

		case <-r2:
			c2++

		case <-time.After(100 * time.Millisecond):
			return
		}
	}
}
