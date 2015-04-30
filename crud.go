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
	"strings"
	"time"
)

const EVENT_KIND = "event"
const USER_KIND = "user"
const EXPENSE = "expense"
const INCOME = "income"

func init() {
	http.HandleFunc("/data", crud)
	http.HandleFunc("/user", userData)
	http.HandleFunc("/", serveApp)
}

func serveApp(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if _, err := authorize(c); err != nil {
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

// Crud is the entry point to CRUD operations. It directs requests based on
// finding an authorized user and the method of the request.
func crud(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u, err := authorize(c)
	if err != nil {
		fmt.Fprintf(w, "Authorization error: %v", err)
		return
	}

	switch r.Method {
	case "GET":
		doGet(w, r, c, u)
	case "POST":
		doPost(w, r, c, u)
	default:
		c.Infof("Method %v not supported!", r.Method)
		fmt.Fprintf(w, "Method %v not supported!", r.Method)
	}
}

func userData(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u, err := authorize(c)
	if err != nil {
		fmt.Fprintf(w, "Authorization error: %v", err)
		return
	}

	switch r.Method {
	case "GET":
		getUser(w, r, c, u)
	case "POST":
		storeUser(w, r, c, u)
	default:
		c.Infof("Method %v not supported!", r.Method)
		fmt.Fprintf(w, "Method %v not supported!", r.Method)
	}
}

func getUser(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
	userId := u.Email
	c.Infof("Getting user %v", userId)
	userKeyId := createUserKeyId(userId, c)
	userKey := datastore.NewKey(c, USER_KIND, userKeyId, 0, nil)
	appUser := new(AppUser)

	err := datastore.Get(c, userKey, appUser)
	if err != nil {
		c.Infof("No user data found for %v. %v", userId, err.Error())
		http.Error(w, "User not found", 404)
		return
	} else {
		resp := UserResponse{true, *appUser}

		if text, err := json.Marshal(resp); err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(text)
		} else {
			handleError(w, err, &c)
		}
	}
}

func storeUser(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
	appUser := new(AppUser)
	if err := readEntity(r, appUser); err != nil {
		handleError(w, err, &c)
	}

	if !isValidUser(appUser) {
		c.Infof("Invalid user object in request: %v", appUser)
		http.Error(w, "Invalid user object!", 400)
		return
	}

	if appUser.Id != u.Email {
		c.Infof("User %v attempted to store data about user %v. Denied.", u.Email, appUser.Id)
		http.Error(w, "Not authorized to change another user!", 403)
		return
	}

	userKeyId := createUserKeyId(appUser.Id, c)
	userKey := datastore.NewKey(c, USER_KIND, userKeyId, 0, nil)
	_, err := datastore.Put(c, userKey, appUser)
	if err != nil {
		handleError(w, err, &c)
	}
	
	c.Infof("Updated user %v", appUser.Id)

	resp := UserResponse{true, *appUser}
	text, err := json.Marshal(resp)
	if err != nil {
		handleError(w, err, &c)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(text)
}

func createUserKeyId(userId string, c appengine.Context) string {
	if userId == "" {
		c.Warningf("Creating a user data key with no userId!")
	}

	return "user:" + userId
}

// DoGet handles GET requests, and returns data to the user based on a query.
func doGet(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
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
	query = applyFilters(qParams, query, c, u)
	query = applySort(qParams, query, c)
	results := make([]Event, 0, 10)
	keys, err := query.GetAll(c, &results)
	if err == nil {
		c.Infof("GetAll returned %v results", len(keys))
	} else {
		handleError(w, err, &c)
	}

	resp := CrudResponse{true, "", "", results}
	if text, err := json.Marshal(resp); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(text)
	} else {
		handleError(w, err, &c)
	}
}

// DoPost handles POST requests.
func doPost(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
	entity := new(Event)
	if err := readEntity(r, entity); err != nil {
		handleError(w, err, &c)
	}

	if err := assertValidKind(entity.Kind, c); err != nil {
		handleError(w, err, &c)
	}

	if entity.Date.IsZero() {
		entity.Date = time.Now()
	}

	entity.User = u.Email

	key := datastore.NewIncompleteKey(c, EVENT_KIND, nil)
	k, err := datastore.Put(c, key, entity)
	if err != nil {
		handleError(w, err, &c)
	}

	resp := CrudResponse{true, k.String(), "", nil}
	text, err := json.Marshal(resp)
	if err != nil {
		handleError(w, err, &c)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(text)
}

// HandleError returns an error for a given HTTP response.
func handleError(w http.ResponseWriter, err error, c *appengine.Context) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// ReadEntity reads a JSON value into entity from a Request body.
// An error is returned if the body cannot be read into entity.
func readEntity(r *http.Request, entity interface{}) error {
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

// ApplyFilters adds filters to query based on the values found in qParams, and potentially
// the current user's Admin role.
func applyFilters(qParams *QParams, query *datastore.Query, c appengine.Context, u *user.User) *datastore.Query {

	query = query.Filter("Kind =", qParams.Kind)
	if user.IsAdmin(c) && qParams.User != "" {
		query = query.Filter("User =", qParams.User)
	} else {
		query = query.Filter("User =", u.Email)
	}

	if qParams.Amount != 0 {
		query = query.Filter("Amount", qParams.Amount)
	}
	if qParams.Comment != "" {
		query = query.Filter("Comment", qParams.Comment)
	}
	if !qParams.MinDate.IsZero() {
		query = query.Filter("Date >=", qParams.MinDate)
	}
	if !qParams.MaxDate.IsZero() {
		query = query.Filter("Date <=", qParams.MaxDate)
	}

	return query
}

// ApplySort adds a sort order to query, based on a parameter found in qParams, if any.
func applySort(qParams *QParams, query *datastore.Query, c appengine.Context) *datastore.Query {
	if qParams.Sort != "" {
		sortField := ""
		desc := false

		if strings.HasPrefix(qParams.Sort, "-") {
			desc = true
			sortField = string([]rune(qParams.Sort)[1:])
		} else {
			sortField = qParams.Sort
		}

		sortField = strings.ToUpper(string([]rune(sortField[:1]))) +
			strings.ToLower(string([]rune(sortField[1:])))

		if desc {
			sortField = "-" + sortField
		}
		c.Infof("Sorting query by %v", sortField)
		query = query.Order(sortField)
	}

	return query
}

// ReadQParams examines request parameters for specific values and stores them
// in qParams.
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

	if r.FormValue("category") != "" {
		qParams.Category = r.FormValue("category")
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

	/*
	 * Read this parameter if present, but will only be used if an admin
	 * is making the request. Normal users will automatically have their
	 * queries filtered to their own events.
	 */
	qParams.User = r.FormValue("user")

	qParams.Sort = r.FormValue("_sort")

	return nil
}

// AssertValidKind verifies that the a 'kind' parameter matches one of
// the valid event types that the application can handle.
func assertValidKind(kind string, c appengine.Context) error {
	if kind != EXPENSE && kind != INCOME {
		return errors.New("Invalid kind: '" + kind + "'")
	} else {
		c.Infof("Data has kind %s", kind)
		return nil
	}
}

func isValidUser(appUser *AppUser) bool {
	return appUser.Id != ""
}

// Authorize verifies that the user making the request should be allowed
// to continue. If the user is not authorized, an error is returned with
// a nil user value.
func authorize(c appengine.Context) (*user.User, error) {
	if u := user.Current(c); u == nil {
		return nil, errors.New("User not logged in!")
	} else {
		return u, nil
	}
}

type Event struct {
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
	Comment  string    `json:"comment"`
	Date     time.Time `json:"date"`
	Kind     string    `json:"kind"`
	User     string    `json:"user"`
}

type QParams struct {
	User     string    `json:"user"`
	Kind     string    `json:"kind"`
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
	Comment  string    `json:"comment"`
	MaxDate  time.Time `json:"before"`
	MinDate  time.Time `json:"after"`
	Sort     string    `json:"_sort"`
}

type AppUser struct {
	Id               string   `json:"id"`
	StartOfWeek      string   `json:"startOfWeek"`
	CustomCategories []string `json:"customCat"`
}

type UserResponse struct {
	Ok   bool    `json:"ok"`
	Data AppUser `json:"data"`
}

type CrudResponse struct {
	Ok    bool    `json:"ok"`
	Key   string  `json:"key"`
	Error string  `json:"error"`
	Data  []Event `json:"data"`
}
