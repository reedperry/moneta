package moneta

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Expense struct {
	Amount  int64     `json:"amount"`
	Comment string    `json:"comment"`
	Date    time.Time `json:"date"`
}

func init() {
	http.HandleFunc("/data", crud)
}

func crud(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		data, err := readRequest(r.Body)
		if err == nil {
			err = json.NewEncoder(w).Encode(data)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "GET":
		fmt.Fprint(w, "Hello!")
		return
	}
}

func readRequest(r io.ReadCloser) (*Expense, error) {
	defer r.Close()
	var expense Expense
	err := json.NewDecoder(r).Decode(&expense)
	return &expense, err
}
