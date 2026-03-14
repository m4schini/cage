package init

import "charm.land/bubbles/v2/list"

func NewBaseList(m *model, styles *styles, delegateKeys *delegateKeyMap) list.Model {
	// Make initial agentList of items.
	const numItems = 24
	items := make([]list.Item, numItems)
	items = []list.Item{
		item{
			title:       "go",
			description: "",
		},
		item{
			title:       "node",
			description: "",
		},
		item{
			title:       "dotnet",
			description: "",
		},
	}

	// Setup agentList.
	delegate := newItemDelegate(delegateKeys, styles, func(i item) {
		m.selectedBase = i.title
	})
	agentList := list.New(items, delegate, 0, 0)
	agentList.Title = "Base"
	agentList.Styles.Title = styles.title
	return agentList
}
