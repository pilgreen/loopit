package tpl

import (
  "bytes"
  "github.com/PuerkitoBio/goquery"
)

func Shim(in bytes.Buffer) (out bytes.Buffer, err error) {
  reader := bytes.NewReader(in.Bytes())
  doc, err := goquery.NewDocumentFromReader(reader)
  if err != nil { return }

  shims := doc.Find("[shim]")
  shims.Each(func(i int, ele *goquery.Selection) {
    placement, _ := ele.Attr("placement")
    query, _ := ele.Attr("shim")
    target := doc.Find(query)

    if target.Length() > 0 {
      // Clean up and place
      ele.RemoveAttr("shim")

      if placement == "after" {
        target.AfterSelection(ele)
      } else {
        target.BeforeSelection(ele)
      }
    } else {
      // New home doesn't exist
      ele.Remove()
    }
  })

  html, err := doc.Html()
  check(err)

  out.WriteString(html)
  return
}
