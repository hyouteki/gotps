package database

import (
	"log"
)

type Device struct {
	Uuid string `json:"uuid"`
	FcmToken string `json:"fcm_token"`
}

func RegisterDevice(device Device) error {
	var query string = `INSERT OR REPLACE INTO devices (uuid, fcm_token) VALUES (?, ?)`
	_, err := db.Exec(query, device.Uuid, device.FcmToken)
	if err != nil {
		log.Printf("error: failed to execute insert device query: %v\n", err)
		return err
	}
	return nil
}
