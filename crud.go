package moneta

import (
	"appengine"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/data", crud)
}

func crud(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprint(w, "Hello!")
	case "POST":
		doPost(w, r)
	}
}

func doPost(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	data, err := getRequestData(r)
	if err == nil {
		for k, v := range *data {
			if k == "kind" {
				kind := v
				c.Infof("Data has kind %s", kind)
			}
		}
	} else {
		c.Errorf(err.Error())
	}
}

func getRequestData(r *http.Request) (*map[string]interface{}, error) {
	defer r.Body.Close()

	var data map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, errors.New("Couldn't get valid JSON object from request body.")
	}

	return &data, nil
}
