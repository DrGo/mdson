# MDSon Markdown Simple Object Notation

## Introduction
- MDSon is a subset of the Markdown (MD) format optimized for serializing arbitaritly nested data structures into Unicode text. It is designed to be easy to read and write for non-programmers, and is particularly suited for defining application's configuration settings. It can be parsed in one pass and is relatively easy to parse despite its use of implictely-defined objects and arrays (objects and array are not enclosed by brackets of any kind). 
- Compared to JSON, MDSon allows comments, do not require using {} and [] brackets to encode nested strcutures and arrays, does not require the use of commas to seprate fields, and has a special raw text field type that can hold a string of arbitary size and contents without the need for escaping any characters.
- Compared to XML, MDSon is a lot simpler (having only a handlful of concepts and rules) less verbose, not requiring opening and closing tags or escaping any characters.
- Compared to YAML, MDSon is not case-sensitive and does not require escaping any characters or use of brackets to define sections. 
- Compared to CSV, MDSon allows comments, permists arbitary nesting of data structures and does not limit the use of commas and quotes in field contents.

## Specifications
- An MDSon source (file or stream) consists of simple text organized into blocks using MD headings.
- Each block corresponds to a data structure (henceforth object) and starts with an MD heading which consists of one or more hashes followed by a unique name, e.g., # Document, ## Section.
- There must be only one root block (object) indicated by single #, e.g., # Document. Every other element (including other blocks) is a child of the root block.
- Any block including the root block, can enclose any number of elements including other blocks, lists of blocks, key-value pairs, lists of strings, numbers or boolean values and free-text elements.
- empty lines and lines starting with '//' are ignored. Useful for adding comments.
- white space is ingored, and can be used for identation.
- Block headings and key names are not case-sensitive.

## Blocks
- Under the root object, objects can be nested to any depth, e.g.,
```
# Family
//this is the root object
name: Simpsons
size: 8

## Address
    //nested within the root block 
    street: 742 
    ### Street
    No : 742
    Name : Evergreen Terrace 
    ### City
    // also nested within the address block
    name: Springfield
    state: ? Kansas 

## Socio-economic status
//nested within the root block 
class : working
```
## lists of blocks
- If a block heading has a suffix of ' List' or ' list', it will be parsed as an array of one or more blocks (objects) of the same type, eg,
```
## Children List
    ### Bart
    gender: male

    ### Lisa
    gender: female

    ### The baby
    gender: female
```
 ## Key-value pairs 
- Key-value pairs consist of an identifier (key) followed by ':' and optionally by a value. They are decoded into fields in the enclosing block object.
- Keys are strings and must start with a character followed optionally by any number of characters or digits other than ':'.
- Values can be of any scalar type: string, integer, float or bool.
- Bool value: 'true' is interpreted as truthy. Anything else including empty values is falsy.
- Quoting is optional but can be used to enforce a string type for the KV pair.

## Scalar lists
- Simple lists of of any scalar type: string, integer, float or bool, which will be parsed an array of the same type.
- A scalar list entry looks like a Key-Value entry, except that
    * the key must have has a suffix of ' List' or ' list'
    * the value is left empty
    * the following lines each start with '-'
- Example:
  ### Lisa
    gender: female
    hobbies:
    - reading
    - playing the saxophone
    - protesting
    voice: eardley Smith

## Raw text fields
- Fields that can hold a string of arbitary length and contents.
- A raw text entry looks like a Key-Value entry, except that the value cannot be empty and must be encolosed in <<>> or <<<>>>.
- Any text enclosed in << >> will be copied verbatim to the corresponding field in the object except that line breaks, if any, are removed. So the entire entry is stored as a single line of text.
- Any text enclosed in <<< >>> will be also copied verbatim to the corresponding field in the object but  line breaks are NOT removed, preserving multi-line text blocks as entered.
- In both cases, << or <<< should be the first character in the text block although it may follow on the same line as the field name (ie it does not need to be in its own line).
- In both cases, >> or >>> should be the last character in the text block and it also does not need to be on its own line
- exmaples 
```
freeText:<< anything here goes >>
```

```
freeText:
<<< `anything` here 'goes'
<head> other \n stuff </head>
some "other" test>>>
```

## Decoding
- Only fields thar are found in the destination type will be decoded
