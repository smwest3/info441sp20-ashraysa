package handlers

import (
	"encoding/json"
	"info441sp20-ashraysa/gateway/models/users"
	"info441sp20-ashraysa/gateway/sessions"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UsersHandler handles requests for the "users" resource
func (c *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		decoder := json.NewDecoder(r.Body)
		newUser := &users.NewUser{}
		err := decoder.Decode(newUser)
		if err != nil {
			http.Error(w, "error decoding response body", http.StatusBadRequest)
			return
		}
		errTwo := newUser.Validate()
		if errTwo != nil {
			http.Error(w, errTwo.Error(), http.StatusBadRequest)
			return
		}
		user, _ := newUser.ToUser()
		var errFour error
		user, errFour = c.UserStore.Insert(user)
		if errFour != nil {
			http.Error(w, errFour.Error(), http.StatusBadRequest)
			return
		}
		newSession := &SessionState{time.Now(), user}
		sessions.BeginSession(c.Key, c.SessionStore, newSession, w)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		enc.Encode(user)
	} else {
		http.Error(w, "request method is not POST", http.StatusMethodNotAllowed)
		return
	}
}

// SpecificUserHandler handles requests for a specific user
func (c *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	currSession := &SessionState{}
	sesID, err := sessions.GetState(r, c.Key, c.SessionStore, currSession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	currID := currSession.User.ID
	urlBase := path.Base(r.URL.String())
	id, _ := strconv.ParseInt(urlBase, 10, 64)

	if r.Method == "GET" {
		var user *users.User
		var er error
		if urlBase != "me" {
			user, er = c.UserStore.GetByID(id)
			if er != nil {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
		} else {
			user = currSession.User
			user, er = c.UserStore.GetByID(user.ID)
			if er != nil {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.Encode(user)

	} else if r.Method == "PATCH" {
		if urlBase != "me" && id != currID {
			http.Error(w, "user is not authenticated", http.StatusForbidden)
			return
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		if urlBase == "me" {
			id = currSession.User.ID
		}
		decoder := json.NewDecoder(r.Body)
		newUpdate := &users.Updates{}
		decoder.Decode(newUpdate)
		errFive := currSession.User.ApplyUpdates(newUpdate)
		if errFive != nil {
			http.Error(w, errFive.Error(), http.StatusNotFound)
			return
		}
		c.SessionStore.Save(sesID, currSession)
		user, errFour := c.UserStore.Update(id, newUpdate)
		if errFour != nil {
			http.Error(w, errFour.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.Encode(user)
	} else {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

// SessionsHandler handles requests for the "sessions" resource, and allows
// clients to begin a new session using an existing user's credentials.
func (c *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		decoder := json.NewDecoder(r.Body)
		newCredentials := &users.Credentials{}
		decoder.Decode(newCredentials)
		user, err := c.UserStore.GetByEmail(newCredentials.Email)
		if err != nil {
			// PASSING SAME TIME AS AUTHENTICATE
			bcrypt.GenerateFromPassword([]byte("doesnt matter"), 13)
			// PASSING TIME DONE

			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		errThree := user.Authenticate(newCredentials.Password)
		if errThree != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		newSession := &SessionState{time.Now(), user}
		sessions.BeginSession(c.Key, c.SessionStore, newSession, w)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		enc.Encode(user)
	} else {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

// SpecificSessionHandler requests related to a specific authenticated session
func (c *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		if path.Base(r.URL.String()) != "mine" {
			http.Error(w, "resource path invalid", http.StatusForbidden)
			return
		}
		_, err := sessions.EndSession(r, c.Key, c.SessionStore)
		if err != nil {
			http.Error(w, "resource path invalid", http.StatusNotFound)
			return
		}
		w.Write([]byte("signed out"))
	} else {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}
}
