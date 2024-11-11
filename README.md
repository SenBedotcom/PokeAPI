# Pokémon API

A simple RESTful API that allows users to interact with Pokémon data. It supports registration, login, fetching specific Pokémon details, abilities, and random Pokémon. JWT authentication is required for certain routes.

## Features
- **User Authentication**: Users can register and log in with JWT authentication.
- **Pokémon Data**: Fetch data for specific Pokémon by name, their abilities, and a random Pokémon.
- **JWT Authentication**: All protected routes require a valid JWT token.

## Project Structure

pokemon-api/
├── main.go                      # Main entry point for the application
├── config/                      # Configuration files
│   └── config.go                # Configuration settings (e.g., JWT secret)
├── controllers/                 # API controllers handling request logic
│   ├── authController.go        # Handles user authentication
│   └── pokemonController.go     # Handles Pokémon data requests
├── models/                      # Data models (e.g., User)
│   └── user.go                  # User model for authentication
├── services/                    # Business logic for the app
│   ├── authService.go           # Handles user-related business logic
│   ├── pokemonService.go        # Handles Pokémon data fetching
│   └── cacheService.go          # Caching logic for Pokémon data
├── middleware/                  # Middleware for authentication and other checks
│   └── authMiddleware.go        # JWT middleware
├── utils/                       # Utility functions
│   ├── hash.go                  # Hash utilities for password hashing
│   └── jwt.go                   # JWT utilities for token creation/validation
├── Dockerfile                   # Docker configuration
├── docker-compose.yml           # Docker Compose configuration for services
└── go.mod                       # Go modules dependencies

## Prerequisites

- Go 1.18+
- Docker (optional, for running the app in containers)

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/yourusername/pokemon-api.git
    cd pokemon-api
    ```

2. Install Go dependencies:

    ```bash
    go mod tidy
    ```

3. Set up environment variables:
    - Create a `.env` file or export the required environment variables:
      - `JWT_SECRET`: Secret key for JWT token generation

4. Run the application:

    ```bash
    go run main.go
    ```

5. Alternatively, you can use Docker to run the application:
   
    ```bash
    docker-compose up --build
    ```

## Endpoints

### Authentication

- **POST `/register`**: Register a new user.
    - Request body:
      ```json
      {
        "username": "example",
        "password": "password123"
      }
      ```

- **POST `/login`**: Log in with existing user credentials to receive a JWT token.
    - Request body:
      ```json
      {
        "username": "example",
        "password": "password123"
      }
      ```
    - Response:
      ```json
      {
        "token": "jwt_token_here"
      }
      ```

### Pokémon Routes (Authenticated)

- **GET `/pokemon/{name}`**: Fetch Pokémon details by name.
    - Headers:
      ```json
      {
        "Authorization": "Bearer {jwt_token_here}"
      }
      ```

- **GET `/pokemon/{name}/ability`**: Fetch Pokémon abilities by name.
    - Headers:
      ```json
      {
        "Authorization": "Bearer {jwt_token_here}"
      }
      ```

- **GET `/pokemon/random`**: Fetch a random Pokémon.
    - Headers:
      ```json
      {
        "Authorization": "Bearer {jwt_token_here}"
      }
      ```

### Error Handling

- **401 Unauthorized**: If the JWT token is missing or invalid, a `401 Unauthorized` error is returned.
- **400 Bad Request**: If the request is malformed or required data is missing, a `400 Bad Request` error is returned.

## Environment Variables

You need to configure the following environment variables for the application:

- `JWT_SECRET`: A secret key used to sign JWT tokens (e.g., `mysecretkey`).

You can set these variables in a `.env` file, or export them directly.

## Testing

1. **Unit Tests**: You can write tests for your services and controllers using Go testing framework.
2. **Integration Tests**: You can integrate with tools like Postman or Swagger to test your API endpoints.

## Docker Usage

To run the application using Docker, make sure you have Docker installed. You can build and run the container with the following commands:

1. **Build and start the containers**:
    ```bash
    docker-compose up --build
    ```

2. **Stop the containers**:
    ```bash
    docker-compose down
    ```

## License

This project is licensed under the MIT License.


