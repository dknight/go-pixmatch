// Package pixmatch is a pixel-level image comparison tool. Heavily inspired
// by https://github.com/mapbox/pixelmatch, but rewritten to idiomatic Go
// language to speed up images comparison.
//
// Go pixmatch has support for PNG, GIF and JPEG (still dirty) formats. This
// tool also accurately  detects anti-aliasing accurately and may (or may not)
// count it as difference.
package pixmatch
