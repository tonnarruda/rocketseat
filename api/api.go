package api

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(db map[string]string) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Post("/api/shorten", handlePost(db))
	r.Get("/{code}", handleGet(db))
	return r
}

type PostBody struct {
	URL string `json:"url"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func handlePost(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body PostBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			SendJSON(
				w,
				Response{Error: "invalid body"},
				http.StatusUnprocessableEntity,
			)
			return
		}
		if _, err := url.Parse(body.URL); err != nil {
			SendJSON(
				w,
				Response{Error: "invalid url"},
				http.StatusBadRequest,
			)
			return
		}
		code := gencode()
		db[code] = body.URL
		SendJSON(
			w,
			Response{Data: code},
			http.StatusCreated,
		)
	}
}

func handleGet(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		data, ok := db[code]

		if !ok {
			http.Error(w, "url nao encontrada", http.StatusNotFound)
		}
		http.Redirect(w, r, data, http.StatusPermanentRedirect)

	}
}

func SendJSON(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal json data", "error", err)
		SendJSON(
			w,
			Response{Error: "something went wrong"},
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response to client", "error", err)
		return
	}
}

const characters = "qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM1234567890"

func gencode() string {
	const n = 8
	bts := make([]byte, n)
	for i := range n {
		bts[i] = characters[rand.Intn(len(characters))]
	}
	return string(bts)
}
