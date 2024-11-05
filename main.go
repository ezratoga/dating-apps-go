package main

import (
    "dating-apps-go/controllers"
    "dating-apps-go/database"
    "dating-apps-go/config"
    // "dating-apps-go/middlewares"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "fmt"
    "github.com/go-playground/validator/v10"
)

// Validator instance
var validate *validator.Validate

// CustomValidator wraps the validator instance to integrate with Echo
type CustomValidator struct {
    validator *validator.Validate
}

// Validate method to satisfy Echo’s Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

func main() {
    // Initialize database
    database.InitDB()

    // Create Echo instance
    echoServer := echo.New()

    // Initialize validator
    validate = validator.New()
    echoServer.Validator = &CustomValidator{validator: validate}

    // Middleware for logging and error handling
    echoServer.Use(middleware.Logger())
    echoServer.Use(middleware.Recover())

    // Auth routes
    echoServer.POST("/signup", controllers.SignUp)
    echoServer.POST("/login", controllers.Login)

    fmt.Println(globalconfig.GetEnvVariable("JWT_KEY"))

    // Profile routes (protected by JWT middleware)
    profileGroup := echoServer.Group("/profiles")
    profileGroup.Use(middleware.JWT([]byte(globalconfig.GetEnvVariable("JWT_KEY"))))
    profileGroup.GET("", controllers.GetProfiles)
    profileGroup.POST("/swipe", controllers.SwipeProfile)
    profileGroup.PUT("/premium/upgrade", controllers.PurchasePremium)
    // profileGroup.PUT("/premium/upgrade", controllers.BuyPremium)


    // Start server
    var initializationPort string = fmt.Sprint(":",globalconfig.GetEnvVariable("PORT"))
    echoServer.Logger.Fatal(echoServer.Start(initializationPort))
}
