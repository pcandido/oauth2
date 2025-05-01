# OAuth2 Example

This project demonstrates a fully functional implementation of the **OAuth 2.0 Authorization Code Flow**, built with **Go** and structured as three independent services:

- 👤 **Client**: A sample application that need user authorization to access protected resources.
- 🧠 **Authorization Server**: Handles authentication, consent, and token generation.
- 🔐 **Resource Server**: An API that serves protected resources.

---

## 📁 Project Structure

```text
oauth2/
├── authorization-server/   # Handles login, consent, and token issuance
├── client/                 # Example client application initiating auth flow
├── resource-server/        # API with protected resources
├── docker-compose.yml      # To run all services together
```

---

## 🚀 Getting Started

### Prerequisites

- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/)

### Environment Setup

1. Clone the repository:

```bash
git clone https://github.com/your-username/oauth2.git
cd oauth2
```

2. Start all services:

```bash
docker-compose up --build
```

---
## 🔐 OAuth 2.0 Authorization Code Flow

This system implements the **OAuth 2.0 Authorization Code Flow**, ideal for secure applications that need to access protected resources on behalf of a user. The diagram below illustrates the entire interaction:

```mermaid
sequenceDiagram
    participant User
    participant Client
    participant AuthorizationServer
    participant ResourceServer

    User->>Client: Accesses the application
    alt User doesn't have an authorization token
        Client->>User: Redirects to /authorize with client_id, redirect_uri, scope, and state
        User->>AuthorizationServer: Request authorization with client_id, redirect_uri, scope, and state
        alt User is not logged in
            AuthorizationServer->>User: Requests login
            User->>AuthorizationServer: Logs in
        end
        AuthorizationServer->>User: Requests consent
        User->>AuthorizationServer: Grants permission
        AuthorizationServer->>User: Redirects to redirect_uri with authorization_code and state
        User->>Client: Calls redirect_uri with authorization_code and state
        Client->>AuthorizationServer: Makes POST to /token with authorization_code, client_id, client_secret, redirect_uri
        AuthorizationServer->>Client: Returns access_token and refresh_token
    end
    Client->>ResourceServer: Makes request with access_token
    ResourceServer->>Client: Returns protected data
    Client->>User: Displays protected data
```

---

### 1. User accesses the application

The user opens the Client App. If no access token is found, the OAuth flow is initiated.

---

### 2. Client redirects to Authorization Server

The client redirects the user to the `/authorize` endpoint on the Authorization Server with:

- `client_id` – pre-registered on the authorization server
- `redirect_uri` – the URL to redirect the user after authorization. This endpoint must handle the authorization response and must also be pre-registered for security reasons, to avoid open redirect attacks.
- `scope` – the actions or permissions requested. All possible scopes should be pre-defined and agreed between client and server.
- `state` – a CSRF protection mechanism. A random value should be generated and persisted securely (in our case, using cookies), and must be compared on the callback to ensure integrity.
- `response_type` – indicates the flow type (here, `code` for Authorization Code flow)

Example:

```
GET /authorize?response_type=code
  &client_id=example-client
  &redirect_uri=http://localhost:8081/callback
  &scope=read
  &state=random123
```

---

### 3. Authorization Server handles login

If the user is not already authenticated, the Authorization Server displays a login form and verifies credentials.

---

### 4. Authorization Server asks for consent

After login, the server asks the user to grant permission for the requested scopes. If granted, the flow continues.

---

### 5. Server redirects with an authorization code

The server redirects the user to the `redirect_uri` with an authorization code and the original `state`.

```
HTTP 302 Found
Location: http://localhost:8081/callback?code=abc123&state=random123
```

If the user denies consent or an error occurs, an `error` parameter will be included in the redirect URI.

---

### 6. Client exchanges code for access token

The client sends a `POST` request to the `/token` endpoint with:

- `client_id` – pre-registered
- `client_secret` – sent securely via backend (not exposed to user). This verifies the authenticity of the client.
- `redirect_uri` – to validate the flow and match the original authorization request
- `authorization_code` – issued by the Authorization Server and stored with context needed to complete the flow
- `grant_type=authorization_code` – specifies the type of OAuth flow

Response:

```json
{
  "access_token": "access-token",
  "refresh_token": "refresh-token",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

The access token is signed using an asymmetric key, allowing other services to validate the token without being able to forge new ones.

---

### 7. Client accesses the protected resource

Using the access token, the client sends a request to the Resource Server:

```
GET /resource
Authorization: Bearer access-token
```

---

### 8. Resource Server validates token and returns data

The Resource Server validates the token using the public key. If valid, it returns the protected data to the client, which is then displayed to the user.

---

## 📌 Endpoints Overview

### Authorization Server

| Method | Path         | Description                   |
|--------|--------------|-------------------------------|
| GET    | `/authorize` | Start OAuth flow              |
| POST   | `/token`     | Exchange code for token       |
| GET    | `/login`     | User login                    |
| POST   | `/consent`   | User consent confirmation     |

---

## 📜 License

MIT License
