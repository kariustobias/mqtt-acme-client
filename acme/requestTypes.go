package acme

type DirectoryReq struct {
	Path string `json:"path"`
}

type NewNonceReq struct {
	Path string `json:"path"`
}

type NewAccountReq struct {
	Contact              []string `json:"contact"`
	TermsOfServiceAgreed bool     `json:"termsOfServiceAgreed"`
}

type NewOrderReq struct {
	Identifiers []Identifier `json:"identifiers"`
}

type FinalizeReq struct {
	CSR string `json:"csr"`
}

type Identifier struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
