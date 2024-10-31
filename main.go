package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type OtpRequest struct {
	Otp string `json:"otp"`
}

type Device struct {
	Uuid     string `json:"uuid"`
	FcmToken string `json:"fcm_token"`
}

var db *sql.DB

func InitDb() error {
    var err error
    db, err = sql.Open("sqlite3", "./devices.db")
    if err != nil {
        return fmt.Errorf("error: failed to open database: %w\n", err)
    }
    
    var createTableQuery string = `
    CREATE TABLE IF NOT EXISTS devices (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid TEXT NOT NULL UNIQUE,
        fcm_token TEXT NOT NULL UNIQUE
    );`
    
    _, err = db.Exec(createTableQuery)
    if err != nil {
        return fmt.Errorf("error: failed to execute `createTableQuery`: %w\n", err)
    }
    
    return nil
}

func JsonResponse(writer http.ResponseWriter, message string, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(map[string]string{"message": message})
}

func ReceiveOtp(writer http.ResponseWriter, req *http.Request) {
	var otpReq OtpRequest
    
    var err error = json.NewDecoder(req.Body).Decode(&otpReq)
	if err != nil || otpReq.Otp == "" {
		fmt.Println("error: no OTP provided")
		JsonResponse(writer, "No OTP provided", http.StatusBadRequest)
		return
	}

	fmt.Printf("info: received OTP: %s\n", otpReq.Otp)
	JsonResponse(writer, "OTP received successfully", http.StatusOK)
}

func RegisterDevice(writer http.ResponseWriter, req *http.Request) {
	var device Device

    var err error = json.NewDecoder(req.Body).Decode(&device)
	if err != nil || device.Uuid == "" || device.FcmToken == "" {
		fmt.Println("error: invalid device registration request")
		JsonResponse(writer, "Invalid device registration request", http.StatusBadRequest)
		return
	}

	var insertQuery string = `INSERT OR REPLACE INTO devices (uuid, fcm_token) VALUES (?, ?)`
	_, err = db.Exec(insertQuery, device.Uuid, device.FcmToken);
	if err != nil {
	    fmt.Errorf("error: failed to execute `insertQuery`: %w\n", err);
	    JsonResponse(writer, "Failed to register device", http.StatusInternalServerError)
        return
	}
	
	fmt.Printf("info: registered device with UUID: %s, FCM token: %s\n",
		device.Uuid, device.FcmToken)
	JsonResponse(writer, "Device registered successfully", http.StatusOK)
}

func main() {
	err := InitDb()
    if err != nil {
        fmt.Println("error: failed to initialize database:", err)
        return
    }
	fmt.Println("info: database initialized")
	defer db.Close()
	
	http.HandleFunc("/receive_otp", ReceiveOtp)
	http.HandleFunc("/register_device", RegisterDevice)

	var ip_port string = "0.0.0.0:3000"
	fmt.Printf("info: server is running at http://%s\n", ip_port)
	if err := http.ListenAndServe(ip_port, nil); err != nil {
		fmt.Println("error: failed to start server: ", err)
	}
}
