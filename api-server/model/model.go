package model

var Config = struct {
	RabbitMqConfig 					RabbitMqConfig
	ApiServerConfig 				ApiServerConfig
}{}

type RabbitMqConfig struct {
	RabbitMqAddress 				string
	RabbitMqPort					int
	RabbitMqUser					string
	RabbitMqPwd 					string
}

type ApiServerConfig struct {
	ApiServerAddress 				string
	ApiServerPort 					int
}





