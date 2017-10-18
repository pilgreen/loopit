package template

import (
  "bytes"
  "github.com/PuerkitoBio/goquery"
)

func Shim(in bytes.Buffer) (out bytes.Buffer, err error) {
  reader := bytes.NewReader(in.Bytes())
  doc, err := goquery.NewDocumentFromReader(reader)
  check(err)

  shims := doc.Find("[shim]")
  shims.Each(func(i int, ele *goquery.Selection) {
    query, _ := ele.Attr("shim")
    placement, _ := ele.Attr("placement")

    if len(query) > 0 {
      target := doc.Find(query)

      if placement == "after" {
        target.AfterSelection(ele)
      } else {
        target.BeforeSelection(ele)
      }
    } else {
      ele.Remove()
    }
  })

  html, err := doc.Html()
  check(err)

  out.WriteString(html)
  return
}
