package tuple

import "github.com/mukeshsingh10/cairndb/common"

// Tuple is a row as the storage sees it: an opaque byte slice plus
// RID where it lives. Interpreting the bytes as a column is the job of the higher
// level (schema and values arrive with the execution project).
type Tuple struct {
	RID  common.RID
	Data []byte
}

// NewTuple() copies the data so that tuple doesn't alias a page's memory.
func NewTuple(rid common.RID, data []byte) Tuple {
	cp := make([]byte, len(data))
	copy(cp, data)
	return Tuple{
		RID:  rid,
		Data: cp,
	}
}

func (t *Tuple) Size() int {
	return len(t.Data)
}
