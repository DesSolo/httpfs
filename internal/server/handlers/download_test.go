package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"httpfs/internal/entities"
	"httpfs/internal/mocks"
	"httpfs/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// withChiUrlParam добавляет параметр в контекст http.Request (для совместимости с chi)
func withChiUrlParam(t *testing.T, r *http.Request, key, value string) *http.Request {
	t.Helper()

	rCtx := chi.NewRouteContext()
	rCtx.URLParams.Add(key, value)

	return r.WithContext(
		context.WithValue(r.Context(), chi.RouteCtxKey, rCtx),
	)
}

func TestHandleDownload(t *testing.T) {
	tests := []struct {
		name             string
		downloaderFunc   func(t *testing.T) Downloader
		reqFunc          func(t *testing.T) *http.Request
		expectStatusCode int
		expectBody       string
	}{
		{
			name: "invalid hash",
			downloaderFunc: func(t *testing.T) Downloader {
				return mocks.NewDownloader(t)
			},
			reqFunc: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			expectStatusCode: http.StatusBadRequest,
			expectBody:       `{"message": "invalid hash"}`,
		},
		{
			name: "file not found",
			downloaderFunc: func(t *testing.T) Downloader {
				d := mocks.NewDownloader(t)

				d.EXPECT().Download(mock.Anything, entities.Hash("123"), mock.Anything).Return(storage.ErrNotFound)

				return d
			},
			reqFunc: func(t *testing.T) *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)

				return withChiUrlParam(t, r, "hash", "123")
			},
			expectStatusCode: http.StatusNotFound,
		},
		{
			name: "fault download (broken downloader)",
			downloaderFunc: func(t *testing.T) Downloader {
				d := mocks.NewDownloader(t)

				d.EXPECT().Download(mock.Anything, entities.Hash("123"), mock.Anything).Return(errors.New("iam broken"))

				return d
			},
			reqFunc: func(t *testing.T) *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)

				return withChiUrlParam(t, r, "hash", "123")
			},
			expectStatusCode: http.StatusInternalServerError,
			expectBody:       `{"message": "fault download file"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := tt.reqFunc(t)
			d := tt.downloaderFunc(t)

			HandleDownload(d)(w, r)

			// проверяем код ответа
			assert.Equal(t, tt.expectStatusCode, w.Result().StatusCode)

			if tt.expectBody != "" {
				var buf bytes.Buffer
				if err := json.Indent(&buf, w.Body.Bytes(), "", " "); err != nil {
					t.Fatal(err)
				}

				// проверяем тело
				assert.JSONEq(t, tt.expectBody, w.Body.String(), buf.String())
			}
		})
	}
}
