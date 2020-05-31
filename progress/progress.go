package main

import (
	"database/sql"
	"encoding/json"
	"info441sp20-ashraysa/gateway/models/users"
	"net/http"
	"strings"
)

// Progress represents a Progress struct
type Progress struct {
	ProgressID int   `json:"progressID"`
	DaysSober  int   `json:"daysSober"`
	UserID     int64 `json:"userID"`
}

// MySQLStore represents a MySQL store
type MySQLStore struct {
	db *sql.DB
}

//NewMySQLStore constructs and returns a new MySQLStore
func NewMySQLStore(DB *sql.DB) *MySQLStore {
	return &MySQLStore{
		db: DB,
	}
}

// ProgressUserHandler tracks the user's sobriety clock and awards points for each day the user clocks-in
func (msq *MySQLStore) ProgressUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-User") != "" {
		decoder := json.NewDecoder(strings.NewReader(r.Header.Get("X-User")))
		user := &users.User{}
		err := decoder.Decode(user)
		if err != nil {
			http.Error(w, "error decoding response body", http.StatusBadRequest)
			return
		}
		// If method is GET
		// If no entry for user, add progrss entry with 0 (port over from pATCH)
		// return the progress struct in response body
		// status OK

		// If method is PATCH
		// Update entry with daysSober + 1
		// return progress struct in response body
		// status oK

		// Returning:
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusOK)
		// enc := json.NewEncoder(w)
		// enc.Encode(user)
		progress := &Progress{}

		if r.Method == "GET" {
			sqlQuery := "select daysSober from Progress where userID = ?"
			res, err := msq.db.Query(sqlQuery, user.ID)
			if err != nil {
				sqlQueryTwo := "insert into Progress(daysSober, userID) values (?, ?)"
				_, errTwo := msq.db.Exec(sqlQueryTwo, 1, user.ID)
				if errTwo != nil {
					http.Error(w, errTwo.Error(), http.StatusInternalServerError)
					return
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			enc := json.NewEncoder(w)
			enc.Encode(progress)
		} else if r.Method == "PATCH" {
			sqlQuery := "select daysSober from Progress where userID = ?"
			res, err := msq.db.Query(sqlQuery, user.ID)
			if err != nil {
				http.Error(w, "User has not logged any days in the sobriety clock", http.StatusBadRequest)
				return
			}
			for res.Next() {
				res.Scan(&progress.DaysSober)
			}

			/*sqlQueryTwo := "insert into Progress(daysSober, userID) values (?, ?)"
			if progress.DaysSober == 0 {
				_, errTwo := msq.db.Exec(sqlQueryTwo, 1, user.ID)
				if errTwo != nil {
					http.Error(w, errTwo.Error(), http.StatusInternalServerError)
					return
				}
			}*/

			// update daysSober and update points for user
			sqlQueryThree := "update Progress set daysSober = ? where userID = ?"
			_, errThree := msq.db.Exec(sqlQueryThree, progress.DaysSober+1, user.ID)
			if errThree != nil {
				http.Error(w, "Error updating daySober for current user", http.StatusInternalServerError)
				return
			}
			sqlQueryFour := "update Users set points = ? where id = ?"
			_, errFour := msq.db.Exec(sqlQueryFour, user.Points+100, user.ID)
			if errFour != nil {
				http.Error(w, "Error adding points for the current user", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			enc := json.NewEncoder(w)
			enc.Encode(progress)
		} else {
			http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
			return
		}
	} else {
		http.Error(w, "User is not authenticated", http.StatusUnauthorized)
		return
	}
}
