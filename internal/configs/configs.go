package configs

// Output returns the configuration of output
type Output struct {
	Path string
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
