#!/usr/bin/env bash

export CC=x86_64-w64-mingw32-gcc
fyne package -os windows -tags=opengl --release
