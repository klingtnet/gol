# gol - a simple golang powered logbook

Pronounced "goal" (/ɡoʊl/).

## installation

```sh
$ make
$ ./main
Listening on http://0.0.0.0:5000
```

If you want to use ssl, you can [generate a certificate](https://devcenter.heroku.com/articles/ssl-certificate-self#generate-private-key-and-certificate-signing-request)
and then start the server using the `-ssl` flag, passing the certificate
and the private key to it:

```sh
$ ./main --ssl server.crt,server.key
Listening on https://0.0.0.0:5000
```

### as docker container

- build the container `docker build -t gol .`
- run the container `docker run -p 5000:5000 gol`

## dependencies

`gol` uses the following libraries (which are awesome):

* [blackfriday](https://github.com/russross/blackfriday) for rendering
    markdown
* [bluemonday](https://godoc.org/github.com/microcosm-cc/bluemonday) to
    sanitize html
* [mux](https://github.com/gorilla/mux) for routing (supports url
    parameters)
* [pflag](https://github.com/ogier/pflag) for posix-style command-line
    flags

Thanks for writing those libraries!

## license

`gol` is licensed under the GNU GPL.
