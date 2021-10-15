// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mobile

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"io"
	"math/big"
	"time"
)

// signPKCS7 does the minimal amount of work necessary to embed an RSA
// signature into a PKCS#7 certificate.
//
// We prepare the certificate using the x509 package, read it back in
// to our custom data type and then write it back out with the signature.
func signPKCS7(rand io.Reader, priv *rsa.PrivateKey, msg []byte) ([]byte, error) {
	const serialNumber = 0x5462c4dd // arbitrary
	name := pkix.Name{CommonName: "gomobile"}

	template := &x509.Certificate{
		SerialNumber:       big.NewInt(serialNumber),
		SignatureAlgorithm: x509.SHA1WithRSA,
		Subject:            name,
		NotAfter:           time.Date(2120, time.January, 1, 12, 0, 0, 0, time.UTC), // more than 50 years for Google Play Store
	}

	b, err := x509.CreateCertificate(rand, template, template, priv.Public(), priv)
	if err != nil {
		return nil, err
	}

	c := certificate{}
	if _, err := asn1.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	h := sha1.New()
	_, err = h.Write(msg)
	if err != nil {
		return nil, err
	}
	hashed := h.Sum(nil)

	signed, err := rsa.SignPKCS1v15(rand, priv, crypto.SHA1, hashed)
	if err != nil {
		return nil, err
	}

	content := pkcs7SignedData{
		ContentType: oidSignedData,
		Content: signedData{
			Version: 1,
			DigestAlgorithms: []pkix.AlgorithmIdentifier{{
				Algorithm:  oidSHA1,
				Parameters: asn1.RawValue{Tag: 5},
			}},
			ContentInfo:  contentInfo{Type: oidData},
			Certificates: c,
			SignerInfos: []signerInfo{{
				Version: 1,
				IssuerAndSerialNumber: issuerAndSerialNumber{
					Issuer:       name.ToRDNSequence(),
					SerialNumber: serialNumber,
				},
				DigestAlgorithm: pkix.AlgorithmIdentifier{
					Algorithm:  oidSHA1,
					Parameters: asn1.RawValue{Tag: 5},
				},
				DigestEncryptionAlgorithm: pkix.AlgorithmIdentifier{
					Algorithm:  oidRSAEncryption,
					Parameters: asn1.RawValue{Tag: 5},
				},
				EncryptedDigest: signed,
			}},
		},
	}

	return asn1.Marshal(content)
}

type pkcs7SignedData struct {
	ContentType asn1.ObjectIdentifier
	Content     signedData `asn1:"tag:0,explicit"`
}

// signedData is defined in rfc2315, section 9.1.
type signedData struct {
	Version          int
	DigestAlgorithms []pkix.AlgorithmIdentifier `asn1:"set"`
	ContentInfo      contentInfo
	Certificates     certificate  `asn1:"tag:0,explicit"`
	SignerInfos      []signerInfo `asn1:"set"`
}

type contentInfo struct {
	Type asn1.ObjectIdentifier
	// Content is optional in PKCS#7 and not provided here.
}

// certificate is defined in rfc2459, section 4.1.
type certificate struct {
	TBSCertificate     tbsCertificate
	SignatureAlgorithm pkix.AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

// tbsCertificate is defined in rfc2459, section 4.1.
type tbsCertificate struct {
	Version      int `asn1:"tag:0,default:2,explicit"`
	SerialNumber int
	Signature    pkix.AlgorithmIdentifier
	Issuer       pkix.RDNSequence // pkix.Name
	Validity     validity
	Subject      pkix.RDNSequence // pkix.Name
	SubjectPKI   subjectPublicKeyInfo
}

// validity is defined in rfc2459, section 4.1.
type validity struct {
	NotBefore time.Time
	NotAfter  time.Time
}

// subjectPublicKeyInfo is defined in rfc2459, section 4.1.
type subjectPublicKeyInfo struct {
	Algorithm        pkix.AlgorithmIdentifier
	SubjectPublicKey asn1.BitString
}

type signerInfo struct {
	Version                   int
	IssuerAndSerialNumber     issuerAndSerialNumber
	DigestAlgorithm           pkix.AlgorithmIdentifier
	DigestEncryptionAlgorithm pkix.AlgorithmIdentifier
	EncryptedDigest           []byte
}

type issuerAndSerialNumber struct {
	Issuer       pkix.RDNSequence // pkix.Name
	SerialNumber int
}

// Various ASN.1 Object Identifies, mostly from rfc3852.
var (
	oidData          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 1}
	oidSignedData    = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 2}
	oidSHA1          = asn1.ObjectIdentifier{1, 3, 14, 3, 2, 26}
	oidRSAEncryption = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
)
