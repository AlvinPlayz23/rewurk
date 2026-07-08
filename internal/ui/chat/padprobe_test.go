package chat

import (
	"fmt"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestPadProbe(t *testing.T) {
	st := lipgloss.NewStyle().PaddingLeft(2)
	out := st.Render("line1\nline2")
	fmt.Printf("QUOTED: %q\n", out)
}
