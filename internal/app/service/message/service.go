package message

import "context"

type MessageService struct {
	messageServiceCache MessageServiceCacheI
}

type MessageServiceParam struct {
	MessageServiceCache MessageServiceCacheI
}

func NewMessageService(_ context.Context, param MessageServiceParam) *MessageService {
	return &MessageService{
		messageServiceCache: param.MessageServiceCache,
	}
}
