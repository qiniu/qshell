package operations

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/qiniu/go-sdk/v7/sandbox"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// SnapshotCreateInfo 保存创建持久快照的参数。
type SnapshotCreateInfo struct {
	SandboxID      string
	Name           string
	Memory         bool
	IdempotencyKey string
	Format         string
}

// SnapshotCreate 创建持久快照，调用方必须保存幂等键以便在结果不确定时继续重试。
func SnapshotCreate(ctx context.Context, info SnapshotCreateInfo) error {
	if info.IdempotencyKey == "" {
		return fmt.Errorf("--idempotency-key is required; reuse the same key and request body after an uncertain outcome")
	}
	if info.Format != sbClient.FormatJSON && info.Format != sbClient.FormatPretty {
		return fmt.Errorf("--format must be pretty or json")
	}
	client, err := sbClient.NewSandboxClient()
	if err != nil {
		return err
	}
	params := sandbox.CreateSnapshotParams{Memory: &info.Memory, IdempotencyKey: info.IdempotencyKey}
	if info.Name != "" {
		params.Name = &info.Name
	}
	snapshot, err := client.CreateSnapshot(ctx, info.SandboxID, params)
	if err != nil {
		return fmt.Errorf("create snapshot failed (idempotency key %s; reuse this key and request body when retrying): %w", info.IdempotencyKey, err)
	}
	if info.Format == sbClient.FormatJSON {
		sbClient.PrintJSON(snapshot)
		return nil
	}
	sbClient.PrintSuccess("Snapshot created: %s", snapshot.SnapshotID)
	fmt.Printf("Names: %s\n", strings.Join(snapshot.Names, ", "))
	return nil
}

// SnapshotListInfo 保存快照过滤与分页参数。
type SnapshotListInfo struct {
	SandboxID string
	Name      string
	Limit     int32
	NextToken string
	Format    string
}

// SnapshotList 查询一页快照，并返回下一页游标。
func SnapshotList(ctx context.Context, info SnapshotListInfo) error {
	if info.Limit < 0 {
		return fmt.Errorf("--limit must be non-negative")
	}
	if info.Format != sbClient.FormatJSON && info.Format != sbClient.FormatPretty {
		return fmt.Errorf("--format must be pretty or json")
	}
	client, err := sbClient.NewSandboxClient()
	if err != nil {
		return err
	}
	params := &sandbox.ListSnapshotsParams{}
	if info.SandboxID != "" {
		params.SandboxID = &info.SandboxID
	}
	if info.Name != "" {
		params.Name = &info.Name
	}
	if info.Limit > 0 {
		params.Limit = &info.Limit
	}
	if info.NextToken != "" {
		params.NextToken = &info.NextToken
	}
	page, err := client.ListSnapshots(ctx, params)
	if err != nil {
		return fmt.Errorf("list snapshots failed: %w", err)
	}
	if info.Format == sbClient.FormatJSON {
		sbClient.PrintJSON(page)
		return nil
	}
	table := sbClient.NewTable(os.Stdout)
	fmt.Fprintln(table, "SNAPSHOT ID\tNAMES")
	for _, item := range page.Items {
		fmt.Fprintf(table, "%s\t%s\n", item.SnapshotID, strings.Join(item.Names, ", "))
	}
	table.Flush()
	if page.NextToken != "" {
		fmt.Printf("Next token: %s\n", page.NextToken)
	}
	return nil
}
