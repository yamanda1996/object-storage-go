package cron

import "context"

func Cron() {
	ctx := context.Background()
	go HeartBeat(ctx)
}