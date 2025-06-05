package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
)

var (
    clientID     = "Ov23liIfbakE2X34uWTR"
    clientSecret = "01dec7d08c1cc283b78d814d057c469fab6d3f27"
    jwtSecret    = []byte("your_jwt_secret_key")
    redirectURI  = "http://localhost:8080/github/callback"
    db           *sql.DB
)

type GitHubUser struct {
    ID        int64  `json:"id"`
    Login     string `json:"login"`
    Name      string `json:"name"`
    AvatarURL string `json:"avatar_url"`
    Email     string `json:"email"`
}

func main() {
    // Initialize DB connection
    err := initDB()
    if err != nil {
        log.Fatalf("Failed to connect to DB: %v", err)
    }
    defer db.Close()

    // Register handlers
    http.HandleFunc("/github/callback", githubCallbackHandler)
	http.HandleFunc("/github/login", githubLoginHandler)

    // Start server
    log.Println("Server started at :8080")
    err = http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}

func initDB() error {
    dsn := "root:12345678@tcp(127.0.0.1:3306)/test_app_github_auth?parseTime=true"
    var err error
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    // Test DB connection
    return db.Ping()
}

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
    githubAuthURL := "https://github.com/login/oauth/authorize"
    params := url.Values{}
    params.Set("client_id", clientID)
    params.Set("redirect_uri", redirectURI)
    params.Set("scope", "read:user user:email") // adjust scopes as needed

    http.Redirect(w, r, githubAuthURL+"?"+params.Encode(), http.StatusFound)
}

func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "Missing code in query", http.StatusBadRequest)
        return
    }

    // Exchange code for access token
    accessToken, err := exchangeCodeForToken(code)
    if err != nil {
        http.Error(w, "Failed to get access token: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Fetch GitHub user info
    user, err := fetchGitHubUser(accessToken)
    if err != nil {
        http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Save or update user in DB
    err = saveOrUpdateUser(user)
    if err != nil {
        log.Println("DB error:", err)
        http.Error(w, "Failed to save user info", http.StatusInternalServerError)
        return
    }

    // Generate JWT token
    jwtToken, err := generateJWT(user.Login)
    if err != nil {
        http.Error(w, "Failed to generate JWT: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Redirect with JWT token in query string
    redirectURL := fmt.Sprintf("http://localhost:5173/landing?token=%s", url.QueryEscape(jwtToken))
    http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func exchangeCodeForToken(code string) (string, error) {
    data := url.Values{}
    data.Set("client_id", clientID)
    data.Set("client_secret", clientSecret)
    data.Set("code", code)
    data.Set("redirect_uri", redirectURI)

    req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
    if err != nil {
        return "", err
    }
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var res struct {
        AccessToken string `json:"access_token"`
        Scope       string `json:"scope"`
        TokenType   string `json:"token_type"`
        Error       string `json:"error"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
        return "", err
    }
    if res.Error != "" {
        return "", fmt.Errorf("oauth error: %s", res.Error)
    }
    return res.AccessToken, nil
}

func fetchGitHubUser(accessToken string) (*GitHubUser, error) {
    req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "token "+accessToken)
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("GitHub API error: %s", string(body))
    }

    var user GitHubUser
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }
    return &user, nil
}

func saveOrUpdateUser(user *GitHubUser) error {
    query := `
    INSERT INTO githubUsers (id, login, name, avatar_url, email)
    VALUES (?, ?, ?, ?, ?)
    ON DUPLICATE KEY UPDATE
      login = VALUES(login),
      name = VALUES(name),
      avatar_url = VALUES(avatar_url),
      email = VALUES(email),
      updated_at = CURRENT_TIMESTAMP
    `
    _, err := db.Exec(query, user.ID, user.Login, user.Name, user.AvatarURL, user.Email)
    return err
}

func generateJWT(username string) (string, error) {
    claims := jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}
