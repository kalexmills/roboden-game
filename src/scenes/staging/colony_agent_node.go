package staging

import (
	"fmt"
	"math"

	"github.com/quasilyte/colony-game/assets"
	"github.com/quasilyte/colony-game/viewport"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
)

const (
	agentFlightHeight float64 = 40.0
	agentPickupSpeed  float64 = 40.0
)

type colonyAgentKind uint8

const (
	agentWorker colonyAgentKind = iota
	agentFreighter
	agentMilitia
	agentFighter
	agentRepeller
	agentRepair
	agentRecharger
	agentGenerator
)

type colonyAgentMode uint8

const (
	agentModeStandby colonyAgentMode = iota
	agentModeAlignStandby
	agentModeCharging
	agentModeMineEssence
	agentModeReturn
	agentModePatrol
	agentModeFollow
	agentModeMakeClone
	agentModeWaitCloning
	agentModePickup
	agentModeResourceTakeoff
	agentModeTakeoff
	agentModeRecycleReturn
	agentModeRecycleLanding
	agentModeMerging
	agentModeMergeTransform
	agentModeBuildBase
)

type agentTraitBits uint64

const (
	traitNeverStop agentTraitBits = 1 << iota
)

type colonyAgentNode struct {
	anim       *ge.Animation
	sprite     *ge.Sprite
	shadow     *ge.Sprite
	diode      *ge.Sprite
	colonyCore *colonyCoreNode

	scene *ge.Scene

	stats *agentStats

	cloningBeam *cloningBeamNode

	pos       gmath.Vec
	spritePos gmath.Vec

	traits agentTraitBits

	mode     colonyAgentMode
	waypoint gmath.Vec
	target   any

	payload    int
	cloneGen   int
	faction    factionTag
	cargoValue float64

	height float64

	attackDelay  float64
	supportDelay float64

	maxHealth  float64
	health     float64
	maxEnergy  float64
	energy     float64
	energyBill float64

	resting bool
	speed   float64

	dist          float64
	waypointsLeft int

	EventDestroyed gsignal.Event[*colonyAgentNode]
}

func newColonyAgentNode(core *colonyCoreNode, stats *agentStats, pos gmath.Vec) *colonyAgentNode {
	a := &colonyAgentNode{
		colonyCore: core,
		stats:      stats,
		pos:        pos,
		height:     40,
	}
	return a
}

func (a *colonyAgentNode) Clone() *colonyAgentNode {
	cloned := newColonyAgentNode(a.colonyCore, a.stats, a.pos)
	cloned.speed = a.speed
	cloned.maxHealth = a.maxHealth
	cloned.maxEnergy = a.maxEnergy
	cloned.traits = a.traits
	cloned.cloneGen = a.cloneGen + 1
	cloned.faction = a.faction
	return cloned
}

func (a *colonyAgentNode) Init(scene *ge.Scene) {
	a.scene = scene

	if a.cloneGen == 0 {
		a.maxHealth = scene.Rand().FloatRange(10, 15)
		a.maxEnergy = scene.Rand().FloatRange(60, 200)
		a.speed = scene.Rand().FloatRange(70, 100)

		if scene.Rand().Chance(0.4) {
			a.traits |= traitNeverStop
		}

		switch a.faction {
		case redFactionTag:
			a.maxHealth *= 1.4
		case blueFactionTag:
			a.maxEnergy *= 1.8
			a.speed *= 1.15
		}
	}

	a.health = a.maxHealth
	a.energy = a.maxEnergy

	a.sprite = scene.NewSprite(a.stats.image)
	a.sprite.Pos.Base = &a.spritePos
	a.colonyCore.world.camera.AddGraphicsAbove(a.sprite)

	if a.faction != neutralFactionTag {
		a.diode = scene.NewSprite(assets.ImageFactionDiode)
		a.diode.Pos.Base = &a.spritePos
		a.diode.Pos.Offset.Y = a.stats.diodeOffset
		var colorScale ge.ColorScale
		colorScale.SetColor(factionByTag(a.faction).color)
		a.diode.SetColorScale(colorScale)
		a.colonyCore.world.camera.AddGraphicsAbove(a.diode)
	}

	shadowImage := assets.ImageSmallShadow
	switch a.stats.size {
	case sizeMedium:
		shadowImage = assets.ImageMediumShadow
	}
	a.shadow = scene.NewSprite(shadowImage)
	a.shadow.Pos.Base = &a.spritePos
	a.colonyCore.world.camera.AddGraphics(a.shadow)

	a.anim = ge.NewRepeatedAnimation(a.sprite, -1)
	a.anim.Tick(scene.Rand().FloatRange(0, 0.7))
}

func (a *colonyAgentNode) IsDisposed() bool { return a.sprite.IsDisposed() }

func (a *colonyAgentNode) AssignMode(mode colonyAgentMode, pos gmath.Vec, target any) bool {
	if a.mode == agentModeMerging && mode != agentModeMergeTransform && mode != agentModeStandby && mode != agentModeAlignStandby {
		panic(fmt.Sprint("changing merging to", mode))
	}

	switch mode {
	case agentModePatrol:
		a.mode = mode
		a.dist = a.colonyCore.radius
		a.waypoint = a.pos.DirectionTo(a.colonyCore.body.Pos).Rotated(0.4).Mulf(a.dist).Add(a.colonyCore.body.Pos)
		a.waypointsLeft = a.scene.Rand().IntRange(40, 70)
		return true

	case agentModeWaitCloning:
		a.mode = mode
		a.target = target
		a.waypoint = gmath.Vec{}
		return true

	case agentModeMakeClone:
		a.mode = mode
		a.target = target
		a.dist = a.scene.Rand().FloatRange(1.2, 2) // cloning time
		a.energyBill += 20
		targetPos := target.(*colonyAgentNode).pos
		a.waypoint = a.pos.DirectionTo(targetPos).Mulf(110).Add(targetPos).Add(a.scene.Rand().Offset(-20, 20))
		return true

	case agentModeMerging:
		a.mode = mode
		a.target = target
		a.dist = a.scene.Rand().FloatRange(5, 6) // merging time
		return true

	case agentModeAlignStandby:
		if a.cloningBeam != nil {
			a.cloningBeam.Dispose()
			a.cloningBeam = nil
		}
		a.mode = mode
		a.waypoint = a.pos.Sub(gmath.Vec{Y: agentFlightHeight - a.height})
		return true

	case agentModeStandby:
		if a.cloningBeam != nil {
			a.cloningBeam.Dispose()
			a.cloningBeam = nil
		}
		a.mode = mode
		maxDist := a.colonyCore.radius * 0.7
		a.dist = a.scene.Rand().FloatRange(40, maxDist)
		a.waypoint = a.pos.DirectionTo(a.colonyCore.body.Pos).Rotated(0.4).Mulf(a.dist).Add(a.colonyCore.body.Pos)
		a.waypointsLeft = a.scene.Rand().IntRange(30, 60)
		return true

	case agentModeFollow:
		a.mode = mode
		a.target = target
		targetPos := target.(*creepNode).pos
		a.waypoint = a.pos.DirectionTo(targetPos).Mulf(100).Add(targetPos).Add(a.scene.Rand().Offset(-40, 40))
		a.waypointsLeft = a.scene.Rand().IntRange(4, 10)
		return true

	case agentModeCharging:
		a.mode = mode
		return true

	case agentModeMineEssence:
		if !a.stats.canGather {
			return false
		}
		energyCost := target.(*essenceSourceNode).pos.DistanceTo(a.pos) / 2
		if energyCost > a.energy {
			return false
		}
		a.energyBill += energyCost
		a.mode = mode
		a.waypoint = target.(*essenceSourceNode).pos.Sub(gmath.Vec{Y: agentFlightHeight}).Add(a.scene.Rand().Offset(-8, 8))
		a.target = target
		return true

	case agentModeTakeoff:
		a.mode = mode
		a.waypoint = a.pos.Sub(gmath.Vec{Y: agentFlightHeight})
		return true

	case agentModePickup:
		a.mode = mode
		a.waypoint = a.pos.Add(gmath.Vec{Y: agentFlightHeight})
		return true

	case agentModeRecycleReturn:
		a.mode = mode
		a.waypoint = a.colonyCore.GetEntrancePos().Sub(gmath.Vec{Y: agentFlightHeight})
		return true

	case agentModeRecycleLanding:
		a.mode = mode
		a.waypoint = a.colonyCore.GetEntrancePos()
		return true

	case agentModeBuildBase:
		construction := target.(*colonyCoreConstructionNode)
		energyCost := construction.pos.DistanceTo(a.pos) * 0.6
		if energyCost > a.energy {
			return false
		}
		a.mode = mode
		a.energyBill += energyCost
		a.dist = a.scene.Rand().FloatRange(5, 7) // build time
		a.target = target
		a.waypoint = gmath.RadToVec(a.scene.Rand().Rad()).Mulf(64.0).Add(construction.pos)
		return true
	}

	return false
}

func (a *colonyAgentNode) Update(delta float64) {
	a.anim.Tick(delta)

	a.shadow.Pos.Offset.Y = a.height + 4
	newShadowAlpha := float32(1.0 - ((a.height / agentFlightHeight) * 0.5))
	a.shadow.SetAlpha(newShadowAlpha)

	// FIXME: this should be fixed in the ge package.
	a.spritePos.X = math.Round(a.pos.X)
	a.spritePos.Y = math.Round(a.pos.Y)

	if a.energyBill != 0 {
		a.energy -= delta * 2
		a.energyBill = gmath.ClampMin(a.energyBill-delta*2, 0)
	}

	if a.resting {
		a.energy = gmath.ClampMax(a.energy+(delta*0.5), a.maxEnergy)
		if a.energy > a.maxEnergy*0.6 {
			a.resting = false
		}
	} else {
		if a.mode != agentModeStandby && a.mode != agentModeCharging && a.energy < a.maxEnergy*0.5 {
			a.resting = true
		}
	}

	a.processAttack(delta)
	a.processSupport(delta)

	switch a.mode {
	case agentModeStandby:
		a.updateStandby(delta)
	case agentModeAlignStandby:
		a.updateAlignStandby(delta)
	case agentModeCharging:
		a.updateCharging(delta)
	case agentModeMineEssence:
		a.updateMineEssence(delta)
	case agentModePickup:
		a.updatePickup(delta)
	case agentModeReturn:
		a.updateReturn(delta)
	case agentModePatrol:
		a.updatePatrol(delta)
	case agentModeFollow:
		a.updateFollow(delta)
	case agentModeWaitCloning:
		a.updateWaitCloning(delta)
	case agentModeMakeClone:
		a.updateMakeClone(delta)
	case agentModeMerging:
		a.updateMerging(delta)
	case agentModeResourceTakeoff:
		a.updateResourceTakeoff(delta)
	case agentModeTakeoff:
		a.updateTakeoff(delta)
	case agentModeRecycleReturn:
		a.updateRecycleReturn(delta)
	case agentModeRecycleLanding:
		a.updateRecycleLanding(delta)
	case agentModeBuildBase:
		a.updateBuildBase(delta)
	}
}

func (a *colonyAgentNode) Dispose() {
	a.sprite.Dispose()
	a.shadow.Dispose()
	if a.diode != nil {
		a.diode.Dispose()
	}
	if a.cloningBeam != nil {
		a.cloningBeam.Dispose()
		a.cloningBeam = nil
	}
}

func (a *colonyAgentNode) Destroy() {
	a.EventDestroyed.Emit(a)
	a.Dispose()
}

func (a *colonyAgentNode) ReceiveEnergyDamage(damage float64) {
	a.energy = gmath.ClampMin(a.energy-damage, 0)
}

func (a *colonyAgentNode) OnDamage(damage damageValue, source gmath.Vec) {
	a.health -= damage.health
	if a.health < 0 {
		playSound(a.scene, a.camera(), assets.AudioAgentDestroyed, a.pos)
		a.colonyCore.actionPriorities.AddWeight(prioritySecurity, 0.05)
		a.colonyCore.actionPriorities.AddWeight(priorityGrowth, 0.01)
		a.Destroy()

		roll := a.scene.Rand().Float()
		if roll < 0.3 {
			createExplosion(a.scene, a.camera(), a.pos)
		} else {
			var scraps *essenceSourceStats
			if roll > 0.6 {
				scraps = smallScrapSource
				if a.stats.size != sizeSmall {
					scraps = scrapSource
				}
			}
			fall := newDroneFallNode(a.colonyCore.world, scraps, a.stats.image, a.shadow.ImageID(), a.pos, a.height)
			a.scene.AddObject(fall)
		}
	}

	a.energy = gmath.ClampMin(a.energy-damage.energy, 0)

	if a.colonyCore.GetSecurityPriority() < 0.65 && a.scene.Rand().Chance(1.0-a.colonyCore.GetSecurityPriority()) {
		a.colonyCore.actionPriorities.AddWeight(prioritySecurity, 0.01)
	}
}

func (a *colonyAgentNode) GetPos() gmath.Vec { return a.pos }

func (a *colonyAgentNode) GetVelocity() gmath.Vec {
	if a.waypoint.IsZero() {
		return gmath.Vec{}
	}
	return a.pos.VecTowards(a.waypoint, a.movementSpeed())
}

func (a *colonyAgentNode) processSupport(delta float64) {
	switch a.stats.kind {
	case agentRepair, agentRecharger:
		// OK
	default:
		return
	}

	a.supportDelay = gmath.ClampMin(a.supportDelay-delta, 0)
	if a.supportDelay != 0 {
		return
	}

	a.supportDelay = a.stats.supportReload * a.scene.Rand().FloatRange(0.7, 1.4)

	switch a.stats.kind {
	case agentRecharger:
		target := a.colonyCore.FindRandomAgent(func(x *colonyAgentNode) bool {
			return x != a &&
				(x.energy+20) < x.maxEnergy &&
				x.pos.DistanceTo(a.pos) < a.stats.supportRange
		})
		if target != nil {
			beam := newBeamNode(a.camera(), ge.Pos{Base: &a.pos}, ge.Pos{Base: &target.pos}, rechargerBeamColor)
			beam.width = 2
			target.energy = gmath.ClampMax(target.energy+10, target.maxEnergy)
			a.scene.AddObject(beam)
			playSound(a.scene, a.camera(), assets.AudioRechargerBeam, a.pos)
		}
	case agentRepair:
		target := a.colonyCore.FindRandomAgent(func(x *colonyAgentNode) bool {
			return x != a &&
				x.health < x.maxHealth &&
				x.pos.DistanceTo(a.pos) < a.stats.supportRange
		})
		if target != nil {
			beam := newBeamNode(a.camera(), ge.Pos{Base: &a.pos}, ge.Pos{Base: &target.pos}, repairBeamColor)
			beam.width = 2
			target.health = gmath.ClampMax(target.health+3, target.maxHealth)
			a.scene.AddObject(beam)
			playSound(a.scene, a.camera(), assets.AudioRepairBeam, a.pos)
		}
	}
}

func (a *colonyAgentNode) processAttack(delta float64) {
	if !a.stats.canPatrol {
		return
	}

	a.attackDelay = gmath.ClampMin(a.attackDelay-delta, 0)
	if a.attackDelay != 0 {
		return
	}

	a.attackDelay = a.scene.Rand().FloatRange(1, 4.5)

	var target *creepNode
	for _, c := range a.colonyCore.world.creeps {
		if c.pos.DistanceTo(a.pos) >= a.stats.attackRange {
			continue
		}
		target = c
		break
	}
	if target == nil {
		return
	}

	switch a.stats.kind {
	case agentMilitia, agentFighter, agentRepeller:
		toPos := snipePos(a.stats.projectileSpeed, a.pos, target.pos, target.GetVelocity())
		p := newProjectileNode(projectileConfig{
			Camera:      a.colonyCore.world.camera,
			Image:       a.stats.projectileImage,
			FromPos:     a.pos,
			ToPos:       toPos,
			Target:      target,
			Speed:       a.stats.projectileSpeed,
			Area:        a.stats.projectileArea,
			RotateSpeed: a.stats.projectileRotateSpeed,
			Damage:      a.stats.projectileDamage,
		})
		a.scene.AddObject(p)
	}

	playSound(a.scene, a.camera(), a.stats.attackSound, a.pos)
}

func (a *colonyAgentNode) movementSpeed() float64 {
	switch a.mode {
	case agentModeTakeoff, agentModeRecycleLanding:
		return 30
	case agentModePickup, agentModeResourceTakeoff, agentModeAlignStandby:
		return agentPickupSpeed
	}
	if a.resting {
		return a.speed * 0.5
	}
	return a.speed
}

func (a *colonyAgentNode) moveTowards(delta float64, pos gmath.Vec) bool {
	travelled := a.movementSpeed() * delta
	if a.pos.DistanceTo(pos) <= travelled {
		a.pos = pos
		return true
	}
	a.pos = a.pos.MoveTowards(pos, travelled)
	return false
}

func (a *colonyAgentNode) updatePatrol(delta float64) {
	if a.moveTowards(delta, a.waypoint) {
		a.dist = a.colonyCore.radius
		a.waypoint = a.pos.DirectionTo(a.colonyCore.body.Pos).Rotated(0.4).Mulf(a.dist).Add(a.colonyCore.body.Pos)
		a.waypointsLeft--
		if a.waypointsLeft == 0 {
			a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
			return
		}
	}
}

func (a *colonyAgentNode) updateWaitCloning(delta float64) {
	cloner := a.target.(*colonyAgentNode)
	if cloner.IsDisposed() {
		a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
		return
	}
}

func (a *colonyAgentNode) updateTakeoff(delta float64) {
	a.height += delta * 30
	if a.moveTowards(delta, a.waypoint) {
		a.height = agentFlightHeight
		a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
	}
}

func (a *colonyAgentNode) updateRecycleReturn(delta float64) {
	if a.moveTowards(delta, a.waypoint) {
		a.colonyCore.openHatchTime = 1.5
		a.AssignMode(agentModeRecycleLanding, gmath.Vec{}, nil)
	}
}

func (a *colonyAgentNode) updateBuildBase(delta float64) {
	target := a.target.(*colonyCoreConstructionNode)
	if target.IsDisposed() {
		if a.cloningBeam != nil {
			a.cloningBeam.Dispose()
		}
		a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
		return
	}
	if !a.waypoint.IsZero() {
		if a.moveTowards(delta, a.waypoint) {
			target.attention += 3.5
			a.waypoint = gmath.Vec{}
			buildPos := ge.Pos{
				Base:   &target.constructPos,
				Offset: gmath.Vec{X: a.scene.Rand().FloatRange(-20, 20)},
			}
			beam := newCloningBeamNode(a.colonyCore.world.camera, false, &a.pos, buildPos)
			a.cloningBeam = beam
			a.scene.AddObject(beam)
			return
		}
		return
	}
	newColony := target.Construct(delta)
	if newColony != nil {
		a.colonyCore.DetachAgent(a)
		newColony.AcceptAgent(a)
		a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
		return
	}
	a.dist -= delta
	if a.dist <= 0 || a.energy < 20 {
		a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
		return
	}
}

func (a *colonyAgentNode) updateRecycleLanding(delta float64) {
	a.height -= delta * 30
	if a.moveTowards(delta, a.waypoint) {
		a.colonyCore.resources.Essence += a.stats.cost
		playSound(a.scene, a.camera(), assets.AudioAgentRecycled, a.pos)
		a.Destroy()
	}
}

func (a *colonyAgentNode) updateMerging(delta float64) {
	target := a.target.(*colonyAgentNode)
	if target.IsDisposed() {
		if a.cloningBeam != nil {
			a.cloningBeam.Dispose()
		}
		a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
		return
	}
	if a.waypoint.IsZero() {
		dist := target.pos.DistanceTo(a.pos)
		if dist > 64 {
			a.waypoint = a.pos.MoveTowards(target.pos, dist-20).Add(a.scene.Rand().Offset(-8, 8))
			return
		}
	}
	if !a.waypoint.IsZero() {
		if a.moveTowards(delta, a.waypoint) {
			a.waypoint = gmath.Vec{}
		}
		return
	}
	if a.cloningBeam == nil {
		beam := newCloningBeamNode(a.colonyCore.world.camera, true, &a.pos, ge.Pos{Base: &target.pos})
		a.cloningBeam = beam
		a.scene.AddObject(beam)
	}
	a.dist -= delta
	if a.pos.DistanceTo(target.pos) > 10 {
		a.pos = a.pos.MoveTowards(target.pos, delta*10)
	} else {
		// Merging is x2 faster when units are next to each other.
		a.dist -= delta
	}
	if a.dist <= 0 {
		a.cloningBeam.Dispose()
		a.cloningBeam = nil
		newStats := mergeAgents(a, target)
		newAgent := a.colonyCore.NewColonyAgentNode(newStats, target.pos)
		var newFaction factionTag
		if newStats.tier == 2 {
			newFaction = a.colonyCore.pickAgentFaction()
		} else {
			newFaction := a.faction
			if newFaction == neutralFactionTag || (target.faction != neutralFactionTag && a.faction != target.faction && a.scene.Rand().Bool()) {
				newFaction = target.faction
			}
		}
		newAgent.faction = newFaction
		a.scene.AddObject(newAgent)
		newAgent.AssignMode(agentModeStandby, gmath.Vec{}, nil)
		target.Destroy()
		a.Destroy()
		return
	}
}

func (a *colonyAgentNode) updateMakeClone(delta float64) {
	target := a.target.(*colonyAgentNode)
	if target.IsDisposed() {
		if a.cloningBeam != nil {
			a.cloningBeam.Dispose()
		}
		a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
		return
	}
	if !a.waypoint.IsZero() {
		if a.moveTowards(delta, a.waypoint) {
			a.waypoint = gmath.Vec{}
			beam := newCloningBeamNode(a.colonyCore.world.camera, false, &a.pos, ge.Pos{Base: &target.pos})
			a.cloningBeam = beam
			a.scene.AddObject(beam)
			return
		}
		return
	}
	a.dist -= delta
	if a.dist <= 0 {
		a.cloningBeam.Dispose()
		a.cloningBeam = nil
		a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
		target.AssignMode(agentModeStandby, gmath.Vec{}, nil)
		clone := a.colonyCore.CloneAgentNode(target)
		a.scene.AddObject(clone)
		clone.AssignMode(agentModeStandby, gmath.Vec{}, nil)
		return
	}
}

func (a *colonyAgentNode) updateFollow(delta float64) {
	if a.moveTowards(delta, a.waypoint) {
		if a.waypointsLeft == 0 {
			a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
			return
		}
		a.waypointsLeft--
		targetPos := a.target.(*creepNode).pos
		a.waypoint = a.pos.DirectionTo(targetPos).Mulf(100).Add(targetPos).Add(a.scene.Rand().Offset(-40, 40))
	}
}

func (a *colonyAgentNode) updateAlignStandby(delta float64) {
	a.height += delta * agentPickupSpeed
	if a.moveTowards(delta, a.waypoint) {
		a.height = agentFlightHeight
		a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
	}
}

func (a *colonyAgentNode) updateStandby(delta float64) {
	a.energy += delta * 0.5
	if a.moveTowards(delta, a.waypoint) {
		if a.waypointsLeft == 0 {
			a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
			return
		}
		a.waypointsLeft--
		a.waypoint = a.pos.DirectionTo(a.colonyCore.body.Pos).Rotated(0.4).Mulf(a.dist).Add(a.colonyCore.body.Pos)
		if !a.hasTrait(traitNeverStop) && a.energy < 40 && a.scene.Rand().Chance(0.2) {
			a.AssignMode(agentModeCharging, gmath.Vec{}, nil)
			return
		}
	}
}

func (a *colonyAgentNode) updateCharging(delta float64) {
	a.energy += delta * 3
	if a.energy >= 50 {
		a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
	}
}

func (a *colonyAgentNode) updateMineEssence(delta float64) {
	if a.moveTowards(delta, a.waypoint) {
		a.AssignMode(agentModePickup, gmath.Vec{}, nil)
	}
}

func (a *colonyAgentNode) updatePickup(delta float64) {
	a.height -= delta * agentPickupSpeed
	if a.moveTowards(delta, a.waypoint) {
		a.height = 0
		a.mode = agentModeResourceTakeoff
		a.waypoint = a.pos.Sub(gmath.Vec{Y: agentFlightHeight})
		source := a.target.(*essenceSourceNode)
		harvested := source.Harvest(a.maxPayload())
		a.payload = harvested
		a.cargoValue += float64(harvested) * source.stats.value
	}
}

func (a *colonyAgentNode) updateResourceTakeoff(delta float64) {
	a.height += delta * agentPickupSpeed
	if a.moveTowards(delta, a.waypoint) {
		a.height = agentFlightHeight
		entranceNum := a.scene.Rand().IntRange(0, 2)
		a.waypoint = a.colonyCore.GetStoragePos().Add(gmath.Vec{Y: float64(entranceNum) * 8})
		a.mode = agentModeReturn
	}
}

func (a *colonyAgentNode) updateReturn(delta float64) {
	if a.moveTowards(delta, a.waypoint) {
		if a.payload != 0 {
			a.colonyCore.resources.Essence += a.cargoValue
			a.payload = 0
			playSound(a.scene, a.camera(), assets.AudioEssenceCollected, a.pos)
		}
		a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
	}
}

func (a *colonyAgentNode) camera() *viewport.Camera {
	return a.colonyCore.world.camera
}

func (a *colonyAgentNode) hasTrait(t agentTraitBits) bool {
	return a.traits&t != 0
}

func (a *colonyAgentNode) maxPayload() int {
	n := a.stats.maxPayload
	if a.faction == yellowFactionTag {
		n++
	}
	return n
}