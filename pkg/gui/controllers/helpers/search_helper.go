package helpers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SearchHelper struct {
	c *HelperCommon
}

func NewSearchHelper(
	c *HelperCommon,
) *SearchHelper {
	return &SearchHelper{
		c: c,
	}
}

func (self *SearchHelper) OpenFilterPrompt(context types.IFilterableContext) error {
	state := self.searchState()

	state.Context = context

	self.searchPrefixView().SetContent(self.c.Tr.FilterPrefix)
	promptView := self.promptView()
	promptView.ClearTextArea()
	promptView.TextArea.TypeString(context.GetFilter())
	promptView.RenderTextArea()

	if err := self.c.PushContext(self.c.Contexts().Search); err != nil {
		return err
	}

	return nil
}

func (self *SearchHelper) OpenSearchPrompt(context types.Context) error {
	state := self.searchState()

	state.Context = context

	self.searchPrefixView().SetContent(self.c.Tr.SearchPrefix)
	promptView := self.promptView()
	// TODO: should we show the currently searched thing here? Perhaps we can store that on the context
	promptView.ClearTextArea()
	promptView.RenderTextArea()

	if err := self.c.PushContext(self.c.Contexts().Search); err != nil {
		return err
	}

	return nil
}

func (self *SearchHelper) DisplayFilterPrompt(context types.IFilterableContext) {
	state := self.searchState()

	state.Context = context
	searchString := context.GetFilter()

	self.searchPrefixView().SetContent(self.c.Tr.FilterPrefix)
	promptView := self.promptView()
	promptView.ClearTextArea()
	promptView.TextArea.TypeString(searchString)
	promptView.RenderTextArea()
}

func (self *SearchHelper) DisplaySearchPrompt(context types.ISearchableContext) {
	state := self.searchState()

	state.Context = context
	searchString := context.GetSearchString()

	self.searchPrefixView().SetContent(self.c.Tr.SearchPrefix)
	promptView := self.promptView()
	promptView.ClearTextArea()
	promptView.TextArea.TypeString(searchString)
	promptView.RenderTextArea()
}

func (self *SearchHelper) searchState() *types.SearchState {
	return self.c.State().GetRepoState().GetSearchState()
}

func (self *SearchHelper) searchPrefixView() *gocui.View {
	return self.c.Views().SearchPrefix
}

func (self *SearchHelper) promptView() *gocui.View {
	return self.c.Contexts().Search.GetView()
}

func (self *SearchHelper) promptContent() string {
	return self.c.Contexts().Search.GetView().TextArea.GetContent()
}

func (self *SearchHelper) Confirm() error {
	state := self.searchState()
	if self.promptContent() == "" {
		return self.CancelPrompt()
	}

	switch state.SearchType() {
	case types.SearchTypeFilter:
		return self.ConfirmFilter()
	case types.SearchTypeSearch:
		return self.ConfirmSearch()
	case types.SearchTypeNone:
		return self.c.PopContext()
	}

	return nil
}

func (self *SearchHelper) ConfirmFilter() error {
	// We also do this on each keypress but we do it here again just in case
	state := self.searchState()

	context, ok := state.Context.(types.IFilterableContext)
	if !ok {
		self.c.Log.Warnf("Context %s is not filterable", state.Context.GetKey())
		return nil
	}

	context.SetFilter(self.promptContent())
	_ = self.c.PostRefreshUpdate(state.Context)

	return self.c.PopContext()
}

func (self *SearchHelper) ConfirmSearch() error {
	state := self.searchState()

	if err := self.c.PopContext(); err != nil {
		return err
	}

	context, ok := state.Context.(types.ISearchableContext)
	if !ok {
		self.c.Log.Warnf("Context %s is searchable", state.Context.GetKey())
		return nil
	}

	searchString := self.promptContent()
	context.SetSearchString(searchString)

	view := context.GetView()

	if err := view.Search(searchString); err != nil {
		return err
	}

	return nil
}

func (self *SearchHelper) CancelPrompt() error {
	self.Cancel()

	return self.c.PopContext()
}

func (self *SearchHelper) Cancel() {
	state := self.searchState()

	switch context := state.Context.(type) {
	case types.IFilterableContext:
		context.SetFilter("")
		_ = self.c.PostRefreshUpdate(context)
	case types.ISearchableContext:
		context.GetView().ClearSearch()
	default:
		// do nothing
	}

	state.Context = nil
}

func (self *SearchHelper) OnPromptContentChanged(searchString string) {
	state := self.searchState()
	switch context := state.Context.(type) {
	case types.IFilterableContext:
		context.SetFilter(searchString)
		_ = self.c.PostRefreshUpdate(context)
	case types.ISearchableContext:
		// do nothing
	default:
		// do nothing (shouldn't land here)
	}
}
