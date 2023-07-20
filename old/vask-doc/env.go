package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
)

// Env is an execution environment
type Env struct {
	fs      *memfs.MemFS
	options *Options
}

// NewEnv creates an execution environment
func NewEnv(options *Options) (*Env, error) {
	env := Env{
		options: options,
		fs:      memfs.Create(),
	}
	env.fs.Mkdir("/tmp", 0777)
	return &env, nil
}

// CreateFile opens a virtual memory file for both reading and writing and create it if it does not exist
func (env *Env) CreateFile(fileName string) (vfs.File, error) {
	return env.fs.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0)
}

func (env *Env) cat(fileName string) error {
	fh, err := env.fs.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	fmt.Println("*********** Start of " + fileName + "************")
	_, err = io.Copy(os.Stdout, fh)
	if err != nil {
		return err
	}
	fmt.Println("*********** End of " + fileName + "************")

	return nil
}

// GetSourceFilePath returns a full path to a source file
func (env *Env) GetSourceFilePath(fileName string) string {
	switch ext := strings.ToLower(filepath.Ext(fileName)); ext {
	case ".gohtml", ".html":
		return env.options.srcDir + "layout/html/" + filepath.Base(fileName)
	case ".css":
		return env.options.srcDir + "layout/css/" + filepath.Base(fileName)
	case ".mdson":
		return env.options.srcDir + "contents/" + filepath.Base(fileName)
	default:
		return ""
	}
}

func (env *Env) CopyFilesFromOS(srcDir string, srcfiles []string, destDir string) error {
	if len(srcfiles) == 0 {
		return fmt.Errorf("must specify at least one source file")
	}
	for _, f := range srcfiles {
		out, _ := env.CreateFile(filepath.Join(destDir, f))
		defer out.Close()
		in, err := os.Open(filepath.Join(srcDir, f))
		if err != nil {
			return fmt.Errorf("cannot read file %s/%s: %s", srcDir, f, err)
		}
		defer in.Close()
		_, err = io.Copy(out, in)
		if err != nil {
			return fmt.Errorf("cannot copy to file %s/%s: %s", destDir, f, err)
		}
	}
	return nil
}

func (env *Env) concatFiles(srcDir string, srcfiles []string, out io.Writer) error {
	if len(srcfiles) == 0 {
		return fmt.Errorf("must specify at least one source file")
	}
	for _, f := range srcfiles {
		buf, err := vfs.ReadFile(env.fs, filepath.Join(srcDir, f))
		if err != nil {
			return fmt.Errorf("cannot read file: %s", err)
		}
		out.Write(buf)
		if err != nil {
			return fmt.Errorf("cannot write to file: %s", err)
		}
	}
	return nil
}

// func ReadDir(dirname string) ([]os.FileInfo, error) {
// 	f, err := os.Open(dirname)
// 	if err != nil {
// 		return nil, err
// 	}
// 	list, err := f.Readdir(-1)
// 	f.Close()
// 	if err != nil {
// 		return nil, err
// 	}
// 	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
// 	return list, nil
// }

// func main() {
// 	files, err := ioutil.ReadDir(".")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, file := range files {
// 		fmt.Println(file.Name())
// 	}
// }
