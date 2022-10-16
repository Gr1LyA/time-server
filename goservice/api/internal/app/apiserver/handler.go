package apiserver

import (
	"fmt"
	"net/http"
	"os"
)


type handler struct {
	http.ServeMux
	cache		Cache
}

func newHandler(cache Cache) *handler {
	s := &handler{
		cache: cache,
	}

	s.HandleFunc("/time", s.getTime)
	return s
}

func (h *handler) getTime (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	content, err := h.cache.Load("responce")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, content)
}