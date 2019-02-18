package installer

import "github.com/pkosiec/terminer/pkg/shell"

func (installer *Installer) SetShell(s shell.Shell) {
	installer.sh = s
}
