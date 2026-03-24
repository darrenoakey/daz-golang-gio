![](banner.jpg)

# daz-golang-gio

Welcome! This is a Go library that makes it easier to build beautiful, cross-platform desktop applications using the [Gio](https://gioui.org) UI toolkit. Whether you're new to GUI programming in Go or just looking for a smoother experience, this library is here to help you get up and running quickly.

## What Is It?

Gio is a powerful framework for building native graphical apps in Go — apps with real windows, buttons, text, and all the things users expect from desktop software. This library sits on top of Gio and gives you a friendlier, more approachable way to work with it.

Think of it as a helpful companion that takes care of some of the repetitive setup so you can focus on building your actual app.

## Getting Started

### Prerequisites

Make sure you have Go installed (version 1.23 or newer). You can check by running:

```bash
go version
```

### Add It to Your Project

In your Go project, run:

```bash
go get github.com/darrenoakey/daz-golang-gio
```

Then import it in your Go code:

```go
import "github.com/darrenoakey/daz-golang-gio"
```

### Try the Example

The fastest way to see it in action is to clone the repository and run the built-in example:

```bash
git clone https://github.com/darrenoakey/daz-golang-gio
cd daz-golang-gio
./run example
```

A window should appear on your screen — that's your first Gio app running!

## Working with the Project

If you're contributing to or exploring the library itself, the `run` script is your main helper. Here's what each command does:

| Command | What it does |
|---|---|
| `./run build` | Compiles the library to check everything is in order |
| `./run test` | Runs all the tests |
| `./run lint` | Checks code formatting and quality |
| `./run check` | Runs build, lint, and tests all together — the full quality check |
| `./run example` | Launches the example application |
| `./run deploy` | Publishes a new release (for maintainers) |

### Running a Specific Test

You can run tests for a specific package by passing a path:

```bash
./run test ./yourpackage/...
```

## Tips & Tricks

- **Start with the example.** The `./example/` folder is the best place to understand how to structure your own app. Copy it as a starting point!

- **Run `./run check` before committing.** It catches formatting issues, build errors, and failing tests all in one go — saves you from surprises later.

- **Gio is immediate mode.** If you're used to other UI frameworks, Gio works a little differently — the UI is redrawn each frame rather than being built once. Once it clicks, it's very satisfying to work with.

- **Cross-platform by default.** Apps built with this library run on Linux, macOS, and Windows without any extra configuration.

- **Keep your dependencies tidy.** Run `go mod tidy` occasionally to make sure your `go.mod` and `go.sum` files stay clean.

## Learn More

- [Gio documentation](https://gioui.org) — the upstream framework this library builds on
- [pkg.go.dev](https://pkg.go.dev/github.com/darrenoakey/daz-golang-gio) — API reference for this library

Happy building! 🎉