package message

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/go-ego/gse"
)

type CreateMessageParm struct {
	// CreateHistogramParams postgres.CreateHistogramParams
}

func (i *MessageService) GetMessageByUser(ctx context.Context, username string) error {
	fmt.Println("MessageService func GetMessageByUser()", username)
	var cache []byte
	err := sonic.Unmarshal(cache, cache)
	if err != nil {
		fmt.Println(err)
	}

	segmentor := &gse.Segmenter{
		AlphaNum: true,
	}
	fmt.Println(segmentor)
	return nil
}

func (i *MessageService) GetMessageByUser2(ctx context.Context, param CreateMessageParm) error {
	fmt.Println("MessageService func GetMessageByUser2()", param)
	return nil
}
