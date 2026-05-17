package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var oauthConfig *oauth2.Config
var allowedUser string

func InitOAuth() {
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"read:user"},
		Endpoint:     github.Endpoint,
	}
	allowedUser = os.Getenv("GITHUB_ALLOWED_USER")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if user := GetUser(r); user != nil {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	if oauthConfig.ClientID == "" {
		http.Error(w, "GitHub OAuth not configured", http.StatusInternalServerError)
		return
	}
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("OAuth exchange error: %v", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	user, err := fetchGitHubUser(token)
	if err != nil {
		log.Printf("Failed to fetch GitHub user: %v", err)
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}

	if allowedUser != "" && user.Username != allowedUser {
		log.Printf("Unauthorized login attempt by: %s", user.Username)
		http.Error(w, "Unauthorized: You are not allowed to access the admin panel", http.StatusForbidden)
		return
	}

	if err := SetUser(w, r, *user); err != nil {
		log.Printf("Failed to set session: %v", err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	_ = ClearUser(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type githubUserResponse struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

func fetchGitHubUser(token *oauth2.Token) (*SessionUser, error) {
	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ghUser githubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return nil, err
	}

	return &SessionUser{
		GitHubID:  ghUser.ID,
		Username:  ghUser.Login,
		AvatarURL: ghUser.AvatarURL,
	}, nil
}
