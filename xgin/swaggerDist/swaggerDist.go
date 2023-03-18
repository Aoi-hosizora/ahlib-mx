package swaggerDist

import (
	_ "embed"
)

// The following files are from:
// - https://github.com/swagger-api/swagger-ui/releases/tag/v3.17.2

//go:embed favicon-16x16.png
var Favicon_16x16_png []byte

//go:embed favicon-32x32.png
var Favicon_32x32_png []byte

//go:embed index.html
var Index_html []byte

//go:embed oauth2-redirect.html
var Oauth2_redirect_html []byte

//go:embed swagger-ui.css
var Swagger_ui_css []byte

//go:embed swagger-ui.js
var Swagger_ui_js []byte

//go:embed swagger-ui-bundle.js
var Swagger_ui_bundle_js []byte

//go:embed swagger-ui-standalone-preset.js
var Swagger_ui_standalone_preset_js []byte
