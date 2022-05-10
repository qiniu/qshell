package provider

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"sync"
	"time"
)

type Provider interface {
	Available() (available bool, err *data.CodeError)
	Provide() (item Item, err *data.CodeError)
	Freeze(item Item)
}

func NewListProvider(items []Item) Provider {
	freezeItems := make([]*freezeItem, 0, len(items))
	for _, i := range items {
		freezeItems = append(freezeItems, &freezeItem{
			item:          i,
			availableTime: nil,
		})
	}
	return &listProvider{
		mu:    sync.Mutex{},
		items: freezeItems,
	}
}

type listProvider struct {
	mu    sync.Mutex
	items []*freezeItem
}

func (l *listProvider) Available() (available bool, err *data.CodeError) {
	if l == nil || len(l.items) == 0 {
		return false, data.NewEmptyError().AppendDesc("no item found")
	}
	i, e := l.Provide()
	return i != nil, e
}

func (l *listProvider) Provide() (item Item, err *data.CodeError) {
	l.mu.Lock()
	for _, i := range l.items {
		if i.Available() {
			item = i.item
			break
		}
	}
	l.mu.Unlock()

	if item == nil {
		err = data.NewEmptyError().AppendDesc("no item found")
	}
	return
}

func (l *listProvider) Freeze(item Item) {
	if item == nil {
		return
	}

	l.mu.Lock()
	for _, i := range l.items {
		if i.item == item {
			i.Freeze()
			break
		}
	}
	l.mu.Unlock()
}

type freezeItem struct {
	item          Item
	availableTime *time.Time
}

func (i *freezeItem) Available() bool {
	if i.availableTime == nil {
		return true
	}
	return time.Now().After(*i.availableTime)
}

func (i *freezeItem) Freeze() {
	t := time.Now().Add(time.Hour * 24 * 365 * 20)
	i.availableTime = &t
}
