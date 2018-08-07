# zwiebelpfanne

Zwiebelpfanne is a horrible compound of the German words for onion (*Zwiebel*)
and socket (*Pfanne*, more like *pan*).

This tool binds a remote
[Tor Hidden Service](https://www.torproject.org/docs/onion-services.html.en)
(also known as *.onion-Domain*) to a local TCP-Port. This enable the use of
Hidden Services for every application, not only those which speak SOCKS5.

A running Tor daemon is required.


## Usage

Bind Facebook's Hidden Service to the localhost on port 1337. The `--onion` flag
is the only required flag.

```
$ ./zwiebelpfanne --onion "facebookcorewwwi.onion:80"
zwiebelpfanne: facebookcorewwwi.onion:80 -> localhost:1337
```

### Flags

- `--tor-socks5` specifies Tor's SOCKS5 daemon (default `localhost:9050`)
- `--listen` specifies where zwiebelpfanne should be bound to
  (default `localhost:1337`)
- `--onion` specifies the Tor Hidden Service zwiebelpfanne should connect to

## Install

Go and a running Tor daemon is required.

```bash
go get github.com/geistesk/zwiebelpfanne
go build github.com/geistesk/zwiebelpfanne

ls ~/go/bin/zwiebelpfanne
```
