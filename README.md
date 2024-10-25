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


Note :
After hitting Signin API copy the token and provide that token in the authorization as a bearer token in access protected route, revoke token, and refresh token.

