# watch

![watch](https://user-images.githubusercontent.com/141232/54884127-c564e080-4e9f-11e9-9116-3f7b72607beb.gif)

_watch_ tool rewritten in go.

## Features

* 1s by default
* `bash -cli` by default
* working aliases

## Usage

```bash
watch [command]
```

Specify command for _watch_ by setting `WATCH_COMMAND`.

```bash
export WATCH_COMMAND=`fish -c`
```

## Install

```bash
go get github.com/antonmedv/watch
```

## License

MIT
