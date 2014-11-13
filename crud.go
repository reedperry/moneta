package moneta

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const EVENT_KIND = "event"

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

	entity := new(event)
	if err := readEntity(r, entity); err == nil {
		if err := assertValidKind(entity, c); err == nil {
			key := datastore.NewIncompleteKey(c, EVENT_KIND, nil)
			if k, err := datastore.Put(c, key, entity); err == nil {
				c.Infof("Stored entity with key.String ", k.String())
				c.Infof("Stored entity with key.StringID ", k.StringID())
				c.Infof("Stored entity with key.IntId ", k.IntID())
				c.Infof("Stored entity with key.Kind ", k.Kind())
				resp := response{k.String(), nil, nil}
				if json, err := json.Marshal(resp); err == nil {
					w.Header().Set("Content-Type", "application/json")
					w.Write(json)
				} else {
					handleError(w, err, &c)
				}
			} else {
				handleError(w, err, &c)
			}
		} else {
			handleError(w, err, &c)
		}
	} else {
		handleError(w, err, &c)
	}
}

func handleError(w http.ResponseWriter, err error, c *appengine.Context) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func readEntity(r *http.Request, entity *event) error {
	defer r.Body.Close()

	var body []byte
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		return readErr
	}

	err := json.Unmarshal(body, entity)
	if err != nil {
		return errors.New("Couldn't get valid JSON object from request body.")
	}

	return nil
}

func assertValidKind(entity *event, c appengine.Context) error {
	kind := entity.Kind
	if kind != "expense" && kind != "credit" {
		return errors.New("Invalid kind")
	} else {
		c.Infof("Data has kind %s", kind)
		return nil
	}
}

type event struct {
	Amount  float64   `json:"amount"`
	Comment string    `json:"comment"`
	Date    time.Time `json:"date"`
	Kind    string    `json:"kind"`
	User    string    `json:"user"`
}

type response struct {
	Key   string                 `json:"key"`
	Error string                 `json:"error"`
	Data  map[string]interface{} `json:"data"`
}
