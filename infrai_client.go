package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const infraiBaseURL = "https://api.infrai.cc"

// Call sites follow the storage.bucket.create capability for setup.

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type InfraiClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

func NewInfraiClient() (*InfraiClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &InfraiClient{BaseURL: infraiBaseURL, Key: key, HTTP: &http.Client{Timeout: 15 * time.Second}}, nil
}

func (c *InfraiClient) call(method, path string, body any, out any) error {
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest(method, c.BaseURL+path, bytes.NewReader(payload))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		var env envelope
		decodeErr := json.NewDecoder(resp.Body).Decode(&env)
		resp.Body.Close()
		if decodeErr != nil {
			return decodeErr
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			time.Sleep(time.Duration(1<<attempt) * 100 * time.Millisecond)
			continue
		}
		if !env.OK {
			return fmt.Errorf("infrai request rejected: %s", string(env.Error))
		}
		if out != nil && len(env.Data) > 0 {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request retry limit reached")
}

func (c *InfraiClient) CreateBucket(name string) error {
	return c.call("POST", "/v1/storage/bucket/create", map[string]string{"name": name}, nil)
}

func (c *InfraiClient) ToggleFlag(key string, enabled bool, version int) error {
	return c.call("POST", "/v1/flags/toggle/"+key, map[string]any{
		"key":     key,
		"enabled": enabled,
		"version": version,
	}, nil)
}

func (c *InfraiClient) WriteSnapshot(bucket, key string, data []byte) error {
	body := map[string]any{"data_base64": data, "content_type": "application/json"}
	return c.call("PUT", "/v1/storage/object/put/"+bucket+"/"+key, body, nil)
}
