package generator

type MessageGenerator interface {
	GetMessageSize() int
}

type UniformMessageGenerator struct {
	MessageSize int
}

func (g UniformMessageGenerator) GetMessageSize() int {
	return g.MessageSize
}

func NewUniformGenerator(msgSize int) *UniformMessageGenerator {
	return &UniformMessageGenerator{
		MessageSize: msgSize,
	}
}
