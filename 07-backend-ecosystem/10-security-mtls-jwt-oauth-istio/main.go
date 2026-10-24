// Shows JWT HS256 sign and verify with the stdlib: signature check,
// expiry check, and why you never trust the header's alg blindly.
// Run: go run ./07-backend-ecosystem/10-security-mtls-jwt-oauth-istio
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var enc = base64.RawURLEncoding

func sign(secret []byte, signingInput string) string {
	m := hmac.New(sha256.New, secret)
	m.Write([]byte(signingInput))
	return enc.EncodeToString(m.Sum(nil))
}

func issue(secret []byte, sub string, ttl time.Duration) string {
	header := enc.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	claims, _ := json.Marshal(map[string]any{"sub": sub, "exp": time.Now().Add(ttl).Unix()})
	input := header + "." + enc.EncodeToString(claims)
	return input + "." + sign(secret, input)
}

func verify(secret []byte, token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errors.New("malformed token")
	}
	hdr, err := enc.DecodeString(parts[0])
	if err != nil || !strings.Contains(string(hdr), `"alg":"HS256"`) {
		return "", errors.New("unexpected alg") // pin the algorithm, ignore "none"
	}
	want := sign(secret, parts[0]+"."+parts[1])
	if !hmac.Equal([]byte(want), []byte(parts[2])) { // constant-time compare
		return "", errors.New("bad signature")
	}
	raw, err := enc.DecodeString(parts[1])
	if err != nil {
		return "", errors.New("bad payload")
	}
	var c struct {
		Sub string `json:"sub"`
		Exp int64  `json:"exp"`
	}
	if err := json.Unmarshal(raw, &c); err != nil {
		return "", err
	}
	if time.Now().Unix() >= c.Exp {
		return "", errors.New("token expired")
	}
	return c.Sub, nil
}

func main() {
	secret := []byte("demo-secret-not-for-production")

	fmt.Println("=== valid token ===")
	tok := issue(secret, "user-42", time.Hour)
	sub, err := verify(secret, tok)
	fmt.Println("  sub:", sub, "err:", err)

	fmt.Println("=== tampered payload ===")
	p := strings.Split(tok, ".")
	evil := enc.EncodeToString([]byte(`{"sub":"admin","exp":9999999999}`))
	_, err = verify(secret, p[0]+"."+evil+"."+p[2])
	fmt.Println("  err:", err)

	fmt.Println("=== expired token ===")
	_, err = verify(secret, issue(secret, "user-42", -time.Minute))
	fmt.Println("  err:", err)

	fmt.Println("=== alg none attack ===")
	none := enc.EncodeToString([]byte(`{"alg":"none"}`)) + "." + p[1] + "."
	_, err = verify(secret, none)
	fmt.Println("  err:", err)
}
