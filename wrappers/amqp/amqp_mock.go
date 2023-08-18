package amqp

type amqpMock struct{}

func NewMock() Amqp {
	return &amqpMock{}
}

func (a *amqpMock) Declare(queue string, msgTTL int, dlx bool) error {
	return nil
}

func (a *amqpMock) Publish(queue string, task string, v any) error {
	return nil
}

func (a *amqpMock) Consume(queue string) error {
	return nil
}

func (a *amqpMock) AddHandler(task string, handler func(payload []byte) error) {
}
