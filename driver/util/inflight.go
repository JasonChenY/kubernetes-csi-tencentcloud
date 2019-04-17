package util

import (
	"sync"
)

// Idempotent is the interface required to manage in flight requests.
type Idempotent interface {
	// The CSI data types are generated using a protobuf.
	// The generated structures are guaranteed to implement the Stringer interface.
	// Example: https://github.com/container-storage-interface/spec/blob/master/lib/go/csi/csi.pb.go#L3508
	// We can use the generated string as the key of our internal inflight database of requests.
	String() string
}

// InFlight is a struct used to manage in flight requests.
type InFlight struct {
	mux      *sync.Mutex
	inFlight map[string]bool
}

// NewInFlight instanciates a InFlight structures.
func NewInFlight() *InFlight {
	return &InFlight{
		mux:      &sync.Mutex{},
		inFlight: make(map[string]bool),
	}
}

// Insert inserts the entry to the current list of inflight requests.
// Returns false when the key already exists.
func (db *InFlight) Insert(entry Idempotent) bool {
	db.mux.Lock()
	defer db.mux.Unlock()

	hash := entry.String()

	_, ok := db.inFlight[hash]
	if ok {
		return false
	}

	db.inFlight[hash] = true
	return true
}

// Delete removes the entry from the inFlight entries map.
// It doesn't return anything, and will do nothing if the specified key doesn't exist.
func (db *InFlight) Delete(h Idempotent) {
	db.mux.Lock()
	defer db.mux.Unlock()

	delete(db.inFlight, h.String())
}
