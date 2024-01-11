package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// 🤢 I tried other ways, but a global variable won out.
var db *pgxpool.Pool

func main() {
	godotenv.Load()

	// Database:
	initDb()
	defer db.Close()

	// Create Gin router
	router := gin.Default()

	isaacApi := router.Group("/isaac-api/api")
	isaacApi.GET("/users/current_user", currentUser)
	isaacApi.GET("/_/experiment", demonstrateBehaviour)

	router.Run()

}

// currentUser returns the current user data, in (roughly) the format of the Java endpoint.
func currentUser(c *gin.Context) {

	user, err := GetCurrentUser(c)
	if err != nil {
		if errors.Is(err, NotLoggedInException{}) {
			c.JSON(http.StatusUnauthorized, gin.H{"errorMessage": "You must be logged in to access this resource."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure"})
		fmt.Printf("ERROR: %s\n", err) // FIXME:print
		return
	}

	c.JSON(http.StatusOK, user)
}

// demonstrateBehaviour experiments with some of the user data in illuminating ways.
func demonstrateBehaviour(c *gin.Context) {
	user, err := GetCurrentUser(c)
	if err != nil {
		if errors.Is(err, NotLoggedInException{}) {
			c.JSON(http.StatusUnauthorized, gin.H{"errorMessage": "You must be logged in to access this resource."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failure"})
		fmt.Printf("ERROR: %s\n", err) // FIXME:print
		return
	}

	resp := gin.H{}

	// null.String behaves much like a Java String. Can't compare directly and must use `.String` property.
	if user.Gender.Valid && user.Gender.String == "PREFER_NOT_TO_SAY" {
		resp["unknownGender"] = true
	}

	// Rounding time only works as expected in UTC.
	// CustomTime dates must be cast before being used.
	yesterday := time.Now().Truncate(24 * time.Hour).Add(-1 * 24 * time.Hour)
	if user.LastSeen != nil && time.Time(*user.LastSeen).After(yesterday) {
		resp["sinceYesterdayMidnight"] = true
	}

	// null.String supports being marshalled to JSON directly without needing anything special.
	resp["gender"] = user.Gender

	// Even simple things like string formatting might return errors, but this simple-statement-if helps.
	if name, err := fmt.Printf("%s %s", user.GivenName, user.FamilyName); err == nil {
		resp["name"] = name
	}

	// Arrays and slices are odd in Go.
	if len(user.RegisteredContexts) > 0 {
		var ctxCopy []string
		for _, ctx := range user.RegisteredContexts {
			s := fmt.Sprintf("Stage: %s @ Examboard: %s", ctx.Stage, ctx.Examboard)
			ctxCopy = append(ctxCopy, s)
		}
		resp["contextSummary"] = ctxCopy
	}

	c.JSON(http.StatusOK, resp)
}
