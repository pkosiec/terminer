package installer

import "github.com/pkosiec/terminer/internal/shell"

func (installer *Installer) SetShell(s shell.Shell) {
	installer.sh = s
}
