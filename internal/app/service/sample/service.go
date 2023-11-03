package sample

import "context"

type SampleService struct {
	sampleServiceDynamodb SampleServiceDynamodbI
}

type SampleServiceParam struct {
	SampleServiceDynamodb SampleServiceDynamodbI
}

func NewSampleService(_ context.Context, param SampleServiceParam) *SampleService {
	return &SampleService{
		sampleServiceDynamodb: param.SampleServiceDynamodb,
	}
}
