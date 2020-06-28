GO := go


build: */*.go */*/*.go
	${GO} build ./...

test:
	${GO} test ./...

clean-dist:
	rm -rf dist

dist: dist-pi dist-linux dist-osx dist-windows

dist-pi: dist/pi/reflector/cwc-reflector dist/pi/station/cwc-station

dist/pi/reflector/cwc-reflector:
	mkdir -p dist/pi/reflector
	GOOS=linux GOARCH=arm GOARM=5 ${GO} build -o dist/pi/reflector/cwc-reflector ./cwc-reflector

dist/pi/station/cwc-station:
	mkdir -p dist/pi/station
	GOOS=linux GOARCH=arm GOARM=5 ${GO} build -o dist/pi/station/cwc-station ./cwc-station

dist-linux: dist/linux/reflector/cwc-reflector.tgz dist/linux/station/cwc-station

dist/linux/reflector/cwc-reflector.tgz:
	mkdir -p dist/linux/reflector
	cp -r cwc-reflector/web dist/linux/reflector/web
	cp -r cwc-reflector/config dist/linux/reflector/config
	GOOS=linux GOARCH=amd64 ${GO} build -o dist/linux/reflector/cwc-reflector ./cwc-reflector
	cd dist/linux && tar cfz reflector/cwc-reflector.tgz reflector/cwc-reflector reflector/config reflector/web

dist/linux/station/cwc-station:
	mkdir -p dist/linux/station
	GOOS=linux GOARCH=amd64 ${GO} build -o dist/linux/station/cwc-station ./cwc-station

dist-osx: dist/osx/reflector/cwc-reflector dist/osx/station/cwc-station

dist/osx/reflector/cwc-reflector:
	mkdir -p dist/osx/reflector
	cp -r cwc-reflector/web dist/osx/reflector/web
	cp -r cwc-reflector/config dist/osx/reflector/config
	GOOS=darwin GOARCH=amd64 ${GO} build -o dist/osx/reflector/cwc-reflector ./cwc-reflector
	cd dist/osx && tar cfz reflector/cwc-reflector.tgz reflector/cwc-reflector reflector/config reflector/web

dist/osx/station/cwc-station:
	mkdir -p dist/osx/station
	GOOS=darwin GOARCH=amd64 ${GO} build -o dist/osx/station/cwc-reflector ./cwc-station

dist-windows: dist/windows/reflector/cwc-reflector.zip dist/windows/cwc-station

dist/windows/reflector/cwc-reflector.zip:
	mkdir -p dist/windows/reflector
	cp -r cwc-reflector/web dist/windows/reflector/web
	cp -r cwc-reflector/config dist/windows/reflector/config
	GOOS=windows GOARCH=386 ${GO} build -o dist/windows/reflector/cwc-reflector.exe ./cwc-reflector
	cd dist/windows && zip -r reflector/cwc-reflector.zip reflector/cwc-reflector.exe reflector/config reflector/web

dist/windows/cwc-station:
	echo "No windows station build due to rpio needing mmap"
#	mkdir -p dist/windows/station
#	GOOS=windows GOARCH=386 ${GO} build -o dist/windows/station/cwc-station.exe ./cwc-station


