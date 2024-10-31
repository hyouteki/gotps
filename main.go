package main

import (
	"log"
	"net/http"

	"gotps/database"
	"gotps/api"
)

func main() {
	err := database.Connect("devices.db")
    if err != nil {
        log.Fatal("error: failed to connect database: ", err)
    }
	log.Println("info: database connected")

	err = database.Init("devices.sql")
	if err != nil {
        log.Fatal("error: failed to initiliaze database: ", err)
    }
	log.Println("info: database initialized")
	
	http.HandleFunc("/receive_otp", api.ReceiveOtpHandler)
	http.HandleFunc("/register_device", api.RegisterDeviceHandler)

	var ipPort string = "0.0.0.0:3000"
	log.Printf("info: server is running at http://%s\n", ipPort)
	err = http.ListenAndServe(ipPort, nil)
	if err != nil {
		log.Fatal("error: failed to start server: ", err)
	}
}
