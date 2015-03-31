package moneta

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const EVENT_KIND = "event"

func init() {
	http.HandleFunc("/data", crud)
	http.HandleFunc("/", serveApp)
}

func serveApp(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if err := authorize(c); err != nil {
		loginURL, _ := user.LoginURL(c, "/")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `You are not signed in! Sign in <a href="%s">here</a>.`, loginURL)
		return
	}

	if r.Method != "GET" {
		fmt.Fprint(w, "Method not supported!")
		return
	}

	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Fprint(w, "index.html not found!")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, string(content))
}

func crud(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if err := authorize(c); err != nil {
		fmt.Fprintf(w, "Authorization error: %v", err)
		return
	}

	switch r.Method {
	case "GET":
		doGet(w, r, c)
	case "POST":
		doPost(w, r, c)
	default:
		fmt.Fprint(w, "Method not supported!")
	}
}

func doGet(w http.ResponseWriter, r *http.Request, c appengine.Context) {
	qParams := new(QParams)
	if err := readQParams(r, qParams); err != nil {
		fmt.Fprintf(w, "Invalid query request: %v", err)
		return
	}

	if err := assertValidKind(qParams.Kind, c); err != nil {
		fmt.Fprint(w, "Invalid kind '%v'", qParams.Kind)
		return
	}

	query := datastore.NewQuery(EVENT_KIND)
	query = applyFilters(qParams, query)
	results := make([]event, 0, 10)
	keys, err := query.GetAll(c, &results)
	if err == nil {
		c.Infof("GetAll returned %v results", len(keys))
	} else {
		handleError(w, err, &c)
	}

	resp := response{"", "", results}
	if json, err := json.Marshal(resp); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	} else {
		handleError(w, err, &c)
	}
}

func doPost(w http.ResponseWriter, r *http.Request, c appengine.Context) {

	entity := new(event)
	if err := readEntity(r, entity); err == nil {
		if err := assertValidKind(entity.Kind, c); err == nil {
			key := datastore.NewIncompleteKey(c, EVENT_KIND, nil)
			if k, err := datastore.Put(c, key, entity); err == nil {
				resp := response{k.String(), "", nil}
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

func applyFilters(qParams *QParams, query *datastore.Query) *datastore.Query {

	query.Filter("Kind", qParams.Kind)

	if qParams.Amount != 0 {
		query = query.Filter("Amount =", qParams.Amount)
	}
	if qParams.Comment != "" {
		query = query.Filter("Comment =", qParams.Comment)
	}
	if !qParams.MinDate.IsZero() {
		query = query.Filter("Date >=", qParams.MinDate)
	}
	if !qParams.MaxDate.IsZero() {
		query = query.Filter("Date <=", qParams.MaxDate)
	}

	return query
}

func readQParams(r *http.Request, qParams *QParams) error {
	var err error

	qParams.Kind = r.FormValue("kind")
	if qParams.Kind == "" {
		return errors.New("Query has no 'kind'!")
	}

	if r.FormValue("amount") != "" {
		qParams.Amount, err = strconv.ParseFloat(r.FormValue("amount"), 64)
		if err != nil {
			return err
		}
	}

	if r.FormValue("comment") != "" {
		qParams.Comment = r.FormValue("comment")
	}

	if r.FormValue("end") != "" {
		maxDate, err := time.Parse(time.RFC3339Nano, r.FormValue("end"))
		if err != nil {
			return err
		} else {
			qParams.MaxDate = maxDate
		}
	}

	if r.FormValue("start") != "" {
		minDate, err := time.Parse(time.RFC3339Nano, r.FormValue("start"))
		if err != nil {
			return err
		} else {
			qParams.MinDate = minDate
		}
	}

	return nil
}

func assertValidKind(kind string, c appengine.Context) error {
	if kind != "expense" && kind != "credit" {
		return errors.New("Invalid kind: '" + kind + "'")
	} else {
		c.Infof("Data has kind %s", kind)
		return nil
	}
}

func authorize(c appengine.Context) error {
	if u := user.Current(c); u == nil {
		return errors.New("User not logged in!")
	}
	return nil
}

type event struct {
	Amount  float64   `json:"amount"`
	Comment string    `json:"comment"`
	Date    time.Time `json:"date"`
	Kind    string    `json:"kind"`
	User    string    `json:"user"`
}

type QParams struct {
	Kind    string    `json:"kind"`
	Amount  float64   `json:"amount"`
	Comment string    `json:"comment"`
	MaxDate time.Time `json:"before"`
	MinDate time.Time `json:"after"`
}

type response struct {
	Key   string  `json:"key"`
	Error string  `json:"error"`
	Data  []event `json:"data"`
}
