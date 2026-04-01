package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type StreamEvent struct {
	Type         string          `json:"type"`
	Index        int             `json:"index,omitempty"`
	Message      json.RawMessage `json:"message,omitempty"`
	Delta        json.RawMessage `json:"delta,omitempty"`
	ContentBlock json.RawMessage `json:"content_block,omitempty"`
	Usage        json.RawMessage `json:"usage,omitempty"`
	Error        error           `json:"-"`
}

func (c *Client) Stream(ctx context.Context, req *Request) (<-chan StreamEvent, error) {
	events := make(chan StreamEvent, 100)

	resp, err := c.doRequest(ctx, "POST", "/messages", req)
	if err != nil {
		close(events)
		return events, fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		close(events)
		return events, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	go func() {
		defer close(events)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

		for scanner.Scan() {
			line := scanner.Text()

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var event StreamEvent
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				continue
			}

			events <- event
		}
	}()

	return events, nil
}
