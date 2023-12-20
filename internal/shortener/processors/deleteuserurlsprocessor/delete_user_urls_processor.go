package deleteurlsprocessor

import (
	"context"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/processors/deleteuserurlsprocessor/message"
	"github.com/anoriar/shortener/internal/shortener/services/deleteuserurls"
)

type DeleteUserURLsProcessor struct {
	deleteUserURLsService deleteuserurls.DeleteUserURLsInterface
	logger                *zap.Logger
	msgChan               chan message.DeleteUserURLsMessage
}

func NewDeleteUserURLsProcessor(deleteUserURLsService deleteuserurls.DeleteUserURLsInterface, logger *zap.Logger) *DeleteUserURLsProcessor {
	instance := &DeleteUserURLsProcessor{
		deleteUserURLsService: deleteUserURLsService,
		logger:                logger,
		msgChan:               make(chan message.DeleteUserURLsMessage, 100),
	}
	go instance.process()

	return instance
}

func (p *DeleteUserURLsProcessor) AddMessage(msg message.DeleteUserURLsMessage) {
	p.msgChan <- msg
}

func (p *DeleteUserURLsProcessor) process() {
	ctx := context.Background()

	for msg := range p.msgChan {
		err := p.deleteUserURLsService.DeleteUserURLs(ctx, msg.UserID, msg.ShortURLs)
		if err != nil {
			p.logger.Error("delete user urls error", zap.String("error", err.Error()))
		}
		continue
	}
}
