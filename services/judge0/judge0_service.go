package judgeServ

type Judge0Result struct {
	Stdout        string `json:"stdout"`
	Stderr        string `json:"stderr"`
	CompileOutput string `json:"compile_output"`
	Time          string `json:"time"`
	Memory        int    `json:"memory"`
	Status        struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"status"`
}

type Judge0Client interface {
	SubmitCode(sourceCode, stdin string, languageID int) (*Judge0Result, error)
}
