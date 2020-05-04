package cron

import "context"

func DoCron() {
	ctx := context.Background()
	go HeartBeat(ctx)
}