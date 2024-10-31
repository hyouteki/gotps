-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Unique identifier for the user's device
    device_uuid TEXT NOT NULL UNIQUE,
    -- Firebase Cloud Messaging (FCM) token for device communication
    device_fcm_token TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS otps (
    otp_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    otp TEXT NOT NULL,
    service_id INTEGER NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (service_id) REFERENCES services(service_id)
);

CREATE TABLE IF NOT EXISTS services (
    service_id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_name TEXT NOT NULL,
    -- Regular expression pattern used to extract OTP from service messages
    regex TEXT NOT NULL,
    -- Example of a message from this service containing an OTP
    sample_otp_message TEXT NOT NULL,
    -- OTP extracted from sample message using the regex (for development/testing purposes)
    sample_otp TEXT NOT NULL,
    -- Expected OTP time-to-live (TTL) in seconds, specific to this service
    expected_otp_ttl INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
    event_id INTEGER PRIMARY KEY AUTOINCREMENT,
    otp_id INTEGER,
    -- Status of the OTP event. Possible values:
    -- 1. "awaiting": Awaiting an OTP from the user (otp_id is NULL in this state)
    -- 2. "received": OTP has been received and is linked via otp_id
    -- 3. "expired": OTP has expired and further action is required.
    --    (May either return to "awaiting" for a new OTP or remove the event)
    event_status TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (otp_id) REFERENCES otps(otp_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);
