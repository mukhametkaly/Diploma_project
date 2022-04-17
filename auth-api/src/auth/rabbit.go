package auth

import (
	"github.com/djumanoff/amqp"
	"github.com/mukhametkaly/Diploma_project/auth-api/src/config"
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
	ResponseQ: "response.fortefood.basket-api",
}

func RabbitMQConStart() (*amqp.Session, error) {
	Loger.Debug("Rabbit configs ", config.AllConfigs.Rabbit)
	cfg := amqp.Config{
		Host:        config.AllConfigs.Rabbit.Host,
		VirtualHost: config.AllConfigs.Rabbit.VirtualHost,
		User:        config.AllConfigs.Rabbit.User,
		Password:    config.AllConfigs.Rabbit.Password,
		Port:        config.AllConfigs.Rabbit.Port,
		LogLevel:    config.AllConfigs.Rabbit.LogLevel,
	}

	sess := amqp.NewSession(cfg)
	err := sess.Connect()
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

func Server() (*amqp.Server, error) {
	sess, err := GetRabbitSession()
	if err != nil {
		Loger.Debugln("GetRabbitSession:", err)
		return nil, err

	}
	if ServerRabbitMQ == nil {
		srv, err := (*sess).Server(srvCfg)
		if err != nil {
			Loger.Debugln("GetRabbitSession:", err)
			return nil, err

		} else {
			ServerRabbitMQ = &srv
			return ServerRabbitMQ, nil
		}
	} else {
		return ServerRabbitMQ, nil
	}
}

func GetRabbitClient() (*amqp.Client, error) {
	sess, err := GetRabbitSession()
	if err != nil {
		Loger.Debugln("GetRabbitSession:", err)
		return nil, err
	}

	if ClientRabbitMQ == nil {
		clt, err := (*sess).Client(cltCfg)
		if err != nil {
			Loger.Debugln("Client:", err)
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
		}
		SessionRabbitMQ = sess
		return SessionRabbitMQ, nil
	}

	return SessionRabbitMQ, nil
}

func GetPublisher() (*amqp.Publisher, error) {
	sess, err := GetRabbitSession()
	if err != nil {
		Loger.Debugln("GetRabbitSession error", err)
		return nil, err
	}
	if PublisherRabbitMQ == nil {
		pb, err := (*sess).Publisher(pubCfg)
		if err != nil {
			Loger.Debugln("GetRabbitSession error", err)
			return nil, err
		} else {
			PublisherRabbitMQ = &pb
			return PublisherRabbitMQ, nil
		}
	} else {
		return PublisherRabbitMQ, nil
	}
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

	err = (*pb).Publish(amqp.Message{Exchange: "X:gateway.in.fanout", RoutingKey: routingKey, Body: body})
	if err != nil {
		Loger.Error("err = ", err)
		Loger.Debug("(*pb).Publish failed")
		return err
	}
	Loger.Debug("SendMessageToCore successfully executed")
	return nil
}

func GetResponseViaRabbit(routingKey string, jsonStr []byte) ([]byte, error) {
	clt, err := GetRabbitClient()
	if err != nil {
		return nil, err
	}

	reply, err := (*clt).Call(routingKey, amqp.Message{Body: jsonStr})
	if err != nil {
		return nil, err
	}

	return reply.Body, nil
}
