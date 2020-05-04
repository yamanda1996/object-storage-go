package model

var Config = struct {
	RabbitMqConfig 					RabbitMqConfig
	ApiServerConfig 				ApiServerConfig
	ElasticSearchConfig 			ElasticSearchConfig
	EtcdServerConfig 				EtcdServerConfig
}{}

type RabbitMqConfig struct {
	RabbitMqAddress 				string
	RabbitMqPort					string
	RabbitMqUser					string
	RabbitMqPwd 					string
}

type ApiServerConfig struct {
	ApiServerAddress 				string
	ApiServerPort 					string
	ApiServerIndex 					string
	ApiServerEtcdPrefix 			string
}

type ElasticSearchConfig struct {
	ElasticSearchAddress			string
	ElasticSearchPort				string
}

type EtcdServerConfig struct {
	EtcdServerAddress				string
	EtcdServerPort					string
}





