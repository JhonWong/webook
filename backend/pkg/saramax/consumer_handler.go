package saramax

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/johnwongx/webook/backend/pkg/logger"
)

type ConsumerHandler[T any] struct {
	fn func(msg *sarama.ConsumerMessage, t T) error
	l  logger.Logger
}

func NewConsumerHandler[T any](fn func(msg *sarama.ConsumerMessage, t T) error, l logger.Logger) sarama.ConsumerGroupHandler {
	return &ConsumerHandler[T]{
		fn: fn,
		l:  l,
	}
}

func (c *ConsumerHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	fmt.Println("setup consumer")
	return nil
}

func (c *ConsumerHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	fmt.Println("cleanup consumer")
	return nil
}

func (c *ConsumerHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, t)
		if err != nil {
			c.l.Error("消费信息格式错误",
				logger.String("源信息", string(msg.Value)),
				logger.Error(err))
			continue
		}

		for i := 0; i < 3; i++ {
			err = c.fn(msg, t)
			if err == nil {
				break
			}
			c.l.Error("处理消息失败",
				logger.Error(err),
				logger.String("topic", msg.Topic),
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset))
		}

		if err != nil {
			c.l.Error("处理消息失败-重试上限",
				logger.Error(err),
				logger.String("topic", msg.Topic),
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset))
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}
