# got [![Go](https://github.com/fedeztk/got/actions/workflows/go.yaml/badge.svg)](https://github.com/fedeztk/got/tree/master/.github/workflows/go.yml)

## Table of Contents

[Demo](#orgab62fc1) -
[Usage](#orgfa2aa9c) -
[Features](#org26baa6c) -
[Testing](#org2744438)

go-translation (shortly `got`), a simple translator build on top of [simplytranslate's APIs](https://codeberg.org/SimpleWeb/SimplyTranslate-Web/src/branch/master/api.md) and the awesome [bubbletea](https://github.com/charmbracelet/bubbletea) tui library.

> Disclaimer: this is my absolute first project in golang, so bugs and clunky code are expected&#x2026;

> The project is still a work-in-progress, breaking changes and heavy refactoring may happen


<a id="orgab62fc1"></a>

## Demo

https://user-images.githubusercontent.com/58485208/165116687-95017b2b-1e1b-4c82-b5ac-dc2cda87afe0.mp4

<a id="orgfa2aa9c"></a>

# Usage

- Install `got`: 

With the `go` tool:
```sh
go install github.com/fedeztk/got/cmd/got
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

- Copy [the sample config](https://github.com/fedeztk/got/blob/master/config.yml) file under ~/.config/got/ as config.yml **or** let the program generate one for you at the first run
- Run it:
```sh
got
```

<a id="org26baa6c"></a>

# Features

-   Interact with google translate easily via the terminal, no need to open a browser!
-   Clean interface with 3 tabs, switch between them with tab/shift-tab:
    -   **text input**: input the sentence you want to translate, press enter to translate
    -   **language selection**: choose between 108 languages, select source language with **s**, target with **t**. Press **h** to show the help menu
    -   **translation**: pager that shows the result of translation
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

