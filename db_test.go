package amasugi

import (
	"testing"
	"time"
)

type User struct {
	Passwd   string    `ami:"password"`
	UpdateAt time.Time `ami:"updated_at"`
}

func (User) GetTableName() string {
	return "user"
}

type Event struct {
	Id        uint64
	CreatedAt time.Time `ami:"created_at"`
}

func (Event) GetTableName() string {
	return "event"
}

type EventRepo struct {
	AbsRepo[Event]
}

type UserRepo struct {
	AbsRepo[User]
}

var userRepo = UserRepo{}
var eventRepo = EventRepo{}

func TestDblib(t *testing.T) {
	// dataQuery := userRepo.Query("id in (?, ?)", 1, 2)
	//
	// dataQuery.Read(func(user User) {
	// 	fmt.Println(user.Passwd)
	// 	fmt.Println(user.UpdateAt)
	// })

	dq2 := eventRepo.Query("1=1")
	dq2.Read(func(event Event) {
		t.Log(event.CreatedAt)
		t.Log(event.Id)
	})

	// stmt, err := db.Prepare("select * from user")
	// if err != nil {
	// 	panic(err)
	// }
	// row, err := stmt.Query()
	// rm := make(map[string]any)
	// row.Next()
	// err = row.Scan(&rm)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// t.Log(rm)
}
