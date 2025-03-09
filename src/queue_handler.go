package main

import (
	"slices"
	"time"
)

type UploadQueueUser struct {
	Ip   string
	Time time.Time
}

var UploadQueue []UploadQueueUser

func IsIPOnQueue(ip string) bool {
	for _, data := range UploadQueue {
		if data.Ip != ip {
			continue
		}

		if data.Time.Before(time.Now()) {
			RemoveExpiredIP()
			return false
		}

		return true
	}

	return false
}

func RemoveExpiredIP() {
	currentTime := time.Now()
	for i := len(UploadQueue) - 1; i >= 0; i-- {
		if UploadQueue[i].Time.Before(currentTime) {
			UploadQueue = slices.Delete(UploadQueue, i, i+1)
		}
	}
}

func AddIPToQueue(ip string, t int) {
	data := UploadQueueUser{
		Ip:   ip,
		Time: time.Unix(time.Now().Unix()+int64(t), 0),
	}

	UploadQueue = append(UploadQueue, data)
}
