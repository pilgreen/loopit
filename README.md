# loopit
A GO program to loop structured data through the template engine

This is a command-line program that pushes a CSV or JSON file through the "text/template" package and echo the results to Stdout. 

Example usage:

```
loopit [options] [path/to/]template.html
```

## Options


#### -data file|url

The data flag will accept a file path on the system or a url to a CSV or JSON file. The program will first determine the appropriate file type by extension and then make a map of the file data to pass to the template.

In the case of a CSV file, loopit will create an object for each row using the first row as the keys for each object. Therefore, the file needs to have a header row with column names that can be used in the dot notation of the "text/template" pacakge. As a general rule stick to single words without punctuation and you should be fine.

In the case of a JSON file, loopit will pass the object or array straight through.


#### -shim

The shim flag is a boolean option. When added, loopit will use [goquery](https://github.com/PuerkitoBio/goquery) to move DOM around after the template has been parsed. This is useful when HTML is provided in a JSON object, but needs to be enhanced with ads or any other content. 

Here is an example that will inject an ad before the 3rd paragraph:

```
<div shim="body p:nth-child(3)">[Ad code goes here]</div>
```



## Template functions

The Go template package is very basic, but also extendable. Below are custom functions that are available to the template.


#### slice *int* *int*

The slice function will let you pull a subset of data. For example, to range over the first five rows of a csv file use the following code.

```
{{ range slice . 0 5 }}
  ... code goes here ...
{{ end }}
```


#### file *string*

The file function will return either a local file or a remote url as a string.

The example below will pull the local file into the document and pass the resulting string to the `minify` function:

```
{{ file "path/to/css/file" | minify "text/css" }}
```


#### minify mimetype *string*

The minify functions sends a string through the [tdewolff/minify](https://github.com/tdewolff/minify) package. 

The following mimetypes are supported, though I haven't fully tested each one yet:

+ "text/css"
+ "text/html"
+ "text/javascript"


#### markdown *string*

This function sends the string through the `blackfriday.MarkdownCommon()` function and returns parsed HTML as a string.


#### ampify *string*

This function converts `<iframe>` tags into `<amp-iframe>` tags. It could be expanded in the future if necessary to incorporate more AMP elements or convert additional native tags.
