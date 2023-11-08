package article

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/johnwongx/webook/backend/internal/repository"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"github.com/johnwongx/webook/backend/pkg/saramax"
	"time"
)

type BatchKafkaConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.Logger
}

func NewBatchKafkaConsumer(client sarama.Client, repo repository.InteractiveRepository, l logger.Logger) *BatchKafkaConsumer {
	return &BatchKafkaConsumer{
		client: client,
		repo:   repo,
		l:      l,
	}
}

func (k *BatchKafkaConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", k.client)
	if err != nil {
		return err
	}

	go func() {
		err := cg.Consume(context.Background(),
			[]string{ReadEventTopic},
			saramax.NewBatchConsumerHandler[ReadEvent](k.Consume, k.l))
		if err != nil {
			k.l.Error("消费循环退出异常", logger.Error(err))
		}
	}()
	return nil
}

func (k *BatchKafkaConsumer) Consume(msg []*sarama.ConsumerMessage, evt []ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	bizs := make([]string, 0, len(evt))
	ids := make([]int64, 0, len(evt))
	for i := 0; i < len(evt); i++ {
		bizs = append(bizs, evt[i].Biz)
		ids = append(ids, evt[i].Aid)
	}
	err := k.repo.BatchIncrReadCnt(ctx, bizs, ids)
	if err != nil {
		k.l.Error("批量增加阅读计数失败",
			logger.Error(err))
	}
	return nil
}
