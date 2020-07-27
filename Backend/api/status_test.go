package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleStatus(t *testing.T) {
	req := httptest.NewRequest("POST", "/games/", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(handleStatus)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("exepcted %v, got %v", http.StatusOK, status)
	}

	expected := `{"ok":true}`
	if body := strings.TrimSpace(rr.Body.String()); body != expected {
		t.Errorf("expected %v [%d], got %v [%d]", expected, len(expected), body, len(body))
	}
}
