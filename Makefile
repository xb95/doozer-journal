unexport GIT_DIR # Needs to be unset for a clean build of external dependencies

LOCAL_GOPATH=${PWD}/.go_path
PKG=github.com/soundcloud/doozer-journal

clean:
	go clean
	git clean -fdx

# The default bazooka target
build:
	mkdir -p $$(dirname $(LOCAL_GOPATH)/src/$(PKG))
	ln -sfn $${PWD} $(LOCAL_GOPATH)/src/$(PKG)
	cd $(LOCAL_GOPATH)/src/$(PKG);\
		GOPATH=$(LOCAL_GOPATH) go get -v -d ./...;\
		GOPATH=$(LOCAL_GOPATH) go build

fmt:
	go fmt ./...
