package model

type Server struct {
	ServerAddress				string
	ServerPort 					string
	Status 						int

}

type DataServer struct {
	Server
}

type ApiServer struct {
	Server
}

type SchedulerServer struct {
	Server

}


