package database

import (
    "dating-apps-go/models"
    "dating-apps-go/config"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "fmt"
)

var DB *gorm.DB

func InitDB() {
    var err error

    dsn := fmt.Sprint("host=",globalconfig.GetEnvVariable("POSTGRE_HOST")," user=",globalconfig.GetEnvVariable("POSTGRE_USER")," password=",globalconfig.GetEnvVariable("POSTGRE_PASSWORD"),        " dbname=",globalconfig.GetEnvVariable("POSTGRE_DATABASE")," sslmode=",globalconfig.GetEnvVariable("POSTGRE_SSL"))
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database")
    }
    fmt.Println("Connection to database postgre sucess")

    // DB.AutoMigrate(&models.User{}, &models.Profile{}, &models.Swipe{})
    DB.AutoMigrate(&models.User{}, &models.Swipe{}, &models.Invoice{})
}
