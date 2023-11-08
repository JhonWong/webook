package saramax

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"time"
)

type BatchConsumerHandler[T any] struct {
	fn           func(msg []*sarama.ConsumerMessage, t []T) error
	l            logger.Logger
	batchSize    int
	timeDuration time.Duration
}

func NewBatchConsumerHandler[T any](fn func(msg []*sarama.ConsumerMessage, t []T) error, l logger.Logger) sarama.ConsumerGroupHandler {
	return &BatchConsumerHandler[T]{
		fn:           fn,
		l:            l,
		batchSize:    10,
		timeDuration: time.Second * 10,
	}
}

func (c *BatchConsumerHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	fmt.Println("setup consumer")
	return nil
}

func (c *BatchConsumerHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	fmt.Println("cleanup consumer")
	return nil
}

func (c *BatchConsumerHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgCh := claim.Messages()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), c.timeDuration)
		msgs := make([]*sarama.ConsumerMessage, 0, c.batchSize)
		ts := make([]T, 0, c.batchSize)
		closed := false
		for i := 0; i < c.batchSize && !closed; i++ {
			select {
			case <-ctx.Done():
				closed = true
			case msg, ok := <-msgCh:
				if !ok {
					cancel()
					return nil
				}
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					c.l.Error("json解析消息数据失败",
						logger.String("topic", msg.Topic),
						logger.Int32("partition", msg.Partition),
						logger.Int64("offset", msg.Offset),
						logger.String("data", string(msg.Value)),
						logger.Error(err))
					continue
				}
				msgs = append(msgs, msg)
				ts = append(ts, t)
			}
		}
		cancel()

		if len(msgs) == 0 {
			continue
		}

		var err error
		for i := 0; i < 3; i++ {
			err = c.fn(msgs, ts)
			if err == nil {
				break
			}
			c.l.Error("消息处理失败", logger.Error(err))
		}

		if err != nil {
			c.l.Error("消息处理失败上限", logger.Error(err))
		}

		for _, msg := range msgs {
			session.MarkMessage(msg, "")
		}
	}

	return nil
}
