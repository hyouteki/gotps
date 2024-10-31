package api

import (
	"encoding/json"
	"log"
	"net/http"

	"gotps/database"
)

type OtpRequest struct {
	Otp string `json:"otp"`
}

func JsonResponse(writer http.ResponseWriter, message string, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(map[string]string{"message": message})
}

func ReceiveOtpHandler(writer http.ResponseWriter, req *http.Request) {
	var otpReq OtpRequest
    
    var err error = json.NewDecoder(req.Body).Decode(&otpReq)
	if err != nil || otpReq.Otp == "" {
		log.Println("error: no OTP provided")
		JsonResponse(writer, "No OTP provided", http.StatusBadRequest)
		return
	}

	log.Printf("info: received OTP: %s\n", otpReq.Otp)
	JsonResponse(writer, "OTP received successfully", http.StatusOK)
}

func RegisterDeviceHandler(writer http.ResponseWriter, req *http.Request) {
	var user database.User

    var err error = json.NewDecoder(req.Body).Decode(&user)
	if err != nil || user.DeviceUuid == "" || user.DeviceFcmToken == "" {
		log.Println("error: invalid user registration request")
		JsonResponse(writer, "Invalid user registration request", http.StatusBadRequest)
		return
	}

	err = database.RegisterDevice(user);
	if err != nil {
	    log.Printf("error: failed to register user: %v\n", err);
	    JsonResponse(writer, "Failed to register user", http.StatusInternalServerError)
        return
	}
	
	log.Printf("info: registered user with UUID: %s, FCM token: %s\n",
		user.DeviceUuid, user.DeviceFcmToken)
	JsonResponse(writer, "User registered successfully", http.StatusOK)
}
