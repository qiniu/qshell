package sandbox

import (
	"context"
	"fmt"

	"github.com/qiniu/go-sdk/v7/sandbox"
)

type listPageTokenKey struct{}

type listPageToken struct {
	value string
}

// SandboxListClient 通过已发布 SDK 查询沙箱，并保留分页响应头。
type SandboxListClient struct {
	*sandbox.Client
}

// NewSandboxListClient 创建支持游标和完整枚举的沙箱客户端。
func NewSandboxListClient() (*SandboxListClient, error) {
	client, err := NewSandboxClient()
	if err != nil {
		return nil, err
	}
	return &SandboxListClient{Client: client}, nil
}

// ListPage 查询一页沙箱，空游标表示已到末页。
func (c *SandboxListClient) ListPage(ctx context.Context, params *sandbox.ListParams) ([]sandbox.ListedSandbox, string, error) {
	token := &listPageToken{}
	items, err := c.Client.List(context.WithValue(ctx, listPageTokenKey{}, token), params)
	return items, token.value, err
}

// ListAll 完整查询后返回去重的沙箱集合，任一页失败时不交付部分结果。
func (c *SandboxListClient) ListAll(ctx context.Context, params *sandbox.ListParams) ([]sandbox.ListedSandbox, error) {
	var query sandbox.ListParams
	if params != nil {
		query = *params
	}
	items := make([]sandbox.ListedSandbox, 0)
	seenIDs := make(map[string]struct{})
	seenTokens := make(map[string]struct{})
	if query.NextToken != nil && *query.NextToken != "" {
		seenTokens[*query.NextToken] = struct{}{}
	}
	for {
		page, nextToken, err := c.ListPage(ctx, &query)
		if err != nil {
			return nil, err
		}
		for _, item := range page {
			if _, exists := seenIDs[item.SandboxID]; exists {
				continue
			}
			seenIDs[item.SandboxID] = struct{}{}
			items = append(items, item)
		}
		if nextToken == "" {
			return items, nil
		}
		if _, repeated := seenTokens[nextToken]; repeated {
			return nil, fmt.Errorf("sandbox list returned a repeated pagination token")
		}
		seenTokens[nextToken] = struct{}{}
		query.NextToken = &nextToken
	}
}
