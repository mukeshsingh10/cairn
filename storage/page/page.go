package page

import (
	"sync"

	"github.com/mukeshsingh10/cairndb/common"
)

type Page struct {
	data     [common.PageSize]byte
	pageID   common.PageID
	pinCount int
	dirty    bool
	latch    sync.RWMutex
}

func newPage() *Page {
	return &Page{pageID: common.InvalidPageID}
}

func (p *Page) Data() []byte {
	return p.data[:]
}

func (p *Page) ID() common.PageID {
	return p.pageID
}

func (p *Page) SetID(id common.PageID) {
	p.pageID = id
}

func (p *Page) PinCount() int {
	return p.pinCount
}

func (p *Page) SetPinCount(n int) {
	p.pinCount = n
}

func (p *Page) IsDirty() bool {
	return p.dirty
}

func (p *Page) SetDirty(dirty bool) {
	p.dirty = dirty
}

func (p *Page) RLatch() {
	p.latch.RLock()
}

func (p *Page) RUnlatch() {
	p.latch.RUnlock()
}

func (p *Page) WLatch() {
	p.latch.Lock()
}

func (p *Page) WUnlatch() {
	p.latch.Unlock()
}

// ResetMemory() zeros the data, as done when frame is reused for another page.
func (p *Page) ResetMemory() {
	p.data = [common.PageSize]byte{}
}
