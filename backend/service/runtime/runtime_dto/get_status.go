package runtime_dto

type GetStatusInput struct{}

type GetStatusOutput struct {
	State      string `json:"state"`
	Version    string `json:"version"`
	InstallDir string `json:"install_dir"`
	NodePath   string `json:"node_path"`
	NpmPath    string `json:"npm_path"`
	ErrorMsg   string `json:"error_msg"`
}
