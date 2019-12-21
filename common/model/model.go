package model

type Server struct {
	ServerAddress				string
	ServerPort 					string

}

type DataServer struct {
	Server
}

type ApiServer struct {
	Server
}

