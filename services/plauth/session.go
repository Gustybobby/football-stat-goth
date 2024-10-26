package plauth

import "encoding/base32"

func GenerateSessionToken() (string, error) {
	tokenBytes, err := generateRandomBytes(20)
	if err != nil {
		return "", err
	}
	token := base32.StdEncoding.EncodeToString(tokenBytes)
	return token, nil
}
