package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/praelatus/backend/models"
	"github.com/praelatus/backend/mw"
	"github.com/praelatus/backend/store"
	"github.com/pressly/chi"
)

func labelRouter() chi.Router {
	router := chi.NewRouter()

	router.Get("/", GetAllLabels)
	router.Post("/", CreateLabel)

	router.Get("/search", SearchLabels)
	router.Get("/:id", GetLabel)
	router.Delete("/:id", DeleteLabel)
	router.Put("/:id", UpdateLabel)

	return router
}

// GetAllLabels will return a JSON array of all labels from the store.
func GetAllLabels(w http.ResponseWriter, r *http.Request) {
	labels, err := Store.Labels().GetAll()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Println(err)
		return
	}

	sendJSON(w, labels)
}

// GetLabel will return a JSON representation of a label
func GetLabel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	lbl := &models.Label{}

	i, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(400)
		w.Write(apiError("invalid id"))
		log.Println(err)
		return
	}

	lbl.ID = int64(i)

	err = Store.Labels().Get(lbl)
	if err != nil {
		w.WriteHeader(500)
		w.Write(apiError(err.Error()))
		return
	}

	sendJSON(w, lbl)
}

// CreateLabel creates a label in the db and return a JSON object of
func CreateLabel(w http.ResponseWriter, r *http.Request) {
	var lbl models.Label

	u := mw.GetUser(r.Context())
	if u == nil {
		w.WriteHeader(403)
		w.Write(apiError("you must be logged in to create a label"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&lbl)
	if err != nil {
		w.WriteHeader(400)
		w.Write(apiError("malformed json"))
		log.Println(err)
		return
	}

	err = Store.Labels().New(&lbl)
	if err != nil {
		w.WriteHeader(500)
		w.Write(apiError(err.Error()))
		log.Println(err)
		return
	}

	sendJSON(w, lbl)
}

// UpdateLabel updates the label in the db and returns a message indicating
// success or failure.
func UpdateLabel(w http.ResponseWriter, r *http.Request) {
	var lbl models.Label

	u := mw.GetUser(r.Context())
	if u == nil {
		w.WriteHeader(403)
		w.Write(apiError("you must be logged in to create a label"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&lbl)
	if err != nil {
		w.WriteHeader(400)
		w.Write(apiError("malformed json"))
		log.Println(err)
		return
	}

	if lbl.ID == 0 {
		id := chi.URLParam(r, "id")
		i, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apiError(http.StatusText(http.StatusBadRequest)))
			return
		}

		lbl.ID = int64(i)
	}

	err = Store.Labels().Save(lbl)
	if err != nil {
		w.WriteHeader(500)
		w.Write(apiError(err.Error()))
		log.Println(err)
		return
	}

	sendJSON(w, lbl)
}

// DeleteLabel deletes labels from the db and returns a repsonse indicating
// success of failure.
func DeleteLabel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	i, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(400)
		w.Write(apiError("invalid id"))
		log.Println(err)
		return
	}

	err = Store.Labels().Remove(models.Label{ID: int64(i)})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("Label successfully deleted"))
}

// SearchLabels will take a url param of query and try to find a label
// with the given name
func SearchLabels(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")

	labels, err := Store.Labels().Search(query)
	if err != nil {
		if err == store.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write(apiError("No labels match that query"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(apiError(err.Error()))
		return
	}

	sendJSON(w, labels)
}
