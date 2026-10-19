package disk

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/mukeshsingh10/cairndb/common"
)

func newDM(t *testing.T) (*DiskManager, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	dm, err := NewDiskManager(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { dm.Close() })
	return dm, path
}

func pageOf(b byte) []byte {
	return bytes.Repeat([]byte{b}, common.PageSize)
}

func TestReadThenWrite(t *testing.T) {
	dm, _ := newDM(t)

	if err := dm.WritePage(0, pageOf(0xAA)); err != nil {
		t.Fatal(err)
	}
	if err := dm.WritePage(3, pageOf(0xBB)); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, common.PageSize)
	for id, want := range map[common.PageID]byte{0: 0xAA, 1: 0, 2: 0, 3: 0xBB} {
		if err := dm.ReadPage(id, buf); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf, pageOf(want)) {
			t.Errorf("page %d: wrong contents", id)
		}
	}
}
