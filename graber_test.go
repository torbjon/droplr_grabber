package main

import (
    "testing"
    "github.com/karlseguin/gspec"
)

func TestVars(t *testing.T) {
    spec := gspec.New(t)
    spec.Expect(site_url).ToEqual("http://d.pr/i/")
    spec.Expect(chars[1]).ToEqual("B")
    spec.Expect(chars[30]).ToEqual("e")
}


func TestResponseCode(t *testing.T) {
    spec := gspec.New(t)
    spec.Expect(response_status_to_code("200 OK")).ToEqual(uint64(200))
    spec.Expect(response_status_to_code("503 Service Temporarily Unavailable")).ToEqual(uint64(503))
}
