package acme

type Path struct {
	Path string `json:"path"`
}

type DirectoryResp struct {
	NewNonce   string `json:"newNonce"`
	NewAccount string `json:"newAccount"`
	NewOrder   string `json:"newOrder"`
	RevokeCert string `json:"revokeCert"`
	KeyChange  string `json:"keyChange"`
}

type NewNonceResp struct {
	ReplayNonce string `json:"replayNonce"`
}

type NewAccountResp struct {
	Status  string   `json:"status"`
	Contact []string `json:"contact"`
	Orders  string   `json:"orders"`
}
