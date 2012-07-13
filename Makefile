unexport GIT_DIR # Needs to be unset for a clean build of external dependencies

include go.mk
PKG=github.com/soundcloud/doozer-journal


clean:
	go clean

# The default bazooka target
build: bundle package

fmt:
	go fmt ./...

########## packaging
update_version:
		sed -i -e "s/const VERSION .*/const VERSION = \"$$(cat VERSION)\"/" main.go

FPM_EXECUTABLE:=$$(dirname $$(dirname $$(gem which fpm)))/bin/fpm
FPM_ARGS=-t deb -m 'doozer-journal authors (see page), Alexander Simmerl <alx@soundcloud.com> (packaging)' --url http://github.com/soundcloud/doozer-journal -s dir
FAKEROOT=fakeroot
RELEASE=$$(cat .release 2>/dev/null || echo "0")

package:
	- mkdir -p $(FAKEROOT)/usr/bin
	cp doozer-journal $(FAKEROOT)/usr/bin
	- mkdir -p $(FAKEROOT)/var/lib/doozer-journal/
	rm -rf *.deb

	$(FPM_EXECUTABLE) -n "doozer-journal" \
		-C $(FAKEROOT) \
		--description "Snapshots, mutation journaling and recovery of doozerd coordinator state." \
		$(FPM_ARGS) -t deb -v $$(cat VERSION) --iteration $(RELEASE) .;
