package operations

import (
	"context"
	"fmt"
	"os"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// TagsInfo 保存模板标签查询参数。
type TagsInfo struct {
	TemplateID string
	Format     string
}

// Tags 查询模板标签的最新构建分配。
func Tags(ctx context.Context, info TagsInfo) error {
	if info.Format != sbClient.FormatJSON && info.Format != sbClient.FormatPretty {
		return fmt.Errorf("--format must be pretty or json")
	}
	client, err := sbClient.NewSandboxClient()
	if err != nil {
		return err
	}
	tags, err := client.ListTemplateTags(ctx, info.TemplateID)
	if err != nil {
		return fmt.Errorf("list template tags failed: %w", err)
	}
	if info.Format == sbClient.FormatJSON {
		sbClient.PrintJSON(tags)
		return nil
	}
	table := sbClient.NewTable(os.Stdout)
	fmt.Fprintln(table, "TAG\tBUILD ID\tCREATED AT")
	for _, tag := range tags {
		fmt.Fprintf(table, "%s\t%s\t%s\n", tag.Tag, tag.BuildID, sbClient.FormatTimestamp(tag.CreatedAt))
	}
	table.Flush()
	return nil
}
