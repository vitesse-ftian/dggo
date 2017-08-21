package ssh

import (
	"github.com/vitesse-ftian/dggo/vitessedata/dglog"
	"os/exec"
	"strings"
)

type ExecResult struct {
	Host string
	Err  error
	Out  string
}

func BinAbs(bin string) string {
	p, err := exec.Command("sh", "-c", "which "+bin).Output()
	dglog.FatalErr(err, "Cannot find path of "+bin)
	return strings.TrimSpace(string(p))
}

func ExecAnyError(res []ExecResult) error {
	for _, r := range res {
		if r.Err != nil {
			return r.Err
		}
	}
	return nil
}

func execOnOneHost(host string, cmd string, c chan ExecResult) {
	result := ExecResult{host, nil, ""}
	result.Out, result.Err = RunCommand(host, cmd)
	c <- result
}

func ExecCmdOnEachHost(hosts []string, cmd string, result chan ExecResult) {
	for _, h := range hosts {
		go execOnOneHost(h, cmd, result)
	}
}

func ExecCmdArrOnEachHost(hosts []string, cmds []string, result chan ExecResult) {
	for i := 0; i < len(hosts); i++ {
		go execOnOneHost(hosts[i], cmds[i], result)
	}
}

func ExecCmdOn(hosts []string, cmd string) []ExecResult {
	ch := make(chan ExecResult)
	res := make([]ExecResult, len(hosts))
	ExecCmdOnEachHost(hosts, cmd, ch)

	for nchecked := 0; nchecked < len(hosts); nchecked++ {
		res[nchecked] = <-ch
	}
	return res
}

func ExecOn(hosts []string, cmds []string) []ExecResult {
	ch := make(chan ExecResult)
	res := make([]ExecResult, len(hosts))
	ExecCmdArrOnEachHost(hosts, cmds, ch)

	for nchecked := 0; nchecked < len(hosts); nchecked++ {
		res[nchecked] = <-ch
	}
	return res
}

func MyHostName() (string, error) {
	h := []string{"localhost"}
	res := ExecCmdOn(h, "hostname")
	if res[0].Err != nil {
		return "", res[0].Err
	}
	return strings.TrimSpace(res[0].Out), nil
}
