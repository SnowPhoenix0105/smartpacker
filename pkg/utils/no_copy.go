package utils

import "sync"

type NoCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (nc NoCopy) Lock()   {}
func (nc NoCopy) Unlock() {}

var _ sync.Locker = NoCopy{}
