// Package pixmatch is a pixel-level image comparison tool. Heavily inspired by
// https://github.com/mapbox/pixelmatch, but rewritten in idiomatic Go to speed
// up images comparison.
//
// Go pixmatch has support for PNG, GIF and JPEG formats. This tool also
// accurately detects anti-aliasing and may (or may not) count it as a
// difference.
package pixmatch
