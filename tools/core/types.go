package core

type RunOptions struct {
	Bin    string   `json:"bin"`
	Args   []string `json:"args"`
	Stdout bool     `json:"output"`
}
