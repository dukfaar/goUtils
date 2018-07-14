package eventbus

type EventHandler func(payload []byte) (interface{}, error)

type EventBus interface {
	Emit(topic string, payload interface{}) error
	On(topic string, channel string, handler EventHandler) error
}

type NsqEventBus struct {
}

func NewNsqEventBus() *NsqEventBus {
	return &NsqEventBus{}
}

func (b *NsqEventBus) Emit(topic string, payload interface{}) error {
	return nil
}

func (b *NsqEventBus) On(topic string, channel string, handler EventHandler) error {
	return nil
}
