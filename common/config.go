package common

import "fmt"

// PageSize is the size of the database page in bytes. Every read and write
// between disk and memory happens in units of this size.
const PageSize = 4096

// PageID identifies a page by its position in database file:
// the page lives in a file at byte offset: PageID * PageSize.
type PageID int32

// FrameID identifies a slot (frame) in a buffer pool's in-memory array.
type FrameID int32

// InvalidPageID marks "no page". For example the end of linked list of pages.
const InvalidPageID PageID = -1

// RID (record id) addresses a tuple physically: which page and which slot
// inside that page's slot array.
type RID struct {
	PageID PageID
	Slot   uint32
}

func (rid *RID) String() string {
	return fmt.Sprintf("page=%d, slot=%d", rid.PageID, rid.Slot)
}
