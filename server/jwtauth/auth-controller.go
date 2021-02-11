package jwtauth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	TokenAuth = jwtauth.New("HS256", []byte("secretkey"), nil)
	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	//_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
	//fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
}

var users = map[string]string{
	"admin": "admin",
	"user":  "admin@jwt.auth",
}

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func Login(w http.ResponseWriter, r *http.Request) {

	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	fmt.Println(r.Body)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]
	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		// w.WriteHeader(http.StatusUnauthorized)  --  will not be processed in login page if using interceptor to break all 401 requests
		// send ok response without token to be handled in ui
		respJSON := make(map[string]interface{})
		respJSON["isLoggedIn"] = false
		res, _ := json.Marshal(respJSON)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	}

	// setup and return token if success
	// TokenAuth = jwtauth.New("HS256", []byte("secretkey"), nil) -- do in init

	//Generate jwt token with claims `user_id:123` here:
	claims := map[string]interface{}{"user_id": creds.Username}
	jwtauth.SetExpiryIn(claims, 10*60*60*1000000000) // set expiry time in nanoseconds -- H*M*S*1000000000 . 1s = 1e+9 nanoseconds. ie 9 zeroes
	_, tokenString, _ := TokenAuth.Encode(claims)
	//fmt.Printf("DEBUG: jwt is %s\n\n", tokenString)

	//var respJson map[string]interface{}
	respJSON := make(map[string]interface{})
	respJSON["isLoggedIn"] = true
	respJSON["token"] = tokenString

	res, _ := json.Marshal(respJSON)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// AuthHello is a sample handler
func AuthHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world.. You are authorised!")
}

// NewRouter returns an HTTP handler that implements the routes for the API
func NewRouter() http.Handler {
	router := chi.NewRouter()

	// Register the API routes
	//Setup Login - No token for login
	router.Post("/login", Login)
	/* To receive token - hit post request _mount-prefix_/login with json body
	{
	"username": "admin",
	"password": "admin"
	}
	*/
	router.Group(func(r chi.Router) {

		// Seek, verify and validate JWT tokens for all calls in the group
		r.Use(jwtauth.Verifier(TokenAuth))
		r.Use(jwtauth.Authenticator) // Handle valid / invalid tokens and accept/reject group.

		// Setup authorized routes which need valid token
		r.Get("/auth/hello", AuthHello) // hit _mount-prefix_/auth/hello with -- Authorization:Bearer token-value'
		// Add another set of authorized subroutes
		// r.Mount("/auth/users", users.NewRouter())
	})
	//End:  JWT - Auth based calls

	return router
}
