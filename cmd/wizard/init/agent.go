package init

import "charm.land/bubbles/v2/list"

func NewAgentList(m *model, styles *styles, delegateKeys *delegateKeyMap) list.Model {
	// Make initial agentList of items.
	const numItems = 24
	items := make([]list.Item, numItems)
	items = []list.Item{
		item{
			title:       "Claude Code",
			description: "by Anthropic",
		},
		item{
			title:       "Codex",
			description: "by OpenAI",
		},
		item{
			title:       "Opencode",
			description: "by Opencode",
		},
	}

	// Setup agentList.
	delegate := newItemDelegate(delegateKeys, styles, func(i item) {
		m.selectedAgent = i.title
	})
	agentList := list.New(items, delegate, 0, 0)
	agentList.Title = "Agent"
	agentList.Styles.Title = styles.title
	return agentList
}
