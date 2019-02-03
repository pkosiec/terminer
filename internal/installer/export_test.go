package installer

import "github.com/pkosiec/terminer/internal/sh"

func (installer *Installer) SetSh(s sh.Sh) {
	installer.sh = s
}
