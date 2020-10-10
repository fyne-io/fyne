# If PREFIX isn't provided, we check for /usr/local and use that if it exists.
# Otherwice we fall back to using /usr.

LOCAL != test -d /usr/local && echo -n "/local" || echo -n ""
LOCAL ?= $(shell test -d /usr/local && echo "/local" || echo "")
PREFIX ?= /usr$(LOCAL)

default:
	# Run "sudo make install" to install the application.

install:
	install -Dm00644 usr/{{.Local}}share/applications/{{.Name}}.desktop $(DESTDIR)$(PREFIX)/share/applications/{{.Name}}.desktop
	install -Dm00755 usr/{{.Local}}bin/{{.Exec}} $(DESTDIR)$(PREFIX)/bin/{{.Exec}}
	install -Dm00644 usr/{{.Local}}share/pixmaps/{{.Icon}} $(DESTDIR)$(PREFIX)/share/pixmaps/{{.Icon}}
