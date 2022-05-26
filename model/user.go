package model

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"log"

	"golang.org/x/crypto/scrypt"
)

const (
	SALT_LENGTH = 8
	SCRYPT_N    = 1 << 0x0e
	SCRYPT_R    = 8
	SCRYPT_P    = 1
	KEY_LENGTH  = 32
)

type UserCredentials struct {
	Salt string `json:"-" bson:"salt"`
	Hash string `json:"-" bson:"hash"`
}

type User struct {
	ID          string          `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Password    string          `json:"password,omitempty" bson:"-"`
	Credentials UserCredentials `json:"-" bson:"credentials"`
}

func (u *User) PrepareToSave() {
	// force MongoDB to generate an ID
	u.ID = ""

	// Hash password
	u.hashPassword()
}

func (u *User) hashPassword() {
	// Generate salt
	salt := make([]byte, SALT_LENGTH)
	if _, err := rand.Read(salt); err != nil {
		log.Fatalln(err)
	}

	// Encode salt
	dk, err := scrypt.Key([]byte(u.Password), salt, SCRYPT_N, SCRYPT_R, SCRYPT_P, KEY_LENGTH)
	if err != nil {
		log.Fatalln(err)
	}

	// Clear plaintext password
	u.Password = ""

	// Store salt and hash
	u.Credentials.Salt = base64.URLEncoding.EncodeToString(salt)
	u.Credentials.Hash = base64.URLEncoding.EncodeToString(dk)
}

func (u *User) CheckPassword(password string) bool {
	// Decode salt
	salt, err := base64.URLEncoding.DecodeString(u.Credentials.Salt)
	if err != nil {
		log.Fatalln(err)
	}

	// Decode hash
	dk, err := base64.URLEncoding.DecodeString(u.Credentials.Hash)
	if err != nil {
		log.Fatalln(err)
	}

	// Hash password
	pdk, err := scrypt.Key([]byte(password), salt, SCRYPT_N, SCRYPT_R, SCRYPT_P, KEY_LENGTH)
	if err != nil {
		return false
	}

	return bytes.Equal(dk, pdk)
}
