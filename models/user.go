package models

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	//ErrUserNotFound err user not found
	ErrUserNotFound = errors.New("user not found")
	//ErrInvalidLogin err invalid login
	ErrInvalidLogin = errors.New("invalid login")
	//ErrUsernameTaken err invalid login
	ErrUsernameTaken = errors.New("username taken")
)

//User struct
type User struct {
	id int64
}

//NewUser constructor
func NewUser(username string, hash []byte) (*User, error) {
	exists, err := client.HExists("user:by-username", username).Result()
	if err != nil {
		return nil, err
	} else if exists {
		return nil, ErrUsernameTaken
	}
	id, err := client.Incr("user:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "username", username)
	pipe.HSet(key, "hash", hash)
	pipe.HSet("user:by-username", username, id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}
	return &User{id}, nil
}

//GetID get id
func (user *User) GetID() (int64, error) {
	return user.id, nil
}

//GetUsername get username
func (user *User) GetUsername() (string, error) {
	key := fmt.Sprintf("user:%d", user.id)
	return client.HGet(key, "username").Result()
}

//GetHash get hash
func (user *User) GetHash() ([]byte, error) {
	key := fmt.Sprintf("user:%d", user.id)
	return client.HGet(key, "hash").Bytes()
}

//Authenticate authenticate
func (user *User) Authenticate(password string) error {
	hash, err := user.GetHash()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidLogin
	}
	return err
}

//GetUserByID get user by id
func GetUserByID(id int64) (*User, error) {
	return &User{id}, nil
}

//GetUserByUsername get user by username
func GetUserByUsername(username string) (*User, error) {
	id, err := client.HGet("user:by-username", username).Int64()
	if err == redis.Nil {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return GetUserByID(id)
}

//AuthenticateUser user authentication
func AuthenticateUser(username, password string) (*User, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, user.Authenticate(password)
}

//RegisterUser register user
func RegisterUser(username, password string) error {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}
	_, err = NewUser(username, hash)
	return err
}
