.PHONY: all clean arm arm64 386 amd64 linux window windowsc darwin ios 

all: android

android: main.go send.go recv.go settings.go theme.go about.go AndroidManifest.xml fdroid-build.sh
	ANDROID_HOME=~/android bash fdroid-build.sh test

clean:
	go clean
	rm crocgui.apk

arm: main.go send.go recv.go settings.go theme.go about.go AndroidManifest.xml
	fyne package -os android/arm --release

arm64: main.go send.go recv.go settings.go theme.go about.go AndroidManifest.xml
	fyne package -os android/arm64 --release

386: main.go send.go recv.go settings.go theme.go about.go AndroidManifest.xml
	fyne package -os android/386 --release

amd64: main.go send.go recv.go settings.go theme.go about.go AndroidManifest.xml
	fyne package -os android/amd64 --release

linux: main.go send.go recv.go settings.go theme.go about.go AndroidManifest.xml
	fyne package -os linux --release

windows: main.go send.go recv.go settings.go theme.go about.go AndroidManifest.xml
	#sudo apt-get install gcc-mingw-w64-x86-64
	CC=x86_64-w64-mingw32-gcc fyne package -os windows --release -tags=opengl

windowsc: main.go send.go recv.go settings.go theme.go about.go AndroidManifest.xml
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOFLAGS=-ldflags=-s go build -tags=opengl

darwin: main.go send.go recv.go settings.go theme.go about.go AndroidManifest.xml
	fyne package -os darwin --release

ios: main.go send.go recv.go settings.go theme.go about.go AndroidManifest.xml
	fyne package -os ios --release
