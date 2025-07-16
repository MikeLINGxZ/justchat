#!/bin/bash

flutter pub run build_runner build
flutter pub run ffigen --config lib/utils/ffi/example_ffi/config.yml
flutter pub run ffigen --config lib/utils/ffi/example_ffi_chat/config.yml