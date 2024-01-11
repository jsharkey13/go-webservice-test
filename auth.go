package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type NotLoggedInException struct{}

func (e NotLoggedInException) Error() string {
	return "You must be logged in to access this resource."
}

// TODO: should this type be UpperCase or lowerCase? Does it matter?
type segueCookie struct {
	Id      string `json:"id"`
	Token   string `json:"token"`
	Expires string `json:"expires"`
	Hmac    string `json:"HMAC"`
}

// This implements the Stringer interface for the segueCookie type.
func (s segueCookie) String() string {
	return fmt.Sprintf("<ID: %s, Token: %s, Expires: %s, HMAC: %s>", s.Id, s.Token, s.Expires, s.Hmac)
}

func calculateHmac(cookie segueCookie) []byte {
	var hmacSaltStr = os.Getenv("HMAC_SALT")
	stringToHmac := fmt.Sprintf("%s|%s|%s", cookie.Id, cookie.Expires, cookie.Token)
	mac := hmac.New(sha256.New, []byte(hmacSaltStr))
	mac.Write([]byte(stringToHmac))
	return mac.Sum(nil)
}

func hasValidHmac(cookie segueCookie) bool {
	ourHmac := calculateHmac(cookie)
	cookieHmac, err := base64.StdEncoding.DecodeString(cookie.Hmac)
	if err != nil {
		return false
	}
	return hmac.Equal(ourHmac, cookieHmac)
}

func isStillValidSession(cookie segueCookie, user RegisteredUser) bool {
	// Check expiry date not passed:
	expiryDate, err := time.Parse(time.RubyDate, cookie.Expires)
	if err != nil {
		return false
	}
	if time.Now().After(expiryDate) {
		return false
	}

	// Check session token:
	if user.sessionToken != cookie.Token {
		return false
	}

	return true
}

func getSegueCookie(c *gin.Context) (segueCookie, error) {
	// This will not return tampered cookies, but can return expired or revoked sessions.
	cookieString, err := c.Cookie("SEGUE_AUTH_COOKIE")
	if err != nil {
		return segueCookie{}, NotLoggedInException{}
	}
	decoded, err := base64.StdEncoding.DecodeString(cookieString)
	if err != nil {
		return segueCookie{}, errors.New("could not decode cookie")
	}

	var cookie segueCookie
	err2 := json.Unmarshal(decoded, &cookie)
	if err2 != nil {
		return segueCookie{}, errors.New("could not parse cookie JSON")
	}

	if !hasValidHmac(cookie) {
		return segueCookie{}, errors.New("invalid cookie HMAC")
	}

	return cookie, nil
}

func getUntrustedUserIdFromCookie(segueCookie segueCookie) (int64, error) {
	return strconv.ParseInt(segueCookie.Id, 10, 64)
}

// GetCurrentUser returns the current user, or an error if no user is logged in.
func GetCurrentUser(c *gin.Context) (RegisteredUser, error) {
	segueCookie, err := getSegueCookie(c)
	if err != nil {
		// TODO: is there a better way of error handling than this? :/
		return RegisteredUser{}, err
	}
	userId, err := getUntrustedUserIdFromCookie(segueCookie)
	if err != nil {
		return RegisteredUser{}, err
	}
	user, err := getUserById(userId)
	if err != nil {
		return RegisteredUser{}, err
	}

	if !isStillValidSession(segueCookie, user) {
		return RegisteredUser{}, NotLoggedInException{}
	}
	return user, nil
}
