// this file is a valid mdson file
// comments will be removed before processing	
// a property (prop) is essentially like a variable; it will be removed from the text
// but made available to the rendering program as part of the DOM
// syntax: dot. followed immediately by a letter or underscore and then colon.
// anything after the colon and until EOL is assigned as a value of type string
// all props belong to the nearest section (see below) except the ones at the beginning
// of the file before any section header is declared which belong to the root section (the document)
.parent: root 
.prop can have long name: some value
.weight: 1
.date: 12July2023
.chapter: My chapter Name 
.book: my great book
.prop with multiline values: are not allowed
 .this is not a prop because the first char is space

// you could refer to any of the above props anywhere in the doc like this {.date}
.today: Today is {date}
date {today}
last date of my book titled {book} was on {date}

-- text==Data .mdson --> booker ---> typst file + typst template --> pdf 
-- html + css

-- table of figures .
== tables of lists 
-- structure + presentation --> 

# introduction

part of my {book}
# {chapter}

// section names are regular markdown headers starting with one or more #
## introduction
// props that belong to this section
.author: Someone
// they can be referred to as {introduction.author}
// refer to the computed .today value above
date {today}

The most important cause of CHF is {chapter.chf[1]}
{chapter1.chf[1]} 
## paragraphs 
>A paragraph is consecutive lines of text with one or more blank lines between them.
>For a line break, add either a backslash \ or two blank spaces at the end of the line.

see @Figxxx --> See Figure 100 
// A list starts with ~ and continues until the next non-list item element 
// this is a list named "Causes of heart failure". The name is all the text
// before the colon and both the name and colon are mandatory
~Causes of heart failure <chf>:
//- list item starts with -
- Atrial fibrillation
- Hypertension
- Myocardial infarction
// This list can be referred to as {.Causes of heart failure}
// its first element can be referred to as {.Causes of heart failure[0]}

### Nested within Introduction
// backticks used to guard preformatted text like code blocks 
// everything else is regular text
~Another List
- item 1
- item 2
- item 3
- 
//Linking to another section or list using markdown link []()
Line 1 in subsection of introduction
Line 2 in subsection of introduction
Line 3 in subsection of introduction


