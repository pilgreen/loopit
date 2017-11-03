# loopit
A GO program to loop structured data through the template engine

This is a command-line program that pipes a CSV or JSON file through the "text/template" package and prints the results to Stdout. 

Example usage:

```
loopit [options] template.html [...template.html]
```

## Command-Line Options

#### -data file|url

The data flag will accept a file path on the system or a url to a CSV or JSON file. Loopit will make a map of the file data and pass it to the template.

In the case of a CSV file, loopit will create an object for each row using the first row as the keys for each object. Therefore, the file needs to have a header row with column names that can be used in the dot notation of the "text/template" pacakge. As a general rule stick to single words without punctuation and you should be fine.

In the case of a JSON file, loopit will pass the object or array straight through.

*Loopit will also accept CSV and JSON data passed in through stdin if -data is ommitted.*

#### -shim

The shim flag is a boolean option. When added, loopit will use [goquery](https://github.com/PuerkitoBio/goquery) to move DOM around after the template has been parsed. This is useful when HTML is provided in a JSON object, but needs to be enhanced with ads or any other content. 

Here is an example that will inject an ad before the 3rd paragraph:

```
<div shim="body p:nth-child(3)">[Ad code goes here]</div>
```

#### -minify

The minify flag will pass the output through the [tdewolff/minify](https://github.com/tdewolff/minify) package. Currently this only works for HTML output.

#### -v

Prints the version.




## Template functions

The Go template package is very basic, but also extendable. Below are custom functions that are available to the template based solely on need so far. Several of these are borrowed from [Hugo](https://gohugo.io/). 

#### dateFormat

This function parses a date using the `time` package and returns a formatted version of your choosing. The Golang [time](https://golang.org/pkg/time/) package is kind of odd, but also very flexible. Make sure to read the documentation fully before working with this. 

The first parameter is the date string to parse, the second is a representation of that date mapped to 'Mon Jan 2 15:04:05 MST 2006' and the third is the desired output format mapped to 'Mon Jan 2 15:04:05 MST 2006'.

```
{{ dateFormat "2017/11/2" "2006/01/02" "Jan 2" }}
```

#### file

The file function will return either a local file or a remote url as a string. The example below will pull the local file into the document and pass the resulting string to the `minify` function:

```
{{ file "path/to/css/file" | minify "text/css" }}
```

#### find 

This function is based on the `Array.find` function in javascript. It returns the first value of a sequence that matches (or doesn't match) the comparator. 

```
{{ $josh := find .employees "name" "Josh" }}
```

or 

```
{{ $awayTeam := find .teams "venue.location" "!=" "home" }}
```

#### inchesToFeet

Simple function that will convert 68 inches to the string 5'8".

```
{{ inchesToFeet .player.height }}
```

#### int

Converts a float to an int type for comparison.

```
{{ range lt (int .points) 65 }} ... {{ end }}
```

#### minify

The minify functions sends a string through the [tdewolff/minify](https://github.com/tdewolff/minify) package. The following mimetypes are supported:

+ "text/css"
+ "text/html"
+ "text/javascript"

```
<style>{{ file "./css/styles.css" | minify "text/css" }}</style>
```

#### markdown

This function sends the string through the `blackfriday.MarkdownCommon()` function and returns parsed HTML as a string.

```
{{ file "./contents.md" | markdown }}
```

#### slice

The slice function will let you pull a subset of data. For example, to range over the first five rows of a csv file use the following code.

```
{{ range slice . 0 5 }} ... {{ end }}
```

#### sort

Sorts a sequence by a path. Ascending and desending order are supported.

```
{{ range sort .people "name.last" }} ... {{ end }}
```

or 

```
{{ range sort .people "age" "desc" }} ... {{ end }}



