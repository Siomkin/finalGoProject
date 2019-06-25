package main

import (
	"main/server"
	"main/server/handlers"

	//"log"
	"main/infrastructure"
	//"time"
)


func main() {

	//ctx := context.Background()

	//login := "ldo"
	//pass := "123"
	//groupName := "2я младшая Б"
	//
	//err := infrastructure.NewGroupRepository().AddGroup(ctx, groupName)
	//
	//group, err := infrastructure.NewGroupRepository().GetGroupByName(ctx, groupName)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(group)

	//infrastructure.AddPrice(ctx, time.Date(2019,06,01, 1,0,0,0, time.Local).Unix(), group.ID, 2.4)
	//infrastructure.AddPrice(ctx, time.Date(2019,06,10, 1,0,0,0, time.Local).Unix(), group.ID, 2.6)

	//infrastructure.NewPricesRepository().ChangePrice(ctx, time.Date(2019,06,10, 1,0,0,0, time.Local).Unix(), group.ID, 2.7)
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

	//res := infrastructure.UserLogin(ctx, login, pass)
	//fmt.Println(res)
	//fmt.Println("Created user with id:", user)


	usersRepository := infrastructure.NewUsersRepository()
	schoolsRepository := infrastructure.NewSchoolRepository()
	groupsRepository := infrastructure.NewGroupRepository()
	childRepository := infrastructure.NewChildrenRepository()
	tabelRepository := infrastructure.NewTabelRepository()
	//
	//authGuard := core.NewAuthGuard(usersRepository)
	//
	//authMiddleware := middleware.NewAuthMiddleware(authGuard)
	//
	usersHandler := handlers.NewUsersHandler(usersRepository, schoolsRepository, groupsRepository, childRepository, tabelRepository)
	//adminsHandler := handlers.NewAdminsHandler(usersRepository)

	//server := server.NewServer(":8080", usersHandler, adminsHandler, authMiddleware)
	srv := server.NewServer(":8080", usersHandler) //, adminsHandler, authMiddleware)
	srv.ConfigureAndRun()

}
