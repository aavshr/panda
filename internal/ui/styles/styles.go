package styles

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	TitleColor             = "#eb5e55"
	TitleSecondaryColor    = "#f9c784"
	ActiveContainerColor   = "#f9c784"
	ListItemColor          = "#c6d8d3"
	ListItemSecondaryColor = "#fdf0d5"
	DescriptionColor       = "#d81e5b"
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

func SetActiveBorder(s *lipgloss.Style) {
	s.BorderForeground(lipgloss.Color(ActiveContainerColor))
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
		Foreground(lipgloss.Color(TitleSecondaryColor))
	return s
}

func DefaultListSelectedStyle() lipgloss.Style {
	s := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(TitleColor))
	return s
}

func DefaultListItemStyle() lipgloss.Style {
	s := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ListItemColor))
	return s
}

func DefaultListItemSecondaryStyle() lipgloss.Style {
	s := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ListItemSecondaryColor))
	return s
}
