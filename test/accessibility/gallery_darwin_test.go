//go:build darwin && cgo && accessibility && axintegration

package accessibility

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"
)

var galleryPath = flag.String("gallery", "", "Path to fyne_accessibility built with -tags accessibility")

var errGalleryForeground = errors.New("the gallery became the foreground application; background tests must not steal focus")

type gallery struct {
	t     *testing.T
	app   axElement
	title string
	done  <-chan struct{}
}

func startGallery(t *testing.T, scenario string) *gallery {
	t.Helper()
	path, err := filepath.Abs(*galleryPath)
	if err != nil {
		t.Fatal(err)
	}
	logPath := filepath.Join(t.TempDir(), "gallery.log")
	log, err := os.Create(logPath)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	cmd := exec.CommandContext(ctx, path, "-scenario", scenario, "-quit-after", "3m", "-background")
	cmd.Stdout, cmd.Stderr = log, log
	cmd.Env = append(os.Environ(), "FYNE_SCALE=1")
	if err := cmd.Start(); err != nil {
		cancel()
		log.Close()
		t.Fatal(err)
	}
	done := make(chan struct{})
	var exitErr error
	go func() {
		exitErr = cmd.Wait()
		close(done)
	}()
	t.Cleanup(func() {
		select {
		case <-done:
			t.Errorf("gallery exited before cleanup: %v", exitErr)
		default:
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
				t.Errorf("stop gallery: %v", err)
			}
			select {
			case <-done:
			case <-time.After(3 * time.Second):
				cancel()
				<-done
			}
		}
		cancel()
		if err := log.Close(); err != nil {
			t.Error(err)
		}
		if t.Failed() {
			output, err := os.ReadFile(logPath)
			if err != nil {
				t.Error(err)
			} else {
				t.Logf("gallery output:\n%s", output)
			}
		}
	})
	app, err := application(cmd.Process.Pid)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(app.close)
	g := &gallery{t: t, app: app, title: "Fyne Accessibility - " + scenario, done: done}
	g.expect("AXStaticText", strings.ToUpper(scenario[:1])+scenario[1:]+" ready", nil)
	role, err := app.text("AXRole")
	if err != nil || role != "AXApplication" {
		t.Fatalf("application role = %q, error = %v", role, err)
	}
	return g
}

func (g *gallery) checkBackground() error {
	g.t.Helper()
	frontmost, err := g.app.text("AXFrontmost")
	if err != nil {
		return err
	}
	if frontmost == "true" {
		return errGalleryForeground
	}
	if frontmost != "false" {
		return fmt.Errorf("AXFrontmost is not available yet")
	}
	return nil
}

// wait reads a fresh tree on every attempt: Fyne may replace native elements
// after any repaint, so an element from a previous action must not be reused.
func (g *gallery) wait(description string, predicate func(*axSnapshot) bool) *axSnapshot {
	g.t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	lastTree := "(no complete accessibility tree)"
	var lastErr error
	for time.Now().Before(deadline) {
		select {
		case <-g.done:
			g.t.Fatalf("gallery exited while waiting for %s", description)
		default:
		}
		err := g.checkBackground()
		var tree *axSnapshot
		if err == nil {
			tree, err = snapshot(g.app, g.title, deadline)
		}
		if err == nil && predicate(tree) {
			err = g.checkBackground()
			if err == nil {
				return tree
			}
		}
		if tree != nil {
			lastTree = tree.String()
			tree.close()
		}
		if errors.Is(err, errGalleryForeground) {
			g.t.Fatal(err)
		}
		lastErr = err
		time.Sleep(50 * time.Millisecond)
	}
	g.t.Fatalf("timed out waiting for %s; last AX error: %v\n%s", description, lastErr, lastTree)
	return nil
}

func (g *gallery) expect(role, label string, attributes map[string]string) {
	g.t.Helper()
	tree := g.wait(fmt.Sprintf("%s %q %v", role, label, attributes), func(tree *axSnapshot) bool {
		node := tree.findWithAttributes(role, label, attributes)
		if node == nil {
			return false
		}
		return node.positioned && node.width > 0 && node.height > 0 && node.attributes["AXRoleDescription"] != ""
	})
	defer tree.close()
	node := tree.findWithAttributes(role, label, attributes)
	for _, attribute := range []string{"AXWindow", "AXTopLevelUIElement"} {
		ok, err := node.element.related(attribute, tree.nodes[0].element)
		if err != nil || !ok {
			g.t.Fatalf("%s %q has incorrect %s: %v\n%s", role, label, attribute, err, tree)
		}
	}
}

func (g *gallery) absent(role, label string) {
	g.t.Helper()
	tree := g.wait(fmt.Sprintf("absence of %s %q", role, label), func(tree *axSnapshot) bool {
		return tree.find(role, label) == nil
	})
	tree.close()
}

func (g *gallery) act(role, label, action string) {
	g.t.Helper()
	tree := g.wait(fmt.Sprintf("%s %q", role, label), func(tree *axSnapshot) bool {
		return tree.find(role, label) != nil
	})
	defer tree.close()
	node := tree.find(role, label)
	actions, err := node.element.actions()
	if err != nil || !slices.Contains(actions, action) {
		g.t.Fatalf("%s %q does not expose %s; actions=%v, error=%v\n%s",
			role, label, action, actions, err, tree)
	}
	if err := node.element.perform(action); err != nil {
		g.t.Fatalf("%s on %s %q: %v\n%s", action, role, label, err, tree)
	}
}

func (g *gallery) setValue(role, label, value string) {
	g.t.Helper()
	tree := g.wait(fmt.Sprintf("%s %q", role, label), func(tree *axSnapshot) bool {
		return tree.find(role, label) != nil
	})
	defer tree.close()
	if err := tree.find(role, label).element.setValue(value); err != nil {
		g.t.Fatalf("set value on %s %q: %v\n%s", role, label, err, tree)
	}
}

func TestAXGallery(t *testing.T) {
	if !trusted() {
		t.Fatal("Accessibility permission is required for this runner. Grant the runner or its launching terminal access in System Settings > Privacy & Security > Accessibility, then rerun. See README.md; no GUI tests were run.")
	}
	if *galleryPath == "" {
		t.Fatal("pass -gallery /path/to/fyne_accessibility (built with -tags accessibility); see README.md")
	}

	t.Run("controls", func(t *testing.T) {
		g := startGallery(t, "controls")
		g.expect("AXButton", "Count presses", map[string]string{"AXEnabled": "true"})
		g.expect("AXButton", "Disabled button", map[string]string{"AXEnabled": "false"})
		g.expect("AXCheckBox", "Already checked", map[string]string{"AXValue": "1"})
		g.expect("AXRadioButton", "Radio Alpha", map[string]string{"AXValue": "1"})
		g.expect("AXProgressIndicator", "", map[string]string{"AXValue": "25%"})
		g.expect("AXProgressIndicator", "", map[string]string{"AXValue": "indeterminate"})
		g.expect("AXProgressIndicator", "", map[string]string{"AXValue": "idle"})
		g.act("AXButton", "Count presses", "AXPress")
		g.expect("AXStaticText", "Pressed 1 times", nil)
		g.act("AXCheckBox", "Enable notifications", "AXPress")
		g.expect("AXCheckBox", "Enable notifications", map[string]string{"AXValue": "1"})
		g.expect("AXStaticText", "Notifications: true", nil)
		g.act("AXRadioButton", "Radio Beta", "AXPress")
		g.expect("AXRadioButton", "Radio Beta", map[string]string{"AXValue": "1"})
		g.expect("AXRadioButton", "Radio Alpha", map[string]string{"AXValue": "0"})
		g.expect("AXStaticText", "Radio: Radio Beta", nil)
		g.act("AXCheckBox", "Group Green", "AXPress")
		g.expect("AXStaticText", "Group: Group Green", nil)
		g.act("AXSlider", "", "AXIncrement")
		g.expect("AXSlider", "", map[string]string{"AXValue": "5"})
		g.expect("AXStaticText", "Slider: 5", nil)
		g.act("AXSlider", "", "AXDecrement")
		g.expect("AXSlider", "", map[string]string{"AXValue": "4"})
		g.setValue("AXSlider", "", "8")
		g.expect("AXStaticText", "Slider: 8", nil)
		g.act("AXButton", "Choose an option", "AXShowMenu")
		g.act("AXButton", "Choice Two", "AXPress")
		g.expect("AXStaticText", "Choice: Choice Two", nil)
	})

	t.Run("text", func(t *testing.T) {
		g := startGallery(t, "text")
		g.expect("AXStaticText", "Plain text", nil)
		g.expect("AXStaticText", "Text grid\nSecond row", nil)
		g.expect("AXTextField", "Display name", map[string]string{"AXValue": "Ada"})
		g.expect("AXTextField", "Read only example", map[string]string{"AXEnabled": "false", "AXValue": "Cannot edit"})
		g.setValue("AXTextField", "Display name", "Grace")
		g.expect("AXTextField", "Display name", map[string]string{"AXValue": "Grace"})
		g.expect("AXStaticText", "Name: Grace", nil)
		g.setValue("AXTextField", "Notes", "Updated\nnotes")
		g.expect("AXTextField", "Notes", map[string]string{"AXValue": "Updated\nnotes"})
		g.expect("AXStaticText", "Notes: Updated\nnotes", nil)
		g.act("AXLink", "Example link", "AXPress")
		g.expect("AXStaticText", "Link activated", nil)
		g.setValue("AXTextField", "Demo password", "updated-password")
		tree := g.wait("password field", func(tree *axSnapshot) bool {
			node := tree.find("AXTextField", "Demo password")
			return node != nil && node.attributes["AXValue"] != ""
		})
		defer tree.close()
		for _, node := range tree.nodes {
			for attribute, value := range node.attributes {
				if strings.Contains(value, "gallery-password") || strings.Contains(value, "updated-password") {
					t.Errorf("password contents exposed in %s on %s", attribute, node.attributes["AXRole"])
				}
			}
		}
		g.setValue("AXTextField", "Required field", "Valid")
		g.act("AXButton", "Submit", "AXPress")
		g.expect("AXStaticText", "Form submitted", nil)
	})

	t.Run("collections", func(t *testing.T) {
		g := startGallery(t, "collections")
		g.expect("AXList", "", nil)
		g.expect("AXOutline", "", nil)
		g.expect("AXTable", "", nil)
		g.expect("AXStaticText", "Cell 1,1", nil)
		g.act("AXRow", "List row 01", "AXPress")
		g.expect("AXStaticText", "List selected: 1", nil)
		g.expect("AXRow", "List row 01", map[string]string{"AXSelected": "true"})
		g.act("AXRow", "Branch", "AXPress")
		g.expect("AXRow", "Branch", map[string]string{"AXExpanded": "true"})
		g.act("AXRow", "Leaf One", "AXPress")
		g.expect("AXStaticText", "Tree selected: Leaf One", nil)
		g.act("AXRow", "Grid item 01", "AXPress")
		g.expect("AXStaticText", "Grid selected: 1", nil)
		g.act("AXButton", "Scroll list to end", "AXPress")
		g.expect("AXRow", "List row 30", nil)
	})

	t.Run("containers", func(t *testing.T) {
		g := startGallery(t, "containers")
		g.expect("AXTabGroup", "", nil)
		g.expect("AXRadioButton", "First tab", map[string]string{"AXSubrole": "AXTabButton", "AXSelected": "true"})
		g.expect("AXStaticText", "First tab content", nil)
		g.absent("AXStaticText", "Second tab content")
		g.act("AXRadioButton", "Second tab", "AXPress")
		g.expect("AXRadioButton", "Second tab", map[string]string{"AXSelected": "true"})
		g.expect("AXStaticText", "Second tab content", nil)
		g.expect("AXStaticText", "Tab: Second tab", nil)
		g.absent("AXStaticText", "First tab content")
		g.act("AXRadioButton", "Document B", "AXPress")
		g.expect("AXStaticText", "Document B content", nil)
		g.absent("AXStaticText", "Document A content")
		g.act("AXButton", "Expand details", "AXPress")
		g.expect("AXStaticText", "Accordion details", nil)
	})

	t.Run("dynamic", func(t *testing.T) {
		g := startGallery(t, "dynamic")
		g.absent("AXStaticText", "Extra content")
		g.act("AXButton", "Toggle extra content", "AXPress")
		g.expect("AXStaticText", "Extra content", nil)
		g.act("AXButton", "Toggle extra content", "AXPress")
		g.expect("AXStaticText", "Extra content hidden", nil)
		g.absent("AXStaticText", "Extra content")
		g.expect("AXButton", "Conditional action", map[string]string{"AXEnabled": "false"})
		g.act("AXCheckBox", "Enable action", "AXPress")
		g.expect("AXButton", "Conditional action", map[string]string{"AXEnabled": "true"})
		g.act("AXButton", "Conditional action", "AXPress")
		g.expect("AXStaticText", "Conditional action activated", nil)
		g.act("AXButton", "Open dialog", "AXPress")
		g.expect("AXStaticText", "Dialog content", nil)
		g.act("AXButton", "Dismiss dialog", "AXPress")
		g.absent("AXStaticText", "Dialog content")
		g.act("AXButton", "Open popup", "AXPress")
		g.expect("AXStaticText", "Popup content", nil)
		g.act("AXButton", "Dismiss popup", "AXPress")
		g.absent("AXStaticText", "Popup content")
	})
}
