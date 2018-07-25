package keyman

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	// TODO: Make sure to handle case where library is on a different path
	OSX_SYSTEM_KEYCHAIN_PATH = "/Library/Keychains/System.keychain"
)

func DeleteTrustedRootByName(commonName string, prompt string) error {
	cmd := elevatedIfNecessary(prompt)("security", "delete-certificate", "-c", commonName, OSX_SYSTEM_KEYCHAIN_PATH)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Unable to run security command: %s\n%s", err, out)
	}
	return nil
}

func (cert *Certificate) AddAsTrustedRoot(prompt string) error {
	tempFileName, err := cert.WriteToTempFile()
	defer func() {
		if err := os.Remove(tempFileName); err != nil {
			log.Debugf("Unable to remove file: %v", err)
		}
	}()
	if err != nil {
		return fmt.Errorf("Unable to create temp file: %s", err)
	}

	cmd := exec.Command("security", "verify-cert", "-c", tempFileName)
	_, err = cmd.CombinedOutput()
	if err == nil {
		// certificate verified successfully so it's already a trusted root, no need
		// to install.
		return nil
	}

	// Add it as a trusted cert
	cmd = elevatedIfNecessary(prompt)("security", "add-trusted-cert", "-d", "-k", OSX_SYSTEM_KEYCHAIN_PATH, tempFileName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Unable to run security command: %s\n%s", err, out)
	} else {
		cmd := exec.Command("security", "verify-cert", "-c", tempFileName)
		out, err := cmd.CombinedOutput()
		log.Debugf("%v: %v", out, err)
		return nil
	}
}

// Checks whether this certificate is install based purely on looking for a cert
// in the system keychain that has the same common name.  This function returns
// true if there are one or more certs in the system keychain whose common name
// matches this cert.
func (cert *Certificate) IsInstalled() (bool, error) {
	cmd := exec.Command("security", "find-certificate", "-c", cert.X509().Subject.CommonName, OSX_SYSTEM_KEYCHAIN_PATH)
	err := cmd.Run()

	found := err == nil
	return found, nil
}
