package analyze

type AnalyzeServiceCacheI interface {
	GetMessageCache(input string) string
}
