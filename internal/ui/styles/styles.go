package styles

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	TitleColor             = lipgloss.Color("#eb5e55")
	TitleSecondaryColor    = lipgloss.Color("#f9c784")
	ActiveContainerColor   = lipgloss.Color("#f9c784")
	ListItemColor          = lipgloss.Color("#c6d8d3")
	ListItemSecondaryColor = lipgloss.Color("#fdf0d5")
	DescriptionColor       = lipgloss.Color("#d81e5b")
)

func MainContainerStyle() lipgloss.Style {
	s := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), true)
	return s
}

func ContainerStyle() lipgloss.Style {
	s := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true)
	return s
}

func SetNormalBorder(s *lipgloss.Style) {
	s.UnsetBorderForeground()
}

func SetSelectedBorder(s *lipgloss.Style) {
	s.BorderForeground(ActiveContainerColor)
}

func SetFocusedBorder(s *lipgloss.Style) {
	s.BorderForeground(TitleColor)
}

func SidebarContainerStyle(width int) lipgloss.Style {
	s := lipgloss.NewStyle().
		Width(width)
	return s
}

func MessagesContainerStyle(width int) lipgloss.Style {
	s := lipgloss.NewStyle().
		Width(width)
	return s
}

func DefaultListStyle() lipgloss.Style {
	s := lipgloss.NewStyle().
		Bold(true).
		Foreground(TitleSecondaryColor)
	return s
}

func DefaultListSelectedStyle() lipgloss.Style {
	s := lipgloss.NewStyle().
		Bold(true).
		Foreground(TitleColor)
	return s
}

func DefaultListItemStyle() lipgloss.Style {
	s := lipgloss.NewStyle().
		Foreground(ListItemColor)
	return s
}

func DefaultListItemSecondaryStyle() lipgloss.Style {
	s := lipgloss.NewStyle().
		Foreground(ListItemSecondaryColor)
	return s
}
