package ttn

import "time"

// TTNUplinkMessage matches the structure of the TTN V3 JSON Webhook.
type TTNUplinkMessage struct {
	EndDeviceIDs struct {
		DeviceID string `json:"device_id"`
	} `json:"end_device_ids"`

	// Root level timestamp
	ReceivedAt time.Time `json:"received_at"`

	UplinkMessage struct {
		// Payload fields
		DecodedPayload map[string]interface{} `json:"decoded_payload"`

		// Metadata: Signal quality (Array of gateways, we usually take the first)
		RxMetadata []struct {
			RSSI        float64 `json:"rssi"`
			ChannelRSSI float64 `json:"channel_rssi"`
			SNR         float64 `json:"snr"`
			GatewayIDs  struct {
				GatewayID string `json:"gateway_id"`
			} `json:"gateway_ids"`
		} `json:"rx_metadata"`

		// Metadata: Transmission settings
		Settings struct {
			DataRate struct {
				Lora struct {
					Bandwidth       float64 `json:"bandwidth"`
					SpreadingFactor float64 `json:"spreading_factor"`
				} `json:"lora"`
			} `json:"data_rate"`
			// Frequency often comes as a string in TTN V3 (e.g. "868000000")
			Frequency interface{} `json:"frequency"`
		} `json:"settings"`
	} `json:"uplink_message"`
}
