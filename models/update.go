package models

import (
	"fmt"
	"strconv"
)

//Update struct
type Update struct {
	id int64
}

//NewUpdate constructor
func NewUpdate(userID int64, body string) (*Update, error) {
	id, err := client.Incr("update:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("update:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "user_id", userID)
	pipe.HSet(key, "body", body)
	pipe.LPush("updates", id)
	pipe.LPush(fmt.Sprintf("user:%d:updates", userID), id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}
	return &Update{id}, nil
}

//GetBody get body from Update
func (update *Update) GetBody() (string, error) {
	key := fmt.Sprintf("update:%d", update.id)
	return client.HGet(key, "body").Result()
}

//GetUser get user data from user_id in Update
func (update *Update) GetUser() (*User, error) {
	key := fmt.Sprintf("update:%d", update.id)
	userID, err := client.HGet(key, "user_id").Int64()
	if err != nil {
		return nil, err
	}
	return GetUserByID(userID)
}

func queryUpdates(key string) ([]*Update, error) {
	updateIDs, err := client.LRange(key, 0, 10).Result()
	if err != nil {
		return nil, err
	}
	updates := make([]*Update, len(updateIDs))
	for i, strID := range updateIDs {
		id, err := strconv.Atoi(strID)
		if err != nil {
			return nil, err
		}
		updates[i] = &Update{int64(id)}
	}
	return updates, nil
}

//GetAllUpdates fetch all updates from redis
func GetAllUpdates() ([]*Update, error) {
	return queryUpdates("updates")
}

//GetUpdates fetch all updates from redis
func GetUpdates(userID int64) ([]*Update, error) {
	key := fmt.Sprintf("user:%d:updates", userID)
	return queryUpdates(key)
}

//PostUpdate post update to redis
func PostUpdate(userID int64, body string) error {
	_, err := NewUpdate(userID, body)
	return err
}
