# oauth2

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