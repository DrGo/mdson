# Vask
A task runner written in Go. Currently can create html (and other format) files from mdson data files and linked templates. It can also concatenate files, inline CSS and create distribution folders (suitable for uploading into an Internet hosting service).

## To compile
- Download the Go toolchain (Go 1.11+ because it uses modules).
- `go build`

## To use

## Caveat
- Work in progress.

## Design
### vask-d
- specialized version to process Alima design docs.
- Expects the following dir structure:
src -|
     |- contents -|
                  |-datasets --> contains dataset1.mdson etc
                  |-routines
                  |-otherstuff
     |- layout   -|
                  |-datasets --> contains datasets.gohtml + any css/js
                  |-routines
                  |-otherstuff
                  | index.html + main.css + main.js + default.gohtml
build --> output files will be copied here

- dirNames = []dirnames under contents/
- foreach dirname: 
  - get gohtmlName from corresponding layout/dirname || default gohtml
  - get a list of all mdson filenames under dirname
  - for each x.mdson:
      - validate contents
      - expand refs
      - execute template into src/dirname/x.html
  - create a toc and save into src/dirname/index.html
  - create src/index.html







