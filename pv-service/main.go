package main

import (
	"fmt"
	"pv-service/database"
)

func main() {
	dbConnection := database.GetDBConnection()
	err := dbConnection.ConnectToDB()
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		return
	}

}
