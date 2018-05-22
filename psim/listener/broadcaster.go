package listener

import (
	"fmt"
	"context"
)

func (s *Service) BroadcastEvents(ctx context.Context) (success bool, err error) {
	source := s.listener.Listen(ctx)
	for event := range source {
		fmt.Println(event)
	}
	return false, nil
}