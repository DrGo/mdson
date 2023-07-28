package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

//go:embed wordlist.txt
var wordsFile []byte

// Generate random mdson files for testing
//Algorithm:
//- generate an mdson struct
//	-
//- output to a file
//

const MaxWord= 5101  //words 
const MaxLineLen = 20 //words
const MaxLineCount = 80 //lines per file 

var wordList = make([]string,0, MaxWord) // loaded from knownWordsFiles
var classlist=[...]string{"upper","middle","lower"}

var count = flag.Int("n", 100, "number of files to gen; default=100")

func main() {
	flag.Parse()
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	for _, w := range read("<wordsFile>", bytes.NewReader(wordsFile), bufio.ScanWords) {
		wordList= append(wordList, w)
	}

	// fmt.Printf("%s\n", wordList)
	if len(flag.Args()) == 0 {
		// add("<stdin>", os.Stdin)
	}
	for i := 0; i < *count; i++ {
		genFile(genFileName())
	}

}


func toss(prob float32) bool {
	return rand.Float32() < prob 
}

const MaxFileName = 6

//TODO: add attribs and attrib refs 

func genFile(filename string) error {
	var w  bytes.Buffer
	genContent(filename, &w, rand.Intn(MaxLineCount))
	return os.WriteFile(filename, w.Bytes(), 0664)
}

func genFileName() string {
	s := ""
	for i := 0; i < MaxFileName ; i ++ {
		s += wordList[rand.Intn(MaxWord)]
	}
	return "gen/"+ s+ ".mdson"
}

func outSentence(w io.Writer, n int, fixed string ) {
	if n==0 {
		return
	}	
	fpos:= rand.Intn(n)
	for i := 0; i  < n; i ++ {
		//ignoring errors
		if i == fpos {
			fmt.Fprint(w, fixed)
		}		
		fmt.Fprint(w, wordList[rand.Intn(MaxWord)], " ")
		// fmt.Printf("i=%d\n", i)
	}
	fmt.Fprintln(w, "")
}

func outHeader(w io.Writer, l int){
	fmt.Fprintln(w, "")	
	fmt.Fprint(w, strings.Repeat("#", l), " ")
	outSentence(w,10,"")
}

func genAttribs(w io.Writer, length int) (attribs []string) {	
	// fmt.Fprintln(w, "")	
	for i := 0; i < length; i++ {
		at := wordList[rand.Intn(MaxWord)] 
		fmt.Fprint(w, "." +at + ":")
		outSentence(w,3,"")
		attribs=append(attribs, at)
	}
	fmt.Fprintln(w, "")	
	return attribs
}	
func genList(w io.Writer, length int) error {	
	fmt.Fprintln(w, "")	
	fmt.Fprint(w, "~" )
	outSentence(w,10,"")
	for i := 0; i < length; i++ {
		fmt.Fprint(w, "  -" )
		outSentence(w, rand.Intn(MaxLineLen), "")
	}
	fmt.Fprintln(w, "")	
	return nil 
}	
func genContent(filename string, w io.Writer, length int) error {
	fmt.Fprintln(w, ".name:", filename)
	fmt.Fprintln(w, ".class:", classlist[rand.Intn(3)])
	attribs:= genAttribs(w, rand.Intn(6))
	attribs=append(attribs, "name","class")
	outHeader(w, 1)
	for i := 0; i < length; i++ {
		// output a header 
		if toss(.05) {
			outHeader(w, 2)
		} 
		if toss(0.1) {
			outHeader(w,3)
		} 
		if toss(.05) {
			genList(w, rand.Intn(MaxLineLen))
		} else {
			// output a text line
			ref := ""
			if toss(.1) {
				ref= "{"+ attribs[rand.Intn(len(attribs))] + "} "	
			}	
			outSentence(w, rand.Intn(MaxLineLen), ref)
			// fmt.PrAintf("i=%d\n", i)
		}
	}
	return nil 
}	



func read(file string, r io.Reader, split bufio.SplitFunc) []string {
	if r == nil {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stdout, "typo: %s\n", err)
			os.Exit(2)
		}
		defer f.Close()
		r = f
	}
	scanner := bufio.NewScanner(r)
	scanner.Split(split)
	words := make([]string, 0, 1000)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stdout, "typo: reading %s: %s\n", file, err)
		os.Exit(2)
	}
	return words
}
