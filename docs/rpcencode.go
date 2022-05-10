package docs

import _ "embed"

//go:embed rpcencode.md
var rpcEncodeDocument string

const RpcEncodeType = "rpcencode"

func init() {
	addCmdDocumentInfo(RpcEncodeType, rpcEncodeDocument)
}
