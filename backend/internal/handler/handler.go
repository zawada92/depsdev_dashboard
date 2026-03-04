package handler

import (
	"encoding/json"
	"net/http"

	"dependency-dashboard/config"
	"dependency-dashboard/internal/domain"
	"dependency-dashboard/internal/model"
	"dependency-dashboard/internal/service"

	_ "dependency-dashboard/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

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
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/docs/*", httpSwagger.Handler(
			httpSwagger.URL("/api/v1/docs/doc.json"),
		))

		r.Route("/dependencies", func(r chi.Router) {
			r.Get("/", h.Get)
			r.Get("/{name}", h.Get)
			r.Post("/{name}", h.Post)
			r.Put("/{name}", h.Put)
			r.Delete("/{name}", h.Delete)
			r.Patch("/{name}", h.Patch)
		})
	})

	return r
}

// Get godoc
// @Summary      Get dependencies
// @Description  Fetches list of all dependencies or a specific one by name with optional score filtering.
// @ID           get-dependencies
// @Produce      json
// @Param        name      path      string  false  "Dependency name (optional)"
// @Param        minScore  query     number  false  "Minimum OpenSSF Score filter"
// @Success      200       {array}   model.Dependency
// @Failure      400       {object}  map[string]string "Invalid input"
// @Failure      404       {object}  map[string]string "Not found"
// @Failure      500       {object}  map[string]string "Internal server error"
// @Router       /api/v1/dependencies [get]
// @Router       /api/v1/dependencies/{name} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	minScore, err := parseMinScore(r)
	if err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	log.Info().
		Str("name", name).
		Float64("minScore", minScore).
		Msg("GET dependencies")

	result, err := h.svc.List(r.Context(), name, minScore)
	if err != nil {
		writeError(w, err)
		return
	}

	if len(result) == 0 && name != "" {
		writeError(w, domain.ErrNotFound)
		return
	}

	if result == nil {
		result = []model.Dependency{}
	}

	writeJSON(w, http.StatusOK, result)
}

// Put godoc
// @Summary      Update or Fetch dependency
// @Description  In UpstreamCache mode: triggers fetch (201). In AuthoritativeDb mode: performs full update (200).
// @ID           put-dependency
// @Accept       json
// @Produce      json
// @Param        name      path      string  true   "Dependency name"
// @Param        body      body      model.PatchDependencyRequest false "Body required only in AuthoritativeDb mode"
// @Success      200       {object}  model.Dependency "Updated object (AuthoritativeDb mode)"
// @Success      201       {nil}     nil "Fetched/Created from upstream (UpstreamCache mode)"
// @Failure      400       {object}  map[string]string "Invalid input"
// @Failure      500       {object}  map[string]string "Internal error"
// @Router       /api/v1/dependencies/{name} [put]
func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if h.cfg.ServiceMode == config.UpstreamCache {
		if err := h.svc.Fetch(r.Context(), name); err != nil {
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, nil)
		return
	}

	var req model.PatchDependencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	if req.Version == nil || req.OpenSSFScore == nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	dep, err := h.svc.Patch(r.Context(), name, req)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dep)
}

// Delete godoc
// @Summary      Delete dependency
// @Description  Removes a dependency from the local repository.
// @ID           delete-dependency
// @Param        name      path      string  true   "Dependency name"
// @Success      204       {nil}     nil "No Content"
// @Failure      500       {object}  map[string]string "Internal error"
// @Router       /api/v1/dependencies/{name} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if err := h.svc.Delete(r.Context(), name); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Post godoc
// @Summary      Fetch dependency from upstream
// @Description  Triggers fetching dependency data from upstream. Returns 405 in UpstreamCache mode.
// @ID           post-dependency
// @Param        name      path      string  true   "Dependency name"
// @Success      201       {nil}     nil "Created/Fetched"
// @Failure      405       {nil}     nil "Method Not Allowed (UpstreamCache mode)"
// @Failure      500       {object}  map[string]string "Internal error"
// @Router       /api/v1/dependencies/{name} [post]
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	if h.cfg.ServiceMode == config.UpstreamCache {
		writeJSON(w, http.StatusMethodNotAllowed, nil)
		return
	}

	name := chi.URLParam(r, "name")

	if err := h.svc.Fetch(r.Context(), name); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, nil)
}

// Patch godoc
// @Summary      Partial update
// @Description  Updates specific fields of a dependency. Returns 405 in UpstreamCache mode.
// @ID           patch-dependency
// @Accept       json
// @Produce      json
// @Param        name      path      string  true   "Dependency name"
// @Param        body      body      model.PatchDependencyRequest true "Partial update fields"
// @Success      200       {object}  model.Dependency
// @Failure      400       {object}  map[string]string "Invalid input"
// @Failure      405       {nil}     nil "Method Not Allowed (UpstreamCache mode)"
// @Failure      500       {object}  map[string]string "Internal error"
// @Router       /api/v1/dependencies/{name} [patch]
func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	if h.cfg.ServiceMode == config.UpstreamCache {
		writeJSON(w, http.StatusMethodNotAllowed, nil)
		return
	}

	name := chi.URLParam(r, "name")

	var req model.PatchDependencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	if req.Version == nil && req.OpenSSFScore == nil {
		writeError(w, domain.ErrInvalidInput)
		return
	}

	dep, err := h.svc.Patch(r.Context(), name, req)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dep)
}
