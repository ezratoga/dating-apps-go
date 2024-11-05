package controllers

import (
    "dating-apps-go/database"
    "dating-apps-go/models"
    "dating-apps-go/config"
    "net/http"
    "github.com/labstack/echo/v4"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt"
    "time"
    "fmt"
)

// SignUp handles user registration
func SignUp(header echo.Context) (err error) {
    user := new(models.SignUpUser) // allocates memory for data type and return pointer of the operand
    userToRegistered := new(models.User)
    err = header.Bind(user); 
    if err != nil {
        return header.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
    }
    // Validate the request
    if err := header.Validate(user); err != nil {
        return header.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed", "details": err.Error()})
    }

    // Parse the birthday string to time.Time format
    birthday, err := time.Parse("2006/01/02", user.Birthday)
    if err != nil {
        return header.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format. Use YYYY-MM-DD."})
    }

    // define all users-column 
    userToRegistered.Username = user.Username
    userToRegistered.Email = user.Email
    userToRegistered.Name = user.Name
    userToRegistered.Gender = user.Gender
    userToRegistered.Birthday = birthday

    // Hash password and save user
    pass := []byte(user.Password)
    hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
    if err != nil {
        return header.JSON(http.StatusBadRequest, map[string]string{"error": "Error creating user. Please try again, or call our Customer Service"})
    }

    userToRegistered.PasswordHash = string(hashedPassword)

    if err := database.DB.Create(&userToRegistered).Error; err != nil {
        return header.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    // return header.JSON(http.StatusOK, user)
    return header.JSON(http.StatusCreated, "User registered successfully")
}

func CreateToken(userId uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userID"] = userId
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix() // Token expires after 1 hour
    fmt.Println(claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(globalconfig.GetEnvVariable("JWT_KEY")))
}

// Login handles user login and returns a JWT token
func Login(header echo.Context) (err error) {
    // Implementation for login with JWT token generation
    user := new(models.LoginUser) // allocates memory for data type and return pointer of the operand
    userData := new(models.User)
    err = header.Bind(user); 
    if err != nil {
        return header.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
    }

    if err := database.DB.Where("username = ?",user.Username).First(&userData).Error; err != nil {
        return header.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong username or password"})
    }

    jwtToken, err := CreateToken(userData.ID)
    if err != nil {
        return header.JSON(http.StatusInternalServerError, err.Error())
    }

    // fmt.Println(userData.PasswordHash);

    if err := bcrypt.CompareHashAndPassword([]byte(userData.PasswordHash), []byte(user.Password)); err != nil {
		return header.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong username or password"})
	}

    // return header.JSON(http.StatusOK, user)
    return header.JSON(http.StatusCreated, map[string]interface{}{
        "token": jwtToken,
        "message": "Login successfully",
    })
}
