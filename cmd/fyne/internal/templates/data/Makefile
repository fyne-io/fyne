# If PREFIX isn't provided, we check for $(DESTDIR)/usr/local and use that if it exists.
# Otherwice we fall back to using /usr.

LOCAL != test -d $(DESTDIR)/usr/local && echo -n "/local" || echo -n ""
LOCAL ?= $(shell test -d $(DESTDIR)/usr/local && echo "/local" || echo "")
PREFIX ?= /usr$(LOCAL)

Name := {{ .Name | printf "%q" }}
Exec := {{ .Exec | printf "%q" }}
Icon := {{ .Icon | printf "%q" }}

default:
	# User install
	# Run "make user-install" to install in ~/.local/
	# Run "make user-uninstall" to uninstall from ~/.local/
	#
	# System install
	# Run "sudo make install" to install the application.
	# Run "sudo make uninstall" to uninstall the application.

install:
	install -Dm00644 usr/{{.Local}}share/applications/$(Name).desktop $(DESTDIR)$(PREFIX)/share/applications/$(Name).desktop
	install -Dm00755 usr/{{.Local}}bin/$(Exec) $(DESTDIR)$(PREFIX)/bin/$(Exec)
	install -Dm00644 usr/{{.Local}}share/pixmaps/$(Icon) $(DESTDIR)$(PREFIX)/share/pixmaps/$(Icon)
uninstall:
	-rm $(DESTDIR)$(PREFIX)/share/applications/$(Name).desktop
	-rm $(DESTDIR)$(PREFIX)/bin/$(Exec)
	-rm $(DESTDIR)$(PREFIX)/share/pixmaps/$(Icon)

user-install:
	install -Dm00644 usr/{{.Local}}share/applications/$(Name).desktop $(DESTDIR)$(HOME)/.local/share/applications/$(Name).desktop
	install -Dm00755 usr/{{.Local}}bin/$(Exec) $(DESTDIR)$(HOME)/.local/bin/$(Exec)
	install -Dm00644 usr/{{.Local}}share/pixmaps/$(Icon) $(DESTDIR)$(HOME)/.local/share/icons/$(Icon)
	sed -i -e "s,Exec=.*,Exec=$(DESTDIR)$(HOME)/.local/bin/$(Exec),g" $(DESTDIR)$(HOME)/.local/share/applications/$(Name).desktop

user-uninstall:
	-rm $(DESTDIR)$(HOME)/.local/share/applications/$(Name).desktop
	-rm $(DESTDIR)$(HOME)/.local/bin/$(Exec)
	-rm $(DESTDIR)$(HOME)/.local/share/icons/$(Icon)
