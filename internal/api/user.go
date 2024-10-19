package api

import (
	"SkinRest/internal/database"
	"SkinRest/pkg/models"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	maxLoginLength    int = 20
	maxPasswordLength int = 32
)

// @BasePath /api/v1

// RegisterHandler godoc
// @Summary Register a new user
// @Description Registers a new user with provided details, returning success message if successful
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.User true "User registration object"
// @Success 200 {object} gin.H {"message": "Success"}
// @Failure 400 {object} gin.H {"error": "Error message"}
// @Failure 500 {object} gin.H {"error": "Internal server error"}
// @Router /register [post]
func RegisterHandler(c *gin.Context) {

	// Get AppContext from this context
	appctx, exists := c.MustGet("appCtx").(*database.AppContext)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var user models.User

	// Get JSON Body
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid fields: " + err.Error()})
		return
	}

	if len(user.Login) > maxLoginLength {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("login must not exceed %d characters", maxLoginLength),
		})
		return
	}

	if len(user.Password) > maxPasswordLength {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("password must not exceed %d characters", maxPasswordLength),
		})
		return
	}

	// Save user to database
	if err := appctx.CreateNewUser(&user); err != nil {
		if err.Error() == models.ErrAlrRegistered.Error() {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		appctx.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// LoginHandler godoc
// @Summary Log in a user
// @Description Authenticates a user by username and password, returning a JWT token if successful
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.User true "User login object"
// @Success 200 {object} gin.H {"token": "JWT Token"}
// @Failure 400 {object} gin.H {"error": "Missing or invalid fields"}
// @Failure 404 {object} gin.H {"error": "This user does not exist"}
// @Failure 500 {object} gin.H {"error": "Internal server error"}
// @Router /login [post]
func LoginHandler(c *gin.Context) {

	// Get AppContext from this context
	appctx, exists := c.MustGet("appCtx").(*database.AppContext)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var user models.User

	// Get JSON Body
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid fields: " + err.Error()})
		return
	}

	if len(user.Login) > maxLoginLength {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("login must not exceed %d characters", maxLoginLength),
		})
		return
	}

	if len(user.Password) > maxPasswordLength {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("password must not exceed %d characters", maxPasswordLength),
		})
		return
	}

	// Fetching User Info if exists from database
	userData, err := appctx.GetInfoUser(&user)

	if err != nil {
		if err.Error() == models.ErrUserNotFound.Error() {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		appctx.Logger.Error(err.Error())
		return
	}

	// validate token
	if err := database.CheckToken(userData.Token); err != nil {
		if err.Error() == models.ErrTokenExpired.Error() {

			newToken, err := appctx.UpdateUserToken(&user)
			if err != nil {
				if err.Error() == models.ErrUserNotFound.Error() {
					c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
					return
				}
				appctx.Logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"token": newToken})
		}
	}

	c.JSON(http.StatusOK, gin.H{"token": userData.Token})

}

// AboutMe godoc
// @Summary Get user information
// @Description Retrieves user information from token including login and skins collection
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.UserInfo "User information object"
// @Failure 404 {object} gin.H {"error": "This user does not exist"}
// @Failure 500 {object} gin.H {"error": "Internal server error"}
// @Router /about [get]
func AboutMe(c *gin.Context) {

	// Get AppContext from this context
	appctx, exists := c.MustGet("appCtx").(*database.AppContext)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Get user data from this context
	userdata, exists := c.MustGet("userData").(*models.UserData)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Get user skins collection from database
	skins, err := appctx.GetUserSkins(userdata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		appctx.Logger.Error(err.Error())
		return
	}

	// Create user information object
	userInfo := models.UserInfo{
		Login: userdata.Login,
		Skins: skins,
	}

	c.JSON(http.StatusOK, userInfo)
}
