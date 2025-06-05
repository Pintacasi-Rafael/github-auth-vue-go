# GitHub OAuth Login - Go + Vue Fullstack App

A full-stack project implementing GitHub OAuth authentication using:

- **Backend**: Go (Standard Library Only)
- **Frontend**: Vue 3 + Vite
- **Database**: MySQL
- **Authentication**: GitHub OAuth + JWT

---

## Features

- GitHub Login with OAuth 2.0
- JWT token generation and verification
- JWT middleware (with expiration check)
- User data stored in MySQL
- Vue frontend with protected landing page

---

---

## Backend Setup (Go)

1. Go to the `back` folder:

```bash
cd back
go mod tidy
go run main.go
```

2. Configure the variables inside main.go

```bash
var (
    clientID     = "YOUR_GITHUB_CLIENT_ID"
    clientSecret = "YOUR_GITHUB_CLIENT_SECRET"
    jwtSecret    = []byte("your_jwt_secret_key")
    redirectURI  = "http://localhost:8080/github/callback"
)
```

3. Create your MySQL githubUsers

```bash
CREATE TABLE IF NOT EXISTS githubUsers (
    id BIGINT PRIMARY KEY,
    login VARCHAR(100) NOT NULL,
    name VARCHAR(255),
    avatar_url VARCHAR(255),
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## Frontend Setup(Vue + Vite)

1. Go to the front folder

```bash
cd front
npm install
```

2. Start the dev server

```bash
npm run dev
```

3. App runs at: http://localhost:5173

- Clicking "Login with GitHub" will redirect to GitHub OAuth
- On success, user is redirected to /landing page with token in the URL
