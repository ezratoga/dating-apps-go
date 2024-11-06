### Dating App Backend
This is the backend system that may be implemented, built with Go and the Echo framework.

### Features
You can register your account, login then swipe the profiles that 

### Tech Stack
Go (Golang): Core backend language
Echo: Web framework for routing and middleware
PostgreSQL: Database for storing user profiles, matches, and messages
JWT: JSON Web Tokens for secure user authentication
Docker: Containerization for easy deployment

### Environment Variables
This is the key that yo can declare in .env file:

Variable	Description
DB_HOST	    Database host address
DB_PORT	    Database port
DB_USER	    Database username
DB_PASSWORD	Database password
JWT_KEY	    Secret key for signing JWT tokens

### Setup and Installation

Clone the repository:
```shell
$ git clone https://github.com/your-username/dating-app-backend.git
$ cd dating-app-backend
```

Install dependencies:
```shell
$ go mod download
```

Set up the environment variables: Create a .env file in the root directory and add the required environment variables as listed above.

Run the application in your local cmd:
```shell
$ go run main.go
```

Run application via docker image:
```shell
$ docker build -t dating-app-go .
$ docker run -p 8080:8080 --env ${ENV_KEY}=${ENV_VALUE} .env dating-app-go
```
