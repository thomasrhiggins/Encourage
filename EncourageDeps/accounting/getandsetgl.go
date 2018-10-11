package accounting

import (
	"fmt"
)

type Ledger struct {
	id int64
}

type Update struct {
	id int64
}

func NewUpdate(userId int64, body string) (*Update, error) {
	id, err := client.Incr("update:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("update:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "user_id", userId)
	pipe.HSet(key, "body", body)
	pipe.LPush("updates", id)
	pipe.LPush(fmt.Sprintf("user:%d:updates", userId), id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}
	return &Update{id}, nil
}
