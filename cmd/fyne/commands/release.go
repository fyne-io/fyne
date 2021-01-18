package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne/internal/mobile"
	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	"fyne.io/fyne/v2/cmd/fyne/internal/util"
)

var (
	macAppStoreCategories = []string{
		"business", "developer-tools", "education", "entertainment", "finance", "games", "action-games",
		"adventure-games", "arcade-games", "board-games", "card-games", "casino-games", "dice-games",
		"educational-games", "family-games", "kids-games", "music-games", "puzzle-games", "racing-games",
		"role-playing-games", "simulation-games", "sports-games", "strategy-games", "trivia-games", "word-games",
		"graphics-design", "healthcare-fitness", "lifestyle", "medical", "music", "news", "photography",
		"productivity", "reference", "social-networking", "sports", "travel", "utilities", "video", "weather",
	}
)

// Declare conformity to Command interface
var _ Command = (*releaser)(nil)

type releaser struct {
	packager

	keyStore     string
	keyStorePass string
	keyPass      string
	developer    string
	password     string
}

// NewReleaser returns a command that can adapt app packages for distribution
func NewReleaser() Command {
	return &releaser{}
}

func (r *releaser) AddFlags() {
	flag.StringVar(&r.os, "os", "", "The operating system to target (android, android/arm, android/arm64, android/amd64, android/386, darwin, freebsd, ios, linux, netbsd, openbsd, windows)")
	flag.StringVar(&r.name, "name", "", "The name of the application, default is the executable file name")
	flag.StringVar(&r.icon, "icon", "Icon.png", "The name of the application icon file")
	flag.StringVar(&r.appID, "appID", "", "For ios or darwin targets an appID is required, for ios this must \nmatch a valid provisioning profile")
	flag.StringVar(&r.appVersion, "appVersion", "", "Version number in the form x, x.y or x.y.z semantic version")
	flag.IntVar(&r.appBuild, "appBuild", 0, "Build number, should be greater than 0 and incremented for each build")
	flag.StringVar(&r.keyStore, "keyStore", "", "Android: location of .keystore file containing signing information")
	flag.StringVar(&r.keyStorePass, "keyStorePass", "", "Android: password for the .keystore file, default take the password from stdin")
	flag.StringVar(&r.keyPass, "keyPass", "", "Android: password for the signer's private key, which is needed if the private key is password-protected. Default take the password from stdin")
	flag.StringVar(&r.certificate, "certificate", "", "iOS/macOS/Windows: name of the certificate to sign the build")
	flag.StringVar(&r.profile, "profile", "", "iOS/macOS: name of the provisioning profile for this release build")
	flag.StringVar(&r.developer, "developer", "", "Windows: the developer identity for your Microsoft store account")
	flag.StringVar(&r.password, "password", "", "Windows: password for the certificate used to sign the build")
	flag.StringVar(&r.tags, "tags", "", "A comma-separated list of build tags")
	flag.StringVar(&r.category, "category", "", "macOS: category of the app for store listing")
}

func (r *releaser) PrintHelp(indent string) {
	fmt.Println(indent, "The release command prepares an application for public distribution.")
}

func (r *releaser) Run(params []string) {
	if err := r.validate(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	r.packager.release = true

	if err := r.beforePackage(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	r.packager.Run(params)
	if err := r.afterPackage(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

func (r *releaser) afterPackage() error {
	if util.IsAndroid(r.os) {
		target := mobile.AppOutputName(r.os, r.packager.name)
		apk := filepath.Join(r.dir, target)
		if err := r.zipAlign(apk); err != nil {
			return err
		}
		return r.signAndroid(apk)
	}
	if r.os == "darwin" {
		return r.packageMacOSRelease()
	}
	if r.os == "ios" {
		return r.packageIOSRelease()
	}
	if r.os == "windows" {
		outName := r.name + ".appx"
		if pos := strings.LastIndex(r.name, ".exe"); pos > 0 {
			outName = r.name[:pos] + ".appx"
		}

		outPath := filepath.Join(r.dir, outName)
		os.Remove(outPath) // MakeAppx will hang if the file exists... ignore result
		if err := r.packageWindowsRelease(outPath); err != nil {
			return err
		}
		return r.signWindows(outPath)
	}
	return nil
}

func (r *releaser) beforePackage() error {
	if util.IsAndroid(r.os) {
		if err := util.RequireAndroidSDK(); err != nil {
			return err
		}
	}

	return nil
}

func (r *releaser) nameFromCertInfo(info string) string {
	// format should be "CN=Company, O=Company, L=City, S=State, C=Country"
	parts := strings.Split(info, ",")
	cn := parts[0]
	pos := strings.Index(strings.ToUpper(cn), "CN=")
	if pos == -1 {
		return cn // not what we were expecting, but should be OK
	}

	return cn[pos+3:]
}

func (r *releaser) packageIOSRelease() error {
	team, err := mobile.DetectIOSTeamID(r.certificate)
	if err != nil {
		return errors.New("failed to determine team ID")
	}

	payload := filepath.Join(r.dir, "Payload")
	os.Mkdir(payload, 0750)
	defer os.RemoveAll(payload)
	appName := mobile.AppOutputName(r.os, r.name)
	payloadAppDir := filepath.Join(payload, appName)
	os.Rename(filepath.Join(r.dir, appName), payloadAppDir)

	entitlementPath := filepath.Join(r.dir, "entitlements.plist")
	entitlements, _ := os.Create(entitlementPath)
	entitlementData := struct{ TeamID, AppID string }{
		TeamID: team,
		AppID:  r.appID,
	}
	err = templates.EntitlementsDarwinMobile.Execute(entitlements, entitlementData)
	entitlements.Close()
	if err != nil {
		return errors.New("failed to write entitlements plist template")
	}

	cmd := exec.Command("codesign", "-f", "-vv", "-s", r.certificate, "--entitlements",
		"entitlements.plist", "Payload/"+appName+"/")
	if err := cmd.Run(); err != nil {
		fyne.LogError("Codesign failed", err)
		return errors.New("unable to codesign application bundle")
	}

	return exec.Command("zip", "-r", appName[:len(appName)-4]+".ipa", "Payload/").Run()
}

func (r *releaser) packageMacOSRelease() error {
	appCert := r.certificate // try to derive two certificates from one name (they will be consistent)
	if strings.Contains(appCert, "Installer") {
		appCert = strings.Replace(appCert, "Installer", "Application", 1)
	}
	installCert := strings.Replace(appCert, "Application", "Installer", 1)
	unsignedPath := r.name + "-unsigned.pkg"

	defer os.RemoveAll(r.name + ".app") // this was the output of package and it can get in the way of future builds
	entitlementPath := filepath.Join(r.dir, "entitlements.plist")
	entitlements, _ := os.Create(entitlementPath)
	err := templates.EntitlementsDarwin.Execute(entitlements, nil)
	entitlements.Close()
	if err != nil {
		return errors.New("failed to write entitlements plist template")
	}
	defer os.Remove(entitlementPath)

	cmd := exec.Command("codesign", "-vfs", appCert, "--entitlement", "entitlements.plist", r.name+".app")
	err = cmd.Run()
	if err != nil {
		fyne.LogError("Codesign failed", err)
		return errors.New("unable to codesign application bundle")
	}

	cmd = exec.Command("productbuild", "--component", r.name+".app", "/Applications/",
		"--product", r.name+".app/Contents/Info.plist", unsignedPath)
	err = cmd.Run()
	if err != nil {
		fyne.LogError("Product build failed", err)
		return errors.New("unable to build macOS app package")
	}
	defer os.Remove(unsignedPath)

	cmd = exec.Command("productsign", "--sign", installCert, unsignedPath, r.name+".pkg")
	return cmd.Run()
}

func (r *releaser) packageWindowsRelease(outFile string) error {
	payload := filepath.Join(r.dir, "Payload")
	os.Mkdir(payload, 0750)
	defer os.RemoveAll(payload)

	manifestPath := filepath.Join(payload, "appxmanifest.xml")
	manifest, err := os.Create(manifestPath)
	if err != nil {
		return err
	}
	manifestData := struct{ AppID, Developer, DeveloperName, Name, Version string }{
		AppID: r.appID,
		// TODO read this info
		Developer:     r.developer,
		DeveloperName: r.nameFromCertInfo(r.developer),
		Name:          r.name,
		Version:       r.combinedVersion(),
	}
	err = templates.AppxManifestWindows.Execute(manifest, manifestData)
	manifest.Close()
	if err != nil {
		return errors.New("failed to write application manifest template")
	}

	util.CopyFile(r.icon, filepath.Join(payload, "Icon.png"))
	util.CopyFile(r.name, filepath.Join(payload, r.name))

	binDir, err := findWindowsSDKBin()
	if err != nil {
		return errors.New("cannot find makeappx.exe, make sure you have installed the Windows SDK")
	}

	cmd := exec.Command(filepath.Join(binDir, "makeappx.exe"), "pack", "/d", payload, "/p", outFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (r *releaser) signAndroid(path string) error {
	signer := filepath.Join(util.AndroidBuildToolsPath(), "/apksigner")

	args := []string{"sign", "--ks", r.keyStore}
	if r.keyStorePass != "" {
		args = append(args, "--ks-pass", "pass:"+r.keyStorePass)
	}
	if r.keyPass != "" {
		args = append(args, "--key-pass", "pass:"+r.keyPass)
	}
	args = append(args, path)

	cmd := exec.Command(signer, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (r *releaser) signWindows(appx string) error {
	binDir, err := findWindowsSDKBin()
	if err != nil {
		return errors.New("cannot find signtool.exe, make sure you have installed the Windows SDK")
	}

	cmd := exec.Command(filepath.Join(binDir, "signtool.exe"),
		"sign", "/a", "/v", "/fd", "SHA256", "/f", r.certificate, "/p", r.password, appx)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (r *releaser) validate() error {
	if r.os == "" {
		r.os = targetOS()
	}
	if util.IsMobile(r.os) || r.os == "windows" {
		if r.appVersion == "" { // Here it is required, if provided then package validate will check format
			return errors.New("missing required -appVersion parameter")
		}
		if r.appBuild <= 0 {
			return errors.New("missing required -appBuild parameter")
		}
	}
	if r.os == "windows" {
		if r.developer == "" {
			return errors.New("missing required -developer parameter for windows release,\n" +
				"use data from Partner Portal, format \"CN=Company, O=Company, L=City, S=State, C=Country\"")
		}
		if r.certificate == "" {
			return errors.New("missing required -certificate parameter for windows release")
		}
		if r.password == "" {
			return errors.New("missing required -password parameter for windows release")
		}
	}
	if util.IsAndroid(r.os) {
		if r.keyStore == "" {
			return errors.New("missing required -keyStore parameter for android release")
		}
	} else if r.os == "darwin" {
		if r.certificate == "" {
			r.certificate = "3rd Party Mac Developer Application"
		}
		if r.category == "" {
			return errors.New("missing required -category parameter for macOS release")
		} else if !isValidMacOSCategory(r.category) {
			return errors.New("category does not match one of the supported list: " +
				strings.Join(macAppStoreCategories, ", "))
		}
	} else if r.os == "ios" {
		if r.certificate == "" {
			r.certificate = "Apple Distribution"
		}
		if r.profile == "" {
			return errors.New("missing required -profile parameter for iOS release")
		}
	}
	return nil
}

func (r *releaser) zipAlign(path string) error {
	unaligned := filepath.Join(filepath.Dir(path), "unaligned.apk")
	err := os.Rename(path, unaligned)
	if err != nil {
		return nil
	}

	cmd := filepath.Join(util.AndroidBuildToolsPath(), "zipalign")
	err = exec.Command(cmd, "4", unaligned, path).Run()
	if err != nil {
		_ = os.Rename(path, unaligned) // ignore error, return previous
		return err
	}
	return os.Remove(unaligned)
}

func findWindowsSDKBin() (string, error) {
	inPath, err := exec.LookPath("makeappx.exe")
	if err == nil {
		return inPath, nil
	}

	matches, err := filepath.Glob("C:\\Program Files (x86)\\Windows Kits\\*\\bin\\*\\*\\makeappx.exe")
	if err != nil || len(matches) == 0 {
		return "", errors.New("failed to look up standard locations for makeappx.exe")
	}

	return filepath.Dir(matches[0]), nil
}

func isValidMacOSCategory(in string) bool {
	found := false
	for _, cat := range macAppStoreCategories {
		if cat == strings.ToLower(in) {
			found = true
			break
		}
	}
	return found
}
