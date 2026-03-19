package rpcclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRawRequestWithContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
		fmt.Fprintln(w, `{"result": null, "error": null, "id": 1}`)
	}))
	defer srv.Close()

	connCfg := &ConnConfig{
		Host:         srv.Listener.Addr().String(),
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	client, err := New(connCfg, nil)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err = client.RawRequestWithContext(ctx, "getblockcount", []json.RawMessage{})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}

	if elapsed > 2*time.Second {
		t.Fatalf("request took %v, expected cancellation within 2s", elapsed)
	}
}
