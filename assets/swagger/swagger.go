package swagger

import (
	"embed"
	"fmt"
	"net/http"
)

//go:embed *
var f embed.FS

// Get returns the embedded swagger assets
func Get() embed.FS {
	return f
}

// Index implements an http.HandlerFunc that provides all links to an internal
// swagger bundle.
func Index(w http.ResponseWriter, _ *http.Request) {
	const html = `
<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
	<meta charset="UTF-8">
	<title>Swagger UI</title>
	<link rel="stylesheet" type="text/css" href="/swagger/swagger-ui.css" />
	<link rel="stylesheet" type="text/css" href="/swagger/index.css" />
	<link rel="icon" type="image/png" href="/swagger/favicon-32x32.png" sizes="32x32" />
	<link rel="icon" type="image/png" href="/swagger/favicon-16x16.png" sizes="16x16" />
  </head>

  <body>
	<div id="swagger-ui"></div>
	<script src="/swagger/swagger-ui-bundle.js" charset="UTF-8"> </script>
	<script src="/swagger/swagger-ui-standalone-preset.js" charset="UTF-8"> </script>
	<script src="/swagger/swagger-initializer.js" charset="UTF-8"> </script>
  </body>
</html>
	`
	fmt.Fprint(w, html)
}
