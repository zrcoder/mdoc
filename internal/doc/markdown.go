package doc

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

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

var (
	ErrNoTitleFound = errors.New("no title found in markdown")
)

func convertMdFile(doc *Doc, parentUrl, localPath string) (err error) {
	body, err := os.ReadFile(localPath)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("read %w", err)
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
		func(cfg *parser.ContextConfig) { cfg.IDs = newIDs() },
	)
	node := md.Parser().Parse(text.NewReader(body), parser.WithContext(ctx))

	// Headings
	tree, err := toc.Inspect(node, body)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("inspect headings %w", err)
	}
	headings := tree.Items
	if len(headings) > 0 {
		headings = headings[0].Items
	}
	doc.Headings = headings

	// Links
	err = inspectLinks(parentUrl, node)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("inspect links %w", err)
	}

	var buf bytes.Buffer
	err = md.Renderer().Render(&buf, body, node)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("render %w", err)
	}
	doc.Content = buf.Bytes()
	meta := gmMeta.Get(ctx)
	weight, ok := meta["weight"].(int)
	if ok {
		doc.Weight = weight
	}
	doc.Title, ok = meta["title"].(string)
	if !ok || doc.Title == "" {
		title, err := parseTitleFromContent(body)
		if err != nil {
			return err
		}
		doc.Title = title
		doc.HideExtraTitle = true
	}
	return nil
}

// parseTitleFromContent parse the markdown content to get the title
func parseTitleFromContent(markdownContent []byte) (string, error) {
	const titlePrefix = "# "
	index := bytes.Index(markdownContent, []byte(titlePrefix))
	if index == -1 {
		return "", ErrNoTitleFound
	}
	markdownContent = markdownContent[index+len(titlePrefix):]
	index = bytes.Index(markdownContent, []byte("\n"))
	if index == -1 {
		return "", ErrNoTitleFound
	}
	line := string(markdownContent[:index])
	line = strings.TrimSpace(line)
	if line == "" {
		return "", ErrNoTitleFound
	}
	if line[0] == '[' {
		index = strings.Index(line, "]")
		if index == -1 {
			return "", ErrNoTitleFound
		}
		return line[1:index], nil
	}
	return line, nil
}

func inspectLinks(parentUrl string, node ast.Node) error {
	return ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
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

		link.Destination = convertRelativeLink(parentUrl, link.Destination)
		return ast.WalkSkipChildren, nil
	})
}

func convertRelativeLink(parentUrl string, link []byte) []byte {
	var anchor []byte
	if i := bytes.IndexByte(link, '#'); i > -1 {
		if i == 0 {
			return link
		}

		anchor = link[i:]
		link = link[:i]
	}

	// _index.md => {parentUrl}
	if bytes.EqualFold(link, []byte(indexMd)) {
		link = append([]byte(parentUrl), anchor...)
		return link
	}

	// xxx.md => xxx
	link = bytes.TrimSuffix(link, []byte(mdFileExtension))

	// ../xxx/_index => ../xxx/
	link = bytes.TrimSuffix(link, []byte(index))

	// Example: ("/docs", "../howto/") => "/docs/howto" TODO: test
	link = []byte(path.Join(parentUrl, string(link)))

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
