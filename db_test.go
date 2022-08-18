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
	AmiRepo[Event]
}

type UserRepo struct {
	AmiRepo[User]
}

type Atlas struct {
	Id           int       `ami:"id"`
	Name         string    `ami:"name"`
	CreateTime   time.Time `ami:"create_time"`
	UpdateTime   time.Time `ami:"update_time"`
	CreateUserId int       `ami:"create_user_id"`
	UpdateUserId int       `ami:"update_user_id"`
	IsDelete     int       `ami:"is_delete"`
}

func (Atlas) GetTableName() string {
	return "xhh_atlas_tag"
}

type AtlasRepo struct {
	AmiRepo[Atlas]
}

var userRepo = UserRepo{}
var eventRepo = EventRepo{}
var atlasRepo = AtlasRepo{}

func TestDblib(t *testing.T) {
	atlasRepo.Insert(&Atlas{
		Id:           0,
		Name:         "liuzx",
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
		CreateUserId: 2,
		UpdateUserId: 3,
		IsDelete:     1,
	})
	// dataQuery := userRepo.Query("id in (?, ?)", 1, 2)
	//
	// dataQuery.Read(func(user User) {
	// 	fmt.Println(user.Passwd)
	// 	fmt.Println(user.UpdateAt)
	// })

	// dq2 := eventRepo.Query("1=1")
	// dq2.Read(func(event Event) {
	// 	t.Log(event.CreatedAt)
	// 	t.Log(event.Id)
	// })

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
