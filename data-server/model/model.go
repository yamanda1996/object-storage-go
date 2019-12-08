package model

var Config = struct {
	RabbitMqConfig 					RabbitMqConfig
	DataServerConfig 				DataServerConfig
}{}

type RabbitMqConfig struct {
	RabbitMqAddress 				string
	RabbitMqPort					int
	RabbitMqUser					string
	RabbitMqPwd 					string
}

type DataServerConfig struct {
	DataServerAddress 				string
	DataServerPort 					int
	StorageRoot						string
}





