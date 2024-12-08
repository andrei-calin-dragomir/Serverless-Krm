package function

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHandle ensures that Handle executes without error and returns the
// HTTP 201 status code indicating no errors.
func TestHandle(t *testing.T) {
	podSpec := `{
		"metadata": {
			"name": "test-pod",
			"namespace": "default",
			"labels": {
				"app": "test-app"
			}
		},
		"spec": {
			"containers": [
				{
					"name": "nginx-container",
					"image": "nginx:latest",
					"ports": [
						{
							"containerPort": 80
						}
					]
				}
			]
		}
	}`
	var (
		w   = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "http://example.com/create-pod", strings.NewReader(podSpec))
		res *http.Response
	)

	Handle(w, req)
	res = w.Result()
	defer res.Body.Close()

	if res.StatusCode != 201 {
		t.Fatalf("unexpected response code: %v", res.StatusCode)
	}
}
