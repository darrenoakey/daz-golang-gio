![](banner.jpg)

# daz-golang-gio

So you want to build a real desktop app in Go — one with windows, buttons, text, and all the things people expect from native software. You've come to the right place.

This library makes it easier to get started with [Gio](https://gioui.org), a fantastic framework for building beautiful, cross-platform desktop applications in Go. Think of this library as a friendly layer on top of Gio that handles some of the repetitive setup so you can focus on your actual app, not the boilerplate.

The result? You spend less time wrestling with setup and more time building something you're proud of. Apps you create with this library run on **Linux, macOS, and Windows** — no extra configuration needed.

---

## What You'll Need

Before anything else, make sure you have **Go 1.23 or newer** installed. Not sure which version you have? Run this in your terminal:

```bash
go version
```

If you need to install or update Go, head over to [go.dev](https://go.dev/dl/) and grab the latest version.

---

## Getting Started

### Add the Library to Your Project

Inside your Go project folder, run:

```bash
go get github.com/darrenoakey/daz-golang-gio
```

Then import it in your code:

```go
import "github.com/darrenoakey/daz-golang-gio"
```

That's it — you're ready to start building!

### See It in Action First

The quickest way to understand what this library can do is to run the built-in example app. Clone the repository and try it out:

```bash
git clone https://github.com/darrenoakey/daz-golang-gio
cd daz-golang-gio
./run example
```

A window will appear on your screen. That's a complete, working desktop application — and it's a great starting point for your own project.

---

## A Guide to Everything You Can Do

All the common tasks in this project are handled by a single helper script called `run`. Here's a friendly walkthrough of everything it can do.

### Building

```bash
./run build
```

This compiles the library to make sure everything is in good shape. Think of it as a quick sanity check — if it builds without errors, you're on solid ground.

### Running Tests

```bash
./run test
```

Runs the full test suite and shows you the results. If you want to test just one specific part of the project, you can pass a path:

```bash
./run test ./yourpackage/...
```

Test output is also saved automatically to `output/testing/last.log`, so you can look back at it any time.

### Checking Code Quality

```bash
./run lint
```

This checks that your code is correctly formatted and follows Go best practices. It uses `gofmt` and `go vet` under the hood, and will also run `golangci-lint` if you have it installed.

### The Full Quality Check

```bash
./run check
```

This is your one-stop confidence booster — it runs the build, the linter, and all the tests in sequence. If everything passes, you're good to go. This is the command to run before committing or sharing your work.

### Running the Example App

```bash
./run example
```

Launches the included example application. It's a great reference whenever you want to see how something is done.

### Deploying a New Release

```bash
./run deploy
```

For project maintainers, this command handles the full release process: runs quality checks, bumps the version number, publishes to GitHub, creates a version tag, and notifies the Go module proxy so the new version is available to everyone.

---

## Tips and Tricks

**Start with the example.** Seriously — the `./example/` folder is the best teacher. Copy it as a starting point for your own app and modify from there. It shows you the right structure from day one.

**Run `./run check` before you commit.** It catches formatting issues, build errors, and test failures all in one go. Making this a habit will save you from surprises later.

**Gio works differently from most UI frameworks.** Instead of building a UI once and updating it, Gio redraws the interface every frame. This is called *immediate mode* rendering. It might feel unusual at first, but once it clicks, it's a wonderfully direct and satisfying way to work.

**Don't worry about platform differences.** Your app will run on Linux, macOS, and Windows without any extra work. Just build and ship.

**Keep your dependencies tidy.** Every now and then, run `go mod tidy` in your project to keep your `go.mod` and `go.sum` files clean and up to date.

**Save test output for later.** The test command automatically saves a full log to `output/testing/last.log`. If a test fails and you want to look at it more carefully, it's right there waiting for you.

---

## Learn More

- [Gio documentation](https://gioui.org) — the upstream framework this library is built on. Great for going deeper.
- [pkg.go.dev reference](https://pkg.go.dev/github.com/darrenoakey/daz-golang-gio) — the full API reference for this library.

Happy building! 🎉

## License

This project is licensed under [CC BY-NC 4.0](https://darren-static.waft.dev/license) - free to use and modify, but no commercial use without permission.