package dockerfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// 基础解析
// ---------------------------------------------------------------------------

func TestParse_BasicDockerfile(t *testing.T) {
	content := `FROM ubuntu:22.04
RUN apt-get update
WORKDIR /app
COPY . /app/
CMD ["python3", "app.py"]
`
	result, err := Parse(content)
	require.NoError(t, err)
	assert.Len(t, result.Instructions, 5)

	assert.Equal(t, "FROM", result.Instructions[0].Name)
	assert.Equal(t, "ubuntu:22.04", result.Instructions[0].Args)
	assert.Equal(t, "RUN", result.Instructions[1].Name)
	assert.Equal(t, "apt-get update", result.Instructions[1].Args)
	assert.Equal(t, "WORKDIR", result.Instructions[2].Name)
	assert.Equal(t, "CMD", result.Instructions[4].Name)
}

func TestParse_NoInstructions(t *testing.T) {
	result, err := Parse("# just a comment\n\n")
	require.NoError(t, err)
	assert.Empty(t, result.Instructions)
}

// ---------------------------------------------------------------------------
// 转义指令
// ---------------------------------------------------------------------------

func TestParse_EscapeDirectiveBacktick(t *testing.T) {
	content := "# escape=`\nFROM alpine\nRUN echo hello `\n    && echo world\n"
	result, err := Parse(content)
	require.NoError(t, err)
	assert.Equal(t, '`', result.EscapeToken)

	run := findInstruction(result, "RUN")
	require.NotNil(t, run)
	assert.Contains(t, run.Args, "echo hello")
	assert.Contains(t, run.Args, "echo world")
}

func TestParse_EscapeDirectiveBackslash(t *testing.T) {
	content := "# escape=\\\nFROM alpine\nRUN echo hello \\\n    && echo world\n"
	result, err := Parse(content)
	require.NoError(t, err)
	assert.Equal(t, '\\', result.EscapeToken)

	run := findInstruction(result, "RUN")
	require.NotNil(t, run)
	assert.Contains(t, run.Args, "echo hello")
	assert.Contains(t, run.Args, "echo world")
}

func TestParse_EscapeBacktickNoJoinOnBackslash(t *testing.T) {
	content := "# escape=`\nFROM alpine\nRUN echo C:\\Users\\test\n"
	result, err := Parse(content)
	require.NoError(t, err)

	run := findInstruction(result, "RUN")
	require.NotNil(t, run)
	assert.Equal(t, `echo C:\Users\test`, run.Args)
}

func TestParse_EscapeDirectiveMustBeFirst(t *testing.T) {
	content := "# this is a comment\n# escape=`\nFROM alpine\nRUN echo hello `\n    && echo world\n"
	result, err := Parse(content)
	require.NoError(t, err)
	// escape= 在非指令注释之后，因此使用默认的 \ 转义符
	assert.Equal(t, '\\', result.EscapeToken)

	run := findInstruction(result, "RUN")
	require.NotNil(t, run)
	assert.Equal(t, "echo hello `", run.Args)
}

// ---------------------------------------------------------------------------
// 行续接
// ---------------------------------------------------------------------------

func TestParse_ContinuationLines(t *testing.T) {
	content := `FROM alpine
RUN apt-get update && \
    apt-get install -y curl && \
    apt-get clean
`
	result, err := Parse(content)
	require.NoError(t, err)

	run := findInstruction(result, "RUN")
	require.NotNil(t, run)
	assert.Contains(t, run.Args, "apt-get update")
	assert.Contains(t, run.Args, "apt-get install -y curl")
	assert.Contains(t, run.Args, "apt-get clean")
}

func TestParse_ContinuationSkipsComments(t *testing.T) {
	content := `FROM alpine
RUN echo hello \
    # this is a comment inside continuation \
    && echo world
`
	result, err := Parse(content)
	require.NoError(t, err)

	run := findInstruction(result, "RUN")
	require.NotNil(t, run)
	assert.Contains(t, run.Args, "echo hello")
	assert.Contains(t, run.Args, "echo world")
	assert.NotContains(t, run.Args, "this is a comment")
}

// ---------------------------------------------------------------------------
// Heredoc
// ---------------------------------------------------------------------------

func TestParse_HeredocRUN(t *testing.T) {
	content := `FROM alpine
RUN <<EOF
echo hello
echo world
EOF
`
	result, err := Parse(content)
	require.NoError(t, err)

	run := findInstruction(result, "RUN")
	require.NotNil(t, run)
	assert.Equal(t, "echo hello\necho world", run.Args)
	assert.Equal(t, "echo hello\necho world", run.Heredoc)
}

func TestParse_HeredocUnterminatedError(t *testing.T) {
	content := `FROM alpine
RUN <<EOF
echo hello
`
	_, err := Parse(content)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unterminated heredoc")
}

func TestParse_HeredocCustomMarker(t *testing.T) {
	content := `FROM alpine
RUN <<SCRIPT
#!/bin/bash
apt-get update
SCRIPT
`
	result, err := Parse(content)
	require.NoError(t, err)

	run := findInstruction(result, "RUN")
	require.NotNil(t, run)
	assert.Contains(t, run.Args, "apt-get update")
}

// ---------------------------------------------------------------------------
// 标志提取
// ---------------------------------------------------------------------------

func TestParse_COPY_Flags(t *testing.T) {
	content := `FROM alpine
COPY --chown=user:group --chmod=755 src/ /dest/
`
	result, err := Parse(content)
	require.NoError(t, err)

	cp := findInstruction(result, "COPY")
	require.NotNil(t, cp)
	assert.Equal(t, "user:group", cp.Flags["chown"])
	assert.Equal(t, "755", cp.Flags["chmod"])
	assert.Equal(t, "src/ /dest/", cp.Args)
}

func TestParse_FROM_Flags(t *testing.T) {
	content := `FROM --platform=linux/amd64 ubuntu:22.04 AS builder
`
	result, err := Parse(content)
	require.NoError(t, err)

	from := findInstruction(result, "FROM")
	require.NotNil(t, from)
	assert.Equal(t, "linux/amd64", from.Flags["platform"])
	assert.Equal(t, "ubuntu:22.04 AS builder", from.Args)
}

// ---------------------------------------------------------------------------
// 警告
// ---------------------------------------------------------------------------

func TestParse_UnknownInstructionWarning(t *testing.T) {
	content := `FROM alpine
FOOBAR something
RUN echo hello
`
	result, err := Parse(content)
	require.NoError(t, err)
	require.Len(t, result.Warnings, 1)
	assert.Contains(t, result.Warnings[0], "unknown instruction")
	assert.Contains(t, result.Warnings[0], "FOOBAR")
}

func TestParse_ONBUILDWarning(t *testing.T) {
	content := `FROM alpine
ONBUILD RUN echo hello
`
	result, err := Parse(content)
	require.NoError(t, err)
	require.Len(t, result.Warnings, 1)
	assert.Contains(t, result.Warnings[0], "ONBUILD")
}

// ---------------------------------------------------------------------------
// 引号感知的 ENV 解析
// ---------------------------------------------------------------------------

func TestParseEnvValues_KeyValue(t *testing.T) {
	result, err := ParseEnvValues(`FOO=bar BAZ="hello world"`, '\\')
	require.NoError(t, err)
	assert.Equal(t, []string{"FOO", "bar", "BAZ", "hello world"}, result)
}

func TestParseEnvValues_DoubleQuotedEscape(t *testing.T) {
	result, err := ParseEnvValues(`MSG="it's \"here\""`, '\\')
	require.NoError(t, err)
	assert.Equal(t, []string{"MSG", `it's "here"`}, result)
}

func TestParseEnvValues_SingleQuotedLiteral(t *testing.T) {
	result, err := ParseEnvValues(`MSG='hello\nworld'`, '\\')
	require.NoError(t, err)
	assert.Equal(t, []string{"MSG", `hello\nworld`}, result)
}

func TestParseEnvValues_Legacy(t *testing.T) {
	result, err := ParseEnvValues("MY_VAR some value with spaces", '\\')
	require.NoError(t, err)
	assert.Equal(t, []string{"MY_VAR", "some value with spaces"}, result)
}

func TestParseEnvValues_Empty(t *testing.T) {
	result, err := ParseEnvValues(`EMPTY=""`, '\\')
	require.NoError(t, err)
	assert.Equal(t, []string{"EMPTY", ""}, result)
}

func TestParseEnvValues_Mixed(t *testing.T) {
	result, err := ParseEnvValues(`A=1 B="two" C='three'`, '\\')
	require.NoError(t, err)
	assert.Equal(t, []string{"A", "1", "B", "two", "C", "three"}, result)
}

// ---------------------------------------------------------------------------
// ParseCommand
// ---------------------------------------------------------------------------

func TestParseCommand_ExecForm(t *testing.T) {
	assert.Equal(t, "node server.js --port 3000", ParseCommand(`["node", "server.js", "--port", "3000"]`))
}

func TestParseCommand_ShellForm(t *testing.T) {
	assert.Equal(t, "echo hello world", ParseCommand("echo hello world"))
}

func TestParseCommand_Empty(t *testing.T) {
	assert.Equal(t, "", ParseCommand(""))
}

// ---------------------------------------------------------------------------
// BOM 和 Windows 换行符
// ---------------------------------------------------------------------------

func TestParse_BOMStripping(t *testing.T) {
	content := "\xef\xbb\xbfFROM alpine\nRUN echo hello\n"
	result, err := Parse(content)
	require.NoError(t, err)
	assert.Equal(t, "FROM", result.Instructions[0].Name)
}

func TestParse_WindowsLineEndings(t *testing.T) {
	content := "FROM alpine\r\nRUN echo hello\r\nWORKDIR /app\r\n"
	result, err := Parse(content)
	require.NoError(t, err)
	assert.Len(t, result.Instructions, 3)
	assert.Equal(t, "WORKDIR", result.Instructions[2].Name)
	assert.Equal(t, "/app", result.Instructions[2].Args)
}

// ---------------------------------------------------------------------------
// StripHeredocMarkers
// ---------------------------------------------------------------------------

func TestStripHeredocMarkers(t *testing.T) {
	assert.Equal(t, "src/ /dest/", StripHeredocMarkers("<<FILE src/ /dest/"))
	assert.Equal(t, "src/ /dest/", StripHeredocMarkers("src/ /dest/"))
}

// ---------------------------------------------------------------------------
// 辅助函数
// ---------------------------------------------------------------------------

func findInstruction(r *ParseResult, name string) *Instruction {
	for i := range r.Instructions {
		if r.Instructions[i].Name == name {
			return &r.Instructions[i]
		}
	}
	return nil
}
