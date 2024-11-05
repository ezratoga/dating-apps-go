package controllers

import (
    "dating-apps-go/database"
    "dating-apps-go/models"
	"github.com/golang-jwt/jwt"
    "github.com/labstack/echo/v4"
    "net/http"
	"fmt"
	"time"
	"github.com/lib/pq"
	"slices"
)

// GetProfiles retrieves up to 10 profiles the user hasn’t seen today
func GetProfiles(header echo.Context) error {
    var profiles []models.Profile
	var users []models.User

	// var profile models.User
    authenticationHeader := header.Get("user").(*jwt.Token)
	userData := authenticationHeader.Claims.(jwt.MapClaims)
	// Retrieve userID from userData
    userID := uint(userData["userID"].(float64)) // Cast to uint


    if err := database.DB.Table("users").Select("*").Where("ID != ?", userID).Limit(10).Find(&users).Error; 
		err != nil {
		return header.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching profiles"})
    }

	profiles = models.MapUsersToProfiles(users)

	for index, value := range users {
		var seenBy = append(users[index].Seen, models.SwipeData{By:userID, At:time.Now()})
		if err := database.DB.Model(&value).Where("ID = ?", value.ID).Update("seen", seenBy).Error; err != nil {
			panic("Error save the profile history seen by")
		}
	}

    return header.JSON(http.StatusOK, profiles)
}

// Check if any match and return the struct to save to DB
func CheckMatchandReturn (whoLike []models.Swipe, swipeUser models.User, targetUser models.User) (models.Swipe, []uint) {
	var swipeMapping models.Swipe
	var whoLikeId []uint
	
	if len(whoLike) > 0 {
		var matchProfile pq.StringArray
		for _, value := range whoLike {
			matchProfile = append(matchProfile, value.Username)
			whoLikeId = append(whoLikeId, value.ID)
		}
	
		swipeMapping = models.Swipe{
			ID: swipeUser.ID,
			Username: swipeUser.Username,
			Like: pq.StringArray{targetUser.Username},
			Pass: pq.StringArray{},
			Match: matchProfile,
		}
	} else {
		swipeMapping = models.Swipe{
			ID: swipeUser.ID,
			Username: swipeUser.Username,
			Like: pq.StringArray{targetUser.Username},
			Pass: pq.StringArray{},
			Match: pq.StringArray{},
		}
	}

	return swipeMapping, whoLikeId
}

// SwipeProfile handles swiping on a profile
func SwipeProfile(header echo.Context) error {
	var SwipePayload models.SwipePayload
    if err := header.Bind(&SwipePayload); err != nil {
        return header.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
    }

	authenticationHeader := header.Get("user").(*jwt.Token)
	userData := authenticationHeader.Claims.(jwt.MapClaims)
	// Retrieve userID from userData
    userID := uint(userData["userID"].(float64)) // Cast to uint

    // Save Last Swipe
	// fmt.Println(userID)
	// fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	var user models.User
	var target models.User
	var swipe models.Swipe
	var swipeMapping models.Swipe
	var whoLikeId []uint
	var whoLike []models.Swipe

	if err := database.DB.Table("users").Model(&user).Where("ID = ?", userID).First(&user).Error; err != nil {
		return header.JSON(http.StatusNotFound, map[string]string{"messsage": "Your account is not valid user"})
	}

	if (user.SwipeCount > 1 && !user.IsPremium) {
		return header.JSON(http.StatusBadRequest, map[string]string{"messsage": "Your account is not allowed to swipe"})
	}

	user.SwipeCount = user.SwipeCount + 1
	user.LastSwipe = time.Now()

	if err := database.DB.Table("users").Model(&target).Where("ID = ?", SwipePayload.TargetId).First(&target).Error; err != nil {
		return header.JSON(http.StatusNotFound, map[string]string{"messsage": "Your target is not valid user"})
	}
	if err := database.DB.Table("users").Select("*").Where("ID = ?", userID).Updates(user).Error; err != nil {
		panic("Error save the profilelast swipe")
	}

	if err := database.DB.Table("swipes").Select("id, username, \"like\"").Where("? = ANY (\"like\")", user.Username).Find(&whoLike).Error; err != nil {
		fmt.Println("No one matches you, ", user.Name)
	}

	if err := database.DB.Table("swipes").Where("ID = ?", userID).First(&swipe).Error; err != nil {
		if SwipePayload.Direction == "right" {
			swipeMapping, whoLikeId = CheckMatchandReturn(whoLike, user, target)
			swipe = swipeMapping
		} else {
			swipe = models.Swipe{
				ID: userID,
				Username: user.Username,
				Like: pq.StringArray{},
				Pass: pq.StringArray{target.Username},
				Match: pq.StringArray{},
			}
		}
		if err:= database.DB.Table("swipes").Create(&swipe).Error; err != nil {
			return header.JSON(http.StatusBadRequest, map[string]string{"message": "Error like/dislike the profile"})
		}
		checkIsMatch := slices.Contains(whoLikeId, SwipePayload.TargetId)
		if checkIsMatch {
			return header.JSON(http.StatusCreated, map[string]string{"message": "You are matching"})
		}
		
		return header.JSON(http.StatusCreated, map[string]string{"message": "Success swipe"})
	}

	if SwipePayload.Direction == "right" {
		swipeMapping, whoLikeId = CheckMatchandReturn(whoLike, user, target)
		swipe.Like = slices.Concat(swipe.Like, swipeMapping.Like)
		swipe.Pass = swipe.Pass

		var savedMatch pq.StringArray

		if (len(swipeMapping.Match) > 0) {
			savedMatch = swipeMapping.Match
		} else {
			savedMatch = slices.Concat(swipe.Match, swipeMapping.Match)
		}

		swipe.Match = savedMatch
	} else {
		var swipePassUpdate pq.StringArray
		swipePassUpdate = slices.Concat(swipePassUpdate, pq.StringArray{target.Username})

		swipe = models.Swipe{
			ID: userID,
			Username: user.Username,
			Like: swipe.Like,
			Pass: swipePassUpdate,
			Match: swipe.Match,
		}
	}

	if err := database.DB.Debug().Model(&swipe).Select("*").Updates(swipe).Error; err != nil {
		return header.JSON(http.StatusBadRequest, map[string]string{"message": "Error like/dislike the profile"})
	}

	checkIsMatch := slices.Contains(whoLikeId, SwipePayload.TargetId)
	if checkIsMatch {
		return header.JSON(http.StatusCreated, map[string]string{"message": "You are matching"})
	}

	return header.JSON(http.StatusCreated, map[string]string{"messsage": "Success swipe"})
}
