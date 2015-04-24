package sessions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/securecookie"
	"github.com/mattbostock/timeoff/model"
	"github.com/nbio/httpcontext"
)

var log *logrus.Logger

type Manager struct {
	sessionName string
	hashKey     []byte
}

func New(name string, hashKey []byte) *Manager {
	return &Manager{name, hashKey}
}

func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(m.sessionName)
	if err != nil {
		return
	}

	s := securecookie.New(m.hashKey, m.hashKey)

	var userID uint64
	err = s.Decode(m.sessionName, cookie.Value, &userID)
	if err != nil {
		log.Debugln(err)
		return
	}

	if userID < 1 {
		log.Errorf("User ID with value %v found in cookie", userID)
		return
	}

	user, err := model.FindUser(userID)
	if err != nil {
		log.Errorf("Cookie user ID %v not found in database", userID)
		return
	}

	httpcontext.Set(r, "user", user)
}

func (m *Manager) SetCookie(w http.ResponseWriter, value interface{}) error {
	// FIXME Consider using non-deterministic ID for added security
	s := securecookie.New(m.hashKey, m.hashKey)
	encoded, err := s.Encode(m.sessionName, value)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not encode cookie with value %v: %s", value, err))
	}

	// Older browsers don't support MaxAge, but it just means the session
	// will end when the browser is closed. That's fine.
	cookie := &http.Cookie{
		Name:     m.sessionName,
		Value:    encoded,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	return nil
}

func (m *Manager) Logout(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     m.sessionName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // delete cookie now
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func SetLogger(l *logrus.Logger) {
	log = l
}
