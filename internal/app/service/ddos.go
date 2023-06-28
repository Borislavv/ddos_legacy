package service

import "sync"

type DDOS struct {
	mu *sync.Mutex
}
