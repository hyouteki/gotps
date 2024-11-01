package main

import (
	"os"
	"log"
	"net/http"

	"gotps/database"
	"gotps/api"
)

func main() {
	database.Constructor(os.Getenv("DATABASE_PATH"), os.Getenv("SQL_FILEPATH"))
	defer database.Deconstructor()
	
	http.HandleFunc("/receive_otp", api.ReceiveOtpHandler)
	http.HandleFunc("/register_device", api.RegisterDeviceHandler)

	var ipPort string = "0.0.0.0:3000"
	log.Printf("info: server is running at http://%s\n", ipPort)
	err := http.ListenAndServe(ipPort, nil)
	if err != nil {
		log.Fatal("error: failed to start server: ", err)
	}
}
