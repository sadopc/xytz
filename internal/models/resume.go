package models

import (
	"sort"

	"github.com/xdagiz/xytz/internal/styles"
	"github.com/xdagiz/xytz/internal/utils"

	"github.com/charmbracelet/bubbles/list"
)

type ResumeItem struct {
	URL      string
	TitleVal string
	FormatID string
}

func (i ResumeItem) Title() string       { return i.TitleVal }
func (i ResumeItem) Description() string { return i.URL }
func (i ResumeItem) FilterValue() string { return i.TitleVal + " " + i.URL }

type ResumeModel struct {
	Visible bool
	List    list.Model
	Width   int
	Height  int
}

func NewResumeModel() ResumeModel {
	dl := styles.NewListDelegate()
	li := list.New([]list.Item{}, dl, 0, 0)
	li.SetShowStatusBar(false)
	li.SetShowTitle(false)
	li.SetShowHelp(false)
	li.KeyMap.Quit.SetKeys("q")
	li.FilterInput.Cursor.Style = li.FilterInput.Cursor.Style.Foreground(styles.MauveColor)
	li.FilterInput.PromptStyle = li.FilterInput.PromptStyle.Foreground(styles.SecondaryColor)

	return ResumeModel{
		Visible: false,
		List:    li,
		Width:   60,
		Height:  10,
	}
}

func (m *ResumeModel) Show() {
	m.Visible = true
	m.LoadItems()
}

func (m *ResumeModel) Hide() {
	m.Visible = false
	m.List.SetItems([]list.Item{})
}

func (m *ResumeModel) LoadItems() {
	items, err := utils.LoadUnfinished()
	if err != nil {
		m.List.SetItems([]list.Item{})
		return
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Timestamp.After(items[j].Timestamp)
	})

	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = ResumeItem{
			URL:      item.URL,
			TitleVal: item.Title,
			FormatID: item.FormatID,
		}
	}

	m.List.SetItems(listItems)
}

func (m *ResumeModel) HandleResize(width, height int) {
	m.Width = width
	m.Height = height
	m.List.SetSize(width, height-7)
}

func (m *ResumeModel) DeleteSelected() {
	if item, ok := m.List.SelectedItem().(ResumeItem); ok {
		utils.RemoveUnfinished(item.URL)
		m.LoadItems()
	}
}

func (m *ResumeModel) SelectedItem() *utils.UnfinishedDownload {
	if item, ok := m.List.SelectedItem().(ResumeItem); ok {
		return &utils.UnfinishedDownload{
			URL:      item.URL,
			Title:    item.TitleVal,
			FormatID: item.FormatID,
		}
	}

	return nil
}

func (m *ResumeModel) View(width, height int) string {
	if !m.Visible {
		return ""
	}

	var headerText string
	if m.List.FilterState() == list.FilterApplied {
		headerText = "Filtered Results"
	} else {
		headerText = "Resume Downloads"
	}

	return styles.SectionHeaderStyle.Render(headerText) + "\n" + styles.ListContainer.Render(m.List.View())
}
