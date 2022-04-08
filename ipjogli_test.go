package ipjogli

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestIP(t *testing.T) {
	tests := map[string]struct {
		headers map[string]string
		env     map[string]string
	}{
		"direct": {headers: nil},
		"xff":    {headers: map[string]string{"X-Forwarded-For": "10.0.0.1, 192.0.2.1"}},
		"fly": {
			headers: map[string]string{
				"Fly-Client-IP":   "192.0.2.1",
				"X-Forwarded-For": "10.0.0.1, 192.0.2.1, 198.51.100.80",
			},
			env: map[string]string{"FLY_APP_NAME": "test"},
		},
	}
	expected := "192.0.2.1\n"

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			for hdr, val := range tt.headers {
				req.Header.Set(hdr, val)
			}
			for env, val := range tt.env {
				os.Setenv(env, val)
			}
			w := httptest.NewRecorder()
			handler(w, req)

			resp := w.Result()
			body, _ := ioutil.ReadAll(resp.Body)

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status %v, got %v", http.StatusOK, resp.StatusCode)
			}

			if string(body) != expected {
				t.Errorf("expected body %q, got %q", expected, string(body))
			}

			// Clean up any environment variables we set to not affect other tests
			for env := range tt.env {
				os.Unsetenv(env)
			}
		})
	}
}
