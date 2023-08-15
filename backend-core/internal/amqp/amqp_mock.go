package amqp

type amqpMock struct{}

func NewMock() Amqp {
	return &amqpMock{}
}

func (a *amqpMock) Publish(task string, v any) error {
	return nil
}
