package message

type MessageServiceCacheI interface {
	GetMessageCache(input string) string
}
