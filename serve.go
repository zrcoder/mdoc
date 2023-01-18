package mdoc

import (
	"encoding/json"
	"fmt"
	tpl "html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/flamego/flamego"
	"github.com/flamego/i18n"
	"github.com/flamego/template"

	"github.com/zrcoder/mdoc/internal/doc"
	"github.com/zrcoder/mdoc/internal/model"
	"github.com/zrcoder/mdoc/internal/static"
	si18n "github.com/zrcoder/mdoc/internal/static/i18n"
	"github.com/zrcoder/mdoc/internal/static/templates"
)

const (
	customTemplatesDir = "templates"
	customI18nDir      = "i18n"
)

func Serve(cfg *model.Config) error {
	docMgr, err := doc.New(cfg)
	if err != nil {
		return err
	}

	f := flamego.New()
	f.Use(flamego.Recovery())
	// custom assets should be served first to support overwrite
	f.Use(customStaticMiddleware(cfg.CustomDirectory))
	// docs must be served
	f.Use(docsDirMiddleware(cfg.DocsBasePath))
	f.Use(builtinStaticMiddleware())
	f.Use(templateMiddleware(cfg))
	languages := cfg.I18nLanguages()
	f.Use(i18nMiddleware(languages))
	f.Use(pageMiddleware(cfg, languages))

	f.Get("/", homeHandler(cfg))
	f.Get(cfg.DocsBasePath+"/?{**}", pageHandler(docMgr))
	f.Any("/webhook", webhookHandler(docMgr))
	f.NotFound(notFound)

	listenAddr := fmt.Sprintf("%s:%s", cfg.HttpAddr, cfg.HttpPort)
	return http.ListenAndServe(listenAddr, f)
}

func customStaticMiddleware(customDirectory string) flamego.Handler {
	option := flamego.StaticOptions{
		Directory: customDirectory,
		SetETag:   true,
	}
	return flamego.Static(option)
}

func docsDirMiddleware(docsDir string) flamego.Handler {
	return flamego.Static(flamego.StaticOptions{Directory: docsDir})
}

func builtinStaticMiddleware() flamego.Handler {
	option := flamego.StaticOptions{
		FileSystem: http.FS(static.Files),
		SetETag:    true,
	}
	return flamego.Static(option)
}

func templateMiddleware(cfg *model.Config) flamego.Handler {
	fs, err := template.EmbedFS(templates.Files, ".", []string{".gohtml"})
	if err != nil {
		return err
	}
	option := template.Options{
		FileSystem:        fs,
		AppendDirectories: []string{filepath.Join(cfg.CustomDirectory, customTemplatesDir)},
		FuncMaps:          []tpl.FuncMap{{"Safe": func(p []byte) tpl.HTML { return tpl.HTML(p) }}},
	}
	return template.Templater(option)
}

func i18nMiddleware(languages []i18n.Language) flamego.Handler {
	option := i18n.Options{
		FileSystem:        http.FS(si18n.Files),
		AppendDirectories: []string{filepath.Join(cfg.CustomDirectory, customI18nDir)},
		Languages:         languages,
	}
	return i18n.I18n(option)
}

func pageMiddleware(cfg *model.Config, languages []i18n.Language) flamego.Handler {
	return func(req *http.Request, data template.Data, locale i18n.Locale) {
		data["Summary"] = cfg
		data["Tr"] = locale.Translate
		data["Lang"] = locale.Lang()
		data["Languages"] = languages
		data["ShowLanguages"] = len(languages) > 1
		data["URL"] = req.URL.Path
	}
}

func homeHandler(cfg *model.Config) flamego.Handler {
	return func(ctx flamego.Context, t template.Template, data template.Data, locale i18n.Locale) {
		if !cfg.HasLandingPage {
			ctx.Redirect(cfg.DocsBasePath)
			return
		}
		data["Title"] = locale.Translate("name") + " - " + locale.Translate("tag_line")
		t.HTML(http.StatusOK, "home")
	}
}

func notFound(t template.Template, data template.Data, locale i18n.Locale) {
	data["Title"] = locale.Translate("status::404")
	t.HTML(http.StatusNotFound, "404")
}

func pageHandler(docMgr *doc.Manager) flamego.Handler {
	return func(ctx flamego.Context, t template.Template, data template.Data, locale i18n.Locale) {
		current := ctx.Param("**")
		if current == "" || current == "/" {
			ctx.Redirect(cfg.DocsBasePath + "/" + docMgr.FirstDocPath())
			return
		}

		if flamego.Env() == flamego.EnvTypeDev {
			err := docMgr.Reload()
			if err != nil {
				panic("reload docs: " + err.Error())
			}
		}

		data["Current"] = current
		data["RootDoc"] = docMgr.Doc(locale.Lang())

		matchedDoc, fallback, err := docMgr.Match(locale.Lang(), current)
		if err != nil {
			notFound(t, data, locale)
			return
		}

		data["Fallback"] = fallback
		data["Category"] = matchedDoc.Category
		data["Title"] = matchedDoc.Title + " - " + locale.Translate("name")
		data["Doc"] = matchedDoc

		if cfg.EditPageLinkFormat != "" {
			blob := strings.TrimPrefix(matchedDoc.LocalPath, cfg.DocsBasePath+"/")
			data["EditLink"] = strings.Replace(cfg.EditPageLinkFormat, "{blob}", blob, 1)
		}
		t.HTML(http.StatusOK, "docs/page")
	}
}

func webhookHandler(docMgr *doc.Manager) flamego.Handler {
	return func(w http.ResponseWriter) {
		err := docMgr.Reload()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": err.Error(),
			})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
