To start the project:

docker-compose up

Curl commands to test:

# Sign up
curl -X POST http://localhost:8080/signup -d '{"email": "user@example.com", "password": "mypassword"}'

# Sign in
curl -X POST http://localhost:8080/signin -d '{"email": "user@example.com", "password": "mypassword"}'

# Access protected route
curl -H "Authorization: Bearer <token>" http://localhost:8080/protected

# Revoke token
curl -X POST -H "Authorization: Bearer <token>" http://localhost:8080/revoke

# Refresh token
curl -X POST -H "Authorization: Bearer <token>" http://localhost:8080/refresh


