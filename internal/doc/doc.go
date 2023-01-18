package doc

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	toc "github.com/abhinav/goldmark-toc"

	"github.com/zrcoder/mdoc/internal/log"
	"github.com/zrcoder/mdoc/internal/osutil"
)

const (
	index           = "_index"
	mdFileExtension = ".md"
	indexMd         = index + mdFileExtension
)

// first key is language name, second key is doc url path
// for parse cache and route matching cases
var docsMemo map[string]map[string]*Doc

// Doc is a node in the documentation hierarchy.
type Doc struct {
	Category  string
	UrlPath   string // the URL path
	LocalPath string // local path with .md extension

	Weight         int    // in md front matter, used to sort docs
	Title          string // in md front matter, used to sort docs, render TOC and so on
	HideExtraTitle bool   // when there is no title in the front matter, we get title in content and not add extra title
	Content        []byte
	Headings       toc.Items // if Title is empty, use Headings[0].Title

	Groups   []*Doc
	Pages    []*Doc
	Previous *PageLink
	Next     *PageLink
}

// PageLink is a link to another page.
type PageLink struct {
	Title   string
	UrlPath string
}

func (d *Doc) memo(m map[string]*Doc) {
	for _, g := range d.Groups {
		m[g.UrlPath] = g
		g.memo(m)
	}
	for _, p := range d.Pages {
		m[p.UrlPath] = p
	}
}

func (d *Doc) sort() {
	less := func(a, b *Doc) bool {
		return a.Weight < b.Weight || a.Weight == b.Weight && a.Title < b.Title
	}
	sort.Slice(d.Pages, func(i, j int) bool {
		return less(d.Pages[i], d.Pages[j])
	})
	sort.Slice(d.Groups, func(i, j int) bool {
		return less(d.Groups[i], d.Groups[j])
	})
	for _, g := range d.Groups {
		g.sort()
	}
}

func (d *Doc) link(baseUrl string) {
	var pre *Doc
	var dfs func(*Doc)
	dfs = func(doc *Doc) {
		for _, g := range doc.Groups {
			if len(g.Content) != 0 {
				pre = g
			}
			dfs(g)
		}
		for _, p := range doc.Pages {
			if pre != nil {
				pre.Next = &PageLink{
					Title:   p.Title,
					UrlPath: p.UrlPath,
				}
				p.Previous = &PageLink{
					Title:   pre.Title,
					UrlPath: pre.UrlPath,
				}
			}
			pre = p
		}
	}
	dfs(d)
}

// parseDocs initializes documentation hierarchy for given languages in the given
// root directory. The language is the key in the returned map.
func parseDocs(root string, languages []string, baseUrl string) (map[string]*Doc, error) {
	res := make(map[string]*Doc, len(languages))
	docsMemo = map[string]map[string]*Doc{}
	for _, language := range languages {
		dir := filepath.Join(root, language)
		if !osutil.IsDir(dir) {
			log.Info("no i18n resources")
			dir = root
		}
		doc := &Doc{LocalPath: dir, UrlPath: ""}
		err := parseDoc(doc, map[string]bool{})
		if err != nil {
			return nil, err
		}
		doc.sort()
		doc.link(baseUrl)
		docsMemo[language] = map[string]*Doc{}
		doc.memo(docsMemo[language])
		res[language] = doc
	}
	return res, nil
}

func parseDoc(doc *Doc, vis map[string]bool) error {
	err := filepath.WalkDir(doc.LocalPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error(path, "===", err)
			return err
		}
		if path == doc.LocalPath {
			return nil
		}
		if vis[path] {
			return nil
		}
		if d.IsDir() {
			vis[path] = true
			group := &Doc{LocalPath: path, UrlPath: joinUrlPath(doc.UrlPath, d.Name())}
			doc.Groups = append(doc.Groups, group)
			indexPath := filepath.Join(path, indexMd)
			if osutil.IsFile(indexPath) {
				err = convertMdFile(group, doc.UrlPath, indexPath)
				if err != nil {
					return err
				}

			}
			return parseDoc(group, vis)
		}

		name := d.Name()
		if name == indexMd || !strings.HasSuffix(name, mdFileExtension) {
			return nil
		}
		vis[path] = true
		urlPath := joinUrlPath(doc.UrlPath, name[:len(name)-len(mdFileExtension)])
		page := &Doc{
			Category:  doc.Title,
			UrlPath:   urlPath,
			LocalPath: path,
		}
		doc.Pages = append(doc.Pages, page)
		return convertMdFile(page, doc.UrlPath, path)
	})

	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func joinUrlPath(base, last string) string {
	last = strings.ReplaceAll(last, " ", "-")
	if base == "" {
		return last
	}
	return strings.Join([]string{base, last}, "/")
}
