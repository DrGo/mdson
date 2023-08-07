// function that gets a chapter as its `body` and formats
// it as a chapter in a scientific book
#let chapter(
    // The Chapter's title.
title: "Chapter title",
// The chapters's authors.
author: "",
// The paper size to use.
paper: "us-letter",
// paletino OS replacement "TeX Gyre Pagella",
body_font: "Minion Pro",
title_font: "Lato",
hyphenate: false,

bibliography-file: none,
// The chapter's content.
body,
) = {

  // Set the document's metadata.
  set document(title: title, author: author)

    // Configure the page properties.
    set page(
        paper: paper,
        margin: (outside: 1.0in, inside: 1.125in, bottom:1.125in+0.4in, top: 1.125in + 0.4in),
        header-ascent: 0.4in,
        footer-descent: 0.3in,

        numbering: "1",
        //TODO: allow optional header + page numbers
        // The header always contains the chapter title on odd pages and
        // the section title on even pages, unless the page is one
        // the starts a chapter (the chapter title is obvious then).
        header: locate(loc => {
            // Are we on an odd page? at() returns an array but we know it has one value
            let i = counter(page).at(loc).first()
            // let pagenum=[#loc.page()]
            if calc.odd(i) {
              return text(0.95em, smallcaps(title))
            }
          // Are we on a page that starts a chapter? (We also check
          // the previous page because some headings contain pagebreaks.)
          let all = query(heading, loc)
            if all.any(it => it.location().page() in (i - 1, i)) {
            return
          }
          // Find the heading of the section we are currently in.
          let before = query(selector(heading).before(loc), loc)
          if before != () {
            align(right, text(0.95em, smallcaps(before.last().body)))
          }
          }),
        )
          //header: locate(loc => {
          //let elems = query(
          //  selector(heading).before(loc),loc,
          //)
          //let even=calc.even(counter(page).at(loc).at(0))
          //let chap_title = smallcaps[
          //  Typst chap_title
          //]
          //let page=[#loc.page()]
          //if even {
          //  page + h(1fr) + chap_title 
          //} else {
          ////if headers are available ie after 1st section
          //  if elems.len()>0 {
          //    let body = smallcaps(counter(heading).display() + " " + elems.last().body) 
          //     emph(body) + h(1fr) + page 
          //  }  
          //}
          //}),   


          // set text(font: "New Computer Modern", lang: "en")
          set text(
              size: 11pt,
              // Set the body font. TeX Gyre Pagella is a free alternative
              // to Palatino.
              font: body_font,
              // font: "New Computer Modern",
              lang: "en",
              hyphenate: hyphenate,
          )
          // MATH
          show math.equation: set text(weight: 400)
          // set math.equation(numbering: "(1.1)") // Currently not directly supported by typst
          set math.equation(numbering: "(1)")
          show math.equation: set block(spacing: 0.65em)

          // Configure lists.
          set enum(indent: 10pt, body-indent: 9pt)
          set list(indent: 10pt, body-indent: 9pt)
          
          // Paragraph settings
          set par(
              justify: true,
              // leading: 0.78em, 
          )
          show par: set block(
            spacing: 0.78em,
            // Set paragraph spacing.
            above: 1.2em,
            below: 1.2em,
            )

          // The first page.
          // page(align(center + horizon)[
          //   #text(2em)[*#title*]
          //   #v(2em, weak: true)
          //   #text(1.6em, author)
          // ])


          // Configure chapter headings.
          show heading.where(level: 1): it => {
            //FIXME: add option to control chapter title page
            // Always start on even pages.
            // pagebreak()
            // counter(page).display(i => if calc.odd(i) {
            //       pagebreak()
            //       })

            // Create the heading numbering.
            let number = if it.numbering != none {
              counter(heading).display(it.numbering)
                h(7pt, weak: true)
            }

            v(5%)
              text(2em, weight: 700, block([#number #it.body]))
              v(1.25em)
          }
          // show heading.where(level: 2): set text(size: 14pt)
  // // Configure headings.
  // show heading: set text(font: "Syne")
  // show heading.where(level: 1): set text(1.1em)
  // show heading.where(level: 1): set par(leading: 0.4em)
  // show heading.where(level: 1): set block(below: 0.8em)
  // show heading: it => {
  //   set text(weight: 600) if it.level > 2
  //   it
  // }
  // show heading: set text(11pt, weight: 400)
    show heading: set text(
      font: title_font,
    )
    set heading(
        numbering: "1.1",    
    )
      // Links should be underlined.
  show link: underline
  // footnotes
  set footnote.entry(
    gap: 0.8em
  )
  // numbering of footnotes
 show footnote.entry: it => {
  let loc = it.note.location()
  numbering(
    "1: ",
    ..counter(footnote).at(loc),
  )
  it.note.body
}

// Configure figures.
  show figure: it => block({
    // Display a backdrop rectangle.
    move(dx: -3%, dy: 1.5%, rect(
      fill: rgb("FF7D79"),
      inset: 0pt,
      move(dx: 3%, dy: -1.5%, it.body)
    ))

    // Display caption.
    if it.has("caption") {
      set align(center)
      set text(font: "Syne")
      v(if it.has("gap") { it.gap } else { 24pt }, weak: true)
      [-- ]
      it.caption
      if it.numbering != none {
        [ (] + counter(figure).display(it.numbering) + [)]
      }
      [ --]
    }

    v(48pt, weak: true)
  })
    body
      // Display bibliography.
  if bibliography-file != none {
    show bibliography: set text(8pt)
    bibliography(bibliography-file, title: text(10pt)[References], style: "ieee")
  }
}

//Title
  // Display the paper's title.
  // v(3pt, weak: true)
  // align(center, text(18pt, title))
  // v(8.35mm, weak: true)
  //
// A stylized block with a quote and its author.
#let blockquote(by, body) = box(inset: (x: 0.4em, y: 12pt), width: 100%, {
  set text(font: "Syne")
  grid(
    columns: (1em, auto, 1em),
    column-gutter: 12pt,
    rows: (1em, auto),
    row-gutter: 8pt,
    text(5em)["],
    line(start: (0pt, 0.45em), length: 100%),
    none, none,
    text(1.4em, align(center, body)),
    none, none,
    v(8pt) + align(right, text(font: "Barlow")[---#by]),
  )
}
)

// A theorum block as ann example of a custom counter
// #theorem[$1 = 1$]
#let c = counter("theorem")
#let theorem(it) = block[
  #c.step()
  *Theorem #c.display():* #it
]

// custom block
// #task(critical: true)[Food today?]
//#task(critical: false)[Work deadline]

#let task(body, critical: false) = {
  set text(red) if critical
  [- #body]
}


