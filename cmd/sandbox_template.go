package cmd

import (
	"github.com/spf13/cobra"

	"github.com/qiniu/qshell/v2/docs"
	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/sandbox/template/operations"
)

var templateCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "template",
		Aliases: []string{"tpl"},
		Short:   "Manage sandbox templates",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTemplateType
			docs.ShowCmdDocument(docs.SandboxTemplateType)
		},
	}
	return cmd
}

var templateListCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.ListInfo{}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List sandbox templates",
		Example: `qshell sandbox template list --format json`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTemplateListType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			operations.List(info)
		},
	}
	cmd.Flags().StringVar(&info.Format, "format", "pretty", "output format: pretty or json")
	return cmd
}

var templateGetCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <templateID>",
		Aliases: []string{"gt"},
		Short:   "Get template details",
		Example: `qshell sandbox template get tmpl-xxxxxxxxxxxx`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTemplateGetType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			if len(args) != 1 {
				_ = cmd.Usage()
				return
			}
			operations.Get(operations.GetInfo{
				TemplateID: args[0],
			})
		},
	}
	return cmd
}

var templateDeleteCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.DeleteInfo{}
	cmd := &cobra.Command{
		Use:     "delete [templateIDs...]",
		Aliases: []string{"dl"},
		Short:   "Delete one or more templates",
		Example: `  qshell sandbox template delete tmpl-xxxxxxxxxxxx -y
  qshell sandbox template delete tmpl-aaa tmpl-bbb -y
  qshell sandbox template delete -s`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTemplateDeleteType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			if len(args) == 0 && !info.Select {
				_ = cmd.Usage()
				return
			}
			info.TemplateIDs = args
			operations.Delete(info)
		},
	}
	cmd.Flags().BoolVarP(&info.Yes, "yes", "y", false, "skip confirmation")
	cmd.Flags().BoolVarP(&info.Select, "select", "s", false, "interactively select templates to delete")
	return cmd
}

var templateBuildsCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "builds <templateID> <buildID>",
		Aliases: []string{"bds"},
		Short:   "View template build status",
		Example: `qshell sandbox template builds tmpl-xxxxxxxxxxxx build-xxxxxxxxxxxx`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTemplateBuildsType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			if len(args) != 2 {
				_ = cmd.Usage()
				return
			}
			operations.Builds(operations.BuildsInfo{
				TemplateID: args[0],
				BuildID:    args[1],
			})
		},
	}
	return cmd
}

var templateBuildCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.BuildInfo{}
	cmd := &cobra.Command{
		Use:   "build",
		Aliases: []string{"bd"},
		Short: "Build a template",
		Long:  "Create a new template and build it, or rebuild an existing template.",
		Example: `  # Create and build a new template from a Docker image
  qshell sandbox template build --name my-template --from-image ubuntu:22.04 --wait

  # Rebuild an existing template
  qshell sandbox template build --template-id tmpl-xxxxxxxxxxxx --from-image ubuntu:22.04

  # Force rebuild without cache
  qshell sandbox template build --template-id tmpl-xxxxxxxxxxxx --no-cache --wait`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTemplateBuildType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			operations.Build(info)
		},
	}
	cmd.Flags().StringVar(&info.Name, "name", "", "template name (for creating a new template)")
	cmd.Flags().StringVar(&info.TemplateID, "template-id", "", "existing template ID (for rebuilding)")
	cmd.Flags().StringVar(&info.FromImage, "from-image", "", "base Docker image")
	cmd.Flags().StringVar(&info.FromTemplate, "from-template", "", "base template")
	cmd.Flags().StringVar(&info.StartCmd, "start-cmd", "", "command to run after build")
	cmd.Flags().StringVar(&info.ReadyCmd, "ready-cmd", "", "readiness check command")
	cmd.Flags().Int32Var(&info.CPUCount, "cpu", 0, "sandbox CPU count")
	cmd.Flags().Int32Var(&info.MemoryMB, "memory", 0, "sandbox memory size in MiB")
	cmd.Flags().BoolVar(&info.Wait, "wait", false, "wait for build to complete")
	cmd.Flags().BoolVar(&info.NoCache, "no-cache", false, "force full rebuild ignoring cache")
	return cmd
}

var templatePublishCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.PublishInfo{Public: true}
	cmd := &cobra.Command{
		Use:     "publish [templateIDs...]",
		Aliases: []string{"pb"},
		Short:   "Publish templates (make public)",
		Example: `  qshell sandbox template publish tmpl-xxxxxxxxxxxx -y
  qshell sandbox template publish -s`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTemplatePublishType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			if len(args) == 0 && !info.Select {
				_ = cmd.Usage()
				return
			}
			info.TemplateIDs = args
			operations.Publish(info)
		},
	}
	cmd.Flags().BoolVarP(&info.Yes, "yes", "y", false, "skip confirmation")
	cmd.Flags().BoolVarP(&info.Select, "select", "s", false, "interactively select templates")
	return cmd
}

var templateUnpublishCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.PublishInfo{Public: false}
	cmd := &cobra.Command{
		Use:     "unpublish [templateIDs...]",
		Aliases: []string{"upb"},
		Short:   "Unpublish templates (make private)",
		Example: `  qshell sandbox template unpublish tmpl-xxxxxxxxxxxx -y
  qshell sandbox template unpublish -s`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTemplateUnpublishType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			if len(args) == 0 && !info.Select {
				_ = cmd.Usage()
				return
			}
			info.TemplateIDs = args
			operations.Publish(info)
		},
	}
	cmd.Flags().BoolVarP(&info.Yes, "yes", "y", false, "skip confirmation")
	cmd.Flags().BoolVarP(&info.Select, "select", "s", false, "interactively select templates")
	return cmd
}

var templateInitCmdBuilder = func(cfg *iqshell.Config) *cobra.Command {
	info := operations.InitInfo{}
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"it"},
		Short:   "Initialize a new template project",
		Long:    "Scaffold a new template project with boilerplate files for the selected language.",
		Example: `  # Interactive mode
  qshell sandbox template init

  # Non-interactive mode
  qshell sandbox template init --name my-template --language go
  qshell sandbox template init --name my-api --language typescript --path ./my-api`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.CmdCfg.CmdId = docs.SandboxTemplateInitType
			if iqshell.ShowDocumentIfNeeded(cfg) {
				return
			}
			operations.Init(info)
		},
	}
	cmd.Flags().StringVar(&info.Name, "name", "", "template project name")
	cmd.Flags().StringVar(&info.Language, "language", "", "programming language (go, typescript, python)")
	cmd.Flags().StringVar(&info.Path, "path", "", "output directory (defaults to ./<name>)")
	return cmd
}

// templateCmdLoader adds the template command and its subcommands to the given parent command.
func templateCmdLoader(parentCmd *cobra.Command, cfg *iqshell.Config) {
	templateCmd := templateCmdBuilder(cfg)
	templateCmd.AddCommand(
		templateListCmdBuilder(cfg),
		templateGetCmdBuilder(cfg),
		templateDeleteCmdBuilder(cfg),
		templateBuildCmdBuilder(cfg),
		templateBuildsCmdBuilder(cfg),
		templatePublishCmdBuilder(cfg),
		templateUnpublishCmdBuilder(cfg),
		templateInitCmdBuilder(cfg),
	)
	parentCmd.AddCommand(templateCmd)
}
