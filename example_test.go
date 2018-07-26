package mds

// func ExampleParse() {
// 	hp := NewParser(DefaultParserOptions())
// 	tokens, err := hp.Parse(strings.NewReader(correct))
// 	if err != nil {
// 		fmt.Printf("error parsing %v \n", err)
// 	}
// 	type Job struct {
// 		Version             string
// 		Command             string
// 		OverWriteOutputFile bool
// 		OutputFile          string
// 		OutputFormat        string
// 		NotInTheMap         bool
// 	}
// 	var job Job
// 	tb, ok := tokens["general"]
// 	if !ok {
// 		fmt.Println("cannot find block general in the map")
// 	}
// 	err = hp.DecodeBlock(&job, tb)
// 	if err != nil {
// 		fmt.Printf("error decoding block %v \n", err)
// 	}
// 	//fmt.Println(fmt.Sprintf("%v\n", err) == correct)
// 	fmt.Println(job.Version)
// 	fmt.Println(job.OutputFile)
// 	fmt.Println(job.OverWriteOutputFile)
// 	fmt.Println(job.NotInTheMap)
// 	// Output:
// 	// 1.0.0
// 	//
// 	// true
// 	// false
// }

const (
	correct = `+++
settings :general
Version: 1.0.0
Command: run
OverWriteOutputFile: true
OutputFileName: fromconfig.docx
//comments: empty lines ignored
+++
`
)
