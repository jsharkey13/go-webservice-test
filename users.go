package main

import (
	"context"
	"database/sql"
	"fmt"

	"gopkg.in/guregu/null.v4"

	"github.com/jsharkey13/go-webservice-test/customtime"
)

type UserContext struct {
	Stage     string `json:"stage"`
	Examboard string `json:"examBoard"`
}

type RegisteredUser struct {
	Id                      int64                  `json:"id"`
	GivenName               string                 `json:"givenName"`
	FamilyName              string                 `json:"familyName"`
	Email                   string                 `json:"email"`
	Role                    string                 `json:"role"`
	Gender                  null.String            `json:"gender"`
	DateOfBirth             *customtime.CustomTime `json:"dateOfBirth"`
	RegistrationDate        *customtime.CustomTime `json:"registrationDate"`
	CountryCode             null.String            `json:"countryCode"`
	SchoolId                null.String            `json:"schoolId"`
	SchoolOther             null.String            `json:"schoolOther"`
	LastUpdated             *customtime.CustomTime `json:"lastUpdated"`
	LastSeen                *customtime.CustomTime `json:"lastSeen"`
	EmailVerificationStatus string                 `json:"emailVerificationStatus"`
	RegisteredContexts      []UserContext          `json:"registeredContexts"`
	// Private (=not exported), since lower case (does this matter?) and no `json:...` format defined.
	sessionToken string
	// Not 'private', but still not in the JSON since `json:"-"` excludes it.
	FakeProperty string `json:"-"`
}

func getUserById(userId int64) (RegisteredUser, error) {
	var user RegisteredUser
	var uDob, uRegDate, uLastUpdated, uLastSeen sql.NullTime

	queryUser := `SELECT id, given_name, family_name, email, role, gender, date_of_birth,
             school_id, school_other,email_verification_status,country_code,registered_contexts,
             registration_date, last_updated, last_seen, session_token FROM users WHERE id=$1;`
	row := db.QueryRow(context.Background(), queryUser, userId)
	err := row.Scan(&user.Id, &user.GivenName, &user.FamilyName, &user.Email, &user.Role, &user.Gender, &uDob,
		&user.SchoolId, &user.SchoolOther, &user.EmailVerificationStatus, &user.CountryCode, &user.RegisteredContexts,
		&uRegDate, &uLastUpdated, &uLastSeen, &user.sessionToken)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err) // FIXME:print
		return RegisteredUser{}, err
	}

	user.DateOfBirth = customtime.NullTimeToCustomTime(uDob)
	user.RegistrationDate = customtime.NullTimeToCustomTime(uRegDate)
	user.LastUpdated = customtime.NullTimeToCustomTime(uLastUpdated)
	user.LastSeen = customtime.NullTimeToCustomTime(uLastSeen)

	return user, nil
}
