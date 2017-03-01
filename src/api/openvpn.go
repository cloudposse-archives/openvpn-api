package api

import (
	"os"
	"io/ioutil"
	"os/exec"
	"strings"
	log "github.com/Sirupsen/logrus"
)

func EnsureUserCerts(name string) (err error) {
	err = nil
	logger := log.WithFields(log.Fields{"class": "openvpn", "method": "EnsureUserCerts"})
	if ! validCertsExits(name) {
		logger.Infof("Certificates from %v nof found", name)
		err = cleanCertsFor(name)
		if err != nil {
			return
		}
		err = generateCertsFor(name)
		if err != nil {
			return
		}
	}

	return
}

func GetClientConfig(name string) (result string, err error) {
	err = nil
	result = ""
	err = recreateAllClientConfigs()
	if err != nil {
		return
	}

	file := "/etc/openvpn/clients/" + name + "/" + name + "-combined.ovpn"

	if _, err = os.Stat(file); os.IsNotExist(err) {
		return
	}

	b, err := ioutil.ReadFile(file)
	result = string(b)

	return
}

func RevokeUser(name string) error {
	// easyrsa revoke {name}
	return nil
}

func validCertsExits(name string) bool {
	logger := log.WithFields(log.Fields{"class": "openvpn", "method": "validCertsExits"})

	cmd := exec.Command("ovpn_listclients")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false
	}

	clients := string(output)
	logger.Debugf("Clients: %s", clients)

	for _, client := range strings.Split(clients, "\n") {
		logger.Debugf("%v", client)
		if strings.Contains(client, name) && strings.Contains(client, "VALID") {
			return true
		}
	}

	return false
}

func cleanCertsFor(name string) (err error) {
	err = nil

	if err = deleteFileSafety("/etc/openvpn/pki/issued/" + name + ".crt"); err != nil {
		return
	}

	if err = deleteFileSafety("/etc/openvpn/pki/reqs/" + name + ".req"); err != nil {
		return
	}

	if err = deleteFileSafety("/etc/openvpn/pki/private/" + name + ".key"); err != nil {
		return
	}

	return
}

func deleteFileSafety(file string) (err error) {
	err = nil
	if _, e := os.Stat(file); os.IsNotExist(e) {
		return
	}

	if err = os.Remove(file); err != nil {
		return
	}

	return

}

func generateCertsFor(user string) error {
	cmd := exec.Command("easyrsa", "build-client-full", user, "nopass")
	return cmd.Run()
}

func recreateAllClientConfigs() error {
	cmd := exec.Command("ovpn_getclient_all")
	return cmd.Run()
}
