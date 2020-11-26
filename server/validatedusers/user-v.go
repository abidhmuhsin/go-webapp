package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        int    `json:"id" validate:"isdefault"`
	Lastname  string `json:"lastname" validate:"required"`
	Firstname string `json:"firstname" validate:"required"`
	Age       int    `json:"age"`
	Email     string `json:"email" validate:"email-blacklist"`
}

var users = []User{}
var id = 0
var validate = validator.New()

// NewRouter returns an HTTP handler that implements the routes for the API
func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Register the API routes
	r.Get("/", getAllUsers)
	r.Post("/", createUser)

	return r
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// {"lastname":"Pu","firstname":"Kak","age":18,"email":"k@te.st"}
	u := User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		responseBody := err.Error()
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		}
		return
	}

	// validation start
	// validate := validator.New() -- use global one after registering custom validator in init function
	err := validate.Struct(u)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		responseBody := map[string]string{"error": validationErrors.Error()}
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		}
		return
	}
	// validation end

	// We don't want an API user to set the ID manually
	// in a production use case this could be an automatically
	// ID in the database
	u.ID = id
	id++

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

func EmailBlacklistValidator(f1 validator.FieldLevel) bool {
	emailBlacklist := []string{"fake@email.com", "spam@email.com"}
	email := f1.Field().String()
	for _, e := range emailBlacklist {
		if email == e {
			return false
		}
	}
	return true
}

func init() {
	fmt.Println("init func in user-v.go - registering custom validator")
	validate.RegisterValidation("email-blacklist", EmailBlacklistValidator)
	// {"lastname":"Pu","firstname":"Kak","age":18,"email":"fake@email.com"}
}
