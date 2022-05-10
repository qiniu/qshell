package docs

import _ "embed"

//go:embed token.md
var tokenDocument string

const TokenType = "token"

func init() {
	addCmdDocumentInfo(TokenType, tokenDocument)
}
