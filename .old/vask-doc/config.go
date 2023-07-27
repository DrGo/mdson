package main

import (
	"fmt"

	"github.com/drgo/core/files"
)

// LoadConfigFile loads info from config file
func LoadConfigFile(configFileName string) error {
	var err error
	if configFileName == "" {
		if configFileName, err = files.GetFullPath("ConfigFileBaseName"); err != nil {
			return err
		}
	}
	//load configuration from config file
	config, err := parse(configFileName)
	// fmt.Println(data)
	if err != nil {
		return fmt.Errorf("failed to parse file '%s': %v", configFileName, err)
	}

	// if job.RosewoodSettings.Debug >= rosewood.DebugAll {
	fmt.Printf("current configuration: \n %s\n", config)
	// }

	// if err = DoRun(job); rosewood.Errors().IsParsingError(err) {
	// 	err = fmt.Errorf("one or more errors occurred during file processing") //do not report again
	// }
	return err
}

// func ExtractOptions(config mdson.Node) (*Options, error) {
// 	config.ChildByName("")
// }
