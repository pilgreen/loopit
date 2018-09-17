package strings

import (
  "text/template"
)

var FuncMap = template.FuncMap {
  "add": Add,
  "dateFormat": DateFormat,
  "findRe": FindRe,
  "findSubRe": FindSubRe,
  "floatToInt": FloatToInt,
  "inchesToFeet": InchesToFeet,
  "lower": ToLower,
  "minify": MinifyCode,
  "markdown": Markdown,
  "matchRe": MatchRe,
  "replace": Replace,
  "replaceRe": ReplaceRe,
  "subtract": Subtract,
  "trim": Trim,
  "unescape": Unescape,
  "upper": ToUpper,
}
