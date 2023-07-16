# mdson
=======
// copyrights

/*
- Section includes one or more chapters; a chapter includes on or more article or div
File format:
file begings with a header (metadata) that includes info on the section/chapter.
standard fields include:
- Type= one of section, page
- Title
- Subtitle
- Description
- draft = false/True; draft files are not included in the build unless --drafts is passed
// Used by the parent section to order its subsections.
// Lower values have higher priority.
- weight = 0
// This sets the number of pages to be displayed per paginated page.
// No pagination will happen if this isn't set or if the value is 0.
- paginate_by = 0
// The date of the post.
// Two formats are allowed: YYYY-MM-DD (2012-10-02) and RFC3339 (2002-10-02T15:00:00Z).
// Do not wrap dates in quotes; the line below only indicates that there is no default date.
// If the section variable `sort_by` is set to `date`, then any page that lacks a `date`
// will not be rendered.
// Setting this overrides a date set in the filename.
- date =

// The last updated date of the post, if different from the date.
// Same format as `date`.
- updated =
// A list of page authors. If a site feed is enabled, the first author (if any)
// will be used as the page's author in the default feed template.
authors = []

// The taxonomies for this page. The keys need to be the same as the taxonomy
// names configured in `config.toml` and the values are an array of String objects. For example,
// tags = ["rust", "web"].
- taxa
- descriptors

eg
Title= environemental health
Authors= 
- Salah Mahmud
- Second author
date= 17July2023
tags=
- health
- ecology

## descrpitors
Agents = Salmonella
Incubation period = 3 days
Infectious period = 1 day before to 3 days after

## Summary

## Introduction

## Diagnosis

## Prevention

## Footnotes


