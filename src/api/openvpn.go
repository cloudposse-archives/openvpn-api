package api

import (
	"os"
	"io/ioutil"
	"os/exec"
	"strings"
	log "github.com/Sirupsen/logrus"
)

// EnsureUserCerts - Ensure that user have valid certificates.
//                   If user certs expired or non exists - create them
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

// GetClientConfig - Regenerate all clients configs
//                   Return client config for user {name}
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

// Check if user certificate exists and it is not expired
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

// Remove user certificates
func cleanCertsFor(name string) (err error) {
	err = nil

	if err = deleteFileSafely("/etc/openvpn/pki/issued/" + name + ".crt"); err != nil {
		return
	}

	if err = deleteFileSafely("/etc/openvpn/pki/reqs/" + name + ".req"); err != nil {
		return
	}

	if err = deleteFileSafely("/etc/openvpn/pki/private/" + name + ".key"); err != nil {
		return
	}

	return
}

// Delete file if exists
func deleteFileSafely(file string) (err error) {
	err = nil
	if _, e := os.Stat(file); os.IsNotExist(e) {
		return
	}

	if err = os.Remove(file); err != nil {
		return
	}

	return

}

// Generate certificates for user {user}
func generateCertsFor(user string) error {
	cmd := exec.Command("easyrsa", "build-client-full", user, "nopass")
	return cmd.Run()
}

// Recreate all users clients configs
func recreateAllClientConfigs() error {
	cmd := exec.Command("ovpn_getclient_all")
	return cmd.Run()
}
