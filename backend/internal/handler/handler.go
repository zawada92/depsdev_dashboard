package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"dependency-dashboard/config"
	"dependency-dashboard/internal/model"
	"dependency-dashboard/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// TODO_TOM logging
type Handler struct {
	cfg *config.Config
	svc *service.Service
}

func New(cfg *config.Config, s *service.Service) *Handler {
	return &Handler{cfg: cfg, svc: s}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{h.cfg.CorsAddress},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/dependencies", h.Get)
		r.Get("/dependencies/{name}", h.Get)
		r.Post("/dependencies/{name}", h.Post)
		r.Put("/dependencies/{name}", h.Put)
		r.Delete("/dependencies/{name}", h.Delete)
		r.Patch("/dependencies/{name}", h.Patch)
	})

	return r
}

// TODO_TOM swagger
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	minScoreStr := r.URL.Query().Get("minScore")

	minScore, _ := strconv.ParseFloat(minScoreStr, 64)

	result, err := h.svc.List(r.Context(), name, minScore)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if len(result) == 0 {
		if name != "" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			result = []model.Dependency{}
			json.NewEncoder(w).Encode(result)
		}
	} else {
		json.NewEncoder(w).Encode(result)
	}
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if err := h.svc.Fetch(r.Context(), name); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if err := h.svc.Delete(r.Context(), name); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	if h.cfg.ServiceMode == config.UpstreamCache {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	name := chi.URLParam(r, "name")

	if err := h.svc.Fetch(r.Context(), name); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	if h.cfg.ServiceMode == config.UpstreamCache {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	name := chi.URLParam(r, "name")

	if err := h.svc.Patch(r.Context(), name); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
