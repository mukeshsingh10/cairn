package disk

import "github.com/mukeshsingh10/cairndb/common"

// DiskRequest is a unit of work for DiskScheduler: read into
// Data or write data to the page
//
// The Caller owns the Data and must not touch in between Schedule
// and Wait: the background goroutines reads or writes it concurrently.
type DiskRequest struct {
	IsWrite bool
	PageID  common.PageID
	Data    []byte

	// done receives exactly one value (nil on sucess) when the request has
	// been carried out. It is buffered so worker never blocks on it.
	done chan error
}

// NewDiskRequest builds a request ready to be passed to Scheduler.
func NewDiskRequest(isWrite bool, pageID common.PageID, data []byte) *DiskRequest {
	return &DiskRequest{
		IsWrite: isWrite,
		PageID:  pageID,
		Data:    data,
		done:    make(chan error),
	}
}

// Wait blocks until the request has completed and returns its IO error
// Call it only once per request.
func (r *DiskRequest) Wait() error {
	return <-r.done
}
