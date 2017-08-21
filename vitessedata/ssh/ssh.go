package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"net"
	"os"
	"os/user"
)

func sshAgent() ssh.AuthMethod {
	if ag, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(ag).Signers)
	}
	return nil
}

func getKeyFile(usr *user.User) (key ssh.Signer, err error) {
	file := usr.HomeDir + "/.ssh/id_rsa"
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return
	}
	return
}

func RunCommand(host string, cmd string) (string, error) {
	host = host + ":22"

	/* Get the current user */
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	/* Set up the ssh config */
	var config *ssh.ClientConfig
	config = &ssh.ClientConfig{
		User:            usr.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.Auth = make([]ssh.AuthMethod, 0, 3)

	/* next option: is there an ssh agent */
	if ag := sshAgent(); ag != nil {
		config.Auth = append(config.Auth, ag)
	}

	/* finally, try to use my public key */
	key, err := getKeyFile(usr)
	if err == nil {
		config.Auth = append(config.Auth, ssh.PublicKeys(key))
	}

	/* if no auth method found, then error out */
	if len(config.Auth) == 0 {
		if err != nil {
			return "", err
		} else {
			return "", fmt.Errorf("ERROR: no authentication method found.")
		}
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return "", err
	}

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.Output(cmd)
	if err != nil {
		return "", err
	}
	return string(output), nil
}
