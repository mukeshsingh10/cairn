package page

import (
	"encoding/binary"

	"github.com/mukeshsingh10/cairndb/common"
)

const (
	HeaderSize = 8
	SlotSize   = 4
)

// TablePage is slotted-page view over one page's bytes. It owns no
// memory: wrap page.Data() (or any PageSize slice) with NewTablePage.
type TablePage struct {
	data []byte
}

func NewTablePage(data []byte) *TablePage {
	if len(data) != common.PageSize {
		panic("NewTablePage: data must be of PageSize bytes")
	}
	return &TablePage{data: data}
}

func (tp *TablePage) Init() {
	tp.setTupleStart(common.PageSize)
	tp.setNumSlots(0)
	tp.SetNextPageID(int32(common.InvalidPageID))
}

func (tp *TablePage) SetNextPageID(id int32) {
	binary.LittleEndian.PutUint32(tp.data[0:4], uint32(id))
}

func (tp *TablePage) NextPageID() int32 {
	return int32(binary.LittleEndian.Uint32(tp.data[0:4]))
}

func (tp *TablePage) NumSlots() int {
	return tp.numSlots()
}

func (tp *TablePage) FreeSpace() int {
	return tp.tupleStart() - (HeaderSize + tp.numSlots()*SlotSize)
}

// NumSlots() returns how many slots are allocated including tombstone.
func (tp *TablePage) numSlots() int {
	return int(binary.LittleEndian.Uint16(tp.data[4:6]))
}

func (tp *TablePage) setNumSlots(n int) {
	binary.LittleEndian.PutUint16(tp.data[4:6], uint16(n))
}

func (tp *TablePage) tupleStart() int {
	return int(binary.LittleEndian.Uint16(tp.data[6:8]))
}

func (tp *TablePage) setTupleStart(off int) {
	binary.LittleEndian.PutUint16(tp.data[6:8], uint16(off))
}
