# pixmatch

[![Tests](https://github.com/dknight/go-pixmatch/actions/workflows/tests.yml/badge.svg)](https://github.com/dknight/go-pixmatch/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dknight/go-pixmatch)](https://goreportcard.com/report/github.com/dknight/go-pixmatch)
[![Go Reference](https://pkg.go.dev/badge/github.com/dknight/go-pixmatch.svg)](https://pkg.go.dev/github.com/dknight/go-pixmatch)

**pixmatch** is a pixel-level image comparison tool. Heavily inspired by
[pixelmatch](https://github.com/mapbox/pixelmatch), but rewritten in idiomatic
Go, **with zero dependencies,** to speed up images comparison.
Go pixmatch has support for **PNG**, **GIF** and **JPEG** formats. This tool
also accurately detects anti-aliasing and may count it as a difference.

## Example output

| Format | Expected                                                                                         | Actual                                                                                           | Difference                                                                                          |
| ------ | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------- |
| JPEG   | ![Hummingbird](https://raw.githubusercontent.com/dknight/go-pixmatch/main/samples/bird-a.jpg)    | ![Hummingbird](https://raw.githubusercontent.com/dknight/go-pixmatch/main/samples/bird-b.jpg)    | ![Hummingbird](https://raw.githubusercontent.com/dknight/go-pixmatch/main/samples/bird-diff.jpg)    |
| GIF    | ![Landscape](https://raw.githubusercontent.com/dknight/go-pixmatch/main/samples/landscape-a.gif) | ![Landscape](https://raw.githubusercontent.com/dknight/go-pixmatch/main/samples/landscape-b.gif) | ![Landscape](https://raw.githubusercontent.com/dknight/go-pixmatch/main/samples/landscape-diff.gif) |
| PNG    | ![Form](https://raw.githubusercontent.com/dknight/go-pixmatch/main/samples/form-a.png)           | ![Form](https://raw.githubusercontent.com/dknight/go-pixmatch/main/samples/form-b.png)           | ![Form](https://raw.githubusercontent.com/dknight/go-pixmatch/main/samples/form-aa-diff.png)        |

## Install

Library:

```sh
go get -u github.com/dknight/go-pixmatch
```

CLI:

```sh
go install github.com/dknight/go-pixmatch/cmd/pixmatch
```

## Library usage

```go
img1, err := NewImageFromPath("./samples/form-a.png")
if err != nil {
    log.Fatalln(err)
}
img2, err := NewImageFromPath("./samples/form-b.png")
if err != nil {
    log.Fatalln(err)
}

// Set some options.
options := NewOptions()
options.SetThreshold(0.05)
options.SetAlpha(0.5)
options.SetDiffColor(color.RGBA{0, 255, 128, 255})
// etc...

diff, err := img1.Compare(img2, options)
if err != nil {
    log.Fatalln(err)
}

fmt.Println(diff)
```

## CLI usage

Usage:

```sh
pixmatch [options] image1.png image2.png
```

Run `pixmatch -h` for the list of supported options.

Example command:

```sh
pixmatch -o diff.png -aa -aacolor=00ffffff -mask ./samples/form-a.png ./samples/form-b.png
```

### Compile binaries

Here is included simple script to compile binaries for some architectures.
If you need something special, you can easily adopt it for your needs.

```sh
./scripts/makebin.sh
```

## Testing and benchmark

Simple tests:

```sh
go test
```

Tests with update diff images:

```sh
UPDATEDIFFS=1 go test
```

Tests with full coverage:

```sh
# Terminal output
UPDATEDIFFS=1 go test -cover

# HTML output
UPDATEDIFFS=1 go test -coverprofile c.out && go tool cover -html c.out
```

Benchmark scripts (outputFile will be written in `logs/` directory):

```sh
./scripts/benchmark.sh <outputFile> [iterations=10]
```

Later, it is easier to analyze it with a cool [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) tool.

## Using with WASM

WASI Preview 1 is very unstable, everything will be changed. This is
an rough example how it might work with golang <= 1.12.6.

Compile the WASM file.


```sh
cd cmd/pixmatch
GOOS=wasip1 GOARCH=wasm go build -o pixmatch.wasm main.go
```

Node.js exanoke file (index.js)

```js
const fs = require('node:fs');
const {WASI} = require('node:wasi');
const imgs = ['form-a.png', 'form-b.png'];
const virtulPathWasiPath = '';
const wasi = new WASI({
  version: 'preview1',
  args: [virtulPathWasiPath, '-mask', '-o', 'diff.png', ...imgs],
  preopens: {
    '/': __dirname,
  },
});

const wasmBuffer = fs.readFileSync('./pixmatch.wasm');
WebAssembly.instantiate(wasmBuffer, wasi.getImportObject()).then(
  (wasmModule) => {
    wasi.start(wasmModule.instance);
  }
);
```

Run node script:

```sh
NODE_NO_WARNINGS=1 node index.js
```

## Known issues, bugs, flaws

- Anti-aliasing detection algorithm can be improved (help appreciated).
- Because of the nature of the JPEG format, comparing them is not a good idea or play with `threshold` parameter.
- I have not tested this tool for 64-bit color images.

## Credits

- To provide 100% compatibility with [pixelmatch](https://github.com/mapbox/pixelmatch).
  Original test files are borrowed from [fixtures](https://github.com/mapbox/pixelmatch/tree/main/test/fixtures).
- [Hummingbird](https://commons.wikimedia.org/wiki/File:Hummingbird.jpg) is taken from Wiki commons by San Diego Zoo.
- Someone on the [Pixilart](https://www.pixilart.com/draw/16x16-6ec491154b5c687) platform created this pixel art girl.
- Form screenshots are made using [PureCSS](https://purecss.io/) framework.

## Contribution

Any help is appreciated. Found a bug, typo, inaccuracy, etc.? Please do not hesitate to make a pull request or file an issue.

## License

MIT 2023
