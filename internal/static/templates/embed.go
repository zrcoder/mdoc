// Copyright 2022 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package templates

import (
	"embed"
)

//go:embed *.gohtml **/*.gohtml
var Files embed.FS
