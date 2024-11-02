package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Ion-Stefan/saas-go-fiber/config"
	"github.com/Ion-Stefan/saas-go-fiber/internal/model"
	"github.com/Ion-Stefan/saas-go-fiber/internal/service"
	"github.com/Ion-Stefan/saas-go-fiber/pkg/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

func SetupOauthRoutes(app fiber.Router, oauthConf *oauth2.Config) {
	app.Get("/oauth/google", func(c *fiber.Ctx) error {
		url := oauthConf.AuthCodeURL("state")
		return c.Redirect(url)
	})

	app.Get("/oauth/redirect", func(c *fiber.Ctx) error {
		// Get code from query params for generating token
		code := c.Query("code")
		if code == "" {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token: code is empty")
		}
		// Exchange code for token
		token, err := oauthConf.Exchange(context.Background(), code)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token: " + err.Error())
		}
		// Set client for getting user information
		client := oauthConf.Client(context.Background(), token)
		// Get user information
		response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info: " + err.Error())
		}

		defer response.Body.Close()
		// Respone user type
		var user util.GoogleUser
		// Read response body
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error reading response body: " + err.Error())
		}
		// Unmarshal user information
		err = json.Unmarshal(bytes, &user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error unmarshal json body " + err.Error())
		}

		dbuser, err := service.GetUserByEmail(user.Email)
		if err != nil {
			return c.Redirect(fmt.Sprintf("%v/login-error", config.Envs.WebsiteURL))
		}

		if dbuser == nil {
			newUser := model.User{
				Email: user.Email,
				Name:  user.Name,
			}

			createErr := service.CreateUser(&newUser)
			if createErr != nil {
				return c.Redirect(fmt.Sprintf("%v/login-error", config.Envs.WebsiteURL))
			}
			jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"user_id": newUser.ID,
				"email":   newUser.Email,
				"admin":   newUser.IsAdmin,
			})
			tokenString, err := jwt_token.SignedString([]byte(config.Envs.JWTSecret))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error generating token: " + err.Error())
			}

			// Set the token in a cookie
			util.SetJWTInCookie(c, tokenString)

			// Redirect to the homepage
			return c.Redirect(fmt.Sprintf("%v/homepage", config.Envs.WebsiteURL))

		}
		// Login the user
		existingUser, err := service.GetUserByEmail(user.Email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error checking user: " + err.Error())
		}
		jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": existingUser.ID,
			"email":   existingUser.Email,
			"admin":   existingUser.IsAdmin,
		})
		tokenString, err := jwt_token.SignedString([]byte(config.Envs.JWTSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error generating token: " + err.Error())
		}

		// Set the token in a cookie
		util.SetJWTInCookie(c, tokenString)

		// Redirect to the homepage
		return c.Redirect(fmt.Sprintf("%v/homepage", config.Envs.WebsiteURL))
	})
}
