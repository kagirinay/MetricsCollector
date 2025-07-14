package agent

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

type fakeClient struct {
	calls int
	urls  []string
}

func (f *fakeClient) Post(url, ct string, body io.Reader) (*http.Response, error) {
	f.calls++
	f.urls = append(f.urls, url)
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer(nil)),
	}, nil
}

func TestAgentReport(t *testing.T) {
	client := &fakeClient{}
	a := New(client, "http://example.com")
	// Запустим один poll и вручную вызовем report
	a.collector.poll()
	a.report()
	if client.calls == 0 {
		t.Fatalf("expected at least one HTTP call")
	}
}
