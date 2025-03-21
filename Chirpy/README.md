# Chirpy

## An API that handles Users Registration Login and Auth Tokens and their Chirpys

Chirpy is an API that implements the minimum requirement for a server. With this project it was possible to learn:

- Understand what web servers are and how they power real-world web applications
- Build a production-style HTTP server in Go, without the use of a framework
- Use JSON, headers, and status codes to communicate with clients via a RESTful API
- Learn what makes Go a great language for building fast web servers
- Use type safe SQL to store and retrieve data from a Postgres database
- Implement a secure authentication/authorization system with well-tested cryptography libraries
- Build and understand webhooks and API keys
- Document the REST API with markdown

## Installation and working with Chirpy

go install
Create a `.env` variable to set:

```
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=secret
POLKA_KEY=polka
```

## Endpoints

The application is handled via `/app/`. However, being this an API, it relies on specific URLs to handle the different usages:
- `GET /admin/metrics`: Retrieves the number of contacts to the API; 
- `POST /admin/reset`: :alarm: Resets the databases. Internal command used only for debugging purposes and that should never exist in production;
- `GET /api/healthz`: Checks the health of the webServer;
- `POST /api/users`: Allows to create a user based on a username and password;
- `PUT /api/users`: Update the email or password information of a authorized user. Requires an authorization token;
- `POST /api/login`: Log's in a user based on the user's Email and Password;
- `POST /api/chirps`: Allowes to post a Chirp for a validated user. The Chirp is passed as part of the request body and it requires the user's authorization token;
- `GET /api/chirps`: Retrieves all the Chirps currently available. It accepts also some queries
    - `?sort`: Reorders the Chirps in order of date of creation. It can either be `asc` or `desc`
    - `?author_id`: Retrieves the Chirps of a specific user. It requires the user uuid.
- `GET /api/chirps/{chirpID}`: Retrieves a specific Chirp based on its Id;
- `DELETE /api/chirps/{chirpID}`: Deletes a single Chirp based on the passed Id. It requires the User token as part of the Athorization http header;
- `POST /api/refresh`: Refreshes an authorization token. It returns to the user a new authorization token and it is based on a refreshing token;
- `POST /api/revoke`: Revokes the access token for a specific user;
- `POST /api/polka/webhooks`: This is a simulation for an external webhook that provides a user with premium access. This is handled via an ApiKey;

In general:
- Authorization tokens should be passed as part of the http.Header and should be of the form "Bearer user_id_auth_token"
- The body should instead contain in a Json form:
    - "email": "example@provider.com"
    - "password": "passwordtobechanged"
    - "body": "This is a Chirp shorter than 140 characters"

## The new horizons of cryptographic algorithms used for password hashing

Why using Bcrypt instead of SHA-256 or MD5 Hashing algorithms? [Why MD5 and SHA are outdeted](https://dev.to/lovestaco/hashing-passwords-why-md5-and-sha-are-outdated-and-why-you-should-use-scrypt-or-bcrypt-48p2)

