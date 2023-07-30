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

func copyFromFSToFile(w io.Writer, fsys fs.FS, path string)(written int64, err error) {
	rc, err := fsys.Open(path)
	if err != nil {
		return -1, err
	}
	defer rc.Close()
	return io.Copy(w, rc)
}

//soruce Ross

func CopyFS(dir string, fsys fs.FS) error {
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
// from wfs
// CopyFS walks the specified root directory on src and copies directories and
// files to dest filesystem.
// func CopyFS2(dest, src fs.FS, root string) error {
// 	return fs.WalkDir(src, root, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil || d == nil {
// 			return err
// 		}
// 		if d.IsDir() {
// 			return MkdirAll(dest, path, d.Type())
// 		}
// 		srcFile, err := src.Open(path)
// 		if err != nil {
// 			return err
// 		}
// 		destFile, err := CreateFile(dest, path, d.Type())
// 		if err != nil {
// 			return err
// 		}
// 		defer destFile.Close()
//
// 		_, err = io.Copy(destFile, srcFile)
// 		return err
// 	})
// }

// // CopyFromFS copies a dir/file from an FS to wfs
// func (fsys *FS) CopyFromFS(dstPath, srcPath string, srcFS fs.FS ) error {
// 	src, err := srcFS.Open(srcPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer src.Close()
//
// 	dst, err := dstFs.Create(dstPath)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	defer func() {
// 		cerr := dst.Close()
// 		if cerr == nil {
// 			err = cerr
// 		}
// 	}()
//
// 	err = dstFs.Chmod(dstPath, info.Mode())
// 	if err != nil {
// 		return nil, err
// 	}
//
// }
//
