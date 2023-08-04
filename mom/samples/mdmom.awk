BEGIN {
    inlisting = 0;
    inlisting2 = 0;
    inlisting3 = 0;
    inparagraph = 0;
    print ".PAPER A4"
    print ".TITLE_STYLE  COLOR YELLOW"
    print ".TITLE \"The Screaming Butterflies\"  \"Flower Power Workgroup\""
    print ".PDF_TITLE \"\\*[$TITLE]"
    print ".\\\" Formatting style, margins"
    print ".PRINTSTYLE TYPESET"
    print ".L_MARGIN   2.5c"
    print ".R_MARGIN   2.5c"
    print ".B_MARGIN   2.5c"
    print ".NEWCOLOR YELLOW #e36c0a"
    print ".HEADING_STYLE 1 COLOR YELLOW"
    print ".HEADING_STYLE 2 COLOR YELLOW"
    print ".\\\" General defaults"
    print ".FAMILY   P"
    print ".FT       R"
    print ".PT_SIZE  11"
    print ".AUTOLEAD 3 3"
    print ".PARA_INDENT 0 \\\" No indent because we're spacing paragraphs."
    print ".HYPHENATION 0 \\\" No hypenantion of words at end of line."
    print ".START"
}
{
    if ( $0 ~ /^# / ) { reset_list(); printf ".HEADING 1 \"%s\"\n", substr($0, 3); }
    else if ( $0 ~ /^## / ) { reset_list(); printf ".HEADING 2 \"%s\"\n", substr($0, 4); }
    else if ( $0 ~ /^### / ) { printf ".HEADING 3 \"%s\"\n", substr($0, 5); }
    else if ( $0 ~ /^#### / ) { printf ".HEADING 4 \"%s\"\n", substr($0, 6); }
    else if ( $0 ~ /^- / ) { 
        if ( inlisting2 == 1 ) {
            inlisting2 = 0;
            print ".LIST OFF"
        }
        if ( inlisting == 0 ) {
            inlisting = 1;
            print ".SP .25v"
            print ".LIST"
        }
        print ".ITEM"; print substr($0, 3) 
    }
    else if ( $0 ~ /^\t- / ) { 
        if ( inlisting == 0 ) {
            inlisting = 1;
            print ".SP .25v"
            print ".LIST"
        }
        if ( inlisting2 == 0 ) {
            inlisting2 = 1;
            print ".LIST"
        }
        print ".ITEM"; print substr($0, 4) 
    }
    else if ( $0 ~ /^\t\t- / ) { 
        if ( inlisting == 0 ) {
            inlisting = 1;
            print ".SP .25v"
            print ".LIST"
        }
        if ( inlisting2 == 0 ) {
            inlisting2 = 1;
            print ".LIST"
        }
        if ( inlisting3 == 0 ) {
            inlisting3 = 1;
            print ".LIST"
        }
        print ".ITEM"; print substr($0, 5) 
    }
    else if ( $0 ~ /^  / ) { 
            if ( inlisting == 1 || inlisting2 == 1  || inlisting3 == 1) {
                gsub("^ *", "", $0)
            }
        }
    else {
        if ( inlisting3 == 1 ) {
            inlisting3 = 0;
            print ".LIST OFF"
        }
        if ( inlisting2 == 1 ) {
            inlisting2 = 0;
            print ".LIST OFF"
        }
        if ( inlisting == 1 ) {
            inlisting = 0;
            print ".LIST OFF"
            print ".SP .25v"
        }
        if ( inparagraph == 0 ) {
            inparagraph = 1;
            print ".PP"
        }
        print $0
    }
}
function reset_list()
{
    if ( inlisting3 == 1 ) {
        inlisting3 = 0;
        print ".LIST OFF"
    }
    if ( inlisting2 == 1 ) {
        inlisting2 = 0;
        print ".LIST OFF"
    }
    if ( inlisting == 1 ) {
        inlisting = 0;
        print ".LIST OFF"
        print ".SP .25v"
    }
}
