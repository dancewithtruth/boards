package amqp

type amqpMock struct{}

func NewMock() Amqp {
	return &amqpMock{}
}

func (a *amqpMock) Publish(queue string, task string, v any) error {
	return nil
}
