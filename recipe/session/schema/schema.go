package schema

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type HandshakeInfo struct {
	JWTSigningPublicKey            string
	AntiCsrf                       string
	AccessTokenBlacklistingEnabled bool
	JWTSigningPublicKeyExpiryTime  uint64
	AccessTokenValidity            uint64
	RefreshTokenValidity           uint64
}

type CreateOrRefreshAPIResponse struct {
	Session struct {
		Handle        string      `json:"handle"`
		UserID        string      `json:"userId"`
		UserDataInJWT interface{} `json:"userDataInJWT"`
	} `json:"session"`
	AccessToken    CreateOrRefreshAPIResponseToken `json:"accessToken"`
	RefreshToken   CreateOrRefreshAPIResponseToken `json:"refreshToken"`
	IDRefreshToken CreateOrRefreshAPIResponseToken `json:"idRefreshToken"`
	AntiCsrfToken  *string                         `json:"antiCsrfToken"`
}

type CreateOrRefreshAPIResponseToken struct {
	Token       string `json:"token"`
	Expiry      uint64 `json:"expiry"`
	CreatedTime uint64 `json:"createdTime"`
}

type TypeInput struct {
	CookieSecure             *bool
	CookieSameSite           *string
	SessionExpiredStatusCode *int
	CookieDomain             *string
	AntiCsrf                 *string
	Override                 *struct {
		Functions func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs      func(originalImplementation APIImplementation) APIImplementation
	}
}

type TypeNormalisedInput struct {
	RefreshTokenPath         supertokens.NormalisedURLPath
	CookieDomain             *string
	CookieSameSite           string
	CookieSecure             bool
	SessionExpiredStatusCode int
	AntiCsrf                 string
	Override                 struct {
		Functions func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs      func(originalImplementation APIImplementation) APIImplementation
	}
}

type SessionContainer struct {
	RevokeSession     func() error
	GetSessionData    func() (interface{}, error)
	UpdateSessionData func(newSessionData interface{}) (interface{}, error)
	GetUserId         func() string
	GetJWTPayload     func() interface{}
	GetHandle         func() string
	GetAccessToken    func() string
	UpdateJWTPayload  func(newJWTPayload interface{}) error
}

type VerifySessionOptions struct {
	AntiCsrfCheck   *bool
	SessionRequired *bool
}

type RecipeImplementation struct {
	Querier                       supertokens.Querier
	Config                        TypeNormalisedInput
	HandshakeInfo                 HandshakeInfo
	GetHandshakeInfo              func() HandshakeInfo
	UpdateJwtSigningPublicKeyInfo func(newKey string, newExpiry uint64)
	CreateNewSession              func(res http.ResponseWriter, userID string, jwtPayload interface{}, sessionData interface{}) (SessionContainer, error)
	GetSession                    func(req *http.Request, res http.ResponseWriter, options *VerifySessionOptions) (*SessionContainer, error)
	RefreshSession                func(req *http.Request, res http.ResponseWriter) (SessionContainer, error)
	RevokeAllSessionsForUser      func(userID string) ([]string, error)
	GetAllSessionHandlesForUser   func(userID string) ([]string, error)
	RevokeSession                 func(sessionHandle string) (bool, error)
	RevokeMultipleSessions        func(sessionHandles []string) ([]string, error)
	GetSessionData                func(sessionHandle string) (interface{}, error)
	UpdateSessionData             func(sessionHandle string, newSessionData interface{}) error
	GetJWTPayload                 func(sessionHandle string) (interface{}, error)
	UpdateJWTPayload              func(sessionHandle string, newJWTPayload interface{}) error
	GetAccessTokenLifeTimeMS      func() (uint64, error)
	GetRefreshTokenLifeTimeMS     func() (uint64, error)
}

type APIOptions struct {
	RecipeImplementation RecipeImplementation
	Config               TypeNormalisedInput
	RecipeID             string
	Req                  *http.Request
	Res                  http.ResponseWriter
	OtherHandler         http.HandlerFunc
}

type APIImplementation struct {
	RefreshPOST   func(options APIOptions) error
	SignOutPOST   func(options APIOptions) (map[string]string, error)
	VerifySession func(verifySessionOptions *VerifySessionOptions, options APIOptions)
}
