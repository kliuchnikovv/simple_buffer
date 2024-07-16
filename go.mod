module github.com/kliuchnikovv/simple_buffer

go 1.22

toolchain go1.22.3

require github.com/kliuchnikovv/edicode v0.0.0-20220328183000-59ed439bb22c

require (
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/wailsapp/wails/v2 v2.1.0 // indirect
	golang.org/x/exp/shiny v0.0.0-20221023144134-a1e5550cf13e // indirect
	golang.org/x/mobile v0.0.0-20210716004757-34ab1303b554 // indirect
)

require (
	github.com/fsnotify/fsnotify v1.5.4
	github.com/leaanthony/slicer v1.6.0 // indirect
	golang.design/x/clipboard v0.6.2
	golang.org/x/image v0.0.0-20211028202545-6944b10bf410 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
)

replace github.com/kliuchnikovv/edicode => ../edicode
