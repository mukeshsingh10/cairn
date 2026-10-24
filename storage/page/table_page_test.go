package page

import (
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
