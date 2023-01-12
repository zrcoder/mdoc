package doc

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path"

	toc "github.com/abhinav/goldmark-toc"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	gmHl "github.com/yuin/goldmark-highlighting"
	gmMeta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmHtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/zrcoder/mdoc/internal/log"
)

func convertMdFile(parentUrl, localPath string) (content []byte, meta map[string]any, headings toc.Items, err error) {
	body, err := os.ReadFile(localPath)
	if err != nil {
		log.Error(err)
		return nil, nil, nil, fmt.Errorf("read %w", err)
	}

	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			gmHtml.WithHardWraps(),
			gmHtml.WithXHTML(),
			gmHtml.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.GFM,
			gmMeta.Meta,
			emoji.Emoji,
			gmHl.NewHighlighting(
				gmHl.WithStyle("base16-snazzy"),
				gmHl.WithGuessLanguage(true),
			),
			extension.NewFootnote(),
		),
	)

	ctx := parser.NewContext(
		func(cfg *parser.ContextConfig) {
			cfg.IDs = newIDs()
		},
	)
	node := md.Parser().Parse(text.NewReader(body), parser.WithContext(ctx))

	// Headings
	tree, err := toc.Inspect(node, body)
	if err != nil {
		log.Error(err)
		return nil, nil, nil, fmt.Errorf("inspect headings %w", err)
	}
	headings = tree.Items
	if len(headings) > 0 {
		headings = headings[0].Items
	}

	// Links
	err = inspectLinks(parentUrl, node)
	if err != nil {
		log.Error(err)
		return nil, nil, nil, fmt.Errorf("inspect links %w", err)
	}

	var buf bytes.Buffer
	err = md.Renderer().Render(&buf, body, node)
	if err != nil {
		log.Error(err)
		return nil, nil, nil, fmt.Errorf("render %w", err)
	}

	return buf.Bytes(), gmMeta.Get(ctx), headings, nil
}

func inspectLinks(pathPrefix string, doc ast.Node) error {
	return ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		link, ok := n.(*ast.Link)
		if !ok {
			return ast.WalkContinue, nil
		}

		dest, err := url.Parse(string(link.Destination))
		if err != nil {
			return ast.WalkContinue, nil
		}

		if dest.Scheme == "http" || dest.Scheme == "https" {
			return ast.WalkSkipChildren, nil
		} else if dest.Scheme != "" {
			return ast.WalkContinue, nil
		}

		//log.Infof("[test]---%s, %s", pathPrefix, link.Destination)
		link.Destination = convertRelativeLink(pathPrefix, link.Destination)
		//log.Infof("[test]===%s\n", link.Destination)
		return ast.WalkSkipChildren, nil
	})
}

func convertRelativeLink(pathPrefix string, link []byte) []byte {
	var anchor []byte
	if i := bytes.IndexByte(link, '#'); i > -1 {
		if i == 0 {
			return link
		}

		anchor = link[i:]
		link = link[:i]
	}

	// Example: _index.md => /docs/introduction
	if bytes.EqualFold(link, []byte(indexMd)) {
		link = append([]byte(pathPrefix), anchor...)
		return link
	}

	// Example: "installation.md" => "installation"
	link = bytes.TrimSuffix(link, []byte(mdFileExtension))

	// Example: "../howto/_index" => "../howto/"
	link = bytes.TrimSuffix(link, []byte(index))

	// Example: ("/docs", "../howto/") => "/docs/howto"
	link = []byte(path.Join(pathPrefix, string(link)))

	link = append(link, anchor...)
	return link
}

// ids is a modified version to allow any non-whitespace characters instead of
// just alphabets or numerics from
// https://github.com/yuin/goldmark/blob/113ae87dd9e662b54012a596671cb38f311a8e9c/parser/parser.go#L65.
type ids struct {
	values map[string]bool
}

func newIDs() parser.IDs {
	return &ids{
		values: map[string]bool{},
	}
}

func (s *ids) Generate(value []byte, kind ast.NodeKind) []byte {
	value = util.TrimLeftSpace(value)
	value = util.TrimRightSpace(value)
	if len(value) == 0 {
		if kind == ast.KindHeading {
			value = []byte("heading")
		} else {
			value = []byte("id")
		}
	}
	if _, ok := s.values[util.BytesToReadOnlyString(value)]; !ok {
		s.values[util.BytesToReadOnlyString(value)] = true
		return value
	}
	for i := 1; ; i++ {
		newResult := fmt.Sprintf("%s-%d", value, i)
		if _, ok := s.values[newResult]; !ok {
			s.values[newResult] = true
			return []byte(newResult)
		}
	}
}

func (s *ids) Put(value []byte) {
	s.values[util.BytesToReadOnlyString(value)] = true
}
