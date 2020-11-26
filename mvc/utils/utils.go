package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			log.Println(err)
			return
		}
	}

	/* or simply use
	if err := json.NewDecoder(r.Body).Decode(x); err != nil {
		log.Println(err)
		return
	}
	log.Println(x)
	*/
	// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
}
