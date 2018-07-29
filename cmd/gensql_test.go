package cmd

import "testing"

func TestJson(t *testing.T) {
	kick([]string{"gensql", "-t", "testdata/*", "-o", "testresult", "-f", "json"})
}
func TestCsv(t *testing.T) {
	kick([]string{"gensql", "-t", "testdata/*", "-o", "testresult", "-f", "csv"})
}

func kick(args []string) {
	RootCmd.SetArgs(args)
	RootCmd.Execute()
}
