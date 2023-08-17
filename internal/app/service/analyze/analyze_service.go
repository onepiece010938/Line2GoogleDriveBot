package analyze

import (
	"context"
	"fmt"
)

type CreateAnalyzeParm struct {
	// CreateHistogramParams postgres.CreateHistogramParams
}

type AnalyzeMessageParm struct {
}

func (i *AnalyzeService) CreateAnalyze(ctx context.Context, param CreateAnalyzeParm) error {
	fmt.Println("AnalyzeService->func CreateAnalyze()")
	// var result message.MessageDomainResult
	// result = message.MessageDomainFunc("aabcccc")
	// fmt.Println(result)
	i.analyzeServiceCache.GetMessageCache("")

	return nil
}

func (i *AnalyzeService) AnalyzeTest(ctx context.Context) (string, error) {
	fmt.Println("AnalyzeService->func AnalyzeTest()")
	// var result message.MessageDomainResult
	// result = message.MessageDomainFunc("aabcccc")
	// fmt.Println(result)
	return "TESTTESTTEST", nil
}
