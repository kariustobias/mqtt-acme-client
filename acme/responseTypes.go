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
	Cert       []byte `json:"cert"`
}

type NewNonceResp struct {
	ReplayNonce string `json:"nonce"`
}

type NewAccountResp struct {
	Status  string   `json:"status"`
	Contact []string `json:"contact"`
	Orders  string   `json:"orders"`
}

type NewOrderResp struct {
	Authorizations []string `json:"authorizations"`
	Finalize       string   `json:"finalize"`
	Certificate    string   `json:"certificate"`
}

type NewAuthorizationResp struct {
	Challenges []Challenges `json:"challenges"`
}

type NewChallengeResp struct {
	Status string `json:"status"`
}

type FinalizeResp struct {
	Certificate string `json:"certificate"`
}

type Challenges struct {
	Typ    string `json:"type"`
	Status string `json:"status"`
	Token  string `json:"token"`
	URL    string `json:"url"`
}
