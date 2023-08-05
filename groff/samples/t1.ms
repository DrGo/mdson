.so macros.ms \"can be used to include all kinds of files.

.\"Debugging
.\"The tmc request is similar to tm1 but does not append a newline (as is done in tm and tm1).
.\" Request: .pev
.\" Print the contents of the current environment and all the currently defined environments (both named and numbered) on stderr.

.\" Request: .pm
.\" Print the entire symbol table on stderr. Names of all defined macros, strings, and diversions are print together with their size in bytes. Since gtroff sometimes adds nodes by itself, the returned size can be larger than expected.


.\" Request: .pnr
.\" Print the names and contents of all currently defined number registers on stderr.

.\" Request: .ptr


.\" Preamble
.\"device papersize=8.5i,11i
.nr PO 3c

.nr PS 10.5p \" Use 10.5-point type.
.ds FAM Minion    \" Use Minion font family.
.\"nr VS 20p  \" vertical size
.P1 \"Prints the header on page 1. Must come before .TL
.TL 
Chapter 1: 
.br
Environmental Health

.NH 1
Introduction
.LP
Must switch to paragraph mode after heading using .LP
Enter sentences in separate lines.
or add two spaces at the end\.  
to indicate a new sentence\. We won\'t do this.

Empty lines separate paragraphs.
So this is another paragraph.

.BIT
Microsoft Word for Windows,
starting with version 1.0 through the current version
(Word 2000).
.BIT
GNU Emacs,
from its first appearance as a standalone editor through
the current version (v20).
See [Bloggs 2002] for details.

.LP
.nr step 0 1  \"create custom var
.nr PI 3n
.NIT
People
.NIT
Animals
.NIT
Things

.\"indented region
.IP \[bu] 2
Lawyers:
.RS \" begin a region indented by #  in the PI register.
.IP \[bu]
Dewey,
.IP \[bu]
Cheatham,
and
.IP \[bu]
and Howe.
.RE \"end indented region
.IP \[bu]
Guns

.\" KS begins a keep, KF a floating keep (eg table that can appear later, and KE ends a keep of either kind.


.\"a keep with a box drawn around it B1 and B2
.LP
.B1
.I Warning:
Happy Fun Ball may suddenly accelerate to dangs`erous
speeds.
.sp \n[PD]/2 \" space by half the inter-paragraph distance
.B2


.QP
This paragraph will be indented by the amount in QI

.PP
Displays are sections of text set off from the surrounding material (typically paragraphs), often differing in indentation, and/or spacing.

.NH 1
Defaults
.LP
The default is 10-point type on 12-point spacing.


.NH 1
Animalia
.NH 2
Arthropoda
.NH 3
Crustacea
.NH 2
Chordata
.NH S 6 6 6
Daimonia
.NH 1
Plantae

.\"If the GROWPS register is set to a value greater than the level of the heading, then the point size of the heading will be increased by PSINCR units over the text size specified by the PS register, for each level by which the heading level is less than the value of GROWPS.
.nr GROWPS 3
.nr PSINCR 1.5p
.
.NH 1
Top Level Heading
.
.NH 2
Second Level Heading
.
.NH 3
Third Level Heading

.SH 2
An Unnumbered Subheading of a ps similar to heading 2

.NH
Margins
.PP
Reg.	Definition	Effective	Default
PO	Page offset (left margin)	next page	1i
LL	Line length	next paragraph	6i
LT	Header/footer length	next paragraph	6i
HM	Top (header) margin	next page	1i
FM	Bottom (footer) margin	next page	1i

.NH 1
Headers and footers
.OH 'left'center'right'
.EH 'left'center'right'

.\"Use groff options to run preprocessors on the input: -e for geqn, -p for gpic, -R for grefer, and -t for gtbl.


.TS H
allbox expand;
cb | cb .
Text      ...of heading...
_
.TH
.T&
l | l .
... the rest of the table follows...
.CW
.TE

