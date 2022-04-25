# got [![Go](https://github.com/fedeztk/got/actions/workflows/go.yaml/badge.svg)](https://github.com/fedeztk/got/tree/master/.github/workflows/go.yml)

## Table of Contents

[Demo](#orgab62fc1)
[Usage](#orgfa2aa9c)
[Features](#org26baa6c)
[Testing](#org2744438)

got(ranslation), a simple translator build on top of [translate-shell](https://github.com/soimort/translate-shell) and the awesome [bubbletea](https://github.com/charmbracelet/bubbletea) tui library.

> Disclaimer: this is my absolute first project in golang, so bugs and clunky code are expected&#x2026;


<a id="orgab62fc1"></a>

## Demo

https://user-images.githubusercontent.com/58485208/165116687-95017b2b-1e1b-4c82-b5ac-dc2cda87afe0.mp4

<a id="orgfa2aa9c"></a>

# Usage

-   install dependencies:

    pacman -S go translate-shell

-   install `got`:

    go get -u https://github.com/fedeztk/got

-   copy [the sample config](https://github.com/fedeztk/got/blob/master/config.yml) file under ~/.config/got/ as config.yml
-   run it:

    got


<a id="org26baa6c"></a>

# Features

-   interact with google translate through translate-shell easily via the terminal
-   clean interface with 3 tabs, switch between them with tab/shift-tab:
    -   **text input**: input the sentence you want to translate, press enter to translate
    -   **language selection**: choose between 121 languages, select source language with **s**, target with **t**
    -   **translation**: pager that shows the result of translation
-   quit anytime with esc or ctrl-c
-   automatically remembers the last languages used


<a id="org2744438"></a>

# Testing

Development is done through `docker`, build the container with:

    docker build -t got .

Check that the build went fine:

    docker images | grep got

Test it with:

    docker run -it -e "TERM=xterm-256color" got

