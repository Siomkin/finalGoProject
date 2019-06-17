package main

import (
	"context"
	"fmt"
	"log"
	"main/infrastructure"
	"time"
)


func main() {

	ctx := context.Background()

	login := "ldo"
	pass := "123"
	groupName := "2я младшая Б"

	infrastructure.AddGroup(ctx, groupName)

	group, err := infrastructure.GetGroupByName(ctx, groupName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(group)

	//infrastructure.AddPrice(ctx, time.Date(2019,06,01, 1,0,0,0, time.Local).Unix(), group.ID, 2.4)
	//infrastructure.AddPrice(ctx, time.Date(2019,06,10, 1,0,0,0, time.Local).Unix(), group.ID, 2.6)

	infrastructure.ChangePrice(ctx, time.Date(2019,06,10, 1,0,0,0, time.Local).Unix(), group.ID, 2.7)
	//date := time.Date(2019,06,10, 1,0,0,0, time.Local).Unix()
	//pr, err := infrastructure.GetPrice(ctx, date, group.ID)
	//fmt.Println(pr)

	//user, err := infrastructure.CreateUser(ctx, login, pass)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//user, err := infrastructure.GetUserByLogin(ctx, login)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}


	res := infrastructure.UserLogin(ctx, login, pass)
	fmt.Println(res)

	//fmt.Println("Created user with id:", user)
}
