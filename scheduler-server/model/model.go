package model

var Config = struct {
	RabbitMqConfig 						RabbitMqConfig
	SchedulerServerConfig 				SchedulerServerConfig
}{}

type RabbitMqConfig struct {
	RabbitMqAddress 					string
	RabbitMqPort						int
	RabbitMqUser						string
	RabbitMqPwd 						string
}

type SchedulerServerConfig struct {
	SchedulerServerAddress 				string
	SchedulerServerPort 				int
}


