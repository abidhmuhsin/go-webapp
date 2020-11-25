package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type User struct {
	ID        string `json:"id"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}

var users = []User{}

// NewRouter returns an HTTP handler that implements the routes for the API
func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Register the API routes
	r.Get("/", getAllUsers)
	r.Post("/", createUser)
	r.Get("/{id}", getUserByID)
	r.Put("/{id}", updateUser)
	r.Delete("/{id}", deleteUser)

	return r
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	//id := mux.Vars(r)["id"]
	id := chi.URLParam(r, "id")
	index := indexByID(users, id)
	if index < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users[index]); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	//id := mux.Vars(r)["id"]
	id := chi.URLParam(r, "id")
	index := indexByID(users, id)
	if index < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	u := User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	users[index] = u
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&u); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	//id := mux.Vars(r)["id"]
	id := chi.URLParam(r, "id")
	index := indexByID(users, id)
	if index < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	users = append(users[:index], users[index+1:]...)
	w.WriteHeader(http.StatusOK)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	u := User{}

	fmt.Println(r.Body, "hi")
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(u, "hi2")
	users = append(users, u)
	response, err := json.Marshal(&u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func indexByID(users []User, id string) int {
	for i := 0; i < len(users); i++ {
		if users[i].ID == id {
			return i
		}
	}
	return -1
}
