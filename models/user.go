package models

import (
    "time"
    "database/sql/driver"
    "encoding/json"
)

// User model
type User struct {
    ID             uint                 `gorm:"primaryKey"`
    Username       string               `gorm:"unique;not null"`
    Email          string               `gorm:"unique;not null"`
    Name           string               `gorm:"not null"`
    Gender         string               `gorm:"not null"`
    Birthday       time.Time            `gorm:"type:date"`
    PasswordHash   string               `gorm:"not null"`
    IsPremium      bool                 `gorm:"default:false"`
    Verified       bool                 `gorm:"default:false"`
    Seen           CollectionSwipeData  `gorm:"type:jsonb;default:'[]';not null"`
    SwipeCount     int                  `gorm:"type:int;default:0;not null"`
    LastSwipe      time.Time            `gorm:"type:timestamp;default:0"`
}

type Profile struct {
    Username       string               `gorm:"unique;not null"`
    Email          string               `gorm:"unique;not null"`
    Name           string               `gorm:"not null"`
    Birthday       time.Time            `gorm:"type:date"`
}

// MapUserToProfile converts a single User to Profile
func MapUserToProfile(user User) Profile {
    return Profile{
        Username: user.Username,
        Email:    user.Email,
        Name:     user.Name,
        Birthday: user.Birthday,
    }
}

// MapUsersToProfiles converts a slice of User structs to a slice of Profile structs
func MapUsersToProfiles(users []User) []Profile {
    profiles := make([]Profile, len(users))
    for i, user := range users {
        profiles[i] = MapUserToProfile(user)
    }
    return profiles
}

type SwipeData struct {
    By  uint        `json:"by"`
    At  time.Time   `json:"at"`
}

type CollectionSwipeData []SwipeData

// Handle JSON type for Insert Swipe Data
func (swipe CollectionSwipeData) Value() (driver.Value, error) {
    return json.Marshal(swipe) // Convert CollectionSwipeData to JSON []byte
}

// handle JSON Input for find or projection Swipe Data
func (swipe *CollectionSwipeData) Scan(value interface{}) error {
    if value == nil {
        *swipe = CollectionSwipeData{}
        return nil
    }
    // Convert the JSON data from []byte to CollectionSwipeData
    return json.Unmarshal(value.([]byte), swipe)
}