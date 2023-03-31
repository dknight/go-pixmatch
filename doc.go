// Package pixmatch is a pixel-level image comparison tool. Heavily inspired by
// [Pixelmatch.js], but rewritten in idiomatic Go to speed up images
// comparison.
//
// Go pixmatch has support for PNG, GIF and JPEG formats. This tool also
// accurately detects anti-aliasing and may (or may not) count it as a
// difference.
// [Pixelmatch.js]: https://github.com/mapbox/pixelmatch
//
// Author:  Dmitri Smirnov (https://www.whoop.ee/)
//
// License: MIT 2023
package pixmatch
