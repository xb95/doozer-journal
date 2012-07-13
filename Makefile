unexport GIT_DIR # Needs to be unset for a clean build of external dependencies

LOCAL_GOPATH=${PWD}/.go_path
PKG=github.com/soundcloud/doozer-journal

bundle: fmt
	mkdir -p $$(dirname $(LOCAL_GOPATH)/src/$(PKG))
	ln -sfn $${PWD} $(LOCAL_GOPATH)/src/$(PKG)
	cd $(LOCAL_GOPATH)/src/$(PKG);\
		GOPATH=$(LOCAL_GOPATH) go get -v -d ./...;\
		GOPATH=$(LOCAL_GOPATH) go build

clean:
	go clean
	git clean -fdx

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
