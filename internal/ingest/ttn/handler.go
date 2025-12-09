package ttn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-ingest-service/internal/cache"
	"go-ingest-service/internal/config"
	"go-ingest-service/internal/ingest/general"
	"go-ingest-service/internal/models"
	// Quota import removed for OSS

	"github.com/gofiber/fiber/v2"
)

// IngestTTNData accepts a webhook from The Things Network.
func IngestTTNData(c *fiber.Ctx) error {
	// HarborDetails loading removed for OSS as SchemaName is not used in the struct
	
	rawBody := c.Body()
	// fmt.Printf("Raw request body: %s\n", string(rawBody)) // Optional logging

	var ttnPayload TTNUplinkMessage
	if err := json.Unmarshal(rawBody, &ttnPayload); err != nil {
		fmt.Printf("[TTN] JSON Unmarshal error: %v\n", err)
		return c.Status(http.StatusBadRequest).JSON(models.APIError{
			Message: "Invalid request body",
			Details: "Could not parse JSON: " + err.Error(),
		})
	}

	// 1. Basic Validation
	if ttnPayload.EndDeviceIDs.DeviceID == "" {
		return c.Status(http.StatusBadRequest).JSON(models.APIError{
			Message: "Validation Error",
			Details: "Missing device_id in TTN payload",
		})
	}

	// ---------------------------------------------------------
	// LOGIC 1: DETERMINE TIME (With Safety Checks)
	// ---------------------------------------------------------
	finalTime := ttnPayload.ReceivedAt

	// Check if 'time' exists in decoded_payload to override
	if val, ok := ttnPayload.UplinkMessage.DecodedPayload["time"]; ok {
		if timeStr, ok := val.(string); ok {
			// Try parsing standard ISO formats
			parsed, err := time.Parse(time.RFC3339, timeStr)
			if err != nil {
				parsed, err = time.Parse("2006-01-02T15:04:05.999Z", timeStr)
			}

			// Only apply if parsing succeeded AND time is reasonable (not > 24h in future)
			if err == nil && parsed.Before(time.Now().Add(24*time.Hour)) {
				finalTime = parsed
			}
		}
		// Remove 'time' so it doesn't get processed as a float value later
		delete(ttnPayload.UplinkMessage.DecodedPayload, "time")
	}

	// ---------------------------------------------------------
	// LOGIC 2: EXTRACT DATA
	// ---------------------------------------------------------
	var batch []general.SensorData
	shipID := ttnPayload.EndDeviceIDs.DeviceID

	// A. Process Decoded Payload (User Data)
	for key, val := range ttnPayload.UplinkMessage.DecodedPayload {
		if floatVal, ok := toFloat64(val); ok {
			batch = append(batch, general.SensorData{
				Time:    finalTime,
				ShipID:  shipID,
				CargoID: key,
				Value:   &floatVal,
				// SchemaName removed for OSS
			})
		}
	}

	// B. Process Metadata (RSSI, SNR, Frequency, etc.)
	// 1. Signal Quality
	if len(ttnPayload.UplinkMessage.RxMetadata) > 0 {
		meta := ttnPayload.UplinkMessage.RxMetadata[0]
		batch = append(batch, makePoint(finalTime, shipID, "rssi", meta.RSSI))
		batch = append(batch, makePoint(finalTime, shipID, "snr", meta.SNR))
		batch = append(batch, makePoint(finalTime, shipID, "channel_rssi", meta.ChannelRSSI))
	}

	// 2. Transmission Settings
	settings := ttnPayload.UplinkMessage.Settings
	if val := settings.DataRate.Lora.Bandwidth; val > 0 {
		batch = append(batch, makePoint(finalTime, shipID, "bandwidth", val))
	}
	if val := settings.DataRate.Lora.SpreadingFactor; val > 0 {
		batch = append(batch, makePoint(finalTime, shipID, "spreading_factor", val))
	}
	// Frequency handling
	if freqVal, ok := toFloat64(settings.Frequency); ok && freqVal > 0 {
		batch = append(batch, makePoint(finalTime, shipID, "frequency", freqVal))
	}

	if len(batch) == 0 {
		// Valid payload structure, but no data we care about. Return 200 OK to stop TTN retrying.
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "No numeric data extracted"})
	}

	// Quota check removed for OSS

	// Queue for Worker
	queuedData := models.QueuedData{
		RetryCount: 0,
		Type:       "ttn",
		Data:       batch,
	}

	dataJSON, err := json.Marshal(queuedData)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to prepare TTN data for queue")
	}

	if err := cache.RedisClient.RPush(c.Context(), config.AppConfig.IngestQueueName, dataJSON).Err(); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to queue TTN data for ingestion")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":      "TTN Payload Processed",
		"data_points": len(batch),
		"ship_id":     shipID,
	})
}

// Helper to create a point quickly (Removed schema param)
func makePoint(t time.Time, ship, cargo string, val float64) general.SensorData {
	return general.SensorData{
		Time:    t,
		ShipID:  ship,
		CargoID: cargo,
		Value:   &val,
	}
}

// Helper to handle JSON number parsing
func toFloat64(unk interface{}) (float64, bool) {
	switch i := unk.(type) {
	case float64:
		return i, true
	case float32:
		return float64(i), true
	case int:
		return float64(i), true
	case int64:
		return float64(i), true
	case string:
		if f, err := strconv.ParseFloat(i, 64); err == nil {
			return f, true
		}
		return 0, false
	default:
		return 0, false
	}
}