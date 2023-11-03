package sample

import (
	"context"
	"fmt"
)

func (s *SampleService) Sample(ctx context.Context, lineID string) (string, error) {
	t, err := s.sampleServiceDynamodb.GetGoogleOAuthToken(lineID)
	if err != nil {
		return "", err
	}
	fmt.Println(t)
	return t.PK, nil
}
