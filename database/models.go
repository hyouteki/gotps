package database

import (
	"log"
)

type User struct {
	DeviceUuid string `json:"uuid"`
	DeviceFcmToken string `json:"fcm_token"`
}

func RegisterDevice(user User) error {
	var query string = `INSERT OR REPLACE INTO users (device_uuid, device_fcm_token) VALUES (?, ?)`
	_, err := db.Exec(query, user.DeviceUuid, user.DeviceFcmToken)
	if err != nil {
		log.Printf("error: failed to insert user into database: %v\n", err)
		return err
	}
	return nil
}

