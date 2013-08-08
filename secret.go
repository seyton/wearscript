package main

import (
	"net/http"
	"crypto/sha1"
	"crypto/rand"
	"io"
	"fmt"
	picarus "github.com/bwhite/picarus/go"
)

func randString() (string, error) {
	nBytes := 12
	b := make([]byte, nBytes)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		fmt.Println("error:", err)
		return "", err
	}
	return string(b), nil
}

func secretHash(secret string) (string) {
	h := sha1.New()
	io.WriteString(h, secret)
	return picarus.UB64Enc(string(h.Sum(nil)))
}

func SecretKeySetupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(400)
		return
	}
	userId, err := userID(r)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	secretType := r.URL.Query().Get(":type")
	if secretType != "raven" && secretType != "borg" && secretType != "notify" {
		w.WriteHeader(400)
		return
	}
	secret, err := randString()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	// Remove previous secret
	prevSecretHash, err := getUserAttribute(userId, "secret_hash_" + secretType)
	if err == nil {
		deleteSecretUser(secretType, prevSecretHash)
	}
	secret = picarus.UB64Enc(secret)
	hash := secretHash(secret)
	setSecretUser(secretType, hash, userId)
	setUserAttribute(userId, "secret_hash_" + secretType, hash)
	io.WriteString(w, secret)
}