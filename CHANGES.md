# next - (not yet released)

- fullscreen editing mode
- saving posts without interupting writing them
- a query interface (with nice syntax in Go, accessible via url
    parameters)
- two new backends:
    - `gol`: proxies to another instance of gol
    - `multi`: writes to multiple storages, for example for backup or
        "remote publishing"
- support for authentication with pluggable providers
    - existing backends: `ldap`, `insecure` (for testing)

# 0.2.0 - Now we're getting fancy...

- markdown preview while editing posts
- sqlite backend works (except sync, which returns nil)

# 0.1.0 - Hello World!

Initial release.

- listing, creating, editing and deleting posts
    * both via the UI and a simple JSON api
- markdown rendering
- mathjax integration
- customizable templates
- pluggable storage backends (`memory` and `json`)
