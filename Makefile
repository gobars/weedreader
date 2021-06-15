all: init install

app=$(notdir $(shell pwd))

init:
	export GOPROXY=https://goproxy.cn

install: init
	go install -ldflags="-s -w" ./...
	ls -lh ~/go/bin/${app}

linux: init
	GOOS=linux GOARCH=amd64 go install -ldflags="-s -w" ./...


upx:
	ls -lh ~/go/bin/${app}
	upx ~/go/bin/${app}
	ls -lh ~/go/bin/${app}
	ls -lh ~/go/bin/linux_amd64/${app}
	upx ~/go/bin/linux_amd64/${app}
	ls -lh ~/go/bin/linux_amd64/${app}

clean:
	rm coverage.out

targz:
	cd .. && tar czvf weedreader.tar.gz --exclude .git --exclude .idea weedreader

run:
	GOLOG_STDOUT=true weedreader -c initassets/weedreader.yml