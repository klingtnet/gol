# gol - a simple golang powered logbook

Pronounced "goal" (/ɡoʊl/).

## installation

```
$ cd gol
$ go get -d .
$ go build main.go
$ ./main
Listening on http://0.0.0.0:5000
```

## dependencies

`gol` uses the following libraries (which are awesome):

* [blackfriday](https://github.com/russross/blackfriday) for rendering
    markdown
* [bluemonday](https://godoc.org/github.com/microcosm-cc/bluemonday) to
    sanitize html
* [mux](https://github.com/gorilla/mux) for routing (supports url
    parameters)

Thanks for writing those libraries!

## license

`gol` is licensed under the GNU GPL.
