package workspace

var (
	// 工作路径
	workspaceDir = ""

	// 当前用户目录
	userDir = ""

	// 当前 job 所在路径
	jobDir = ""
)

func GetWorkspace() string {
	return workspaceDir
}

func GetUserDir() string {
	return userDir
}

func GetJobDir() string {
	return jobDir
}

