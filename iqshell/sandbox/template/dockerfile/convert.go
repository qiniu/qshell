// Package dockerfile 提供将通用 Dockerfile 解析结果转换为
// sandbox.TemplateStep 格式的适配器。
package dockerfile

import (
	"fmt"
	"strings"

	"github.com/qiniu/go-sdk/v7/sandbox"

	dfparser "github.com/qiniu/qshell/v2/iqshell/common/dockerfile"
)

// ConvertResult 保存 Dockerfile 转换为 sandbox 格式的输出。
type ConvertResult struct {
	// BaseImage 是 FROM 指令中的镜像名。
	BaseImage string

	// Steps 是从 Dockerfile 指令转换后的构建步骤列表。
	Steps []sandbox.TemplateStep

	// StartCmd 是从 CMD 或 ENTRYPOINT 指令中提取的启动命令。
	StartCmd string

	// ReadyCmd 在存在 CMD/ENTRYPOINT 时默认为 "sleep 20"（匹配 e2b 行为）。
	ReadyCmd string

	// Warnings 包含非致命的解析问题。
	Warnings []string
}

// Convert 解析 Dockerfile 内容并将指令转换为 sandbox TemplateStep。
// 匹配 e2b v2 构建系统行为：
//   - 在前面插入 USER root 和 WORKDIR / 步骤
//   - 如果未找到 USER 指令，在末尾追加 USER user
//   - 如果未找到 WORKDIR 指令，在末尾追加 WORKDIR /home/user
func Convert(content string) (*ConvertResult, error) {
	parsed, err := dfparser.Parse(content)
	if err != nil {
		return nil, err
	}

	result := &ConvertResult{
		Warnings: parsed.Warnings,
	}
	var steps []sandbox.TemplateStep
	hasUser := false
	hasWorkdir := false
	fromCount := 0

	// 插入默认步骤（匹配 e2b 行为）
	steps = append(steps, makeStep("USER", "root"))
	steps = append(steps, makeStep("WORKDIR", "/"))

	for _, inst := range parsed.Instructions {
		switch inst.Name {
		case "FROM":
			fromCount++
			if fromCount > 1 {
				result.Warnings = append(result.Warnings,
					"multi-stage build detected; using the last FROM stage as the runtime base image")
			}
			// 提取镜像名（第一个非 AS 的 token），多阶段构建取最后一个 FROM
			result.BaseImage = extractImage(inst.Args)

		case "RUN":
			if inst.Args == "" {
				return nil, fmt.Errorf("line %d: empty RUN instruction", inst.Line)
			}
			args := []string{inst.Args}
			steps = append(steps, sandbox.TemplateStep{
				Type: "RUN",
				Args: &args,
			})

		case "COPY", "ADD":
			user, src, dest, err := parseCopyArgs(inst.Args, inst.Flags)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid %s instruction: %w", inst.Line, inst.Name, err)
			}
			args := []string{src, dest, user, ""}
			steps = append(steps, sandbox.TemplateStep{
				Type: "COPY",
				Args: &args,
			})

		case "WORKDIR":
			if inst.Args == "" {
				return nil, fmt.Errorf("line %d: empty WORKDIR instruction", inst.Line)
			}
			hasWorkdir = true
			args := []string{inst.Args}
			steps = append(steps, sandbox.TemplateStep{
				Type: "WORKDIR",
				Args: &args,
			})

		case "USER":
			if inst.Args == "" {
				return nil, fmt.Errorf("line %d: empty USER instruction", inst.Line)
			}
			hasUser = true
			args := []string{inst.Args}
			steps = append(steps, sandbox.TemplateStep{
				Type: "USER",
				Args: &args,
			})

		case "ENV":
			envArgs, err := dfparser.ParseEnvValues(inst.Args, parsed.EscapeToken)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid ENV instruction: %w", inst.Line, err)
			}
			steps = append(steps, sandbox.TemplateStep{
				Type: "ENV",
				Args: &envArgs,
			})

		case "ARG":
			argArgs, hasDefault := parseArgValues(inst.Args)
			if hasDefault {
				steps = append(steps, sandbox.TemplateStep{
					Type: "ENV",
					Args: &argArgs,
				})
			}
			// 无默认值的 ARG 仅为构建时变量，不生成 ENV 步骤

		case "CMD":
			result.StartCmd = dfparser.ParseCommand(inst.Args)
			result.ReadyCmd = "sleep 20"

		case "ENTRYPOINT":
			result.StartCmd = dfparser.ParseCommand(inst.Args)
			result.ReadyCmd = "sleep 20"
		}
	}

	// 如果未显式设置，追加默认值（匹配 e2b 行为）
	if !hasUser {
		steps = append(steps, makeStep("USER", "user"))
	}
	if !hasWorkdir {
		steps = append(steps, makeStep("WORKDIR", "/home/user"))
	}

	if result.BaseImage == "" {
		return nil, fmt.Errorf("no FROM instruction found in Dockerfile")
	}

	result.Steps = steps
	return result, nil
}

// extractImage 从 FROM 参数中提取镜像名，忽略 AS 别名。
func extractImage(args string) string {
	for f := range strings.FieldsSeq(args) {
		if strings.ToUpper(f) == "AS" {
			break
		}
		return f
	}
	return args
}

// parseCopyArgs 从 COPY/ADD 指令中提取 user、src 和 dest。
func parseCopyArgs(args string, flags map[string]string) (user, src, dest string, err error) {
	// 从 --chown 标志中提取用户
	if chown, ok := flags["chown"]; ok {
		if u, _, found := strings.Cut(chown, ":"); found {
			user = u
		} else {
			user = chown
		}
	}

	// 去除 heredoc 标记
	args = dfparser.StripHeredocMarkers(args)

	fields := strings.Fields(args)
	if len(fields) < 2 {
		return "", "", "", fmt.Errorf("COPY/ADD requires at least source and destination")
	}

	dest = fields[len(fields)-1]
	src = strings.Join(fields[:len(fields)-1], " ")
	return user, src, dest, nil
}

// parseArgValues 将 ARG name[=default_value] 解析为 ["name", "value"]。
// 返回解析结果和是否包含默认值（即是否有 = 号）。
func parseArgValues(rest string) ([]string, bool) {
	key, value, hasDefault := strings.Cut(rest, "=")
	return []string{strings.TrimSpace(key), strings.TrimSpace(value)}, hasDefault
}

// makeStep 创建一个简单的 TemplateStep。
func makeStep(typ string, args ...string) sandbox.TemplateStep {
	a := make([]string, len(args))
	copy(a, args)
	return sandbox.TemplateStep{
		Type: typ,
		Args: &a,
	}
}
