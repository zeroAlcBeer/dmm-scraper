package config

// Output returns the configuration of output
type Output struct {
	Path    string
	NeedCut bool
}

// Proxy returns the configuration of proxy
type Proxy struct {
	Enable bool
	Socket string
}

// DMMApi returns the configuration of DMMApi
type DMMApi struct {
	ApiId       string
	AffiliateId string
}

// Configs ...
type Configs struct {
	Output Output
	Proxy  Proxy
	DMMApi DMMApi
}

// Default ...
func Default() *Configs {
	return &Configs{
		Output: Output{
			Path:    "output/{year}/{num}",
			NeedCut: true,
		},
		Proxy: Proxy{
			Enable: false,
			Socket: "",
		},
		DMMApi: DMMApi{
			ApiId:       "",
			AffiliateId: "",
		},
	}
}
