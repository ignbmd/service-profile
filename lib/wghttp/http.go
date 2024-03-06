package wghttp

import "sync"

var HttpWG *sync.WaitGroup

func NewHttpWg() {
	if HttpWG == nil {
		HttpWG = new(sync.WaitGroup)
	}
}
