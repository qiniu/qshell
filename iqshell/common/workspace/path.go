package workspace

var (
	// 工作路径
	workspacePath = ""

	// 当前用户目录
	userPath = ""
)

func GetWorkspace() string {
	return workspacePath
}

func GetUserPath() string {
	return userPath
}
