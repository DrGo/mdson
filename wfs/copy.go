package wfs

//

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// ReadFile reads the file named by path from fs and returns the contents.
func ReadFile(fsys fs.FS, path string) ([]byte, error) {
	rc, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

func copyFromFSToWriter(w io.Writer, fsys fs.FS, path string)(written int64, err error) {
	rc, err := fsys.Open(path)
	if err != nil {
		return -1, err
	}
	defer rc.Close()
	return io.Copy(w, rc)
}

//soruce Ross
//FlushFS copy an fs to a dir in an OS filesystem
// it will create dir if it does not exist and overwrite
// any files with same path 
func FlushFS(dir string, fsys fs.FS) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		targ := filepath.Join(dir, filepath.FromSlash(path))// full dest path=dir+ this file's path 
		if d.IsDir() {
			if err := os.MkdirAll(targ, 0777); err != nil {
				return err
			}
			return nil
		}
		r, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()
		info, err := r.Stat()
		if err != nil {
			return err
		}
		w, err := os.OpenFile(targ, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666|info.Mode()&0777)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, r); err != nil {
			w.Close()
			return fmt.Errorf("copying %s: %v", path, err)
		}
		if err := w.Close(); err != nil {
			return err
		}
		return nil
	})
}
