package page

import (
	"encoding/binary"
	"errors"

	"github.com/mukeshsingh10/cairndb/common"
)

const (
	HeaderSize = 8
	SlotSize   = 4
)

var (
	// ErrNoSpace means the tuple does not fit in this page's free space.
	ErrNoSpace = errors.New("table page: not enough free space")
	// ErrInvalidSlot means the slot number was never allocated.
	ErrInvalidSlot = errors.New("table page: invalid slot")
	// ErrTupleDeleted means the slot exists but its tuple was deleted.
	ErrTupleDeleted = errors.New("table page: tuple deleted")
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

func (tp *TablePage) InsertTuple(data []byte) (uint32, error) {
	need := len(data) + SlotSize
	if need > tp.FreeSpace() {
		return 0, ErrNoSpace
	}
	newStart := tp.tupleStart() - len(data)
	copy(tp.data[newStart:], data)
	slot := tp.numSlots()
	tp.setSlot(slot, newStart, len(data))
	tp.setNumSlots(slot + 1)
	tp.setTupleStart(newStart)
	return uint32(slot), nil
}

func (tp *TablePage) GetTuple(slot uint32) ([]byte, error) {
	if int(slot) >= tp.numSlots() {
		return nil, ErrInvalidSlot
	}
	offset, size := tp.getSlot(int(slot))
	if offset == 0 {
		return nil, ErrTupleDeleted
	}
	out := make([]byte, size)
	copy(out, tp.data[offset:offset+size])
	return out, nil
}

func (tp *TablePage) DeleteTuple(slot uint32) error {
	if int(slot) >= tp.numSlots() {
		return ErrInvalidSlot
	}
	offset, size := tp.getSlot(int(slot))
	if offset == 0 {
		return ErrTupleDeleted
	}
	tp.setSlot(int(slot), 0, size)
	return nil
}

func (tp *TablePage) setSlot(i, offset, size int) {
	pos := tp.slotPos(i)
	binary.LittleEndian.PutUint16(tp.data[pos:pos+2], uint16(offset))
	binary.LittleEndian.PutUint16(tp.data[pos+2:pos+4], uint16(size))
}

func (tp *TablePage) getSlot(i int) (offset, size int) {
	pos := tp.slotPos(i)
	offset = int(binary.LittleEndian.Uint16(tp.data[pos : pos+2]))
	size = int(binary.LittleEndian.Uint16(tp.data[pos+2 : pos+4]))
	return offset, size
}

func (tp *TablePage) slotPos(i int) int {
	return HeaderSize + i*SlotSize
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
