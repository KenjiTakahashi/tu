[![Build Status](https://travis-ci.org/KenjiTakahashi/tu.png?branch=master)](https://travis-ci.org/KenjiTakahashi/tu)

**tu** is a batch music tagging helper.

This is basically a [tagutil](https://github.com/kAworu/tagutil) wrapper with some additional features and convenience shortcuts. And concurrency, because... Go!

## installation

First, you have to [get Go](http://golang.org/doc/install). Note that version >= 1.1 is required.

Second, you have to [get tagutil](https://github.com/kAworu/tagutil). Note that JSON support is required.

Then, just

```bash
$ go get github.com/KenjiTakahashi/tu
```

should get you going.

## usage

#### w

```bash
$ tu r PATTERN FILES...
```

Writes tags to files based on their filenames. Pattern conforms to [tagutil](https://github.com/kAworu/tagutil#renaming-files)'s definition.

#### e

```bash
$ tu e FILES...
```

Opens interactive editing session for each file. It works by opening a YAML formatted tag list in your `$EDITOR`.

#### t

```bash
$ tu t [-t TAGS] FILES...
```

Applies TitleCase transformation to specified files. If `-t` flag is present, it should contain a comma separated list of tag names to transform. Otherwise, all found tags are transformed.

#### r

```bash
$ tu r [-Y] PATTERN FILES...
```

Renames files based on their tags. This actually only calls `tagutil -p rename:PATTERN FILES...`.

If `-Y` flag is present, all questions are answered YES. **Note:** If applying a pattern on two different files results in the same filename, this option may eat your files. So be careful.

## titlecase

There is also a package here named `titlecase`, which is more or less a rewrite of [Stuart Coville](http://muffinresearch.co.uk)'s Python library (available [here](https://github.com/ppannuto/python-titlecase)).

It is used to capitalize song names based on NY Times Manual of Style. Meaning, it generally capitalizes first letter of every word, but tries to get proper on "small words" and other corner cases which should not be capitalized.

It was moved to a separate package, so that others can make use of it. Documentation is available through [Godoc](http://godoc.org/github.com/KenjiTakahashi/tu/titlecase).
