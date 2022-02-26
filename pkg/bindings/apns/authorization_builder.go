package apns

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/bhojpur/service/pkg/utils/logger"
)

// The "issued at" timestamp in the JWT must be within one hour from the
// APNS server time. I set the expiration time at 55 minutes to ensure that
// a new certificate gets generated before it gets too close and risking a
// failure.
const expirationMinutes = time.Minute * 55

type authorizationBuilder struct {
	logger              logger.Logger
	mutex               sync.RWMutex
	authorizationHeader string
	tokenExpiresAt      time.Time
	keyID               string
	teamID              string
	privateKey          interface{}
}

func (a *authorizationBuilder) getAuthorizationHeader() (string, error) {
	authorizationHeader, ok := a.readAuthorizationHeader()
	if ok {
		return authorizationHeader, nil
	}

	return a.generateAuthorizationHeader()
}

func (a *authorizationBuilder) readAuthorizationHeader() (string, bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if time.Now().After(a.tokenExpiresAt) {
		return "", false
	}

	return a.authorizationHeader, true
}

func (a *authorizationBuilder) generateAuthorizationHeader() (string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.logger.Debug("Authorization token expired; generating new token")

	now := time.Now()
	claims := jwt.StandardClaims{
		IssuedAt: time.Now().Unix(),
		Issuer:   a.teamID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = a.keyID
	signedToken, err := token.SignedString(a.privateKey)
	if err != nil {
		return "", err
	}

	a.authorizationHeader = "bearer " + signedToken
	a.tokenExpiresAt = now.Add(expirationMinutes)

	return a.authorizationHeader, nil
}
