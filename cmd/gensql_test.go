package cmd

import "testing"

func TestBasic(t *testing.T) {
	kick([]string{"gensql", "-t", "testdata/*", "-o", "testresult", "-f", "json"})
}

func kick(args []string) {
	RootCmd.SetArgs(args)
	RootCmd.Execute()
}
