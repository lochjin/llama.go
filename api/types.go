package api

type PropsResponse struct {
	BuildInfo  string     `json:"build_info"`
	ModelPath  string     `json:"model_path"`
	NCtx       int64      `json:"n_ctx"`
	Modalities Modalities `json:"modalities"`
}

type Modalities struct {
	Vision bool `json:"vision"`
	Audio  bool `json:"audio"`
}
