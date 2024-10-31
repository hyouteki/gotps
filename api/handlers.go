package api

import (
    "bytes"
    "os"
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

func SendFCMNotification(fcmToken string, message string) error {
    url := "https://fcm.googleapis.com/fcm/send"

    payload := map[string]interface{}{
        "to": fcmToken,
        "notification": map[string]string{
            "title": "Sample FCM Notification",
            "body":  message,
        },
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        log.Printf("error: failed to marshal payload: %v\n", err)
        return err
    }
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
    if err != nil {
        log.Printf("error: failed to create request: %v\n", err)
        return err
    }

    fcmServerKey := os.Getenv("FCM_TOKEN")
    req.Header.Set("Authorization", "key="+fcmServerKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("error: failed to send request: %v\n", err)
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("error: received non-200 response: %d\n", resp.StatusCode)
    }

    log.Println("info: FCM notification sent successfully")
    return nil
}
