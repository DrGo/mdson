# MDSon Markdown Simple Object Notation

## Introduction
- MDSon is a subset of the Markdown (MD) format optimized for defining application configration/settings. It is designed to be easy to read and write for non-programmers, although it might be a bit more challenging to parse, mainly because objects and arrays are implictely defined, e.g.,objects and array are not enclosed by brackets of any kind. 
- An MDSon source (file or stream) consists of simple text organized into blocks using MD headings.
- Each block corresponds to an object and starts with an MD heading which consists of one or more hashes followed by a unique name, e.g., # Document, ## Section.
- There must be only one root block (object) indicated by single #, e.g., # Document. Every other element (including other blocks) is a child of the root block.
- Any block including the root block, can enclose any number of elements including other blocks or lists of blocks, key-value pairs, list (or maps) of strings, numbers or boolean values and free-text elements.
- empty lines and lines starting with '//' are ignored. Useful for adding comments.
- white space is ingored
- Block headings and key names are not case-sensitive

## Blocks
- Under the root object, objects can be nested to any depth, e.g.,
```
# root
key: value
## nested within root
key: value
## also nested within root
key: value
### nested within the above 2nd level heading object
key: value
### also nested within the above 2nd level heading object
key: value
## nested within root 
key: value
```
## lists of objects
- If a block heading has a suffix of ' List' or ' list', it will be parsed as an array of one or more blocks (objects) of the same type, eg 
```
## Sections List
### Section1
### Section2
```
- Key-value pairs consist of an identifier (key) followed by : and optionally by a value. They are decoded into fields in the enclosing block object. eg
- keys are strings and must start with a character followed optionally by any number of characters or digits.
- values can be string, integer, float or bool
- bool: true is interpreted as truthy anything else including empty values is falsy.
- quoting is optional but can be used to enforce a string type for the KV pair
- simple lists of primitves: if a KV entry consists only of a key followed by : (ie no value) and the following lines each start with '-', the KV entry is decode into an array consisting of all the items with '-'. 
- simple maps of primtives: If a block heading has a suffix of ' List' or ' List', it will be parsed as an array of all blocks    

## Free text
- Any enclosed in << >> will be copied verbatim to the corresponding field in the object except that line breaks if any are removed. So the entire entry is stored as a single line of text.
- Any enclosed in <<< >>> will be also copied verbatim to the corresponding field in the object but  line breaks are NOT removed, preserving multi-line text blocks as entered.
- In both cases << or <<< should be the first character in the text block although it may follow on the same line as the field name (ie it does not need to be in its own line.
- In both cases >> or >>> should be the last character in the text block and it also does not need to be on its own line
- exmaples 
```
freeText:
<<< `anything` here 'goes'
<head> other \n stuff </head>
some "other" test>>>
```
```
freeText:<< anything here goes >>
```


## Decoding
- Only fields thar are found in the destination type will be decoded

.DisallowUnknownFields 


- The json package only accesses the exported fields of struct types (those that begin with an uppercase letter). Therefore only the the exported fields of a struct will be present in the JSON output.
- How does Unmarshal identify the fields in which to store the decoded data? For a given JSON key "Foo", Unmarshal will look through the destination struct's fields to find (in order of preference):

An exported field with a tag of "Foo" (see the Go spec for more on struct tags),
An exported field named "Foo", or
An exported field named "FOO" or "FoO" or some other case-insensitive match of "Foo".
What happens when the structure of the JSON data doesn't exactly match the Go type?

b := []byte(`{"Name":"Bob","Food":"Pickle"}`)
var m Message
err := json.Unmarshal(b, &m)
Unmarshal will decode only the fields that it can find in the destination type. In this case, only the Name field of m will be populated, and the Food field will be ignored. This behavior is particularly useful when you wish to pick only a few specific fields out of a large JSON blob. It also means that any unexported fields in the destination struct will be unaffected by Unmarshal.

Unmarshaling that data into a FamilyMember value works as expected, but if we look closely we can see a remarkable thing has happened. With the var statement we allocated a FamilyMember struct, and then provided a pointer to that value to Unmarshal, but at that time the Parents field was a nil slice value. To populate the Parents field, Unmarshal allocated a new slice behind the scenes. This is typical of how Unmarshal works with the supported reference types (pointers, slices, and maps).

Consider unmarshaling into this data structure:

type Foo struct {
    Bar *Bar
}
If there were a Bar field in the JSON object, Unmarshal would allocate a new Bar and populate it. If not, Bar would be left as a nil pointer.