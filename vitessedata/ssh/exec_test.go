package ssh

import (
	"fmt"
	"testing"
)

func runCmd(hosts []string, cmd string) bool {
	ok := true
	ch := make(chan ExecResult)
	ExecCmdOnEachHost(hosts, cmd, ch)
	for i, _ := range hosts {
		r := <-ch
		if r.Err != nil {
			ok = false
			fmt.Printf("FAIL: Host # %d, %s %s\n", i, r.Host, r.Err.Error())
		} else {
			fmt.Printf("OK! Host # %d, %v\n", i, r)
		}
	}
	return ok
}

func TestExec(t *testing.T) {
	mlcmd := `printf '
this is a multi line cmd
line2
line3
' `
	t.Run("exec", func(t *testing.T) {
		hosts := []string{"localhost", "localhost", "localhost"}
		ok := runCmd(hosts, "hostname; hostname")
		if !ok {
			t.Error("Running hostname failed")
		}

		ok = runCmd(hosts, mlcmd)
		if !ok {
			t.Error("printf failed.")
		}
	})
}
