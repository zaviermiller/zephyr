# Zephyr
Zephyr is a reactive frontend Go framework based on core ideas from React and Vue 3. The purpose of this project is to provide an easy to use WebAssembly reactive framework API, to further test the technologies features and capabilities. Some experienced Gophers may have more technical questions, and to them I extend an invitation to read the [FAQ]() on the Zephyr website.

This package uses the `unsafe` and `syscall/js` APIs, both of which are subject to change. All part of the library that use

-Zavier Miller (creator of Zephyr)

### Acknowledgements

---

## Design Philosophies and Goals
 - **Beginner-friendly** - Easy for people coming over from JS and beginners alike
 - **Accessible and responsive first** - Default APIs provide size breakpoints, 
 - **Make it Go** - The Zephyr templating language is meant to be written in a way that mimics how Go code is written. There are a few abstractions to make frontend development easier, but for the most part, just write Go as you normally would. The entire templating language and build process is optional, if you prefer to use just Go, go ahead!
 - **Robust tooling**
 - **Intuitive abstractions**
 ---
## Documentation
---

### Installation
To install Zephyr, just go get it: `go get -u github.com/zaviermiller/zephyr && go install zephyr`. This will install the zephyr package and then install the Zephyr CLI tool.

There are two main ways to develop an app with Zephyr: pure Go and the Zephyr dev environment. 
## Benchmarks
---
## Features and Roadmap
These are the concrete features currently offered and planned by the Zephyr API:
- [x] Zephyr core runtime
- [x] Reactive data API
- [x] Component API
- [x] Asynchronous runtime
- [x] One-way data flow between components as props
- [x] TinyGo support
- [x] Vue-based DOM updates (initial pass)
- [ ] Two-way data flow with DOM and Component Events
- [ ] Lifecycle hooks (OnInit, OnMount, OnRerender, Ondeactivate, On)
- [ ] Zephyr CLI tool (create, dev, build, test)
- [ ] Width breakpoint APIs
- [ ] Testing framework
- [ ] Zephyr dev environment
- [ ] Vue-based `.zefr` compiler
- [ ] Vue-based slots
- [ ] [Accessible-only build]()
- [ ] Easy Dockerization
- [ ] Support using any Go package in Components
- [ ] Plugins
- [ ] Routing plugin (separate import)
<!-- - [ ] Global state plugin (separate import) -->
- [ ] GraphQL plugin (seperate import)
- [ ] Content plugin (seperate import, requires SSG or SSR)
- [ ] Very simple Material UI plugin (separate import)
- [ ] Support build targets: SSR, SPA, SSG
- [ ] Chrome dev tools
- - [ ] Component inspector
- - [ ] State manager
- - [ ] More
- - [ ] & more! [Accepting suggestions]()




---
## Issues
