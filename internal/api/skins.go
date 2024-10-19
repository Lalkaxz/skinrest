package api

import (
	"SkinRest/internal/database"
	"SkinRest/pkg/models"
	"fmt"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	maxSkinNameLength   int = 30
	maxSkinSourceLength int = 255
)

// @BasePath /api/v1

// AddNewSkin godoc
// @Summary Add a new skin
// @Description Adds a new skin for the authenticated user, returning the created skin data
// @Tags skins
// @Accept json
// @Produce json
// @Param skin body models.Skin true "Skin object"
// @Success 201 {object} models.Skin "Created skin data"
// @Failure 400 {object} gin.H {"error": "Missing or invalid fields"}
// @Failure 404 {object} gin.H {"error": "This user does not exist"}
// @Failure 500 {object} gin.H {"error": "Internal server error"}
// @Router /skins/add [post]
func AddNewSkin(c *gin.Context) {

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

	var skin models.Skin

	// Get JSON Body
	if err := c.ShouldBindJSON(&skin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid fields: " + err.Error()})
		return
	}

	// Validation "skin type" field
	if skin.Type != "Classic" && skin.Type != "Slim" {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidSkinType.Error()})
		return
	}

	if len(skin.Name) > maxSkinNameLength {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("skin name must not exceed %d characters", maxSkinNameLength),
		})
		return
	}

	if len(skin.Src) > maxSkinSourceLength {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("skin source must not exceed %d characters", maxSkinSourceLength),
		})
		return
	}

	// Save skin to database
	skinData, err := appctx.AddNewSkin(userdata, &skin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		appctx.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusCreated, skinData)
}

// GetSkinsCollection godoc
// @Summary Retrieve user's skin collection
// @Description Gets the collection of skins for the authenticated user
// @Tags skins
// @Accept json
// @Produce json
// @Success 200 {array} models.Skin "List of user's skins"
// @Failure 404 {object} gin.H {"error": "This user does not exist"}
// @Failure 500 {object} gin.H {"error": "Internal server error"}
// @Router /skins [get]
func GetSkinsCollection(c *gin.Context) {

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

	c.JSON(http.StatusOK, skins)
}

// GetSkin godoc
// @Summary Retrieve a specific skin by ID
// @Description Gets the details of a specific skin for the authenticated user using the skin ID
// @Tags skins
// @Accept json
// @Produce json
// @Param id path int true "Skin ID"
// @Success 200 {object} models.Skin "Details of the requested skin"
// @Failure 400 {object} gin.H {"error": "Error message"}
// @Failure 404 {object} gin.H {"error": "Skin not found"}
// @Failure 404 {object} gin.H {"error": "This user does not exist"}
// @Failure 500 {object} gin.H {"error": "Internal server error"}
// @Router /skins/{id} [get]
func GetSkin(c *gin.Context) {

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

	// Get query param
	idParam := c.Param("id")

	// Convert query param to integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidIdFormat.Error()})
		return
	}

	if id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be greater than or equal to 1"})
		return
	}

	// Get skin from database
	skinData, err := appctx.GetUserSkin(userdata, id)
	if err != nil {
		if err == models.ErrSkinNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		appctx.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, skinData)
}

// DeleteSkin godoc
// @Summary Delete a specific skin by ID
// @Description Deletes the specified skin for the authenticated user using the skin ID
// @Tags skins
// @Accept json
// @Produce json
// @Param id path int true "Skin ID"
// @Success 200 {object} gin.H {"status": "Success"}
// @Failure 400 {object} gin.H {"error": "Error message"}
// @Failure 404 {object} gin.H {"error": "This user does not exist"}
// @Failure 404 {object} gin.H {"error": "Skin not found"}
// @Failure 500 {object} gin.H {"error": "Internal server error"}
// @Router /skins/{id} [delete]
func DeleteSkin(c *gin.Context) {
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
	// Get Query param
	idParam := c.Param("id")

	// Convert query param to integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidIdFormat.Error()})
		return
	}

	if id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be greater than or equal to 1"})
		return
	}

	// Remove skin from database
	err = appctx.DeleteUserSkin(userdata, id)
	if err != nil {
		if err == models.ErrSkinNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		appctx.Logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Success"})
}
