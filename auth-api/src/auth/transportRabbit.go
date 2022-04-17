package auth

import (
	"context"
	"encoding/json"
	"github.com/djumanoff/amqp"
	"github.com/go-kit/kit/endpoint"
	"github.com/mukhametkaly/Diploma_project/auth-api/src/models"
)

type RBServer struct {
	e   endpoint.Endpoint
	dec RBDecoderRequestFunc
	enc RBEncodeResponseFunc
}

func getEndpointsMap(ss Service, routingKey string) (RBServer, error) {
	switch routingKey {
	case "request.auth.get.userinfo":
		return RBServer{e: makeAuthEndpoint(ss), dec: decodeGetUserInfoRBMQRequest, enc: encodeGetUserInfoRBMQResponse}, nil
	default:
		return RBServer{}, newErrorString(0, "no such routing key "+routingKey)
	}
}

type RBDecoderRequestFunc func(context.Context, amqp.Message) (request interface{}, err error)

type RBEncodeResponseFunc func(context.Context, interface{}, string) (*amqp.Message, error)

func MakeAuthServiceRabbitMQ(s Service, d amqp.Message) *amqp.Message {
	Loger.Debug("In MakeAuthServiceRabbitMQ")
	serv, err := getEndpointsMap(s, d.RoutingKey)
	request, err := serv.dec(nil, d)
	if err != nil {
		Loger.Error(err)
		return nil
	}
	response, err := serv.e(nil, request)
	if err != nil {
		Loger.Error(err)
		return nil
	}
	result, err := serv.enc(nil, response, d.ReplyTo)
	if err != nil {
		Loger.Error(err)
		return nil
	}
	return result
}

func decodeGetUserInfoRBMQRequest(_ context.Context, mes amqp.Message) (interface{}, error) {
	Loger.Debug("IN decodeStockSkusRBMQRequest")

	var token string
	err := json.Unmarshal(mes.Body, &token)
	if err != nil {
		Loger.Debugln("umarshaling error")
		return nil, err
	}

	return token, nil
}

func encodeGetUserInfoRBMQResponse(_ context.Context, obj interface{}, routingKey string) (*amqp.Message, error) {
	request := obj.(*models.User)
	orderByte, err := json.Marshal(*request)
	if err != nil {
		Loger.Debugln("marshalling error")
		return nil, err
	}

	return &amqp.Message{Body: orderByte}, nil
}
