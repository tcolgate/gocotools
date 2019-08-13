package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
)

//Badge represents a coverage badge
type Badge struct {
	Coverage float64
	Color    string
}

var colors = map[string]string{
	"brightgreen": "#44cc11",
	"green":       "#97ca00",
	"yellow":      "#dfb317",
	"orange":      "#fe7d37",
	"red":         "#e05d44",
}

var _badgeTemplate string = `
<svg xmlns="http://www.w3.org/2000/svg" width="96" height="20">
	<title>{{.Coverage}}</title>
	<desc>Generated with covbadger (https://github.com/imsky/covbadger)</desc>
	<linearGradient id="smooth" x2="0" y2="100%">
		<stop offset="0" stop-color="#bbb" stop-opacity=".1" />
		<stop offset="1" stop-opacity=".1" />
	</linearGradient>
	<rect rx="3" width="96" height="20" fill="#555" />
	<rect rx="3" x="60" width="36" height="20" fill="{{.Color}}" />
	<rect x="60" width="4" height="20" fill="{{.Color}}" />
	<rect rx="3" width="96" height="20" fill="url(#smooth)" />
	<g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,sans-serif" font-size="11">
		<text x="30" y="15" fill="#010101" fill-opacity=".3">coverage</text>
		<text x="30" y="14">coverage</text>
		<text x="78" y="15" fill="#010101" fill-opacity=".3">{{.Coverage}}%</text>
		<text x="78" y="14">{{.Coverage}}%</text>
	</g>
</svg>`

func RenderBadge(coverage float64) (string, error) {
	if coverage < 0 || coverage > 100 {
		return "", fmt.Errorf("Invalid coverage: %f", coverage)
	}

	var buffer bytes.Buffer
	badgeTemplate, _ := template.New("badge").Parse(_badgeTemplate)

	color := colors["red"]

	if coverage > 95 {
		color = colors["brightgreen"]
	} else if coverage > 80 {
		color = colors["green"]
	} else if coverage > 60 {
		color = colors["yellow"]
	} else if coverage > 40 {
		color = colors["orange"]
	}

	_ = badgeTemplate.Execute(&buffer, &Badge{coverage, color})
	return buffer.String(), nil
}

/*
<coverage line-rate="0.028990632" branch-rate="0" version="" timestamp="1565726391679" lines-covered="492" lines-valid="16971" branches-covered="0" branches-valid="0" complexity="0">
*/

type cobSummary struct {
	XMLName  xml.Name `xml:"coverage"`
	LineRate float64  `xml:"line-rate,attr"`
}

func main() {
	path := filepath.Join("coverage.xml")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	cob := cobSummary{}

	err = xml.Unmarshal([]byte(data), &cob)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	log.Printf("cob: %#v", cob)

	fmt.Print(RenderBadge(cob.LineRate))
}
