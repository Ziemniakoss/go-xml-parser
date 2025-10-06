package goxmlparser_test

import (
	"os"
	"strings"
	"testing"

	goxmlparser "github.com/Ziemniakoss/go-xml-parser"
)

func TestGetXmlNodes(t *testing.T) {
	content, err := os.ReadFile("testdata/spring.xml")
	if err != nil {
		t.Error("Could not read test file: " + err.Error())
	}
	xmlDoc, err := goxmlparser.GetXmlNodes(string(content))
	if err != nil {
		t.Error("Could not read doc: " + err.Error())
	}
	if xmlDoc.RootTag.TagName != "beans" {
		t.Error("Wrong root tag " + xmlDoc.RootTag.TagName)
	}
	if len(xmlDoc.RootTag.Children) != 2 {
		t.Error("Root tag should have 2 children, got", len(xmlDoc.RootTag.Children))
	}

	secondNode := xmlDoc.RootTag.Children[1]
	if strings.TrimSpace(secondNode.TextContent) != "Content" {
		t.Error("Wrongly parsed text content of tag:", secondNode.TextContent)
	}
	expectedSecondRange := goxmlparser.Range{
		Start: goxmlparser.Position{
			Line:      10,
			Character: 4,
		},
		End: goxmlparser.Position{
			Line:      12,
			Character: 10,
		},
	}
	if expectedSecondRange != secondNode.Range {
		t.Error("Wrong second element range", expectedSecondRange, secondNode.Range)
	}
}

func TestGetXmlNodes_NoNodes(t *testing.T) {
	content, err := os.ReadFile("testdata/springNoBeans.xml")
	if err != nil {
		t.Error("Could not read test file: " + err.Error())
	}
	xmlDoc, err := goxmlparser.GetXmlNodes(string(content))
	if err != nil {
		t.Error("Could not read doc: " + err.Error())
	}
	if xmlDoc.RootTag.TagName != "beans" {
		t.Error("Wrong root tag " + xmlDoc.RootTag.TagName)
	}
	if len(xmlDoc.RootTag.Children) != 0 {
		t.Error("Children shuld be empty list got")
	}
}
