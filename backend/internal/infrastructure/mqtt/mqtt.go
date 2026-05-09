package mqtt

import (
	"fmt"
	"sync"
	"time"

	"github.com/danisbagus/genset-monitoring/backend/internal/config"
	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

var (
	mqttClient pahomqtt.Client
	mqttOnce   sync.Once
)

// NewMQTT creates and returns a singleton MQTT client with auto-reconnect.
func NewMQTT(cfg config.MQTTConfig, log *zap.Logger) (pahomqtt.Client, error) {
	var initErr error

	mqttOnce.Do(func() {
		broker := fmt.Sprintf("tcp://%s:%d", cfg.Broker, cfg.Port)

		opts := pahomqtt.NewClientOptions().
			AddBroker(broker).
			SetClientID(cfg.ClientID).
			SetUsername(cfg.Username).
			SetPassword(cfg.Password).
			SetAutoReconnect(true).
			SetMaxReconnectInterval(30 * time.Second).
			SetConnectRetry(true).
			SetConnectRetryInterval(5 * time.Second).
			SetKeepAlive(60 * time.Second).
			SetPingTimeout(10 * time.Second).
			SetOnConnectHandler(func(c pahomqtt.Client) {
				log.Info("MQTT connected", zap.String("broker", broker))
				registerSubscriptions(c, log)
			}).
			SetConnectionLostHandler(func(c pahomqtt.Client, err error) {
				log.Warn("MQTT connection lost", zap.Error(err))
			}).
			SetReconnectingHandler(func(c pahomqtt.Client, opts *pahomqtt.ClientOptions) {
				log.Info("MQTT reconnecting...")
			})

		client := pahomqtt.NewClient(opts)

		token := client.Connect()
		token.WaitTimeout(10 * time.Second)
		if token.Error() != nil {
			initErr = fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
			return
		}

		mqttClient = client
	})

	return mqttClient, initErr
}

// IsConnected returns true if the MQTT client is connected.
func IsConnected(client pahomqtt.Client) bool {
	if client == nil {
		return false
	}
	return client.IsConnected()
}

// registerSubscriptions sets up all topic subscriptions.
// Extend this function to add real subscriptions as the system grows.
func registerSubscriptions(client pahomqtt.Client, log *zap.Logger) {
	topics := map[string]byte{
		"genset/+/telemetry": 1,
		"genset/+/status":    1,
		"genset/+/alert":     1,
	}

	for topic, qos := range topics {
		t := client.Subscribe(topic, qos, buildMessageHandler(log, topic))
		t.WaitTimeout(5 * time.Second)
		if t.Error() != nil {
			log.Error("MQTT subscribe failed",
				zap.String("topic", topic),
				zap.Error(t.Error()),
			)
			continue
		}
		log.Info("MQTT subscribed", zap.String("topic", topic))
	}
}

// buildMessageHandler returns a message handler for the given topic.
func buildMessageHandler(log *zap.Logger, topic string) pahomqtt.MessageHandler {
	return func(client pahomqtt.Client, msg pahomqtt.Message) {
		log.Debug("MQTT message received",
			zap.String("topic", msg.Topic()),
			zap.ByteString("payload", msg.Payload()),
		)
		// TODO: dispatch to service layer based on topic
	}
}
