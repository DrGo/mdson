package booker

//
// type WritableFS interface {
// 	fs.FS
// 	OpenFile(name string, flag int, perm fs.FileMode) (WritableFile, error)
// }
//
// type WritableFile interface {
// 	fs.File
// 	io.Writer
// }
//
// func Create(fsys WritableFS, name string) (WritableFile, error)
//
// func WriteFile(fsys WritableFS, name string, data []byte, perm fs.FileMode) error
//
//
// type WriteFile interface {
// 	io.ByteWriter
// 	fs.File
// }
//
//
//
// func Create(fsys FS, name string) (WriteFile, error) { }
//
// type WFile interface {
//     Stat() (FileInfo, error)
//     Write(p []byte) (n int, err error)
//     Close() error
// }
//
// type WriteFS interface {
//     OpenFile(name string, flag int, perm FileMode) (WFile, error)
// }
//
// type MkDirFS interface {
//     MkDir(name string, perm FileMode) error
// }
//
// func Create(fsys WriteFS, name string) (WFile, error) {
//     // Use fsys.OpenFile ...
// }
//
// func WriteFile(fsys WriteFS, name string, data []byte, perm FileMode) error {
//     // Use fsys.OpenFile, Write, and Close ...
// }
//
// func MkDirAll(fsys MkDirFS, path string, perm FileMode) error {
//     // Use fsys.MkDir to do the work.
//     // Also requires either Stat or Open to check for parents.
//     // I'm not sure how to structure that either/or requirement.
// }
