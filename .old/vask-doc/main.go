package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Options stores configuration
type Options struct {
	srcDir   string
	buildDir string
	layout   string
	contents string
	Title    string //website title
}

// const ConfigFileBaseName = "vask.mdson"

const testPath = "test-files/"

var options = &Options{
	srcDir:   testPath + "src/",
	buildDir: testPath + "build/",
	layout:   testPath + "src/layout/",
	contents: testPath + "src/contents/",
}

func main() {
	// log.Println("started")
	// Create a fully writable filesystem in memory
	env, _ := NewEnv(options)
	_ = env

	files, err := ioutil.ReadDir(options.contents)
	check(err)
	dirNames := []string{}
	for _, file := range files {
		if file.IsDir() {
			dirNames = append(dirNames, filepath.Base(file.Name()))
		}
	}

	for _, d := range dirNames {
		// get mdson files from every directory under contents
		srcDir := filepath.Join(options.contents, d)
		fmt.Println("processing:", srcDir)
		c, err := NewCollection(srcDir + "/*.mdson")
		check(err)
		// create corresponding build folder under options.build
		buildDir := filepath.Join(options.buildDir, d)
		fmt.Println("create build dir", buildDir)
		err = os.MkdirAll(buildDir, 666)
		check(err)
		c.Values["buildDir"] = buildDir

		// parse the corresponding template file name for this directory
		templateFileName := filepath.Join(options.layout, d, d+".html")
		t, err := template.New(filepath.Base(templateFileName)).Funcs(funcList).ParseFiles(templateFileName)
		check(err)
		c.Values["template"] = t
		rs, err := c.ForEach(render)
		check(err)
		_, err = io.Copy(os.Stdout, rs[0])
		check(err)
	}

	// if err := runTemplate(env, "team.gohtml", "team.mdson"); err != nil {
	// 	log.Fatalln("failed to parse template", err)
	// }
	// if err := runTemplate(env, "collaborators.gohtml", "collaborators.mdson"); err != nil {
	// 	log.Fatalln("failed to parse template", err)
	// }
	// // if err := env.cat("/tmp/collaborators.html"); err != nil {}
	// // inline css
	// s, err := inlineCSS(env.GetSourceFilePath("header.html"), env.GetSourceFilePath("styles.css"))
	// if err != nil {
	// 	log.Fatalln("failed to inline css", err)
	// }

	// if err := vfs.WriteFile(env.fs, "tmp/header.html", s, 0644); err != nil {
	// 	log.Fatalln("failed to inline css", err)
	// }
	// srcFiles := []string{"banner.html", "footer.html"}
	// // copy rest of files as is
	// err = env.CopyFilesFromOS(options.srcHTML, srcFiles, "/tmp/")
	// if err != nil {
	// 	log.Fatalln(fmt.Errorf("failed to copy files to tmp: %v", err))
	// }

	// // concat all html files into an index.html in the build dir
	// buildFileName := filepath.Join(options.destDir, "index.html")
	// buildOut, err := os.Create(buildFileName)
	// if err != nil {
	// 	log.Fatalln("failed to create build files", err)
	// }
	// defer buildOut.Close()

	// // indexFileName := "/tmp/index.html"
	// // out, err := env.CreateFile(indexFileName)
	// // if err != nil {
	// // 	log.Fatalln(fmt.Errorf("failed to create output file '%s': %v", indexFileName, err))
	// // }
	// // defer out.Close()
	// srcFiles = []string{"header.html", "banner.html", "team.html", "collaborators.html", "footer.html"}
	// if err := env.concatFiles("/tmp/", srcFiles, buildOut); err != nil {
	// 	log.Fatalln("failed to concat files", err)
	// }
	// // // copy /tmp/index.html to build/index.html

	// // _, err = io.Copy(buildOut, out)
	// // if err != nil {
	// // 	log.Fatalln("failed to create build files", err)
	// // }
	// // out.Sync()
	// // out.Close()
	// buildOut.Sync()
	// err = buildOut.Close()
	// if err != nil {
	// 	log.Fatalln("failed to create build files", err)
	// }
	// //env.cat("/tmp/index.html")
	// // outFileName := options.destDir + "index.html"
	// // out, err := os.Create(outFileName)
	// // if err != nil {
	// // 	log.Fatalln("failed to open output file", options.srcDir, err)
	// // }
	// // defer out.Close()

	// log.Println("success")
}

// func main() {
// 	if err := RunApp(); err != nil {
// 		log.Fatalln(err)
// 	}
// }

//RunApp has all program logic; entry point for all tests
//WARNING: not thread-safe; this is the only function allowed to change the job configuration
func RunApp() error {
	// if len(os.Args) == 1 { //no command line arguments
	// 	return DoFromConfigFile("")
	// }
	// exeName := os.Args[0]
	// // job, err := LoadConfigFromCommandLine()
	// // if err != nil {
	// // 	return err
	// // }
	// // if job.RosewoodSettings.Debug == rosewood.DebugAll {
	// // 	fmt.Printf("current settings:\n%s\n", job)
	// // }
	// switch job.Command { //TODO: check command is case insensitive
	// case "do":
	// 	if len(job.RwFileNames) == 0 {
	// 		return fmt.Errorf("must specify an MDSon configuration file")
	// 	}
	// 	if err = DoFromConfigFile(job.RwFileNames[0]); err != nil {
	// 		return err
	// 	}
	// case "check":
	// 	job.RosewoodSettings.CheckSyntaxOnly = true
	// 	fallthrough
	// case "run":
	// 	job.Command = "process" //change to print nicer messages
	// 	//FIXME: this check is not working
	// 	if err = DoRun(job); rosewood.Errors().IsParsingError(err) {
	// 		err = fmt.Errorf("one or more errors occurred during file processing") //do not report again
	// 	}
	// case "init":
	// 	// configFilename, err := DoInit(job)
	// 	// if err == nil && job.RosewoodSettings.Debug >= rosewood.DebugUpdates {
	// 	// 	fmt.Printf("configuration saved as '%s'\n", configFilename)
	// 	// }
	// 	return err
	// 	// case "version":
	// 	// 	fmt.Println(getVersion())
	// 	// case "help", "h":
	// 	// 	helpMessage(job.RwFileNames, getVersion())
	// 	// default:
	// 	// 	helpMessage(nil, getVersion())
	// 	// 	return fmt.Errorf(ErrWrongCommand, exeName)
	// }
	return nil
}

// //LoadConfigFromCommandLine creates a object using command line arguments
// func LoadConfigFromCommandLine() (*rosewood.Job, error) {
// 	job := rosewood.DefaultJob(rosewood.DefaultSettings()) //TODO: ensure all defaults are reasonable
// 	flgSets, _ := setupCommandFlag(job)
// 	flg, err := args.ParseCommandLine(flgSets[0], flgSets[1:]...)

// 	if err != nil {
// 		return nil, err
// 	}
// 	job.Command = flg.Name()
// 	if len(flg.Args()) == 0 {
// 		return job, nil
// 	}
// 	switch runtime.GOOS {
// 	case "windows":
// 		job.RwFileNames, err = WinArgsToFileNames(flg.Args()[0])
// 		if err != nil {
// 			return nil, err
// 		}
// 	default:
// 		for _, fileName := range flg.Args() {
// 			job.RwFileNames = append(job.RwFileNames, fileName)
// 		}
// 	}
// 	return job, nil
// }

func WinArgsToFileNames(args string) ([]string, error) {
	if !strings.ContainsAny(args, "*?") {
		return []string{args}, nil
	}
	return filepath.Glob(args)
}

func cat(c *Collection, r io.Reader, index int) (io.Reader, error) {
	fmt.Println(c.GetFileName(index))
	var buf bytes.Buffer
	w := io.MultiWriter(&buf, os.Stdout)
	_, err := io.Copy(w, r)
	return &buf, err
}

func render(c *Collection, r io.Reader, index int) (io.Reader, error) {
	srcFileName := c.GetFileName(index)
	fmt.Println("rendering", srcFileName)
	// buildDir := c.Values["buildDir"].(string)
	// destFileName := filepath.Join(buildDir, strings.TrimSuffix(srcFileName, filepath.Ext(srcFileName))+".html")
	data, err := parse(srcFileName)
	// fmt.Println(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file '%s': %v", srcFileName, err)
	}
	t, ok := c.Values["template"].(*template.Template)
	if !ok {
		return nil, fmt.Errorf("template is not set")
	}
	// _ = t
	// _ = data
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render file '%s': %v", srcFileName, err)
	}
	return &buf, err
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
