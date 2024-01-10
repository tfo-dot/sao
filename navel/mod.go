package navel

import "encoding/json"

type Navel struct {
	Source       NavelScript
	ID           *string
	Markings     []string
	BranchStates map[string]bool
}

func NewScript(content string) Navel {
	var script NavelScript

	err := json.Unmarshal([]byte(content), &script)

	if err != nil {
		panic(err)
	}

	return Navel{
		Source:       script,
		ID:           nil,
		Markings:     make([]string, 0),
		BranchStates: map[string]bool{},
	}
}

func (n *Navel) SetOption(id string) *NavelRegular {
	return nil
}

type NavelScript struct {
	Entrypoint string
	Objects    []NavelObject
}

type NavelType string

var (
	NavelTypeRegular *NavelType = nil
	NavelTypeProxy   NavelType  = "proxy"
	NavelTypeBranch  NavelType  = "dependant"
)

type NavelObject struct {
	ID       string  `json:"id"`
	Goto     string  `json:"goto"`
	Type     *string `json:"type"`
	Content  string  `json:"content"`
	Response struct {
		Title   string        `json:"title"`
		Options []NavelOption `json:"options"`
	} `json:"response"`
	Branch struct {
		Condition []string     `json:"condition"`
		Then      NavelRegular `json:"then"`
		Else      NavelRegular `json:"else"`
	} `json:"branch"`
}

type NavelProxy struct {
	ID   string  `json:"id"`
	Goto string  `json:"goto"`
	Type *string `json:"type"`
}

type NavelRegular struct {
	ID       string  `json:"id"`
	Content  string  `json:"content"`
	Type     *string `json:"type"`
	Response struct {
		Title   string        `json:"title"`
		Options []NavelOption `json:"options"`
	} `json:"response"`
}

type NavelBranch struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Branch struct {
		Condition []string     `json:"condition"`
		Then      NavelRegular `json:"then"`
		Else      NavelRegular `json:"else"`
	} `json:"branch"`
}

type NavelOption struct {
	Label     string   `json:"label"`
	Goto      string   `json:"goto"`
	Input     string   `json:"input"`
	Mark      string   `json:"mark"`
	BlockedBy []string `json:"blockedBy"`
}
