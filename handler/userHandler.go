package handler

import (
	"log"
	"miniproject/config"
	"miniproject/entity"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new user
// @Description Register a new user with the provided email and password
// @ID register-user
// @Accept json
// @Produce json
// @Param request body entity.User true "User registration request body"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 500 {object} map[string]interface{} "Failed to hash password" "Failed to register user" "Failed to send registration email"
// @Router /register [post]
func RegisterUser(c echo.Context) error {
	input := new(entity.User)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid request data"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to hash password"})
	}

	// Simpan user ke database
	user := entity.User{
		Email:         input.Email,
		Password:      string(hashedPassword),
		DepositAmount: input.DepositAmount,
	}
	if err := config.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to register user"})
	}

	if err := sendRegistrationEmail(user.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to send registration email"})
	}

	// Menghilangkan data sensitif dari respons
	user.Password = ""

	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "success register", "user": user})
}

// @Summary Login a user
// @Description Log in with the provided email and password
// @ID login-user
// @Accept json
// @Produce json
// @Param request body entity.User true "User login request body"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Failed to generate token"
// @Router /login [post]
func LoginUser(c echo.Context) error {
	input := new(entity.User)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid request data"})
	}

	// Cari user berdasarkan email
	var user entity.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "Invalid credentials"})
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "Invalid credentials"})
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.Id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to generate token"})
	}

	// Menghilangkan data sensitif dari respons
	user.Password = ""

	return c.JSON(http.StatusOK, map[string]interface{}{"token": tokenString, "message": "login success"})
}

// Middleware untuk memeriksa token JWT
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "Token is missing"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("12345"), nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "Invalid token"})
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := int(claims["sub"].(float64))
		log.Printf("User ID from token: %d", userID)

		c.Set("user", userID)

		return next(c)
	}
}
