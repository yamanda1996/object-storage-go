package model

var Config = struct {
	DataServerConfig DataServerConfig
}{}

type DataServerConfig struct {
	RabbitMqAddress 				string
	RabbitMqPort					int
	RabbitMqUser					string
	RabbitMqPwd 					string
}




