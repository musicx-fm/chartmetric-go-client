package chartmetric

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func buildJSONBody(body any) (io.Reader, error) {
	if body == nil {
		return http.NoBody, nil
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	return bytes.NewBuffer(jsonBody), nil
}

func readResponseBody(respBody io.ReadCloser) ([]byte, error) {
	var buf bytes.Buffer

	if _, err := io.Copy(&buf, respBody); err != nil {
		return nil, fmt.Errorf("copy bytes from response body: %w", err)
	}

	return buf.Bytes(), nil
}

func isStatusSuccess(resp *http.Response) bool {
	statusCode := resp.StatusCode

	return statusCode >= 200 && statusCode < 300
}

func addQueryParams(req *http.Request, params map[string]any) {
	if len(params) == 0 {
		return
	}

	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, fmt.Sprintf("%v", val))
	}

	req.URL.RawQuery = q.Encode()
}
