package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OtpRequest struct {
	Otp string `json:"otp"`
}

type DeviceRegistration struct {
	UUID     string `json:"uuid"`
	FCMToken string `json:"fcm_token"`
}

var deviceTokens = make(map[string]string)

func receiveOtp(writer http.ResponseWriter, req *http.Request) {
	var otpReq OtpRequest

	err := json.NewDecoder(req.Body).Decode(&otpReq)
	if err != nil || otpReq.Otp == "" {
		fmt.Println("error: no OTP provided")
		http.Error(writer, "No OTP provided", http.StatusBadRequest)
		return
	}

	fmt.Printf("info: received OTP: %s\n", otpReq.Otp)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{"message": "OTP received successfully"})
}

func registerDevice(writer http.ResponseWriter, req *http.Request) {
	var registration DeviceRegistration

	err := json.NewDecoder(req.Body).Decode(&registration)
	if err != nil || registration.UUID == "" || registration.FCMToken == "" {
		fmt.Println("error: invalid device registration")
		http.Error(writer, "Invalid device registration", http.StatusBadRequest)
		return
	}

	deviceTokens[registration.UUID] = registration.FCMToken
	fmt.Printf("info: registered device with UUID: %s, FCM token: %s\n",
		registration.UUID, registration.FCMToken)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{"message": "Device registered successfully"})
}

func main() {
	http.HandleFunc("/receive_otp", receiveOtp)
	http.HandleFunc("/register_device", registerDevice)

	var ip_port string = "0.0.0.0:3000"
	fmt.Printf("info: server is running at http://%s\n", ip_port)
	if err := http.ListenAndServe(ip_port, nil); err != nil {
		fmt.Println("error: failed to start server: ", err)
	}
}
