package funcs

import (
  "bytes"
  "github.com/PuerkitoBio/goquery"
)

func Shim(in bytes.Buffer) (out bytes.Buffer, err error) {
  reader := bytes.NewReader(in.Bytes())
  doc, err := goquery.NewDocumentFromReader(reader)
  Check(err)

  shims := doc.Find("[shim]")
  shims.Each(func(i int, ele *goquery.Selection) {
    query, _ := ele.Attr("shim")

    if len(query) > 0 {
      target := doc.Find(query)
      target.BeforeSelection(ele)
    } else {
      ele.Remove()
    }
  })

  html, err := doc.Html()
  Check(err)

  out.WriteString(html)
  return
}
