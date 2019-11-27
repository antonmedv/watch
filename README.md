<img src="https://user-images.githubusercontent.com/141232/54884127-c564e080-4e9f-11e9-9116-3f7b72607beb.gif" width="500" align="right" alt="watch">

# ⏰ watch

_watch_ tool rewritten in go.

## Features

* 2s interval
* working aliases
* configurable shell
* windows support

## Usage

```bash
watch [command]
```

Specify command for _watch_ by setting `WATCH_COMMAND` (`bash -cli` by default).

```bash
export WATCH_COMMAND=`fish -c`
```

## Example

```bash
watch git status
```

```bash
watch curl wttr.in
```

```bash
watch 'll | grep .go'
```

## Install

```bash
go get github.com/antonmedv/watch
```

## License

MIT
