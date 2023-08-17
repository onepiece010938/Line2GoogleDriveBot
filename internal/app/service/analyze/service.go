package analyze

import "context"

type AnalyzeService struct {
	analyzeServiceCache AnalyzeServiceCacheI
}

type AnalyzeServiceParam struct {
	AnalyzeServiceCache AnalyzeServiceCacheI
}

func NewAnalyzeService(_ context.Context, param AnalyzeServiceParam) *AnalyzeService {
	return &AnalyzeService{
		analyzeServiceCache: param.AnalyzeServiceCache,
	}
}
