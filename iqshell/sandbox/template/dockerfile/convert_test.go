package dockerfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// 基础测试
// ---------------------------------------------------------------------------

func TestConvert_BasicDockerfile(t *testing.T) {
	content := `FROM ubuntu:22.04
RUN apt-get update && apt-get install -y curl
WORKDIR /app
COPY . /app/
CMD ["python3", "app.py"]
`
	result, err := Convert(content)
	require.NoError(t, err)

	assert.Equal(t, "ubuntu:22.04", result.BaseImage)
	assert.Equal(t, "python3 app.py", result.StartCmd)
	assert.Equal(t, "sleep 20", result.ReadyCmd)

	// 步骤：USER root, WORKDIR /, RUN, WORKDIR /app, COPY, USER user
	// （因为已设置 WORKDIR，末尾不追加 WORKDIR /home/user）
	assert.Len(t, result.Steps, 6)
	assert.Equal(t, "USER", result.Steps[0].Type)
	assert.Equal(t, []string{"root"}, *result.Steps[0].Args)
	assert.Equal(t, "WORKDIR", result.Steps[1].Type)
	assert.Equal(t, []string{"/"}, *result.Steps[1].Args)
	assert.Equal(t, "RUN", result.Steps[2].Type)
	assert.Equal(t, []string{"apt-get update && apt-get install -y curl"}, *result.Steps[2].Args)
	assert.Equal(t, "WORKDIR", result.Steps[3].Type)
	assert.Equal(t, []string{"/app"}, *result.Steps[3].Args)
	assert.Equal(t, "COPY", result.Steps[4].Type)
	assert.Equal(t, []string{".", "/app/", "", ""}, *result.Steps[4].Args)
	assert.Equal(t, "USER", result.Steps[5].Type)
	assert.Equal(t, []string{"user"}, *result.Steps[5].Args)
}

func TestConvert_NoFROM(t *testing.T) {
	_, err := Convert("RUN echo hello")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no FROM instruction")
}

func TestConvert_DefaultUserAndWorkdir(t *testing.T) {
	content := `FROM alpine:3.18
RUN echo hello
`
	result, err := Convert(content)
	require.NoError(t, err)

	last2 := result.Steps[len(result.Steps)-2:]
	assert.Equal(t, "USER", last2[0].Type)
	assert.Equal(t, []string{"user"}, *last2[0].Args)
	assert.Equal(t, "WORKDIR", last2[1].Type)
	assert.Equal(t, []string{"/home/user"}, *last2[1].Args)
}

func TestConvert_ExplicitUserAndWorkdir(t *testing.T) {
	content := `FROM alpine:3.18
USER myuser
WORKDIR /opt/app
`
	result, err := Convert(content)
	require.NoError(t, err)

	lastStep := result.Steps[len(result.Steps)-1]
	assert.Equal(t, "WORKDIR", lastStep.Type)
	assert.Equal(t, []string{"/opt/app"}, *lastStep.Args)
}

func TestConvert_ENV_KeyValueForm(t *testing.T) {
	content := `FROM alpine
ENV FOO=bar BAZ="hello world"
`
	result, err := Convert(content)
	require.NoError(t, err)

	envStep := result.Steps[2]
	assert.Equal(t, "ENV", envStep.Type)
	assert.Equal(t, []string{"FOO", "bar", "BAZ", "hello world"}, *envStep.Args)
}

func TestConvert_ENV_LegacyForm(t *testing.T) {
	content := `FROM alpine
ENV MY_VAR some value with spaces
`
	result, err := Convert(content)
	require.NoError(t, err)

	envStep := result.Steps[2]
	assert.Equal(t, "ENV", envStep.Type)
	assert.Equal(t, []string{"MY_VAR", "some value with spaces"}, *envStep.Args)
}

func TestConvert_ARG(t *testing.T) {
	content := `FROM alpine
ARG VERSION=1.0
ARG NAME
`
	result, err := Convert(content)
	require.NoError(t, err)

	// ARG VERSION=1.0 有默认值，生成 ENV 步骤
	argStep1 := result.Steps[2]
	assert.Equal(t, "ENV", argStep1.Type)
	assert.Equal(t, []string{"VERSION", "1.0"}, *argStep1.Args)

	// ARG NAME 无默认值，不生成 ENV 步骤，下一个应该是默认追加的 USER user
	argStep2 := result.Steps[3]
	assert.Equal(t, "USER", argStep2.Type)
	assert.Equal(t, []string{"user"}, *argStep2.Args)
}

func TestConvert_CMD_ShellForm(t *testing.T) {
	content := `FROM alpine
CMD echo hello world
`
	result, err := Convert(content)
	require.NoError(t, err)
	assert.Equal(t, "echo hello world", result.StartCmd)
	assert.Equal(t, "sleep 20", result.ReadyCmd)
}

func TestConvert_CMD_ExecForm(t *testing.T) {
	content := `FROM alpine
CMD ["node", "server.js", "--port", "3000"]
`
	result, err := Convert(content)
	require.NoError(t, err)
	assert.Equal(t, "node server.js --port 3000", result.StartCmd)
}

func TestConvert_ENTRYPOINT(t *testing.T) {
	content := `FROM alpine
ENTRYPOINT ["python3", "-m", "flask"]
`
	result, err := Convert(content)
	require.NoError(t, err)
	assert.Equal(t, "python3 -m flask", result.StartCmd)
	assert.Equal(t, "sleep 20", result.ReadyCmd)
}

func TestConvert_COPY_WithChown(t *testing.T) {
	content := `FROM alpine
COPY --chown=myuser:mygroup src/ /dest/
`
	result, err := Convert(content)
	require.NoError(t, err)

	copyStep := result.Steps[2]
	assert.Equal(t, "COPY", copyStep.Type)
	assert.Equal(t, []string{"src/", "/dest/", "myuser", ""}, *copyStep.Args)
}

func TestConvert_ADD(t *testing.T) {
	content := `FROM alpine
ADD app.tar.gz /opt/
`
	result, err := Convert(content)
	require.NoError(t, err)

	addStep := result.Steps[2]
	assert.Equal(t, "COPY", addStep.Type)
	assert.Equal(t, []string{"app.tar.gz", "/opt/", "", ""}, *addStep.Args)
}

func TestConvert_ContinuationLines(t *testing.T) {
	content := `FROM alpine
RUN apt-get update && \
    apt-get install -y curl && \
    apt-get clean
`
	result, err := Convert(content)
	require.NoError(t, err)

	runStep := result.Steps[2]
	assert.Equal(t, "RUN", runStep.Type)
	assert.Contains(t, (*runStep.Args)[0], "apt-get update")
	assert.Contains(t, (*runStep.Args)[0], "apt-get install -y curl")
	assert.Contains(t, (*runStep.Args)[0], "apt-get clean")
}

func TestConvert_IgnoredInstructions(t *testing.T) {
	content := `FROM alpine
EXPOSE 8080
VOLUME /data
LABEL maintainer="test"
RUN echo done
`
	result, err := Convert(content)
	require.NoError(t, err)

	// 仅包含 USER root, WORKDIR /, RUN, USER user, WORKDIR /home/user
	assert.Len(t, result.Steps, 5)
}

func TestConvert_Comments(t *testing.T) {
	content := `# This is a comment
FROM alpine
# Another comment
RUN echo hello
`
	result, err := Convert(content)
	require.NoError(t, err)
	assert.Equal(t, "alpine", result.BaseImage)
	assert.Len(t, result.Steps, 5)
}

func TestConvert_FROMWithPlatform(t *testing.T) {
	content := `FROM --platform=linux/amd64 ubuntu:22.04
RUN echo test
`
	result, err := Convert(content)
	require.NoError(t, err)
	assert.Equal(t, "ubuntu:22.04", result.BaseImage)
}

func TestConvert_FROMWithAlias(t *testing.T) {
	content := `FROM ubuntu:22.04 AS builder
RUN echo test
`
	result, err := Convert(content)
	require.NoError(t, err)
	assert.Equal(t, "ubuntu:22.04", result.BaseImage)
}

func TestConvert_NoCmdNoStartCmd(t *testing.T) {
	content := `FROM alpine
RUN echo hello
`
	result, err := Convert(content)
	require.NoError(t, err)
	assert.Equal(t, "", result.StartCmd)
	assert.Equal(t, "", result.ReadyCmd)
}

// ---------------------------------------------------------------------------
// 转义指令测试
// ---------------------------------------------------------------------------

func TestConvert_EscapeDirectiveBacktick(t *testing.T) {
	content := "# escape=`\nFROM alpine\nRUN echo hello `\n    && echo world\n"
	result, err := Convert(content)
	require.NoError(t, err)

	runStep := result.Steps[2]
	assert.Contains(t, (*runStep.Args)[0], "echo hello")
	assert.Contains(t, (*runStep.Args)[0], "echo world")
}

func TestConvert_EscapeDirectiveBackslash(t *testing.T) {
	// 显式 # escape=\ 应与默认行为一致
	content := "# escape=\\\nFROM alpine\nRUN echo hello \\\n    && echo world\n"
	result, err := Convert(content)
	require.NoError(t, err)

	runStep := result.Steps[2]
	assert.Contains(t, (*runStep.Args)[0], "echo hello")
	assert.Contains(t, (*runStep.Args)[0], "echo world")
}

func TestConvert_EscapeDirectiveBacktickNoJoinOnBackslash(t *testing.T) {
	// 使用反引号转义时，反斜杠不应作为续行符
	content := "# escape=`\nFROM alpine\nRUN echo C:\\Users\\test\n"
	result, err := Convert(content)
	require.NoError(t, err)

	runStep := result.Steps[2]
	assert.Equal(t, `echo C:\Users\test`, (*runStep.Args)[0])
}

func TestConvert_EscapeDirectiveMustBeFirst(t *testing.T) {
	// 非指令注释在 escape= 之前，意味着 escape 指令被忽略
	content := "# this is a comment\n# escape=`\nFROM alpine\nRUN echo hello `\n    && echo world\n"
	result, err := Convert(content)
	require.NoError(t, err)

	// 反引号续行不应生效（escape= 在非指令注释之后）
	// RUN 行应该只是 "echo hello `"
	runStep := result.Steps[2]
	assert.Equal(t, "echo hello `", (*runStep.Args)[0])
}

// ---------------------------------------------------------------------------
// Heredoc 测试
// ---------------------------------------------------------------------------

func TestConvert_HeredocRUN(t *testing.T) {
	content := `FROM alpine
RUN <<EOF
echo hello
echo world
EOF
`
	result, err := Convert(content)
	require.NoError(t, err)

	runStep := result.Steps[2]
	assert.Equal(t, "RUN", runStep.Type)
	assert.Equal(t, "echo hello\necho world", (*runStep.Args)[0])
}

func TestConvert_HeredocUnterminatedError(t *testing.T) {
	content := `FROM alpine
RUN <<EOF
echo hello
echo world
`
	_, err := Convert(content)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unterminated heredoc")
}

func TestConvert_HeredocCustomMarker(t *testing.T) {
	content := `FROM alpine
RUN <<SCRIPT
#!/bin/bash
apt-get update
apt-get install -y curl
SCRIPT
`
	result, err := Convert(content)
	require.NoError(t, err)

	runStep := result.Steps[2]
	assert.Contains(t, (*runStep.Args)[0], "apt-get update")
	assert.Contains(t, (*runStep.Args)[0], "apt-get install -y curl")
}

// ---------------------------------------------------------------------------
// 引号和转义处理测试
// ---------------------------------------------------------------------------

func TestConvert_ENV_DoubleQuotedEscape(t *testing.T) {
	// 双引号值内的转义引号
	content := `FROM alpine
ENV MSG="it's \"here\""
`
	result, err := Convert(content)
	require.NoError(t, err)

	envStep := result.Steps[2]
	assert.Equal(t, "ENV", envStep.Type)
	assert.Equal(t, []string{"MSG", `it's "here"`}, *envStep.Args)
}

func TestConvert_ENV_SingleQuotedLiteral(t *testing.T) {
	// 单引号：所有内容为字面量，不处理转义
	content := `FROM alpine
ENV MSG='hello\nworld'
`
	result, err := Convert(content)
	require.NoError(t, err)

	envStep := result.Steps[2]
	assert.Equal(t, []string{"MSG", `hello\nworld`}, *envStep.Args)
}

func TestConvert_ENV_MultiplePairsMixed(t *testing.T) {
	content := `FROM alpine
ENV A=1 B="two" C='three'
`
	result, err := Convert(content)
	require.NoError(t, err)

	envStep := result.Steps[2]
	assert.Equal(t, []string{"A", "1", "B", "two", "C", "three"}, *envStep.Args)
}

func TestConvert_ENV_EmptyValue(t *testing.T) {
	content := `FROM alpine
ENV EMPTY=""
`
	result, err := Convert(content)
	require.NoError(t, err)

	envStep := result.Steps[2]
	assert.Equal(t, []string{"EMPTY", ""}, *envStep.Args)
}

// ---------------------------------------------------------------------------
// 多阶段构建警告测试
// ---------------------------------------------------------------------------

func TestConvert_MultiStageWarning(t *testing.T) {
	content := `FROM golang:1.21 AS builder
RUN go build -o /app
FROM alpine:3.18
COPY --from=builder /app /app
CMD ["/app"]
`
	result, err := Convert(content)
	require.NoError(t, err)

	// 应使用最后一个 FROM 作为基础镜像
	assert.Equal(t, "alpine:3.18", result.BaseImage)
	// 应有多阶段构建的警告
	assert.Len(t, result.Warnings, 1)
	assert.Contains(t, result.Warnings[0], "multi-stage")
}

// ---------------------------------------------------------------------------
// 未知指令警告测试
// ---------------------------------------------------------------------------

func TestConvert_UnknownInstructionWarning(t *testing.T) {
	content := `FROM alpine
FOOBAR something
RUN echo hello
`
	result, err := Convert(content)
	require.NoError(t, err)

	assert.Len(t, result.Warnings, 1)
	assert.Contains(t, result.Warnings[0], "unknown instruction")
	assert.Contains(t, result.Warnings[0], "FOOBAR")
}

func TestConvert_ONBUILDWarning(t *testing.T) {
	content := `FROM alpine
ONBUILD RUN echo hello
RUN echo world
`
	result, err := Convert(content)
	require.NoError(t, err)

	assert.Len(t, result.Warnings, 1)
	assert.Contains(t, result.Warnings[0], "ONBUILD")
}

// ---------------------------------------------------------------------------
// BOM 处理测试
// ---------------------------------------------------------------------------

func TestConvert_BOMStripping(t *testing.T) {
	// UTF-8 BOM 前缀
	content := "\xef\xbb\xbfFROM alpine\nRUN echo hello\n"
	result, err := Convert(content)
	require.NoError(t, err)
	assert.Equal(t, "alpine", result.BaseImage)
}

// ---------------------------------------------------------------------------
// 续行中的注释
// ---------------------------------------------------------------------------

func TestConvert_ContinuationSkipsComments(t *testing.T) {
	content := `FROM alpine
RUN echo hello \
    # this is a comment inside continuation \
    && echo world
`
	result, err := Convert(content)
	require.NoError(t, err)

	runStep := result.Steps[2]
	cmd := (*runStep.Args)[0]
	assert.Contains(t, cmd, "echo hello")
	assert.Contains(t, cmd, "echo world")
	assert.NotContains(t, cmd, "this is a comment")
}

// ---------------------------------------------------------------------------
// 错误行号测试
// ---------------------------------------------------------------------------

func TestConvert_ErrorIncludesLineNumber(t *testing.T) {
	content := `FROM alpine
RUN echo hello
WORKDIR
`
	_, err := Convert(content)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "line")
	assert.Contains(t, err.Error(), "empty WORKDIR")
}

// ---------------------------------------------------------------------------
// 带多个标志的 COPY
// ---------------------------------------------------------------------------

func TestConvert_COPY_WithFromAndChmod(t *testing.T) {
	content := `FROM alpine
COPY --from=builder --chmod=755 /app /app
`
	result, err := Convert(content)
	require.NoError(t, err)

	copyStep := result.Steps[2]
	assert.Equal(t, "COPY", copyStep.Type)
	assert.Equal(t, []string{"/app", "/app", "", ""}, *copyStep.Args)
}

// ---------------------------------------------------------------------------
// Windows 风格换行符
// ---------------------------------------------------------------------------

func TestConvert_WindowsLineEndings(t *testing.T) {
	content := "FROM alpine\r\nRUN echo hello\r\nWORKDIR /app\r\n"
	result, err := Convert(content)
	require.NoError(t, err)
	assert.Equal(t, "alpine", result.BaseImage)
	assert.Equal(t, "WORKDIR", result.Steps[3].Type)
	assert.Equal(t, []string{"/app"}, *result.Steps[3].Args)
}

// ---------------------------------------------------------------------------
// Dockerfile 构建场景测试
// ---------------------------------------------------------------------------

func TestConvert_TypicalWebApp(t *testing.T) {
	// 典型的 Web 应用 Dockerfile，包含 RUN + COPY + WORKDIR + CMD
	content := `FROM node:18-alpine
RUN apk add --no-cache python3 make g++
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --production
COPY . .
CMD ["node", "server.js"]
`
	result, err := Convert(content)
	require.NoError(t, err)

	assert.Equal(t, "node:18-alpine", result.BaseImage)
	assert.Equal(t, "node server.js", result.StartCmd)
	assert.Equal(t, "sleep 20", result.ReadyCmd)

	// 步骤：USER root, WORKDIR /, RUN apk, WORKDIR /app, COPY package*, RUN npm, COPY ., USER user
	assert.Len(t, result.Steps, 8)
	assert.Equal(t, "USER", result.Steps[0].Type)
	assert.Equal(t, []string{"root"}, *result.Steps[0].Args)
	assert.Equal(t, "WORKDIR", result.Steps[1].Type)
	assert.Equal(t, []string{"/"}, *result.Steps[1].Args)
	assert.Equal(t, "RUN", result.Steps[2].Type)
	assert.Equal(t, "WORKDIR", result.Steps[3].Type)
	assert.Equal(t, "COPY", result.Steps[4].Type)
	assert.Equal(t, "RUN", result.Steps[5].Type)
	assert.Equal(t, "COPY", result.Steps[6].Type)
	assert.Equal(t, "USER", result.Steps[7].Type)
	assert.Equal(t, []string{"user"}, *result.Steps[7].Args)
}

func TestConvert_PythonAppWithENV(t *testing.T) {
	// Python 应用，包含 ENV、ARG、COPY、ENTRYPOINT
	content := `FROM python:3.11-slim
ARG APP_VERSION=1.0.0
ENV PYTHONUNBUFFERED=1 APP_ENV="production"
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
ENTRYPOINT ["python", "-m", "uvicorn", "main:app", "--host", "0.0.0.0"]
`
	result, err := Convert(content)
	require.NoError(t, err)

	assert.Equal(t, "python:3.11-slim", result.BaseImage)
	assert.Equal(t, "python -m uvicorn main:app --host 0.0.0.0", result.StartCmd)
	assert.Equal(t, "sleep 20", result.ReadyCmd)

	// 步骤：USER root, WORKDIR /, ARG→ENV, ENV, WORKDIR /app, COPY, RUN, COPY, USER user
	assert.Len(t, result.Steps, 9)

	// ARG 转换为 ENV
	assert.Equal(t, "ENV", result.Steps[2].Type)
	assert.Equal(t, []string{"APP_VERSION", "1.0.0"}, *result.Steps[2].Args)

	// ENV 多键值对
	assert.Equal(t, "ENV", result.Steps[3].Type)
	assert.Equal(t, []string{"PYTHONUNBUFFERED", "1", "APP_ENV", "production"}, *result.Steps[3].Args)
}

func TestConvert_GoAppWithUserAndWorkdir(t *testing.T) {
	// Go 应用，显式设置 USER 和 WORKDIR，不应追加默认值
	content := `FROM golang:1.24-alpine
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/server .
USER appuser
WORKDIR /app
CMD ["/app/server"]
`
	result, err := Convert(content)
	require.NoError(t, err)

	assert.Equal(t, "golang:1.24-alpine", result.BaseImage)
	assert.Equal(t, "/app/server", result.StartCmd)

	// 最后两步应该是用户显式设置的 USER 和 WORKDIR，不追加默认值
	lastStep := result.Steps[len(result.Steps)-1]
	assert.Equal(t, "WORKDIR", lastStep.Type)
	assert.Equal(t, []string{"/app"}, *lastStep.Args)

	secondLast := result.Steps[len(result.Steps)-2]
	assert.Equal(t, "USER", secondLast.Type)
	assert.Equal(t, []string{"appuser"}, *secondLast.Args)
}

func TestConvert_DockerfileWithCOPYChownAndADD(t *testing.T) {
	// 混合 COPY（带 --chown）和 ADD 指令
	content := `FROM ubuntu:22.04
RUN useradd -m appuser
COPY --chown=appuser:appuser config/ /etc/app/
ADD scripts.tar.gz /opt/scripts/
USER appuser
WORKDIR /home/appuser
CMD ["bash"]
`
	result, err := Convert(content)
	require.NoError(t, err)

	assert.Equal(t, "ubuntu:22.04", result.BaseImage)

	// COPY 带 --chown
	copyStep := result.Steps[3]
	assert.Equal(t, "COPY", copyStep.Type)
	assert.Equal(t, []string{"config/", "/etc/app/", "appuser", ""}, *copyStep.Args)

	// ADD 转换为 COPY
	addStep := result.Steps[4]
	assert.Equal(t, "COPY", addStep.Type)
	assert.Equal(t, []string{"scripts.tar.gz", "/opt/scripts/", "", ""}, *addStep.Args)
}

func TestConvert_DockerfileOnlyFROMAndRUN(t *testing.T) {
	// 最简 Dockerfile：仅 FROM + RUN，无 CMD
	content := `FROM alpine:3.19
RUN apk add --no-cache curl wget
`
	result, err := Convert(content)
	require.NoError(t, err)

	assert.Equal(t, "alpine:3.19", result.BaseImage)
	assert.Equal(t, "", result.StartCmd)
	assert.Equal(t, "", result.ReadyCmd)

	// 步骤：USER root, WORKDIR /, RUN, USER user, WORKDIR /home/user
	assert.Len(t, result.Steps, 5)
}

func TestConvert_DockerfileWithMultipleRUN(t *testing.T) {
	// 多条 RUN 指令
	content := `FROM ubuntu:22.04
RUN apt-get update
RUN apt-get install -y curl
RUN apt-get clean && rm -rf /var/lib/apt/lists/*
CMD ["bash"]
`
	result, err := Convert(content)
	require.NoError(t, err)

	assert.Equal(t, "ubuntu:22.04", result.BaseImage)
	assert.Equal(t, "bash", result.StartCmd)

	// 步骤：USER root, WORKDIR /, RUN×3, USER user, WORKDIR /home/user
	assert.Len(t, result.Steps, 7)
	assert.Equal(t, "RUN", result.Steps[2].Type)
	assert.Equal(t, []string{"apt-get update"}, *result.Steps[2].Args)
	assert.Equal(t, "RUN", result.Steps[3].Type)
	assert.Equal(t, []string{"apt-get install -y curl"}, *result.Steps[3].Args)
	assert.Equal(t, "RUN", result.Steps[4].Type)
	assert.Contains(t, (*result.Steps[4].Args)[0], "apt-get clean")
}

func TestConvert_DockerfileWithHeredocAndENV(t *testing.T) {
	// Heredoc RUN + ENV 组合
	content := `FROM python:3.11
ENV DEBIAN_FRONTEND=noninteractive
RUN <<EOF
apt-get update
apt-get install -y --no-install-recommends gcc
apt-get clean
rm -rf /var/lib/apt/lists/*
EOF
COPY . /app
WORKDIR /app
CMD ["python", "main.py"]
`
	result, err := Convert(content)
	require.NoError(t, err)

	assert.Equal(t, "python:3.11", result.BaseImage)
	assert.Equal(t, "python main.py", result.StartCmd)

	// ENV 步骤
	envStep := result.Steps[2]
	assert.Equal(t, "ENV", envStep.Type)
	assert.Equal(t, []string{"DEBIAN_FRONTEND", "noninteractive"}, *envStep.Args)

	// Heredoc RUN 步骤，内容包含多行命令
	runStep := result.Steps[3]
	assert.Equal(t, "RUN", runStep.Type)
	assert.Contains(t, (*runStep.Args)[0], "apt-get update")
	assert.Contains(t, (*runStep.Args)[0], "apt-get install")
	assert.Contains(t, (*runStep.Args)[0], "apt-get clean")
}

func TestConvert_IgnoredInstructionsInFullDockerfile(t *testing.T) {
	// 包含 EXPOSE、VOLUME、LABEL 等被忽略指令的完整 Dockerfile
	content := `FROM nginx:alpine
LABEL maintainer="test@example.com"
EXPOSE 80 443
VOLUME ["/var/cache/nginx"]
COPY nginx.conf /etc/nginx/nginx.conf
COPY html/ /usr/share/nginx/html/
HEALTHCHECK --interval=30s --timeout=3s CMD curl -f http://localhost/ || exit 1
CMD ["nginx", "-g", "daemon off;"]
`
	result, err := Convert(content)
	require.NoError(t, err)

	assert.Equal(t, "nginx:alpine", result.BaseImage)
	assert.Equal(t, "nginx -g 'daemon off;'", result.StartCmd)

	// EXPOSE、VOLUME、LABEL、HEALTHCHECK 被忽略，不生成步骤
	// 步骤：USER root, WORKDIR /, COPY×2, USER user, WORKDIR /home/user
	assert.Len(t, result.Steps, 6)
	for _, step := range result.Steps {
		assert.NotEqual(t, "EXPOSE", step.Type)
		assert.NotEqual(t, "VOLUME", step.Type)
		assert.NotEqual(t, "LABEL", step.Type)
		assert.NotEqual(t, "HEALTHCHECK", step.Type)
	}
}

func TestConvert_CMDOverridesENTRYPOINT(t *testing.T) {
	// CMD 在 ENTRYPOINT 之后出现，应以最后一个为准
	content := `FROM alpine
ENTRYPOINT ["python3"]
CMD ["-m", "flask", "run"]
`
	result, err := Convert(content)
	require.NoError(t, err)

	// 最后出现的是 CMD，所以 StartCmd 应该是 CMD 的值
	assert.Equal(t, "-m flask run", result.StartCmd)
	assert.Equal(t, "sleep 20", result.ReadyCmd)
}

func TestConvert_ENTRYPOINTOverridesCMD(t *testing.T) {
	// ENTRYPOINT 在 CMD 之后出现
	content := `FROM alpine
CMD ["default-cmd"]
ENTRYPOINT ["/entrypoint.sh"]
`
	result, err := Convert(content)
	require.NoError(t, err)

	// 最后出现的是 ENTRYPOINT，所以 StartCmd 应该是 ENTRYPOINT 的值
	assert.Equal(t, "/entrypoint.sh", result.StartCmd)
	assert.Equal(t, "sleep 20", result.ReadyCmd)
}
