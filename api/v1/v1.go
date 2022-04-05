package v1

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed api.yaml
var api string

// Get returns api yaml file
func API(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, api)
}
