package docs

import _ "embed"

//go:embed rpcdecode.md
var rpcDecodeDocument string

const RpcDecodeType = "rpcdecode"

func init() {
	addCmdDocumentInfo(RpcDecodeType, rpcDecodeDocument)
}
