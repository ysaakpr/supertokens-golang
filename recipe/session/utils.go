package session

import (
	"errors"
	"net/url"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/session/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
	"golang.org/x/net/publicsuffix"
)

func ValidateAndNormaliseUserInput(recipeInstance *SessionRecipe, appInfo supertokens.NormalisedAppinfo, config *schema.TypeInput) (schema.TypeNormalisedInput, error) {
	typeNormalisedInput := MakeTypeNormalisedInput(appInfo)

	topLevelAPIDomain, err := GetTopLevelDomainForSameSiteResolution(appInfo.APIDomain.GetAsStringDangerous())
	if err != nil {
		return schema.TypeNormalisedInput{}, err
	}
	topLevelWebsiteDomain, err := GetTopLevelDomainForSameSiteResolution(appInfo.WebsiteDomain.GetAsStringDangerous())
	if err != nil {
		return schema.TypeNormalisedInput{}, err
	}

	cookieSameSite := CookieSameSite_LAX
	if topLevelAPIDomain != topLevelWebsiteDomain {
		cookieSameSite = CookieSameSite_NONE
	}
	typeNormalisedInput.CookieSameSite = cookieSameSite

	antiCsrf := AntiCSRF_NONE
	if cookieSameSite == CookieSameSite_NONE {
		antiCsrf = AntiCSRF_VIA_CUSTOM_HEADER
	}
	typeNormalisedInput.AntiCsrf = antiCsrf

	if config.CookieDomain != nil {
		cookieDomain, err := NormaliseSessionScopeOrThrowError(*config.CookieDomain)
		if err != nil {
			return schema.TypeNormalisedInput{}, err
		}
		typeNormalisedInput.CookieDomain = &cookieDomain
	}

	if config.CookieSameSite != nil {
		typeNormalisedInput.CookieSameSite = *config.CookieSameSite
	}

	if config.CookieSecure != nil {
		typeNormalisedInput.CookieSecure = *config.CookieSecure
	}

	if config.SessionExpiredStatusCode != nil {
		typeNormalisedInput.SessionExpiredStatusCode = *config.SessionExpiredStatusCode
	}

	if config.AntiCsrf != nil {
		if *config.AntiCsrf != AntiCSRF_NONE || *config.AntiCsrf != AntiCSRF_VIA_TOKEN || *config.AntiCsrf != AntiCSRF_VIA_TOKEN {
			return typeNormalisedInput, errors.New("antiCsrf config must be one of 'NONE' or 'VIA_CUSTOM_HEADER' or 'VIA_TOKEN'")
		}
		typeNormalisedInput.AntiCsrf = *config.AntiCsrf
	}

	IsAnIPAPIDomain, err := supertokens.IsAnIPAddress(topLevelAPIDomain)

	IsAnIPWebsiteDomain, err := supertokens.IsAnIPAddress(topLevelWebsiteDomain)
	if typeNormalisedInput.CookieSameSite == CookieSameSite_NONE &&
		!typeNormalisedInput.CookieSecure &&
		!(topLevelAPIDomain == "localhost" || IsAnIPAPIDomain) &&
		!(topLevelWebsiteDomain == "localhost" || IsAnIPWebsiteDomain) {
		return typeNormalisedInput, errors.New("Since your API and website domain are different, for sessions to work, please use https on your apiDomain and dont set cookieSecure to false.")
	}

	if config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
	}

	refreshAPIPath, err := supertokens.NewNormalisedURLPath(RefreshAPIPath)
	if err != nil {
		return schema.TypeNormalisedInput{}, err
	}
	typeNormalisedInput.RefreshTokenPath = appInfo.APIBasePath.AppendPath(*refreshAPIPath)

	return typeNormalisedInput, nil
}

func MakeTypeNormalisedInput(appInfo supertokens.NormalisedAppinfo) schema.TypeNormalisedInput {
	return schema.TypeNormalisedInput{
		CookieDomain:             nil,
		CookieSameSite:           CookieSameSite_LAX,
		CookieSecure:             strings.HasPrefix(appInfo.APIDomain.GetAsStringDangerous(), "https"),
		SessionExpiredStatusCode: 401,
		AntiCsrf:                 AntiCSRF_NONE,
		Override: struct {
			Functions func(originalImplementation schema.RecipeImplementation) schema.RecipeImplementation
			APIs      func(originalImplementation schema.APIImplementation) schema.APIImplementation
		}{
			Functions: func(originalImplementation schema.RecipeImplementation) schema.RecipeImplementation {
				return originalImplementation
			},
			APIs: func(originalImplementation schema.APIImplementation) schema.APIImplementation {
				return originalImplementation
			},
		},
	}
}

func GetTopLevelDomainForSameSiteResolution(URL string) (string, error) {
	urlObj, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	hostname := urlObj.Host
	isAnIP, err := supertokens.IsAnIPAddress(hostname)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(hostname, "localhost") || strings.HasPrefix(hostname, "localhost.org") || isAnIP {
		return "localhost", nil
	}
	parsedURL, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return "", errors.New("Please make sure that the apiDomain and websiteDomain have correct values")
	}
	return parsedURL, nil
}

func NormaliseSessionScopeOrThrowError(sessionScope string) (string, error) {
	sessionScope = strings.TrimSpace(sessionScope)
	sessionScope = strings.ToLower(sessionScope)

	if strings.HasPrefix(sessionScope, ".") {
		sessionScope = sessionScope[1:]
	}

	if !strings.HasPrefix(sessionScope, "http://") && !strings.HasPrefix(sessionScope, "https://") {
		sessionScope = "http://" + sessionScope
	}

	urlObj, err := url.Parse(sessionScope)
	if err != nil {
		return "", errors.New("Please provide a valid sessionScope")
	}
	sessionScope = urlObj.Host
	if strings.HasPrefix(sessionScope, ".") {
		sessionScope = sessionScope[1:]
	}

	noDotNormalised := sessionScope

	isAnIP, err := supertokens.IsAnIPAddress(sessionScope)
	if err != nil {
		return "", err
	}
	if sessionScope == "localhost" || isAnIP {
		noDotNormalised = sessionScope
	}
	if strings.HasPrefix(sessionScope, ".") {
		noDotNormalised = "." + sessionScope
	}
	return noDotNormalised, nil
}
