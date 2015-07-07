package elasticsearch

import (
	"net/url"
)

type API interface {
	Path() string
	Query() url.Values
}

type baseAPI struct {
	Queries url.Values
}

func NewBaseAPI() *baseAPI {
	return &baseAPI{url.Values{}}
}

func (api *baseAPI) Path() string {
	return ""
}

func (api *baseAPI) Query() url.Values {
	return api.Queries
}

func (api *baseAPI) singleOption(key, value string) {
	api.Queries[key] = []string{value}
}

func (api *baseAPI) Pretty() {
	api.singleOption("pretty", "true")
}

func (api *baseAPI) YAML() {
	api.singleOption("format", "yaml")
}

func (api *baseAPI) MachineReadable() {
	api.singleOption("human", "false")
}

func (api *baseAPI) AddFilterPath(paths ...string) {
	api.Queries["filter_path"] = append(api.Queries["filter_path"], paths...)
}

func (api *baseAPI) FlatSettings() {
	api.singleOption("flat_settings", "true")
}

func (api *baseAPI) CamelCase() {
	api.singleOption("camelCase", "true")
}
