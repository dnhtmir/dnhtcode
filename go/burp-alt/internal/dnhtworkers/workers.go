package dnhtworkers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"intruder/internal/dnhtclient"
)

const (
	ph = "{PLACEHOLDER}"
)

type Worker struct {
	Client  dnhtclient.HTTPClient
	Method  string
	URL     string
	Headers http.Header
	Body    string
	RPS     int
	Jitter  float64
}

func (w *Worker) StartWorker(id int, jobs <-chan string, results chan<- string) {
	for job := range jobs {
		url := strings.ReplaceAll(w.URL, ph, job)
		body := strings.ReplaceAll(w.Body, ph, job)

		headers := make(http.Header)
		for key, values := range w.Headers {
			newKey := strings.ReplaceAll(key, ph, job)
			for _, value := range values {
				newValue := strings.ReplaceAll(value, ph, job)
				headers.Add(newKey, newValue)
			}
		}

		respBody, err := dnhtclient.ExecuteRequest(w.Client, w.Method, url, headers, body)
		if err != nil {
			results <- fmt.Sprintf("Worker %d: Error - %s", id, err)
			continue
		}
		results <- fmt.Sprintf("Worker %d: URL: %s Len: %s", id, url, respBody) // TODO: respBody is taking the len

		if w.RPS > 0 {
			delay := time.Second / time.Duration(w.RPS)
			if w.Jitter > 0 {
				jitterValue := time.Duration(rand.Float64() * float64(delay) * w.Jitter)
				delay += jitterValue
			}
			time.Sleep(delay)
		}
	}
}
