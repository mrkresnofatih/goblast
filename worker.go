package goblast

import "sync"

type BackgroundWorker struct {
	RunState *sync.WaitGroup
}
