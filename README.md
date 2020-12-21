# loopit

A GO program to loop structured data through the template engine. This is a command-line program that pipes a CSV or JSON file through the "text/template" package and prints the results to Stdout. If you need to make a full blown website or blog, and have stumbled onto this page, I highly recommend using [Hugo](https://gohugo.io/) instead. I often have to make small snippets of data-based code, portions of a page, or single pages in multiple environments. With these type of projects, the file structure required by Hugo sometimes prevents me from using it.

Example usage:

```
loopit [options] template.html [...template.html]
```

# Command-Line options

### -data file|url

The data flag will accept a file path on the system or a url to a CSV, JSON or RSS file. Loopit will make a map of the file data and pass it to the template.

In the case of a CSV file, loopit will create an object for each row using the first row as the keys for each object. Therefore, the file needs to have a header row with column names that can be used in the dot notation of the "text/template" pacakge. As a general rule stick to single words without punctuation and you should be fine.

In the case of a JSON file, loopit will pass the object or array straight through.

In the case of an RSS file, loopit will convert most standard elements.

*If you do not specify a template file, loopit will print the parsed data to sdout. It will also accept data piped in through stdin if -data is ommitted.*

### -markdown

The markdown flag passes the output through the [russross/blackfriday](https://github.com/russross/blackfriday) package.

### -minify

The minify flag will pass the output through the [tdewolff/minify](https://github.com/tdewolff/minify) package. Currently this only works for HTML output.

### -out filename

Write the markup to a file instead of stdout. This uses `ioutil.WriteFile()` and replaces the file each time.

### -watch

After the initial render, loopit will enter watch mode and re-run the render function any time a template is written. It also recursively walks the current directory and adds any `.js` and `.css` files it finds to the list.

### -shim

The shim flag is a boolean option. When added, loopit will use [goquery](https://github.com/PuerkitoBio/goquery) to move DOM around after the template has been parsed. This is useful when HTML is provided in a JSON object, but needs to be enhanced with ads or any other content. 

Here is an example that will inject an ad before the 3rd paragraph:

```
<div shim="body p:nth-child(3)">[Ad code goes here]</div>
```

### -version

Prints the version.



# Template functions

The Go template package is very basic, but also extendable. Below are custom functions that are available in the template based solely on need so far. Several of these are borrowed from [Hugo](https://gohugo.io/). 

## Strings

### add

Adds an int to an int.

```
{{ add 1 $index }}
```

### dateFormat

This function parses a date using the `time` package and returns a formatted version of your choosing. The Golang [time](https://golang.org/pkg/time/) package is kind of odd, but also very flexible. Make sure to read the documentation fully before working with this. 

The first parameter is the date string to parse, the second is a representation of that date mapped to 'Mon Jan 2 15:04:05 MST 2006' and the third is the desired output format mapped to 'Mon Jan 2 15:04:05 MST 2006'. 

The fourth parameter is optional, and allows you to alter the Timezone using the time.LoadLocation() function.

```
{{ dateFormat "2017/11/2" "2006/01/02" "Jan 2" "America/Chicago" }}
```

### findRe

This function returns the first string match using the regexp package.

```
{{ findRe "^foo" "foobar" }}
```

### findSubRe

This is similar to findRe, but it returns a slice of strings that match sub-patterns.

```
{{ findSubRe "^foo(.*)$" "foobar" }}

This would return: ["foobar", "bar"]
```

### floatToInt

Straightforward but sometimes necessary when comparing data types from JSON feeds that are casted with json.UnMarshal()

```
{{ floatToInt .floatPropertyInFeed }}
```

### inchesToFeet

Simple function that will convert 68 inches to the string 5'8".

```
{{ inchesToFeet .player.height }}
```

### lower

Returns a string converted to lower case.

```
{{ lower "John Smith" }}
```

### minify

The minify functions sends a string through the [tdewolff/minify](https://github.com/tdewolff/minify) package. The following mimetypes are supported:

+ "text/css"
+ "text/html"
+ "text/javascript"

```
<style>{{ file "./css/styles.css" | minify "text/css" }}</style>
```

### markdown

This function sends the string through the `blackfriday.MarkdownCommon()` function and returns parsed HTML as a string.

```
{{ file "./contents.md" | markdown }}
```

### matchRe

Just like findRe but returns a boolean instead of the matching string

```
{{ findRe "^foo" "foobar" }}
```

### replace

Passes a string through the `strings.Replace` function swapping the first match.

```
{{ replace "http" "https" .url }}
```

### replaceRe

More advanced replacing using the regexp package. This replaces all matching strings in the pattern.

```
{{ replaceRe "^<div xmlns.*?>" "" .body }}
```

### subtract

Subtracts an int from an int.

```
{{ subtract 1 $index }}
```

### trim

Trims characters and space from both sides of the string.

```
{{ trim "/" .link }}
```

### unescape

Runs a string through the url.QueryUnescape() function, which converts each 3-byte encoded substring of the form "%AB" into the hex-decoded byte 0xAB. Useful for urls in feeds. 

**Note: escape is already available in the Go Template package.**

```
{{ unescape .link }}
```

### upper

Converts a string to uppercase.

```
{{ upper "huge" }}
```

## Collections

### find 

This function is based on the `Array.find` function in javascript. It returns the first value of a sequence that matches (or doesn't match) the comparator. 

```
{{ $josh := find .employees "name" "Josh" }}
```

or 

```
{{ $awayTeam := find .teams "venue.location" "!=" "home" }}
```

### join

Joins a slice using a separator and returns the string.

```
{{ join "," .authors }}
```

### pluck

A port of Underscore's pluck function. Uses dot notation to grab nested objects and returns a slice.

```
{{ pluck "album.title" .relatedWork }}
```

### slice

Pulls a subset of data by index. For example, to range over the first five rows of a csv file use the following code.

```
{{ range slice . 0 5 }} ... {{ end }}
```

### sort

Sorts a sequence by a path. Ascending and desending order are supported.

```
{{ range sort .people "name.last" }} ... {{ end }}
```

or 

```
{{ range sort .people "age" "desc" }} ... {{ end }}
```

### where

A port of Hugo's where function, this returns a filtered slice and is desinged to work with a range. Dot notation is possible. 

The default comparator is `==`, but you can use any of the following: `=, ==, eq, !=, <>, ne, <, lt, >, gt, <=, lte, >=, gte, matches`.

`matches` runs a pattern through the regexp.MatchString() function.

```
{{ range (where .friends ".name.first" "ne" "John") }}
```

## Scratch

This is a copy of a concept from Hugo. The Go template package is intentionally dumb by design, but that doesn't always work when a JSON feed is created by a third party and time to publication is critical. Variable scope is very tight, so to loosen it up a bit you can use scratch. This lets you maintain global variables and alter them from within loops and conditionals.

It works very much like local storage in the browser, with key/string pairings. There are three functions to get, set and delete a pairing. Unlike Hugo, only strings are supported.

```
{{ $scratch := scratch }}
{{ $scratch.Set("name", "Jane") }}

{{ if [condition] }}
  {{ $scratch.Set("name", "Gloria") }}
{{ end }}

<h2>Hello, {{ $scratch.Get("name") }}</h2>
```

## Partials

Loopit is designed to make small snippets quickly, but often those snippets can turn into full blown pages that need some DRY tools. Partials allows you to run data through an external template and get a string, which can be further manipulated by the functions in this template.

### file

Simply returns file contents as a string. Relative paths are compared to the working directory when the loopit command is run, not relative to the template. This is a tad confusing but allows for maximum flexibility, like common atomic structures for multiple projects. For consistency, I personally use a Makefile.

```
{{ file "./css/atoms.css" | minify "text/css" }}
```

### partial

Returns contents of a file as a string, but runs data through loopit and includes all functions (except nested partials I haven't figured that one out yet). Just like `file`, relative paths are compared to the working directory.

```
{{ partial "./amp-components/amp-img.tmpl" .photo }}
```
