package inventory

import (
	"encoding/json"
	"github.com/djumanoff/amqp"
	"github.com/mukhametkaly/Diploma/document-api/src/config"
	"github.com/mukhametkaly/Diploma/document-api/src/models"
	"github.com/sony/sonyflake"
	"strconv"
	"sync"
)

var srvCfg = amqp.ServerConfig{
	ResponseX: "X:gateway.out.fanout",
	RequestX:  "X:routing.topic",
}
var pubCfg = amqp.PublisherConfig{}
var SessionRabbitMQ *amqp.Session
var PublisherRabbitMQ *amqp.Publisher
var ServerRabbitMQ *amqp.Server
var ClientRabbitMQ *amqp.Client
var cltCfg = amqp.ClientConfig{
	ResponseX: "X:routing.topic",
	RequestX:  "X:gateway.in.fanout",
	ResponseQ: "response.darbiz.transfer-api",
}

func RabbitMQConStart() (*amqp.Session, error) {
	cfg := amqp.Config{
		Host:        config.AllConfigs.Rabbit.Host,
		VirtualHost: config.AllConfigs.Rabbit.VirtualHost,
		User:        config.AllConfigs.Rabbit.User,
		Password:    config.AllConfigs.Rabbit.Password,
		Port:        config.AllConfigs.Rabbit.Port,
		LogLevel:    config.AllConfigs.Rabbit.LogLevel,
	}

	once := sync.Once{}

	once.Do(func() {
		flake := sonyflake.NewSonyflake(sonyflake.Settings{})
		id, err := flake.NextID()
		if err != nil {
			return
		}
		s := strconv.FormatUint(id, 16)

		cltCfg.ResponseQ += "." + s
	})

	sess := amqp.NewSession(cfg)
	err := sess.Connect()
	if err != nil {
		return nil, err
	}

	return &sess, nil
}

func GetRabbitClient() (*amqp.Client, error) {
	sess, err := GetRabbitSession()
	if err != nil {
		return nil, err
	}

	if ClientRabbitMQ == nil {
		clt, err := (*sess).Client(cltCfg)
		if err != nil {
			return nil, err

		} else {
			ClientRabbitMQ = &clt
			return ClientRabbitMQ, nil
		}
	} else {
		return ClientRabbitMQ, nil
	}
}

func GetRabbitSession() (*amqp.Session, error) {
	if SessionRabbitMQ == nil {
		sess, err := RabbitMQConStart()
		if err != nil {
			return nil, err
		} else {
			SessionRabbitMQ = sess
			return SessionRabbitMQ, nil
		}
	} else {
		return SessionRabbitMQ, nil
	}
}

func GetPublisher() (*amqp.Publisher, error) {
	sess, err := GetRabbitSession()
	if err != nil {
		return nil, err
	}
	if PublisherRabbitMQ == nil {
		pb, err := (*sess).Publisher(pubCfg)
		if err != nil {
			return nil, err
		} else {
			PublisherRabbitMQ = &pb
			return PublisherRabbitMQ, nil
		}
	} else {
		return PublisherRabbitMQ, nil
	}
}

func Server() (*amqp.Server, error) {
	sess, err := GetRabbitSession()
	if err != nil {
		return nil, err

	}
	if ServerRabbitMQ == nil {
		srv, err := (*sess).Server(srvCfg)
		if err != nil {
			return nil, err

		} else {
			ServerRabbitMQ = &srv
			return ServerRabbitMQ, nil
		}
	} else {
		return ServerRabbitMQ, nil
	}
}

type message struct {
	JsonClass  string          `json:"jsonClass"`
	Body       json.RawMessage `json:"body"`
	Headers    amqp.Table      `json:"headers"`
	RoutingKey string          `json:"routingKey"`
	ReplyTo    string          `json:"replyTo,omitempty"`
}

func SendMessageToCore(body []byte, routingKey string) error {
	Loger.Debug("SendMessageToCore started")
	Loger.Debug("routingKey = ", routingKey)
	pb, err := GetPublisher()
	if err != nil {
		Loger.Error("err = ", err)
		Loger.Debug("GetPublisher failed")
		return err
	}
	var event message
	event.JsonClass = "Event"
	event.RoutingKey = routingKey
	event.Body = body
	event.Headers = make(map[string]interface{})
	event.Headers["Accept-Language"] = "ru"
	eventByte, err := json.Marshal(event)
	if err != nil {
		Loger.Error("err = ", err)
		Loger.Debug("json.Marshal(event) failed")
		return err
	}
	err = (*pb).Publish(amqp.Message{Exchange: "X:gateway.in.fanout", RoutingKey: event.RoutingKey, Body: eventByte})
	if err != nil {
		Loger.Error("err = ", err)
		Loger.Debug("(*pb).Publish failed")
		return err
	}
	Loger.Debug("SendMessageToCore successfully executed")
	return nil
}

func SendProductsToUpdateCount(products []models.InventoryProduct, merchantId string) error {
	sn := make([]ShortNomen, 0, len(products))

	for i := range products {
		sn = append(sn, ShortNomen{
			Barcode: products[i].Barcode,
			Amount:  products[i].ActualAmount,
		})
	}

	req := UpdateProductsCountRequest{
		MerchantId: merchantId,
		Action:     "UPD",
		Nomens:     sn,
	}

	reqByte, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return SendMessageToCore(reqByte, "event.catalog.update.nomenclatures")
}

type UpdateProductsCountRequest struct {
	MerchantId string       `json:"merchant_id"`
	Action     string       `json:"action"`
	Nomens     []ShortNomen `json:"nomens"`
}

type ShortNomen struct {
	Barcode string
	Amount  float64
}
