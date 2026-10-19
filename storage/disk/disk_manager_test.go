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

func TestReadBeyondEOFIsZeros(t *testing.T) {
	dm, _ := newDM(t)
	buf := pageOf(0xFF)
	if err := dm.ReadPage(10, buf); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, pageOf(0)) {
		t.Errorf("expected zero-filled page")
	}
}

func TestPersistAccrossReopen(t *testing.T) {
	dm, path := newDM(t)
	if err := dm.WritePage(1, pageOf(0x42)); err != nil {
		t.Error(err)
	}
	dm.Close()

	diskmanager, err := NewDiskManager(path)
	if err != nil {
		t.Error(err)
	}
	buf := make([]byte, common.PageSize)
	if err := diskmanager.ReadPage(1, buf); err != nil {
		t.Error(err)
	}
	if !bytes.Equal(buf, pageOf(0x42)) {
		t.Error("data lost after repoen")
	}
}

func TestRejectBadArguments(t *testing.T) {
	dm, _ := newDM(t)
	if err := dm.WritePage(1, make([]byte, 10)); err == nil {
		t.Error("short write buffer should fail")
	}
	if err := dm.ReadPage(-1, make([]byte, common.PageSize)); err == nil {
		t.Error("negetive page id should fail")
	}
}

func TestIOCounters(t *testing.T) {
	dm, _ := newDM(t)
	_ = dm.WritePage(0, pageOf(0x11))
	_ = dm.ReadPage(0, make([]byte, common.PageSize))
	if dm.NumWrites() != 1 || dm.NumReads() != 1 {
		t.Errorf("writes=%d reads=%d, want 1 and 1", dm.NumWrites(), dm.NumReads())
	}
}
