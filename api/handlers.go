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
	var device database.Device

    var err error = json.NewDecoder(req.Body).Decode(&device)
	if err != nil || device.Uuid == "" || device.FcmToken == "" {
		log.Println("error: invalid device registration request")
		JsonResponse(writer, "Invalid device registration request", http.StatusBadRequest)
		return
	}

	err = database.RegisterDevice(device);
	if err != nil {
	    log.Printf("error: failed to register device: %v\n", err);
	    JsonResponse(writer, "Failed to register device", http.StatusInternalServerError)
        return
	}
	
	log.Printf("info: registered device with UUID: %s, FCM token: %s\n",
		device.Uuid, device.FcmToken)
	JsonResponse(writer, "Device registered successfully", http.StatusOK)
}
