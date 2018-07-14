package eventbus

import (
	"sync"
	"testing"
)

func BenchmarkNsqEventBusReceive(b *testing.B) {
	eventbus := NewNsqEventBus("localhost:4150", "localhost:4161")

	for i := 0; i < b.N; i++ {
		eventbus.Emit("testtopic", "hi, i'm a cookie")
	}

	var wg sync.WaitGroup
	wg.Add(b.N)
	b.ResetTimer()
	consumer := eventbus.On("testtopic", "test", func(payload []byte) error {
		wg.Done()
		return nil
	})
	defer consumer.Stop()

	wg.Wait()
}
