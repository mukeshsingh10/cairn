package page

import (
	"bytes"
	"errors"
	"testing"

	"github.com/mukeshsingh10/cairndb/common"
)

func newTP(t *testing.T) *TablePage {
	t.Helper()
	tp := NewTablePage(make([]byte, common.PageSize))
	tp.Init()
	return tp
}

func TestInit(t *testing.T) {
	tp := newTP(t)
	if tp.NumSlots() != 0 {
		t.Errorf("NumSlots=%d, want 0", tp.NumSlots())
	}
	if tp.NextPageID() != int32(common.InvalidPageID) {
		t.Errorf("NextPageID=%d, want InvalidPageID", tp.NextPageID())
	}
	if got, want := tp.FreeSpace(), common.PageSize-HeaderSize; got != want {
		t.Errorf("got=%d, want=%d", got, want)
	}
}

func TestNextPageID(t *testing.T) {
	tp := newTP(t)
	tp.SetNextPageID(43)
	if tp.NextPageID() != 43 {
		t.Errorf("NextPageID=%d, want=43", tp.NextPageID())
	}
}

func TestInsertAndGet(t *testing.T) {
	tp := newTP(t)
	tuples := [][]byte{[]byte("foo"), []byte("bar"), []byte("foo bar baz")}
	for i, tuple := range tuples {
		slot, err := tp.InsertTuple(tuple)
		if err != nil {
			t.Error(err)
		}
		if slot != uint32(i) {
			t.Errorf("slot=%d want=%d (slots are sequential)", slot, i)
		}
	}

	for i, want := range tuples {
		got, err := tp.GetTuple(uint32(i))
		if err != nil || !bytes.Equal(got, want) {
			t.Errorf("tuple %d = %q, %v want  %q", i, got, err, want)
		}
	}
}

func TestGetReturnsCopy(t *testing.T) {
	tp := newTP(t)
	slot, _ := tp.InsertTuple([]byte("hello"))
	got, _ := tp.GetTuple(slot)
	got[0] = 'X'
	again, _ := tp.GetTuple(slot)
	if string(again) != "hello" {
		t.Error("GetTuple must return a copy, not alias page memory")
	}
}

func TestInsertCopiesInput(t *testing.T) {
	tp := newTP(t)
	in := []byte("hello")
	slot, _ := tp.InsertTuple(in)
	in[0] = 'X'
	got, _ := tp.GetTuple(slot)
	if string(got) != "hello" {
		t.Error("InsertTuple must copy its input")
	}
}

func TestDelete(t *testing.T) {
	tp := newTP(t)
	s0, _ := tp.InsertTuple([]byte("keep"))
	s1, _ := tp.InsertTuple([]byte("drop"))
	if err := tp.DeleteTuple(s1); err != nil {
		t.Error(err)
	}
	if _, err := tp.GetTuple(s1); !errors.Is(err, ErrTupleDeleted) {
		t.Errorf("GetTuple(deleted) err = %v want ErrTupleDeleted", err)
	}
	if got, err := tp.GetTuple(s0); err != nil || string(got) != "keep" {
		t.Errorf("other tuple damaged: %q, %v", got, err)
	}

	// RID must stay stable: the next insert get a new slot, not the freed one.
	s2, _ := tp.InsertTuple([]byte("new"))
	if s2 != 2 {
		t.Errorf("slot after delete = %d want 2 (slot are never reused)", s2)
	}
}

func TestInvalidSlot(t *testing.T) {
	tp := newTP(t)
	if _, err := tp.GetTuple(0); !errors.Is(err, ErrInvalidSlot) {
		t.Errorf("GetTuple on empty page err = %v want ErrInvalidSlot", err)
	}
	if err := tp.DeleteTuple(7); !errors.Is(err, ErrInvalidSlot) {
		t.Errorf("DeleteTuple err = %v want ErrInvalidSlot", err)
	}
}

func TestFillPageExactly(t *testing.T) {
	tp := newTP(t)
	const size = 100
	tuple := bytes.Repeat([]byte{7}, size)
	want := (common.PageSize - HeaderSize) / (size + SlotSize)

	n := 0
	for {
		_, err := tp.InsertTuple(tuple)
		if errors.Is(err, ErrNoSpace) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		n++
		if n > common.PageSize {
			t.Error("never ran out of space")
		}
	}
	if n != want {
		t.Errorf("fit %d tuples of %d bytes, want %d", n, size, want)
	}

	for i := 0; i < n; i++ {
		got, err := tp.GetTuple(uint32(i))
		if err != nil || !bytes.Equal(got, tuple) {
			t.Errorf("tuple %d corrupted %v", i, err)
		}
	}
}

func TestTupleLargerThanPage(t *testing.T) {
	tp := newTP(t)
	if _, err := tp.InsertTuple(make([]byte, common.PageSize)); !errors.Is(err, ErrNoSpace) {
		t.Errorf("err = %v want ErrNoSpace", err)
	}
	if tp.NumSlots() != 0 {
		t.Error("failed insert must not allocate a slot")
	}
}

func TestFailedInsertLeavePageIntact(t *testing.T) {
	tp := newTP(t)
	_, _ = tp.InsertTuple([]byte("first"))
	before := tp.FreeSpace()
	_, _ = tp.InsertTuple(make([]byte, common.PageSize))
	if tp.FreeSpace() != before {
		t.Error("failed insert change FreeSpace")
	}
}

func TestSurvivesRewrapOfSameBytes(t *testing.T) {
	// The page's state must live entirely in its byte. So re-wraping the same
	// buffer (as happens after a page is written to disk and read back) works.
	buf := make([]byte, common.PageSize)
	tp := NewTablePage(buf)
	tp.Init()
	slot, _ := tp.InsertTuple([]byte("persist me"))

	tp2 := NewTablePage(buf)
	got, err := tp2.GetTuple(slot)
	if err != nil || string(got) != "persist me" {
		t.Errorf("got %q, %v after re-wrap", got, err)
	}
}
