// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mobile

// APK is the archival format used for Android apps. It is a ZIP archive with
// three extra files:
//
//	META-INF/MANIFEST.MF
//	META-INF/CERT.SF
//	META-INF/CERT.RSA
//
// The MANIFEST.MF comes from the Java JAR archive format. It is a list of
// files included in the archive along with a SHA1 hash, for example:
//
//	Name: lib/armeabi/libbasic.so
//	SHA1-Digest: ntLSc1eLCS2Tq1oB4Vw6jvkranw=
//
// For debugging, the equivalent SHA1-Digest can be generated with OpenSSL:
//
//	cat lib/armeabi/libbasic.so | openssl sha1 -binary | openssl base64
//
// CERT.SF is a similar manifest. It begins with a SHA1 digest of the entire
// manifest file:
//
//	Signature-Version: 1.0
//	Created-By: 1.0 (Android)
//	SHA1-Digest-Manifest: aJw+u+10C3Enbg8XRCN6jepluYA=
//
// Then for each entry in the manifest it has a SHA1 digest of the manfiest's
// hash combined with the file name:
//
//	Name: lib/armeabi/libbasic.so
//	SHA1-Digest: Q7NAS6uzrJr6WjePXSGT+vvmdiw=
//
// This can also be generated with openssl:
//
//	echo -en "Name: lib/armeabi/libbasic.so\r\nSHA1-Digest: ntLSc1eLCS2Tq1oB4Vw6jvkranw=\r\n\r\n" | openssl sha1 -binary | openssl base64
//
// Note the \r\n line breaks.
//
// CERT.RSA is an RSA signature block made of CERT.SF. Verify it with:
//
//	openssl smime -verify -in CERT.RSA -inform DER -content CERT.SF cert.pem
//
// The APK format imposes two extra restrictions on the ZIP format. First,
// it is uncompressed. Second, each contained file is 4-byte aligned. This
// allows the Android OS to mmap contents without unpacking the archive.

// Note: to make life a little harder, Android Studio stores the RSA key used
// for signing in an Oracle Java proprietary keystore format, JKS. For example,
// the generated debug key is in ~/.android/debug.keystore, and can be
// extracted using the JDK's keytool utility:
//
//	keytool -importkeystore -srckeystore ~/.android/debug.keystore -destkeystore ~/.android/debug.p12 -deststoretype PKCS12
//
// Once in standard PKCS12, the key can be converted to PEM for use in the
// Go crypto packages:
//
//	openssl pkcs12 -in ~/.android/debug.p12 -nocerts -nodes -out ~/.android/debug.pem
//
// Fortunately for debug builds, all that matters is that the APK is signed.
// The choice of key is unimportant, so we can generate one for normal builds.
// For production builds, we can ask users to provide a PEM file.

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
)

// NewWriter returns a new Writer writing an APK file to w.
// The APK will be signed with key.
func NewWriter(w io.Writer, priv *rsa.PrivateKey) *Writer {
	apkw := &Writer{priv: priv}
	apkw.w = zip.NewWriter(&countWriter{apkw: apkw, w: w})
	return apkw
}

// Writer implements an APK file writer.
type Writer struct {
	offset   int
	w        *zip.Writer
	priv     *rsa.PrivateKey
	manifest []manifestEntry
	cur      *fileWriter
}

// Create adds a file to the APK archive using the provided name.
//
// The name must be a relative path. The file's contents must be written to
// the returned io.Writer before the next call to Create or Close.
func (w *Writer) Create(name string) (io.Writer, error) {
	if err := w.clearCur(); err != nil {
		return nil, fmt.Errorf("apk: Create(%s): %v", name, err)
	}
	res, err := w.create(name)
	if err != nil {
		return nil, fmt.Errorf("apk: Create(%s): %v", name, err)
	}
	return res, nil
}

func (w *Writer) create(name string) (io.Writer, error) {
	// Align start of file contents by using Extra as padding.
	if err := w.w.Flush(); err != nil { // for exact offset
		return nil, err
	}
	const fileHeaderLen = 30 // + filename + extra
	start := w.offset + fileHeaderLen + len(name)
	extra := start % 4

	zipfw, err := w.w.CreateHeader(&zip.FileHeader{
		Name:  name,
		Extra: make([]byte, extra),
	})
	if err != nil {
		return nil, err
	}
	w.cur = &fileWriter{
		name: name,
		w:    zipfw,
		sha1: sha1.New(),
	}
	return w.cur, nil
}

// Close finishes writing the APK. This includes writing the manifest and
// signing the archive, and writing the ZIP central directory.
//
// It does not close the underlying writer.
func (w *Writer) Close() error {
	if err := w.clearCur(); err != nil {
		return fmt.Errorf("apk: %v", err)
	}

	hasDex := false
	for _, entry := range w.manifest {
		if entry.name == "classes.dex" {
			hasDex = true
			break
		}
	}

	manifest := new(bytes.Buffer)
	if hasDex {
		fmt.Fprint(manifest, manifestDexHeader)
	} else {
		fmt.Fprint(manifest, manifestHeader)
	}
	certBody := new(bytes.Buffer)

	for _, entry := range w.manifest {
		n := entry.name
		h := base64.StdEncoding.EncodeToString(entry.sha1.Sum(nil))
		fmt.Fprintf(manifest, "Name: %s\nSHA1-Digest: %s\n\n", n, h)
		cHash := sha1.New()
		fmt.Fprintf(cHash, "Name: %s\r\nSHA1-Digest: %s\r\n\r\n", n, h)
		ch := base64.StdEncoding.EncodeToString(cHash.Sum(nil))
		fmt.Fprintf(certBody, "Name: %s\nSHA1-Digest: %s\n\n", n, ch)
	}

	mHash := sha1.New()
	_, err := mHash.Write(manifest.Bytes())
	if err != nil {
		return err
	}
	cert := new(bytes.Buffer)
	fmt.Fprint(cert, certHeader)
	fmt.Fprintf(cert, "SHA1-Digest-Manifest: %s\n\n", base64.StdEncoding.EncodeToString(mHash.Sum(nil)))
	cert.Write(certBody.Bytes())

	mw, err := w.Create("META-INF/MANIFEST.MF")
	if err != nil {
		return err
	}
	if _, err := mw.Write(manifest.Bytes()); err != nil {
		return err
	}

	cw, err := w.Create("META-INF/CERT.SF")
	if err != nil {
		return err
	}
	if _, err := cw.Write(cert.Bytes()); err != nil {
		return err
	}

	rsa, err := signPKCS7(rand.Reader, w.priv, cert.Bytes())
	if err != nil {
		return fmt.Errorf("apk: %v", err)
	}
	rw, err := w.Create("META-INF/CERT.RSA")
	if err != nil {
		return err
	}
	if _, err := rw.Write(rsa); err != nil {
		return err
	}

	return w.w.Close()
}

const manifestHeader = `Manifest-Version: 1.0
Created-By: 1.0 (Go)

`

const manifestDexHeader = `Manifest-Version: 1.0
Dex-Location: classes.dex
Created-By: 1.0 (Go)

`

const certHeader = `Signature-Version: 1.0
Created-By: 1.0 (Go)
`

func (w *Writer) clearCur() error {
	if w.cur == nil {
		return nil
	}
	w.manifest = append(w.manifest, manifestEntry{
		name: w.cur.name,
		sha1: w.cur.sha1,
	})
	w.cur.closed = true
	w.cur = nil
	return nil
}

type manifestEntry struct {
	name string
	sha1 hash.Hash
}

type countWriter struct {
	apkw *Writer
	w    io.Writer
}

func (c *countWriter) Write(p []byte) (n int, err error) {
	n, err = c.w.Write(p)
	c.apkw.offset += n
	return n, err
}

type fileWriter struct {
	name   string
	w      io.Writer
	sha1   hash.Hash
	closed bool
}

func (w *fileWriter) Write(p []byte) (n int, err error) {
	if w.closed {
		return 0, fmt.Errorf("apk: write to closed file %q", w.name)
	}
	_, err = w.sha1.Write(p)
	if err != nil {
		return 0, fmt.Errorf("apk: sha1 write %s", err)
	}
	n, err = w.w.Write(p)
	if err != nil {
		return 0, fmt.Errorf("apk: %v", err)
	}
	return n, err
}
