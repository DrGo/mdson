// implements an in-memory fs
package wfs

// Copyright 2020 The Go Authors. All rights reserved.
// Modified by Salah Mahmud to support bigger files
import (
	"bytes"
	"io"
	"io/fs"
	"path"
	"sort"
	"strings"
	"time"
)

// A FS is a simple in-memory file system for use in tests,
// represented as a map from path names (arguments to Open)
// to information about the files or directories they represent.
//
// The map need not include parent directories for files contained
// in the map; those will be synthesized if needed.
// But a directory can still be included by setting the MapFile.Mode's ModeDir bit;
// this may be necessary for detailed control over the directory's FileInfo
// or to create an empty directory.
//
// File system operations read directly from the map,
// so that the file system can be changed by editing the map as needed.
// An implication is that file system operations must not run concurrently
// with changes to the map, which would be a race.
// Another implication is that opening or reading a directory requires
// iterating over the entire map, so a FS should typically be used with not more
// than a few hundred entries or directory reads.
type FS struct {
	m map[string]*WFile
}

func NewFS() FS {
	return FS{
		m: make(map[string]*WFile),
	}
}

// A WFile describes a single file in a MapFS.
type WFile struct {
	bytes.Buffer             // file content
	Mode         fs.FileMode // FileInfo.Mode
	ModTime      time.Time   // FileInfo.ModTime
	Sys          any         // FileInfo.Sys
}

var _ fs.FS = FS(FS{})
var _ fs.File = (*openWFile)(nil)
var _ io.WriteCloser = (*openWFile)(nil)

// Open opens the named file.
func (fsys FS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	file := fsys.m[name]
	if file != nil && file.Mode&fs.ModeDir == 0 {
		// Ordinary file
		return &openWFile{name, wFileInfo{path.Base(name), file}, 0}, nil

	}

	// Directory, possibly synthesized.
	// Note that file can be nil here: the map need not contain explicit parent directories for all its files.
	// But file can also be non-nil, in case the user wants to set metadata for the directory explicitly.
	// Either way, we need to construct the list of children of this directory.
	var list []wFileInfo
	var elem string
	var need = make(map[string]bool)
	if name == "." {
		elem = "."
		for fname, f := range fsys.m {
			i := strings.Index(fname, "/")
			if i < 0 {
				if fname != "." {
					list = append(list, wFileInfo{fname, f})
				}
			} else {
				need[fname[:i]] = true
			}
		}
	} else {
		elem = name[strings.LastIndex(name, "/")+1:]
		prefix := name + "/"
		for fname, f := range fsys.m {
			if strings.HasPrefix(fname, prefix) {
				felem := fname[len(prefix):]
				i := strings.Index(felem, "/")
				if i < 0 {
					list = append(list, wFileInfo{felem, f})
				} else {
					need[fname[len(prefix):len(prefix)+i]] = true
				}
			}
		}
		// If the directory name is not in the map,
		// and there are no children of the name in the map,
		// then the directory is treated as not existing.
		if file == nil && list == nil && len(need) == 0 {
			return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
		}
	}
	for _, fi := range list {
		delete(need, fi.name)
	}
	for name := range need {
		list = append(list, wFileInfo{name, &WFile{Mode: fs.ModeDir}})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].name < list[j].name
	})

	if file == nil {
		file = &WFile{Mode: fs.ModeDir}
	}
	return &wDir{name, wFileInfo{elem, file}, list, 0}, nil
}

// fsOnly is a wrapper that hides all but the fs.FS methods,
// to avoid an infinite recursion when implementing special
// methods in terms of helpers that would use them.
// (In general, implementing these methods using the package fs helpers
// is redundant and unnecessary, but having the methods may make
// MapFS exercise more code paths when used in tests.)
type fsOnly struct{ fs.FS }

func (fsys FS) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(fsOnly{fsys}, name)
}

func (fsys FS) Stat(name string) (fs.FileInfo, error) {
	return fs.Stat(fsOnly{fsys}, name)
}

func (fsys FS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(fsOnly{fsys}, name)
}

func (fsys FS) Glob(pattern string) ([]string, error) {
	return fs.Glob(fsOnly{fsys}, pattern)
}

// WriteFile writes data to the named file, creating it if necessary.
// If the file does not exist, WriteFile creates it with permissions perm (before umask);
// otherwise WriteFile truncates it before writing, without changing permissions.
func (fsys FS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	f := &WFile{
		//Buffer needs no initialization
		Mode:    perm,
		ModTime: time.Now().Local(),
	}
	f.Write(data)
	fsys.m[name] = f
	return nil
}

type noSub struct {
	FS
}

func (noSub) Sub() {} // not the fs.SubFS signature

func (fsys FS) Sub(dir string) (fs.FS, error) {
	return fs.Sub(noSub{fsys}, dir)
}

// A wFileInfo implements fs.FileInfo and fs.DirEntry for a given map file.
type wFileInfo struct {
	name string
	f    *WFile
}

func (i *wFileInfo) Name() string               { return i.name }
func (i *wFileInfo) Size() int64                { return int64(i.f.Len()) }
func (i *wFileInfo) Mode() fs.FileMode          { return i.f.Mode }
func (i *wFileInfo) Type() fs.FileMode          { return i.f.Mode.Type() }
func (i *wFileInfo) ModTime() time.Time         { return i.f.ModTime }
func (i *wFileInfo) IsDir() bool                { return i.f.Mode&fs.ModeDir != 0 }
func (i *wFileInfo) Sys() any                   { return i.f.Sys }
func (i *wFileInfo) Info() (fs.FileInfo, error) { return i, nil }
func (i *wFileInfo) String() string             { return fs.FormatFileInfo(i) }

// An openWFile is a regular (non-directory) fs.File open for reading.
type openWFile struct {
	path string
	wFileInfo
	offset int64
}

func (f *openWFile) Stat() (fs.FileInfo, error) { return &f.wFileInfo, nil }

func (f *openWFile) Close() error { return nil }

func (f *openWFile) Read(b []byte) (int, error) {
	if f.offset >= int64(f.f.Len()) {
		return 0, io.EOF
	}
	if f.offset < 0 {
		return 0, &fs.PathError{Op: "read", Path: f.path, Err: fs.ErrInvalid}
	}
	n := copy(b, f.f.Bytes()[f.offset:])
	f.offset += int64(n)
	return n, nil
}

func (f *openWFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		// offset += 0
	case 1:
		offset += f.offset
	case 2:
		offset += int64(f.f.Len())
	}
	if offset < 0 || offset > int64(f.f.Len()) {
		return 0, &fs.PathError{Op: "seek", Path: f.path, Err: fs.ErrInvalid}
	}
	f.offset = offset
	return offset, nil
}

func (f *openWFile) ReadAt(b []byte, offset int64) (int, error) {
	if offset < 0 || offset > int64(f.f.Len()) {
		return 0, &fs.PathError{Op: "read", Path: f.path, Err: fs.ErrInvalid}
	}
	n := copy(b, f.f.Bytes()[offset:])
	if n < len(b) {
		return n, io.EOF
	}
	return n, nil
}
func (f *openWFile) Write(p []byte) (int, error) {
	return f.f.Write(p)
}

// A wDir is a directory fs.File (so also an fs.ReadDirFile) open for reading.
type wDir struct {
	path string
	wFileInfo
	entry  []wFileInfo
	offset int
}

func (d *wDir) Stat() (fs.FileInfo, error) { return &d.wFileInfo, nil }
func (d *wDir) Close() error               { return nil }
func (d *wDir) Read(b []byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.path, Err: fs.ErrInvalid}
}

func (d *wDir) ReadDir(count int) ([]fs.DirEntry, error) {
	n := len(d.entry) - d.offset
	if n == 0 && count > 0 {
		return nil, io.EOF
	}
	if count > 0 && n > count {
		n = count
	}
	list := make([]fs.DirEntry, n)
	for i := range list {
		list[i] = &d.entry[d.offset+i]
	}
	d.offset += n
	return list, nil
}
