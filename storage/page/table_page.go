package page

import "github.com/mukeshsingh10/cairndb/common"

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
