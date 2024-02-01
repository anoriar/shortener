package deleteurlsprocessor

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/processors/deleteuserurlsprocessor/message"
	"github.com/anoriar/shortener/internal/shortener/services/deleteuserurls"
)

// DeleteUserURLsProcessor missing godoc.
type DeleteUserURLsProcessor struct {
	deleteUserURLsService deleteuserurls.DeleteUserURLsInterface
	logger                *zap.Logger
	msgChan               chan message.DeleteUserURLsMessage
}

// NewDeleteUserURLsProcessor missing godoc.
func NewDeleteUserURLsProcessor(deleteUserURLsService deleteuserurls.DeleteUserURLsInterface, logger *zap.Logger) *DeleteUserURLsProcessor {
	instance := &DeleteUserURLsProcessor{
		deleteUserURLsService: deleteUserURLsService,
		logger:                logger,
		msgChan:               make(chan message.DeleteUserURLsMessage, 100),
	}

	return instance
}

// GetMessageChan missing godoc.
func (p *DeleteUserURLsProcessor) GetMessageChan() chan message.DeleteUserURLsMessage {
	return p.msgChan
}

// AddMessage missing godoc.
func (p *DeleteUserURLsProcessor) AddMessage(msg message.DeleteUserURLsMessage) {
	p.msgChan <- msg
}

// Start missing godoc.
func (p *DeleteUserURLsProcessor) Start(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Delete user URLs task: cancel")
			close(p.msgChan)
			return
		case msg, ok := <-p.msgChan:
			if !ok {
				p.logger.Info("Delete user URLs task: channel is closed")
				return
			}
			p.process(ctx, msg)
		}
	}

}

func (p *DeleteUserURLsProcessor) process(ctx context.Context, msg message.DeleteUserURLsMessage) {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		p.logger.Error("Delete user URLs task: json marshal error", zap.String("error", err.Error()))
	}
	p.logger.Info("Delete user URLs task: received message:", zap.String("msg", string(msgJSON)))

	err = p.deleteUserURLsService.DeleteUserURLs(ctx, msg.UserID, msg.ShortURLs)
	if err != nil {
		p.logger.Error("Delete user URLs task: error", zap.String("error", err.Error()))
	}
}
