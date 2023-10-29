package filesystem

import "time"

type File interface {
	Disk(disk string) File
	Extension() (string, error)
	File() string
	GetClientOriginalName() string
	GetClientOriginalExtension() string
	HashName(path ...string) string
	LastModified() (time.Time, error)
	MimeType() (string, error)
	Size() (int64, error)
	Store(path string) (string, error)
	StoreAs(path string, name string) (string, error)
}
