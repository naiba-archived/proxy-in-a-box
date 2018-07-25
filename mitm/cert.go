package mitm

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/getlantern/keyman"
)

const (
	day         = 24 * time.Hour
	doubleWeeks = day * 14
	month       = 1
	year        = 1
)

//GenerateCertForClient rt
func (m *MITM) GenerateCertForClient() error {
	var err error
	if m.pk, err = keyman.LoadPKFromFile(m.TLSConf.PrivateKeyFile); err != nil {
		m.pk, err = keyman.GeneratePK(2048)
		if err != nil {
			return fmt.Errorf("Unable to generate private key: %s", err)
		}
		m.pk.WriteToFile(m.TLSConf.PrivateKeyFile)
	}
	m.pkPem = m.pk.PEMEncoded()
	m.issuingCert, err = keyman.LoadCertificateFromFile(m.TLSConf.CertFile)
	if err != nil || m.issuingCert.ExpiresBefore(time.Now().AddDate(0, month, 0)) {
		m.issuingCert, err = m.pk.TLSCertificateFor(
			time.Now().AddDate(year, 0, 0),
			true,
			nil,
			m.TLSConf.Organization,
			m.TLSConf.CommonName,
		)
		if err != nil {
			return fmt.Errorf("Unable to generate self-signed issuing certificate: %s", err)
		}
		m.issuingCert.WriteToFile(m.TLSConf.CertFile)
	}
	m.issuingCertPem = m.issuingCert.PEMEncoded()
	return nil
}

//FakeCert rt
func (m *MITM) FakeCert(domain string) (*tls.Certificate, error) {
	cert, has := m.cache.Get("DC" + domain)
	if has {
		fmt.Println("CertCache:", domain)
		return cert.(*tls.Certificate), nil
	}

	//create certificate
	certTTL := doubleWeeks
	generatedCert, err := m.pk.TLSCertificateFor(
		time.Now().Add(certTTL),
		false,
		m.issuingCert,
		m.TLSConf.Organization,
		domain)
	if err != nil {
		return nil, fmt.Errorf("Unable to issue certificate: %s", err)
	}
	keyPair, err := tls.X509KeyPair(generatedCert.PEMEncoded(), m.pkPem)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse keypair for tls: %s", err)
	}

	m.cache.Set("DC"+domain, &keyPair, time.Hour*48)
	return &keyPair, nil
}
