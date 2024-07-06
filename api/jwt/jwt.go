package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	Sub string `json:"sub"`
	Id  int64  `json:"id"`
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
}

func base64URLEncode(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(jsonBytes), nil
}

func signMessage(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func CreateJWT(secret string, id int64) (string, error) {
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	payload := Payload{
		Sub: "1234567890",
		Id:  id,
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Hour * 24).Unix(),
	}

	encodedHeader, err := base64URLEncode(header)
	if err != nil {
		return "", err
	}

	encodedPayload, err := base64URLEncode(payload)
	if err != nil {
		return "", err
	}

	message := encodedHeader + "." + encodedPayload
	signature := signMessage(message, secret)

	jwt := message + "." + signature
	return jwt, nil
}

func VerifyJWT(token, secret string) (Payload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Payload{}, fmt.Errorf("invalid token format")
	}

	message := parts[0] + "." + parts[1]
	receivedSignature := parts[2]
	expectedSignature := signMessage(message, secret)

	if receivedSignature != expectedSignature {
		return Payload{}, fmt.Errorf("invalid token signature")
	}

	decodedPayload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Payload{}, err
	}

	var payload Payload
	err = json.Unmarshal(decodedPayload, &payload)
	if err != nil {
		return Payload{}, err
	}

	if time.Now().Unix() > payload.Exp {
		return Payload{}, fmt.Errorf("token has expired")
	}

	return payload, nil
}
