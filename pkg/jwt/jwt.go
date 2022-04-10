package jwt

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/dgrijalva/jwt-go"

	"github.com/google/uuid"
)

func NewTokenizer(hmacSecret []byte, validity time.Duration) *Tokenizer {
	return &Tokenizer{hmacSecret: hmacSecret, validity: validity}
}

type Tokenizer struct {
	hmacSecret []byte
	validity   time.Duration
}

func (t *Tokenizer) Generate(u *model.User) (string, error) {
	now := time.Now().UTC()
	teamID := ""
	if u.TeamID != nil {
		teamID = *u.TeamID
	}
	claims := &UserClaims{
		// Refers to model.User
		UserID: u.ID,
		// Refers to model.Team
		TeamID: teamID,
		// https://tools.ietf.org/html/rfc7519#section-4.1
		StandardClaims: jwt.StandardClaims{
			// The "iat" (issued at) claim identifies the time at which the JWT was
			// issued.  This claim can be used to determine the age of the JWT. Its
			// value MUST be a number containing a NumericDate value. Use of this
			// claim is OPTIONAL.
			IssuedAt: now.Unix(),
			// The "exp" (expiration time) claim identifies the expiration time on
			// or after which the JWT MUST NOT be accepted for processing. The
			// processing of the "exp" claim requires that the current date/time
			// MUST be before the expiration date/time listed in the "exp" claim.
			ExpiresAt: now.Add(t.validity).Unix(),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.hmacSecret)
}

func (t *Tokenizer) Valid(token string) (userID, teamID string, err error) {
	claims := &UserClaims{}
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will
	// return an error if the token is invalid (if it has expired according to
	// the expiry time we set on sign in), or if the signature does not match
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return t.hmacSecret, nil
	})
	return claims.UserID, claims.TeamID, err
}

type UserClaims struct {
	jwt.StandardClaims

	UserID string
	TeamID string
}

func (u *UserClaims) Valid() error {
	if _, err := uuid.Parse(u.UserID); err != nil {
		return err
	}
	return u.StandardClaims.Valid()
}

// Example:
// Authorization: Bearer <token>
func ExtractToken(r *http.Request) (string, error) {
	headline := r.Header.Get("Authorization")
	if headline == "" {
		return "", errors.New("missing field 'Authorization' in header")
	}
	splitToken := strings.Split(headline, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("could not split token")
	}
	return splitToken[1], nil
}
