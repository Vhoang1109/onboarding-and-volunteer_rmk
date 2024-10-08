package transport

import (
	"net/http"
	"github.com/cesc1802/onboarding-and-volunteer-service/feature/sign_in/dto"
	"github.com/cesc1802/onboarding-and-volunteer-service/feature/sign_in/usecase"
	"github.com/gin-gonic/gin"
	"log"
)

type SignInHandler struct {
	signInUsecase usecase.SignInUsecase
}

// NewSignInHandler creates a new instance of SignInHandler.
func NewSignInHandler(signInUsecase usecase.SignInUsecase) *SignInHandler {
	return &SignInHandler{
		signInUsecase: signInUsecase,
	}
}

// SignIn handles the HTTP request for user sign-in.
func (h *SignInHandler) SignIn(c *gin.Context) {
	var signInDTO dto.SignInRequestDTO
	if err := c.ShouldBindJSON(&signInDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	response, err := h.signInUsecase.SignIn(signInDTO)
	if err != nil {
		log.Printf("SignIn error: %v", err) // Logging the error
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SignUp handles the HTTP request for user sign-up.
func (h *SignInHandler) SignUp(c *gin.Context) {
	var signUpDTO dto.SignUpRequestDTO
	if err := c.ShouldBindJSON(&signUpDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	response, err := h.signInUsecase.SignUp(signUpDTO)
	if err != nil {
		log.Printf("SignUp error: %v", err) // Logging the error
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign up"})
		}
		return
	}

	c.JSON(http.StatusCreated, response)
}

// RegisterRoutes sets up the routes for the sign-in and sign-up handlers.
func (h *SignInHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/sign-in", h.SignIn)
	router.POST("/sign-up", h.SignUp)
}
