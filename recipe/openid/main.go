/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package openid

import (
	"github.com/supertokens/supertokens-golang/recipe/jwt/jwtmodels"
	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *openidmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func CreateJWTWithContext(payload map[string]interface{}, validitySecondsPointer *uint64, userContext supertokens.UserContext, useStaticSigningKey *bool) (jwtmodels.CreateJWTResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return jwtmodels.CreateJWTResponse{}, err
	}
	return (*instance.RecipeImpl.CreateJWT)(payload, validitySecondsPointer, useStaticSigningKey, userContext)
}

func GetJWKSWithContext(userContext supertokens.UserContext) (jwtmodels.GetJWKSResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return jwtmodels.GetJWKSResponse{}, err
	}
	return (*instance.RecipeImpl.GetJWKS)(userContext)
}

func GetOpenIdDiscoveryConfigurationWithContext(userContext supertokens.UserContext) (openidmodels.GetOpenIdDiscoveryConfigurationResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return openidmodels.GetOpenIdDiscoveryConfigurationResponse{}, err
	}
	return (*instance.RecipeImpl.GetOpenIdDiscoveryConfiguration)(userContext)
}

func CreateJWT(payload map[string]interface{}, validitySecondsPointer *uint64, useStaticSigningKey *bool) (jwtmodels.CreateJWTResponse, error) {
	return CreateJWTWithContext(payload, validitySecondsPointer, &map[string]interface{}{}, useStaticSigningKey)
}

func GetJWKS() (jwtmodels.GetJWKSResponse, error) {
	return GetJWKSWithContext(&map[string]interface{}{})
}

func GetOpenIdDiscoveryConfiguration() (openidmodels.GetOpenIdDiscoveryConfigurationResponse, error) {
	return GetOpenIdDiscoveryConfigurationWithContext(&map[string]interface{}{})
}
