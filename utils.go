package goxmlparser

import (
	"errors"

	antlr "github.com/antlr4-go/antlr/v4"
)

// Position in a text
type Position struct {
	// Line position in a document (zero-based).
	Line uint `json:"line"`

	// Character offset on a line in a document (zero-based). Assuming that the line is
	// represented as a string, the `character` value represents the gap between the
	// `character` and `character + 1`.
	//
	// If the character value is greater than the line length it defaults back to the
	// line length.
	Character uint `json:"character"`
}

/**
 * A range in a text document expressed as (zero-based) start and end positions.
 *
 * If you want to specify a range that contains a line including the line ending
 * character(s) then use an end position denoting the start of the next line.
 * For example:
 * ```ts
 * {
 *     start: { line: 5, character: 23 }
 *     end : { line 6, character : 0 }
 * }
 * ```
 */
type Range struct {
	// The range's start position
	Start Position
	// The range's end position.
	End Position
}

type XmlNode struct {
	TagName     string
	Range       Range
	Children    []XmlNode
	TextContent string
}

type XmlDocument struct {
	RootTag XmlNode
}

func ExtractRange(ctx antlr.ParserRuleContext) Range {
	return Range{
		Start: Position{
			Character: uint(ctx.GetStart().GetColumn()),
			Line:      uint(ctx.GetStart().GetLine()) - 1,
		},
		End: Position{
			Character: uint(ctx.GetStop().GetColumn()),
			Line:      uint(ctx.GetStop().GetLine()) - 1,
		},
	}
}

func GetXmlNodes(xmlString string) (XmlDocument, error) {
	is := antlr.NewInputStream(xmlString)
	lexer := NewXMLLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewXMLParser(stream)
	parser.RemoveErrorListeners()
	doc := parser.Document()
	var nodes []XmlNode
	for _, child := range doc.GetChildren() {
		eleCtx, ok := child.(*ElementContext)
		if ok {
			if child, err := convertToXmlNode(eleCtx); err == nil {
				nodes = append(nodes, child)
			}
		}
	}
	if len(nodes) > 0 {
		return XmlDocument{RootTag: nodes[0]}, nil
	}
	return XmlDocument{}, errors.New("Empty Document")
}

func convertToXmlNode(ctx *ElementContext) (XmlNode, error) {
	names := ctx.AllName()
	if len(names) == 0 {
		return XmlNode{}, errors.New("Tag has no name")
	}
	tagName := names[0].GetText()
	content := ""
	var contentCtx = ctx.Content()
	var children []XmlNode
	if contentCtx != nil {
		content = contentCtx.GetText()
		for _, child := range contentCtx.GetChildren() {
			eleCtx, ok := child.(*ElementContext)
			if ok {
				if child, err := convertToXmlNode(eleCtx); err == nil {
					children = append(children, child)
				}
			}
		}
	}
	return XmlNode{
		Range:       ExtractRange(ctx),
		TagName:     tagName,
		Children:    children,
		TextContent: content,
	}, nil
}
