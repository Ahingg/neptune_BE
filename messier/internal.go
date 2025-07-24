package messier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"neptune/backend/pkg/utils"
	"net/http"
	"time"
)

func SendRequest(ctx context.Context, method, url string, reqBody interface{}, resp interface{}, authToken string) error {
	fmt.Printf("Sending %s request to: %s\n", method, url)
	var bodyReader io.Reader

	if reqBody != nil {
		jsonBytes, err := json.Marshal(reqBody)
		if err != nil {
			utils.CheckPanic(err)
			return err
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		utils.CheckPanic(fmt.Errorf("failed to marshal request body: %w", err))
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	respHTTP, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		utils.CheckPanic(fmt.Errorf("failed to send request: %w", err))
		return fmt.Errorf("failed to send request to messier: %w", err)
	}

	defer respHTTP.Body.Close()

	if respHTTP.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(respHTTP.Body)
		return fmt.Errorf("external API returned status %d: %s", respHTTP.StatusCode, string(bodyBytes))
	}
	if resp != nil {
		if err := json.NewDecoder(respHTTP.Body).Decode(resp); err != nil {
			return fmt.Errorf("failed to decode external API response: %w", err)
		}
	}
	return nil
}
