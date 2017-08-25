# loopit
A GO program to loop structured data through the template engine

This is a command-line program that pushes a CSV or JSON file through the "text/template" package. All options are passed in through flags.

#### -csv file|url

The csv flag will accept either a file path on the system or a publicly available url. The program will make a map of the file data creating an object for each row. The first row will be used as the keys for each object. Therefore, the file needs to have a header row with column names that can be used in the dot notation of the "text/template" pacakge. As a general rule stick to single words without punctuation and you should be fine.

#### -json file|url

The json flag simply passes the file contents through to the template. You can use JSON arrays or objects, along with nested arrays and objects.

#### -template file

If you omit the template flag, the program will echo the JSON equivalent of the resulting map. This is a good way to see what objects are available for the template. Including the template parameter will pass the map through and echo the evaluated code. For more information on what is available with the "text/template" package, refer to the [online documentation](https://golang.org/pkg/text/template/).

## Template functions

The Go template package is very basic, but also extendable. Below are custom functions that are available inside the template.

#### slice([]interface{}) start(int) end(int)

The slice function will let you pull a subset of data. For example, to range over the first five rows of a csv file use the following code.

```
{{ range slice . 0 5 }}
  ... code goes here ...
{{ end }}
```

#### file path(string)

The file function will return either a local file or a remote url as a string.

#### minify mimetype(string) contents(string)

The minify functions sends a string through the [tdewolff/minify](https://github.com/tdewolff/minify) package. 

The following mimetypes are supported, though I haven't fully tested each one yet:

+ "text/css"
+ "text/html"
+ "text/javascript"

#### markdown text(string)

This function just sends the string through the `blackfriday.MarkdownCommon()` function and returns a string.

#### ampify html(string)

This function converts `<iframe>` tags into `<amp-iframe>` tags. It could be expanded in the future if necessary to incorporate more amp elemens or convert additional native tags.
