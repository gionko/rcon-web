# rcon-web

[![Build Status](https://travis-ci.org/dieselburner/rcon-web.svg)](https://travis-ci.org/dieselburner/rcon-web)
[![License](https://img.shields.io/github/license/dieselburner/rcon-web.svg)](https://github.com/dieselburner/rcon-web/blob/master/LICENSE.md)

Web interface for Source engine based game server administration

<!-- TOC -->
- [Overview](#overview)
- [Screenshots](#screenshots)
- [Downloads](#downloads)
- [Installation](#installation)
- [Configuration](#configuration)

## Overview

`rcon-web` is a web administration interface for Source engine based game servers. Interface is made mobile-first and is totally responsive.

## Screenshots

All screenshots are made from an actual phone, Android 8.0 on Samsung A5 2017.

Login screen:

<img src="https://raw.githubusercontent.com/dieselburner/rcon-web/master/documentation/images/login.png" width="256">

Main screen:

<p float="left">
  <img src="https://raw.githubusercontent.com/dieselburner/rcon-web/master/documentation/images/main-players.png" width="256">
  <img src="https://raw.githubusercontent.com/dieselburner/rcon-web/master/documentation/images/main-maps.png" width="256">
  <img src="https://raw.githubusercontent.com/dieselburner/rcon-web/master/documentation/images/main-config.png" width="256">
</p>

Actions:

<p float="left">
  <img src="https://raw.githubusercontent.com/dieselburner/rcon-web/master/documentation/images/click-user.png" width="256">
  <img src="https://raw.githubusercontent.com/dieselburner/rcon-web/master/documentation/images/click-map.png" width="256">
  <img src="https://raw.githubusercontent.com/dieselburner/rcon-web/master/documentation/images/click-bot-count.png" width="256">
  <img src="https://raw.githubusercontent.com/dieselburner/rcon-web/master/documentation/images/click-bot-damage.png" width="256">
  <img src="https://raw.githubusercontent.com/dieselburner/rcon-web/master/documentation/images/click-bot-difficulty.png" width="256">
</p>

## Downloads

Latest release is available [here](https://github.com/dieselburner/rcon-web/releases/latest).

## Installation

Installation instructions are located in a separate documentation file, [here](https://github.com/dieselburner/rcon-web/blob/master/INSTALL.md). Configuring `rcon-web` for external access is not covered here, since this often is quite individual topic. Hopefully, Internet is full of information about port forwarding, CNAME's, dynamic DNS'es, etc.

## Configuration

There are two configuration files that needs to be configured.

`uwsgi.ini` contains web application specific configuration, and it might need to be adjusted for your environment, like paths or user/group id's.

`rcon-web.conf` contains game server configuration. Game server address and password goes here. Please note, your login password is the same as in this configuration file.
Additionally, it is possible to fine-tune couple of options via this configuration file, in particular, prefix for user ban/kick messages, and provide readable map names.
