package doc

import (
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"

	"github.com/zrcoder/mdoc/internal/log"
	"github.com/zrcoder/mdoc/internal/model"
	"github.com/zrcoder/mdoc/internal/osutil"
)

// Manager is a store maintaining documentation hierarchies for multiple languages.
type Manager struct {
	// The list of config values
	languages   []string
	baseURLPath string
	rootDir     string

	// The list of inferred values
	defaultLanguage string
	docs            atomic.Value
	docsMemo        atomic.Value
	reloadLock      sync.Mutex
}

// New initializes the documentation store from given config.
func New(cfg *model.Config) (*Manager, error) {
	i18nLanguages := cfg.I18nLanguages()
	languages := make([]string, len(i18nLanguages))
	for i, v := range i18nLanguages {
		languages[i] = v.Name
	}
	s := &Manager{
		languages:       languages,
		baseURLPath:     cfg.DocsBasePath,
		defaultLanguage: languages[0],
		rootDir:         cfg.DocsDirectory,
	}
	err := s.Reload()
	if err != nil {
		return nil, errors.Wrap(err, "reload")
	}
	return s, nil
}

// Reload re-initializes the documentation store.
func (s *Manager) Reload() error {
	s.reloadLock.Lock()
	defer s.reloadLock.Unlock()
	if !osutil.IsDir(s.rootDir) {
		log.Error("not exist nor is dir:", s.rootDir)
		return errors.Errorf("directory root %q does not exist", s.rootDir)
	}

	docs, err := parseDocs(s.rootDir, s.languages, s.baseURLPath)
	if err != nil {
		return errors.Wrap(err, "init doc")
	}

	s.setDocs(docs)
	s.setDocsMemo(docsMemo)
	return nil
}

func (s *Manager) getDocs() map[string]*Doc {
	return s.docs.Load().(map[string]*Doc)
}

func (s *Manager) setDocs(docs map[string]*Doc) {
	s.docs.Store(docs)
}
func (s *Manager) setDocsMemo(memo map[string]map[string]*Doc) {
	s.docsMemo.Store(memo)
}
func (s *Manager) getDocsMemo() map[string]map[string]*Doc {
	return s.docsMemo.Load().(map[string]map[string]*Doc)
}

// FirstDocPath returns the URL path of the first doc that has content in the
// default language.
func (s *Manager) FirstDocPath() string {
	return firstDocPath(s.getDocs()[s.defaultLanguage])
}

func firstDocPath(doc *Doc) string {
	if len(doc.Groups) == 0 {
		if len(doc.Pages) == 0 {
			return "404"
		}
		return doc.Pages[0].UrlPath
	}
	if doc.Groups[0].Content != nil {
		return doc.Groups[0].UrlPath
	}
	return firstDocPath(doc.Groups[0])
}

// Doc returns the Doc of the given language. It returns the Doc of the default
// language if the given language is not found.
func (s *Manager) Doc(language string) *Doc {
	docs := s.getDocs()
	doc, ok := docs[language]
	if !ok {
		return docs[s.defaultLanguage]
	}
	return doc
}

var ErrNoMatch = errors.New("no match for the path")

// Match matches a node with given path in given language. If there is no such
// node exists or the node content is empty, it falls back to use the node with
// the same path in default language.
func (s *Manager) Match(language, path string) (d *Doc, fallback bool, err error) {
	docMemo := s.getDocsMemo()
	memo := docMemo[language]
	if len(memo) == 0 {
		memo = docMemo[s.defaultLanguage]
		if len(memo) == 0 {
			return nil, false, ErrNoMatch
		}
	}

	d, ok := memo[path]
	if ok && len(d.Content) > 0 {
		return d, false, nil
	}

	return nil, false, ErrNoMatch
}
