package ir

const Version = 1

type Program struct {
	Version   int        `json:"version"`
	Files     []*File    `json:"files"`
	Types     []*Type    `json:"types"`
	Symbols   []*Symbol  `json:"symbols"`
	Globals   []*Global  `json:"globals"`
	Functions []*Function `json:"functions"`
}

type File struct {
	Path     string `json:"path"`
	Source   string `json:"source,omitempty"`
	RootNode *Node  `json:"rootNode"`
}

type Node struct {
	Kind     string  `json:"kind"`
	Pos      int     `json:"pos"`
	End      int     `json:"end"`
	Children []*Node `json:"children,omitempty"`
	Symbol   string  `json:"symbol,omitempty"`
	Type     string  `json:"type,omitempty"`
}

type Type struct {
	Id          string   `json:"id"`
	Kind        string   `json:"kind"`
	Flags       uint32   `json:"flags"`
	Name        string   `json:"name,omitempty"`
	Members     []string `json:"members,omitempty"`
	Properties  []string `json:"properties,omitempty"`
	Signatures  []string `json:"signatures,omitempty"`
	TypeArgs    []string `json:"typeArgs,omitempty"`
	Target      string   `json:"target,omitempty"`
	Value       any      `json:"value,omitempty"`
}

type Symbol struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Flags        uint32   `json:"flags"`
	Kind         string   `json:"kind"`
	Type         string   `json:"type,omitempty"`
	Declarations []string `json:"declarations,omitempty"`
	Members      []string `json:"members,omitempty"`
	Exports      []string `json:"exports,omitempty"`
}

type Global struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Type   string `json:"type"`
}

type Function struct {
	Name       string     `json:"name"`
	Symbol     string     `json:"symbol"`
	Signature  string     `json:"signature"`
	Parameters []Param    `json:"parameters"`
	ReturnType string     `json:"returnType"`
	Body       *Block     `json:"body,omitempty"`
}

type Param struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Block struct {
	Statements []Statement `json:"statements"`
}

type Statement struct {
	Kind string `json:"kind"`
	Data any    `json:"data"`
}
