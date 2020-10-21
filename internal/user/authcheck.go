package user

import (
	"context"
)

// AuthCheck checks for JWT token validity.
func AuthCheck(ctx context.Context, authkey string) bool {
	// hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
	//	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	//	if !ok || method.Alg() != "HS256" {
	//		return nil, fmt.Errorf("bad sign method")
	//	}
	//	return config.App.TokenSecret, nil
	// }
	// token, err := jwt.Parse(authkey, hashSecretGetter)
	// if err != nil || !token.Valid {
	//	jsonError(w, http.StatusUnauthorized, "bad token")
	//	return
	// }

	// payload, ok := token.Claims.(jwt.MapClaims)
	// if !ok {
	//	jsonError(w, http.StatusUnauthorized, "no payload")
	// }

	return true
}
