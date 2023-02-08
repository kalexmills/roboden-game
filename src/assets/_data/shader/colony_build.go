//go:build ignore
// +build ignore

package main

var Time float

func gridPattern(v, colorMult vec4, hash int, p, pixSize, originTexPos vec2) vec4 {
	posHash := int(p.x+p.y) * int(p.y*5)
	state := posHash % hash
	if state == int(1) {
		p.x += 1.0
	} else if state == int(2) {
		p.x -= 1.0
	} else if state == int(3) {
		p.y += 1.0
	} else if state == int(4) {
		p.y -= 1.0
	} else {
		return v
	}
	return imageSrc0At(p/pixSize+originTexPos) * colorMult
}

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	v := imageSrc0UnsafeAt(texCoord)
	if v.a == 0 {
		return v
	}

	pixSize := imageSrcTextureSize()
	originTexPos, _ := imageSrcRegionOnTexture()
	actualTexPos := vec2(texCoord.x-originTexPos.x, texCoord.y-originTexPos.y)
	actualPixPos := actualTexPos * pixSize

	initialY := -2.0
	offsetY := 10.0 * Time

	dist := distance(actualPixPos, vec2(32, initialY-offsetY))
	if dist > (60.0 * (1.4 - Time)) {
		return v
	}
	if dist > (54.0 * (1.4 - Time)) {
		return gridPattern(v, vec4(1, 1.1, 1.3, 1.0), 15, actualPixPos, pixSize, originTexPos)
	}
	if dist > (48.0 * (1.4 - Time)) {
		return gridPattern(v, vec4(0.9, 1.2, 1.6, 1.0), 12, actualPixPos, pixSize, originTexPos)
	}
	if dist > (42.0 * (1.4 - Time)) {
		return gridPattern(v, vec4(0.8, 1.4, 2.0, 1.0), 8, actualPixPos, pixSize, originTexPos)
	}
	if dist > (40.0 * (1.4 - Time)) {
		v = gridPattern(v, vec4(0.7, 0.7, 0.7, 1.0), 7, actualPixPos, pixSize, originTexPos)
		v.xyz *= 0.4
		return v
	}

	return vec4(0, 0, 0, 0)
}