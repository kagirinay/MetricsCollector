package agent

import (
	"io"
	"net/http"
)

// Sender описывает минимальный HTTP-клиент, который необходим
// непосредственно логике агента. Благодаря этому в тесттах можно
// подменять стандартный *http.Client "заглушкой".
type Sender interface {
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}
