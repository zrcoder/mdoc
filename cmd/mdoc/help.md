# mdoc

A dead simple document server

## usage

```shell
mdoc <directory>
```

## config

you can make a `doc.yaml` file to config your document, a simple of the config file is like below:

```yaml
# The local address to listen on, use "0.0.0.0" to listen on all network interfaces
httpAddr: localhost
# The local port to listen on
httpPort: 9999
# The description of the site
description: mdoc is a dead simple web server for markdown documentation
# The local directory that contains Markdown docs, if is empty, current directory will be used
docsDirectory: "docs"
# The local directory that contains custom assets (e.g. templates, images, CSS, JavaScript, local files)
customDirectory: static
# Whether the site has a landing page, set to "false" to redirect users for documentation directly
hasLandingPage: true
# Whether the pages have a navigation bar
hasNavBar: true
# The url base path for documentation, start with '/', can be empty
docsBasePath: /docs
# The list of languages that is supported
languages:
  - name: en-US
    description: English
  - name: zh-CN
    description: 简体中文
# The format to construct a edit page link, leave it empty to disable, e.g.
# https://github.com/zrcoder/mdoc/blob/main/docs/{blob}
editPageLinkFormat:
```
