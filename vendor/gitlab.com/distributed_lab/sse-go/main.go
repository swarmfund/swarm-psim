package sse

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	//ErrNilChan will be returned by Notify if it is passed a nil channel
	ErrNilChan = fmt.Errorf("nil channel given")
)

type RequestBuilder func() (*http.Request, error)

type Listener struct {
	Request RequestBuilder
	client  http.Client
	retry   *time.Ticker
}

func NewListener(request RequestBuilder) *Listener {
	return &Listener{
		Request: request,
		client:  http.Client{},
		retry:   time.NewTicker(1 * time.Second),
	}
}

func (l *Listener) Events() <-chan Event {
	events := make(chan Event)
	go func() {
		for {
			err := l.Subscribe(events)
			if err != nil {
				events <- Event{
					Err: err,
				}
			}
			// wait for some time before retry
			<-l.retry.C
		}
	}()
	return events
}

func (l *Listener) Subscribe(events chan<- Event) error {
	if events == nil {
		return ErrNilChan
	}

	request, err := l.Request()
	if err != nil {
		return fmt.Errorf("error getting sse request: %v", err)
	}
	request.Header.Set("Accept", "text/event-stream")

	response, err := l.client.Do(request)
	if err != nil {
		return fmt.Errorf("error performing request %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error getting resource: %d", response.StatusCode)
	}

	br := bufio.NewReader(response.Body)

	delim := []byte{':', ' '}

	var event Event

	for {
		bs, err := br.ReadBytes('\n')

		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF {
			break
		}

		if len(bs) < 2 {
			continue
		}

		spl := bytes.SplitN(bs, delim, 2)

		if len(spl) < 2 {
			continue
		}

		event = Event{}
		switch string(spl[0]) {
		case "data":
			event.Data = bytes.NewBuffer(bytes.TrimSpace(spl[1]))
			events <- event
		}
	}

	return nil
}
