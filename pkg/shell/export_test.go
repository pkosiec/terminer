package shell

func ExposeInternalShell() *shell {
	return &shell{}
}

func (s *shell) IsCommandAvailable(cmdName string) bool {
	return s.isCommandAvailable(cmdName)
}
