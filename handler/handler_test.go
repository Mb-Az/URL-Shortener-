package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	mapstruct "urlShortener/mapStruct"
	"urlShortener/shortener"

	"github.com/gin-gonic/gin"
)

type ShortenRes struct {
	ShortLink string `json:"short_link"`
}

type ShortenErr struct {
	Error string `json:"error"`
}

type RequestMap struct {
	longUrl   string
	ShortLink string
}

var (
	r      *gin.Engine
	reqMap = &RequestMap{}
)

func StartServer() {
	gin.SetMode(gin.TestMode)
	r = gin.New()
	mapS := mapstruct.New("sync")

	r.POST("/shorten", Shorten(mapS))
	r.GET("/link/:id", Link(mapS))
}

func TestMain(m *testing.M) {
	StartServer()
	os.Exit(m.Run())
}

func TestPostShortenOk(t *testing.T) {

	longUrl, err := shortner.RandomString(20)
	// reqMap.longUrl = longUrl

	if err != nil {
		t.Fatalf("Error in generating long url %v, got", err)
	}

	var body []byte
	body = fmt.Appendf(body, `{"url":"%v"}`, longUrl)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d",
			http.StatusOK, w.Code)
	}

	var res ShortenRes
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	// reqMap.ShortLink = res.ShortLink

}

func TestPostShortenErr(t *testing.T) {

	body := []byte(`{"url":""}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest, w.Code)
	}

	var res ShortenErr
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

}

func TestGetLink(t *testing.T) {

	longUrl, err := shortner.RandomString(20)

	if err != nil {
		t.Fatalf("Error in generating long url %v, got", err)
	}

	var body []byte
	body = fmt.Appendf(body, `{"url":"%v"}`, longUrl)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d",
			http.StatusOK, w.Code)
	}

	var res ShortenRes
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	// fmt.Println(res.ShortLink)

	parts := strings.Split(res.ShortLink, "/")
	short := parts[len(parts)-1]

	// fmt.Println(short)

	// fmt.Println("/link/" + short)

	req = httptest.NewRequest(http.MethodGet, "/link/"+short, nil)

	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d",
			http.StatusSeeOther, w.Code)
	}

	location := w.Header().Get("Location")
	if location != "/link/"+longUrl {
		t.Errorf("expected redirect to %v, got %v", "/link/"+longUrl, location)
	}

}

func TestGetLinkErr(t *testing.T) {

	fakeUrl, err := shortner.RandomString(8)

	if err != nil {
		t.Fatalf("Error in generating long url %v, got", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf(`/link/%v`, fakeUrl), nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d",
			http.StatusNotFound, w.Code)
	}

}

func TestPostShortenConcurrent(t *testing.T) {
	const numRequests = 20000

	var wg sync.WaitGroup

	errors := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			body, err := json.Marshal(map[string]string{
				"url": fmt.Sprintf("https://example.com/%d", i),
			})
			if err != nil {
				errors <- err
				return
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/shorten",
				bytes.NewReader(body),
			)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				errors <- fmt.Errorf(
					"request %d: expected status %d, got %d; body: %s",
					i,
					http.StatusOK,
					w.Code,
					w.Body.String(),
				)
				return
			}

			var res ShortenRes

			if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
				errors <- fmt.Errorf(
					"request %d: failed to decode response: %w",
					i,
					err,
				)
				return
			}

			if res.ShortLink == "" {
				errors <- fmt.Errorf(
					"request %d: empty short link",
					i,
				)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}
