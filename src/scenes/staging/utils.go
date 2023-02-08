package staging

import (
	"github.com/quasilyte/colony-game/assets"
	"github.com/quasilyte/colony-game/viewport"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

func posIsFree(world *worldState, skipColony *colonyCoreNode, pos gmath.Vec, radius float64) bool {
	for _, source := range world.essenceSources {
		if source.pos.DistanceTo(pos) < (radius + source.stats.size) {
			return false
		}
	}
	for _, construction := range world.coreConstructions {
		if construction.pos.DistanceTo(pos) < (radius + 40) {
			return false
		}
	}
	for _, colony := range world.colonies {
		if colony == skipColony {
			continue
		}
		if colony.body.Pos.DistanceTo(pos) < (radius + 40) {
			return false
		}
	}
	return true
}

func createExplosion(scene *ge.Scene, camera *viewport.Camera, pos gmath.Vec) {
	explosion := newEffectNode(camera, pos, assets.ImageSmallExplosion1)
	scene.AddObject(explosion)

	explosionSoundIndex := scene.Rand().IntRange(0, 4)
	explosionSound := resource.AudioID(int(assets.AudioExplosion1) + explosionSoundIndex)
	playSound(scene, camera, explosionSound, pos)
}

func roundedPos(pos gmath.Vec) gmath.Vec {
	x := int(pos.X)
	y := int(pos.Y)
	if x%2 != 0 {
		x++
	}
	if y%2 != 0 {
		y++
	}
	return gmath.Vec{X: float64(x), Y: float64(y)}
}

func correctedPos(sector gmath.Rect, pos gmath.Vec, pad float64) gmath.Vec {
	if pos.X < (pad + sector.Min.X) {
		pos.X = pad + sector.Min.X
	} else if pos.X > (sector.Max.X - pad) {
		pos.X = sector.Max.X - pad
	}
	if pos.Y < (pad + sector.Min.Y) {
		pos.Y = pad + sector.Min.Y
	} else if pos.Y > (sector.Max.Y - pad) {
		pos.Y = sector.Max.Y - pad
	}
	return pos
}

func snipePos(projectileSpeed float64, fireFrom, targetPos, targetVelocity gmath.Vec) gmath.Vec {
	if targetVelocity.IsZero() || projectileSpeed == 0 {
		return targetPos
	}
	dist := targetPos.DistanceTo(fireFrom)
	predictedPos := targetPos.Add(targetVelocity.Mulf(dist / projectileSpeed))
	return predictedPos
}

func randFind[T any](rand *gmath.Rand, slice []T, f func(x T) bool) T {
	var result T
	randWalk(rand, slice, func(x T) bool {
		if f(x) {
			result = x
			return false
		}
		return true
	})
	return result
}

func randWalk[T any](rand *gmath.Rand, slice []T, f func(x T) bool) {
	if len(slice) == 0 {
		return
	}
	var slider gmath.Slider
	slider.SetBounds(0, len(slice)-1)
	slider.TrySetValue(rand.IntRange(0, len(slice)-1))
	for i := 0; i < len(slice); i++ {
		if !f(slice[slider.Value()]) {
			break
		}
		slider.Inc()
	}
}

func playSound(scene *ge.Scene, camera *viewport.Camera, id resource.AudioID, pos gmath.Vec) {
	if camera.ContainsPos(pos) {
		scene.Audio().PlaySound(id)
	}
}