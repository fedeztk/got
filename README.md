# got
[![GO](https://github.com/fedeztk/got/actions/workflows/go.yaml/badge.svg)](https://github.com/fedeztk/got/tree/master/.github/workflows/go.yml) [![GHCR](https://github.com/fedeztk/got/actions/workflows/deploy.yaml/badge.svg)](https://github.com/fedeztk/got/tree/release/.github/workflows/deploy.yml) [![AUR](https://img.shields.io/aur/version/go-translation-git?logo=archlinux)](https://aur.archlinux.org/packages/go-translation-git) [![Go Report Card](https://goreportcard.com/badge/github.com/fedeztk/got)](https://goreportcard.com/report/github.com/fedeztk/got)

## Table of Contents

[Demo](#orgab62fc1) -
[Usage](#orgfa2aa9c) -
[Features](#org26baa6c) -
[Testing](#org2744438)

go-translation (shortly `got`), a simple translator and text-to-speech app built on top of [simplytranslate](https://codeberg.org/SimpleWeb/SimplyTranslate-Web/src/branch/master/api.md) and [lingvatranslate](https://github.com/thedaviddelta/lingva-translate) APIs. The interface is made with the awesome [bubbletea](https://github.com/charmbracelet/bubbletea) tui library.

> :warning: simplytranslate is currently down. I still kept it as a valid backend in the hope of a future comeback, by default `got` uses now lingvatranslate

Screenshots [here](#org26baa6c)

> Disclaimer: this is my absolute first project in golang, so bugs and clunky code are expected&#x2026;

> The project is still a work-in-progress, breaking changes and heavy refactoring may happen


<a id="orgab62fc1"></a>

## Demo


https://user-images.githubusercontent.com/58485208/166154625-7c5556bd-74aa-4425-a046-160ba793b792.mp4


<a id="orgfa2aa9c"></a>

# Usage

- Install `got`: 

With the `go` tool:
```sh
go install github.com/fedeztk/got/cmd/got@latest
```
**Or** from source:
```sh
# clone the repo
git clone https://github.com/fedeztk/got.git
# install manually 
make install
```
In both cases make sure that you have the go `bin` directory in your path:
```sh
export PATH="$HOME/go/bin:$PATH"
```
If you are an Arch user there is also an AUR package available:
```sh
paru -S go-translation-git
```
- Copy [the sample config](https://github.com/fedeztk/got/blob/master/config.yml) file under ~/.config/got/ as config.yml **or** let the program generate one for you at the first run
- Run it interactively:
```sh
got            # use last used engine, default is google
got -e reverso # change engine to reverso
```
-  Or in oneshot mode:
```sh
got -o -s en -t it "Hello World"          # use default (google)
got -o -e libre -s en -t it "Hello World" # use libre-translate
```
For more information check the help (`got -h`)
<a id="org26baa6c"></a>

# Features

-   Interact with various translation engines easily via the terminal, no need to open a browser!
-   Clean interface with 3 tabs, switch between them with tab/shift-tab:
	-   **text input**: input the sentence you want to translate, press **enter** to translate
![image](https://user-images.githubusercontent.com/58485208/173687247-2a1ad240-44f8-46ff-b8de-c55b3eccc4c4.png)
	-   **language selection**: choose between 108 languages, select source language with **s**, target with **t** and **i** to invert the target with the source. Press **?** to show the full help menu
	![image](https://user-images.githubusercontent.com/58485208/173687797-6325ccc9-5745-43af-b9a8-35b97bd94675.png)
	Full help:
	![image](https://user-images.githubusercontent.com/58485208/173687516-33d48c4c-206a-4b85-9678-ee6684ba71e4.png)
	-   **translation**: pager that shows the result of translation. Copy translation with **y**, listen the translation with **p**
	![image](https://user-images.githubusercontent.com/58485208/173687675-5d073c2c-428a-4a27-9cb2-4b0c803a8a5e.png)
- **engines** (only available with simplytranslate backend): choose between google, libre-translate, reverso and iciba (deepl is not working yet)
-   quit anytime with **esc** or **ctrl-c**
-   automatically remembers the last languages used


<a id="org2744438"></a>

# Testing

Development is done through `docker`, build the container with:

    make docker-build

Check that the build went fine:

    docker images | grep got

Test it with:

    make docker-run

Pre-built Docker image available [here](https://github.com/fedeztk/got/pkgs/container/got)

