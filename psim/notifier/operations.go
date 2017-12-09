package notifier

import (
	"context"

	"io/ioutil"

	"encoding/json"
	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	sse "gitlab.com/distributed_lab/sse-go"
	"gitlab.com/tokend/psim/psim/notifier/internal/operations"
)

func (s *Service) listenOperations(ctx context.Context) {
	if s.Operations == nil {
		s.logger.Warn("operations listener is not enabled")
		return
	}

	listener := sse.NewListener(s.paymentsRequest)
	events := listener.Events()

	var err error
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Received stop signal")
			//ToDo: add sse connection closer
			return

		case event, ok := <-events:
			s.logger.Info("Received sse event")
			if !ok {
				s.logger.Info("SSE connection closed")
				return
			}
			if event.Err != nil {
				s.errors <- errors.Wrap(event.Err, "Failed to get event")
				continue
			}

			err = s.processOperationEvent(event)
			if err != nil {
				s.errors <- err
			}
		}
	}
}

func (s *Service) processOperationEvent(event sse.Event) (err error) {
	defer func() error {
		if err != nil {
			err = errors.Wrap(err, "Process event failed")
		}
		return err
	}()

	rawOperation, err := ioutil.ReadAll(event.Data)
	if err != nil {
		err = errors.Wrap(err, "Failed to unmarshal op")
		return
	}

	err = s.processOperation(rawOperation)
	logFields := logan.F{
		"success": err == nil,
		"cursor":  s.Operations.Cursor,
		"error":   err,
	}

	// Wrapping of the nil error return nil
	err = errors.Wrap(err, "Failed to process operation", logFields)
	s.logger.WithFields(logFields).Info("Event processing completed")
	return err
}

func (s *Service) processOperation(rawOperation []byte) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = errors.FromPanic(rec)
			s.logger.WithStack(err).WithError(err).Error("Recover of the operation processing")
		}
	}()

	base := new(operations.Base)
	err = json.Unmarshal(rawOperation, base)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal op_base")
	}

	s.Operations.Cursor = fmt.Sprintf("%d", base.ID)
	logger := s.logger.WithFields(base.LogFields())
	logger.Debug("Got new operation")

	op, err := operations.ParseOperation(base, rawOperation)
	if err != nil {
		cause := errors.Cause(err)
		if cause == operations.ErrorUnsupportedOpType {
			logger.WithError(err).Debug("Unsupported operation. Skip")
			return nil
		}
		return errors.Wrap(err, "failed to parse operation", base.LogFields())
	}

	participantsMap, err := s.getParticipants(op.ParticipantsRequest())
	if err != nil {
		return err
	}
	if len(participantsMap[base.ID]) == 0 {
		return nil
	}

	op.UpdateParticipants(participantsMap[base.ID])
	letters, err := op.CraftLetters(s.ProjectName)
	if err != nil {
		return errors.Wrap(err, "failed to create letters", base.LogFields())
	}

	for _, letter := range letters {
		err = s.sendEmail(letter, s.Operations.PayloadID)
		if err != nil {
			return errors.Wrap(err, "Failed to send offer letter", base.LogFields())
		}
	}

	return nil
}
