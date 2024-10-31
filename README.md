> Go OTP Server for integration with the [gibotp](https://github.com/hyouteki/gibotp) Android application.
> This server registers devices and requests OTPs from specified services by sending notifications to the registered device.

# API documentation
Server exposes two endpoints.

## 1. Receive OTP
- **URL**: `/receive_otp`
- **Method**: `POST`
- **Description**: Receives an OTP sent from the client and validates its presence in the request. If valid, a success message is returned.

### Request Body
| Field | Type   | Description            | Required |
|-------|--------|------------------------|----------|
| otp   | String | The OTP to be verified | Yes      |

### Example
```bash
curl -X POST http://0.0.0.0:3000/receive_otp -H "Content-Type: application/json" -d '{"otp":"123456"}'
```
```json
{
    "otp": "123456"
}
```

### Responses
- `200 OK`: The OTP was received successfully.
- `400 Bad Request`: No OTP was provided in the request.

## 2. Register Device
- **URL**: `/register_device`
- **Method**: `POST`
- **Description**: Registers a device using the provided UUID and FCMToken. If successful, the server stores the device's FCM token using the UUID as a key.

### Request Body
| Field	    | Type	 | Description	                                     | Required |
|-----------|--------|---------------------------------------------------|----------|
| uuid	    | String | The unique identifier for the device              | Yes      |
| fcm_token	| String | The Firebase Cloud Messaging token for the device | Yes      |

### Example
```bash
curl -X POST http://0.0.0.0:3000/register_device -H "Content-Type: application/json" -d '{"uuid":"device-123", "fcm_token":"abcd1234fcmToken"}'
```
```json
{
    "uuid": "device-123",
    "fcm_token": "abcd1234fcmToken"
}
```

### Responses
- `200 OK`: The device was registered successfully.
- `400 Bad Request`: The request is missing required fields or has invalid data.
- `500 Internal Server Error`: Unable to process the request. Possible reasons: Error occurred while inserting data into the database for user registration.
