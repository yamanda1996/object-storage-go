package model

type HeartBeat struct {
	DataServerAddress 	string
	Timestamp 			int64
	CpuUsage 			float64
	MemUsage 			float64
	DiskUsage 			float64
}



