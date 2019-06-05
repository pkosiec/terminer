![Terminer](./assets/logo.png)

# Terminer

[![Build Status](https://travis-ci.com/pkosiec/terminer.svg?branch=master)](https://travis-ci.com/pkosiec/terminer) [![Go Report Card](https://goreportcard.com/badge/github.com/pkosiec/terminer)](https://goreportcard.com/report/github.com/pkosiec/terminer) [![codecov](https://codecov.io/gh/pkosiec/terminer/branch/master/graph/badge.svg)](https://codecov.io/gh/pkosiec/terminer) ![GitHub release](https://img.shields.io/github/release/pkosiec/terminer.svg)

Upgrade your terminal experience with a single command.

> **:warning: WARNING**: This project is in a very early stage. Use it at your own risk.

Terminer is an cross-platform installer for terminal presets. Install Fish or ZSH shell packed with useful plugins and sleek prompts. Use one of starter recipes or make yours.

![Screenshot](./assets/screenshot.png)

## Table of contents

- [Motivation](#motivation)
- [Installation](#installation)
- [Usage](#usage)
  - [Quick start](#quick-start)
  - [Available recipes](#available-recipes)
  - [What is a Recipe](#what-is-a-recipe)
- [Available commands](#available-commands)
  - [`install`](#install)
  - [`rollback`](#rollback)
  - [`version`](#version)

## Motivation

Command line interface (CLI) is a great way to access various operating system functions. It allows you to automate time-consuming tasks with a single command. Also, not all features are available with graphical user interface (GUI).

The command line interface is available through a program called terminal emulator, which launches shell. Shell is a program which processes user input and returns output.

Usually, default shell setup is bare-bone. Luckily, there is a way to upgrade it: install different shell along with useful plugins. They can introduce features such as autocompletion or syntax highlighting. You can also install prompts, which can display useful pieces of information, like current git branch. All these additions heavily improve user experience and save your time.

There is one downside of the shell experience customization. It's a time-consuming task. Installing different shell is easy, but to configure it well, you can easily spend hours for searching useful plugins, prompts, themes and fonts.

Not anymore. Use Terminer. Bootstrap your complete shell configuration in a moment.

## Installation

To install this tool into `/usr/local/bin` directory, run the following command:

```bash
curl -sfL https://raw.githubusercontent.com/pkosiec/terminer/master/install.sh | sh -s -- -b /usr/local/bin
```


> **Note:** Make sure that the directory in your `$PATH`. If not, run:
```bash
export PATH=$PATH:/usr/local/bin
```

Include the line above in you `~/.bashrc` file or other configuration file for your current shell.

## Usage

Terminer operates on recipes, which consist of shell commands.
The most essential commands are `install` and `rollback`.

### Quick start

To install a recipe from official repository, run:

```bash
terminer install [recipe name]
```

To rollback a recipe from official repository, run:

```bash
terminer rollback [recipe name]
```

### Available recipes

To see all maintained recipes, navigate to the [**`recipes`**](./recipes) directory.

### What is a Recipe

Recipe is a YAML or JSON file with shell commands put in a proper order. Recipe consists of stages, which contain steps. Every step is a different shell command.

This is an example recipe, which just prints messages for all steps in all stages - not only during install, but also for rollback operation:

```yaml
os: darwin
metadata:
  name: Recipe
  description: Recipe Description

stages:
  - metadata:
      name: Stage 1
      description: Stage 1 description
      url: https://stage1.example.com
    steps:
      - metadata:
          name: Step 1
          url: https://step1.stage1.example.com
        execute:
          run:
            - echo "Step 1 of Stage 1"
        rollback:
          run:
            - echo "Rollback of Step 1 of Stage 1"
      - metadata:
          name: Step 2
          url: https://step2.stage1.example.com
        execute:
          run:
            - echo "Step 2 of Stage 1"
        rollback:
          run:
            - echo "Rollback of Step 2 of Stage 1"
  - metadata:
      name: Stage 2
      description: Stage 2 description
      url: https://stage2.example.com
    steps:
      - metadata:
          name: Step 1
          url: https://step1.stage2.example.com
        execute:
          run:
            - echo "Step 1 of Stage 1"
          shell: sh
        rollback:
          run:
            - echo "Rollback of Step 1 of Stage 2"
```

## Available commands

The following section describes all available commands in Terminer CLI.

### `install`

Install command installs a recipe from the official recipe repository. You can use additional flags to install a recipe from a local or remote file.

**Usage**

```bash
terminer install [recipe name]
```

**Flags**

```
-f, --filepath string   Recipe file path
-h, --help              help for install
-u, --url string        Recipe URL
```

**Examples**

```
terminer install zsh-starter
terminer install -f ./recipe.yaml
terminer install --file /Users/sample-user/recipe.yml
terminer install -u https://example.com/recipe.yaml
terminer install --url http://foo.bar/recipe.yml
```

### `rollback`

Rollback command uninstalls a recipe from the official recipe repository. You can use additional flags to rollback a recipe from a local or remote file.

**Usage**

```bash
terminer rollback [recipe name]
```

**Flags**

```
-f, --filepath string   Recipe file path
-h, --help              help for install
-u, --url string        Recipe URL
```

**Examples**

```bash
terminer rollback zsh-starter
terminer rollback -f ./recipe.yaml
terminer rollback --file /Users/sample-user/recipe.yml
terminer rollback -u https://example.com/recipe.yaml
terminer rollback --url http://foo.bar/recipe.yml
```

### `version`

Prints the application version

**Usage**

```bash
terminer version
```
