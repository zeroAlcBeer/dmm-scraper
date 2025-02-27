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

// Configs ...
type Configs struct {
	Output Output
	Proxy  Proxy
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
	}
}
