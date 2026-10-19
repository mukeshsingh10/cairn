package disk

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/mukeshsingh10/cairndb/common"
)

type DiskManager struct {
	file      *os.File
	numWrites atomic.Int64
	numReads  atomic.Int64
}

func NewDiskManager(path string) (*DiskManager, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("open db file: %w", err)
	}

	return &DiskManager{file: f}, nil
}

func (d *DiskManager) WritePage(id common.PageID, data []byte) error {
	if id < 0 {
		return fmt.Errorf("write page: invalid page id %d", id)
	}
	if len(data) != common.PageSize {
		return fmt.Errorf("write page: got %d bytes want %d", len(data), common.PageSize)
	}
	if _, err := d.file.WriteAt(data, int64(id)*common.PageSize); err != nil {
		return fmt.Errorf("write page: %d: %w", id, err)
	}
	if err := d.file.Sync(); err != nil {
		return fmt.Errorf("write page: %d: %w", id, err)
	}
	d.numWrites.Add(1)
	return nil
}

func (d *DiskManager) ReadPage(id common.PageID, buf []byte) error {
	if id < 0 {
		return fmt.Errorf("read page: invalid page id %d", id)
	}
	if len(buf) != common.PageSize {
		return fmt.Errorf("read page: buffer is %d want %d", len(buf), common.PageSize)
	}
	n, err := d.file.ReadAt(buf, int64(id)*common.PageSize)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("read page: %d: %w", id, err)
	}
	for i := n; i < len(buf); i++ {
		buf[i] = 0
	}
	d.numReads.Add(1)
	return nil
}

func (d *DiskManager) NumWrites() int64 {
	return d.numWrites.Load()
}

func (d *DiskManager) NumReads() int64 {
	return d.numReads.Load()
}

func (d *DiskManager) Close() error {
	return d.file.Close()
}
