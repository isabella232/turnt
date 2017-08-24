package lib

import (
	"crypto"
	"crypto/hmac"
	"fmt"
	"encoding/base64"

	"github.com/rapid7/turnt/utils"
)

const AUTHN_PROTOCOL = "Rapid7-HMAC-V1"

// Generate the digest header
func GenerateDigest(algorithm crypto.Hash, body string) string {
	h := algorithm.New()
	h.Write([]byte(body))

	strAlg := utils.GetAlgorithmAsString(algorithm)
	return fmt.Sprintf("%s=%s", strAlg, base64.StdEncoding.EncodeToString(h.Sum(nil)))
}

// Generate the authorization header
func GenerateAuthorization(algorithm crypto.Hash, identity string, signature string) string {
	strAlg := utils.GetAlgorithmAsString(algorithm)
	strAuth := fmt.Sprintf("%s:%s", identity, signature)
	authorization := base64.StdEncoding.EncodeToString([]byte(strAuth))

	return fmt.Sprintf("%s-%s %s", AUTHN_PROTOCOL, strAlg, authorization)
}

// Generate the signature header
func GenerateSignature(algorithm crypto.Hash, identity string, secret string, digest string, method string, uri string, host string, date int64) string {
	mac := hmac.New(algorithm.New, []byte(secret))

	mac.Write([]byte(fmt.Sprintf("%s %s\n", method, uri)))
	mac.Write([]byte(fmt.Sprintf("%s\n", host)))
	mac.Write([]byte(fmt.Sprintf("%d\n", date)))
	mac.Write([]byte(fmt.Sprintf("%s\n", identity)))
	mac.Write([]byte(fmt.Sprintf("%s\n", digest)))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
