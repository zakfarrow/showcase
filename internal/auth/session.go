package auth

import (
	"encoding/gob"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const (
	sessionName    = "showcase-session"
	userSessionKey = "user"
)

var store *sessions.CookieStore

func InitStore() {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production-32chars"
	}
	store = sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	}
}

type SessionUser struct {
	GitHubID  int64
	Username  string
	AvatarURL string
}

func GetUser(r *http.Request) *SessionUser {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return nil
	}
	user, ok := session.Values[userSessionKey].(SessionUser)
	if !ok {
		return nil
	}
	return &user
}

func SetUser(w http.ResponseWriter, r *http.Request, user SessionUser) error {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Values[userSessionKey] = user
	return session.Save(r, w)
}

func ClearUser(w http.ResponseWriter, r *http.Request) error {
	session, _ := store.Get(r, sessionName) // Ignore error - clear anyway
	// Clear all session values
	for key := range session.Values {
		delete(session.Values, key)
	}
	// Set MaxAge to -1 to delete the cookie
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	}
	return session.Save(r, w)
}

func init() {
	// Register the SessionUser type for gob encoding
	gob.Register(SessionUser{})
}
