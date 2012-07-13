# Makefile to include for Go projects
#
#   Installs dependencies in local GOPATH (.gopath),
#   and binary in $PWD/bin.
#
# usage:
#
#   PKG=github.com/org/project-name
#
#   include go.mk
#
#   build: bundle
#   clean: bundle-clean
#
LOCAL_GOPATH=${PWD}/.gopath
PKG_GOPATH=$(LOCAL_GOPATH)/src/$(PKG)

# this is needed for bazooka builds
unexport GIT_DIR

bundle: bundle-check
	mkdir -p bin $$(dirname $(PKG_GOPATH))
	ln -sfn $${PWD} $(PKG_GOPATH)
	cd $(PKG_GOPATH);\
	  GOPATH=$(LOCAL_GOPATH) go get -v -d;\
	  GOPATH=$(LOCAL_GOPATH) GOBIN=$${PWD}/bin go install -v $(PKG)

bundle-check:
	@test -n "$(PKG)" || { echo "PKG variable must be set to package name" && exit 1; }

bundle-clean:
	rm -rf $(LOCAL_GOPATH)
