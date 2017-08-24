package utils

import (
	"reflect"
	"crypto"
	"strings"
)

var supported_algorithms = map[string]crypto.Hash {
	"SHA256": crypto.SHA256,
	"SHA512": crypto.SHA512,
}

var supported_algorithms_lookup = map[crypto.Hash]string {
	crypto.SHA256: "SHA256",
	crypto.SHA512: "SHA512",
}

// Determine if the hashing algorithm is supported
func AlgorithmIsSupported(algorithm string) bool {
	_, ok := supported_algorithms[strings.ToUpper(algorithm)]
	return ok
}

// Get a list of supported algorithms
func GetSupportedAlgorithms() string {
	return strings.Join(mapSupportedAlgorithmsKeys(supported_algorithms), ", ")
}

// Get the algorithm as a crypto.Hash
func GetAlgorithmType(algorithm string) crypto.Hash {
	return supported_algorithms[strings.ToUpper(algorithm)]
}

// Get the algorithm as a string
func GetAlgorithmAsString(algorithm crypto.Hash) string {
	return strings.ToUpper(supported_algorithms_lookup[algorithm])
}

func mapSupportedAlgorithmsKeys(m map[string]crypto.Hash) []string {
	keys := reflect.ValueOf(m).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	return strkeys
}
