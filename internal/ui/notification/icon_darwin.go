//go:build darwin

package notification

import _ "embed"

// Icon is the PNG data for the Crush icon, used for OSC 99 notifications.
//
//go:embed crush-icon.png
var Icon []byte
