package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/bekzat-kamen/RestApiFromZero/internal/database"
	"github.com/bekzat-kamen/RestApiFromZero/internal/models"
)

type Handlers struct {
	store *database.TaskStore
}

func NewHandlers(store *database.TaskStore) *Handlers {
	return &Handlers{store}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func (h *Handlers) getAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handlers) getTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	task, err := h.store.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) createTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}
}
