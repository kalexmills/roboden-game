package staging

import (
	"fmt"
	"math"
	"sort"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/roboden-game/gamedata"
	"github.com/quasilyte/roboden-game/pathing"
)

const (
	maxWallSegments int     = 16
	wallTileSize    float64 = 32
)

type levelGenerator struct {
	scene        *ge.Scene
	rng          gmath.Rand
	world        *worldState
	playerSpawn  gmath.Vec
	sectors      []gmath.Rect
	sectorSlider gmath.Slider

	pendingResources []*essenceSourceNode
}

func newLevelGenerator(scene *ge.Scene, world *worldState) *levelGenerator {
	g := &levelGenerator{
		scene: scene,
		world: world,
	}
	g.rng.SetSeed(world.config.Seed)
	g.sectors = []gmath.Rect{
		{Min: gmath.Vec{X: 0, Y: 0}, Max: gmath.Vec{X: g.world.width / 2, Y: g.world.height / 2}},
		{Min: gmath.Vec{X: g.world.width / 2, Y: 0}, Max: gmath.Vec{X: g.world.width, Y: g.world.height / 2}},
		{Min: gmath.Vec{X: 0, Y: g.world.height / 2}, Max: gmath.Vec{X: g.world.width / 2, Y: g.world.height}},
		{Min: gmath.Vec{X: g.world.width / 2, Y: g.world.height / 2}, Max: gmath.Vec{X: g.world.width, Y: g.world.height}},
	}
	g.sectorSlider.SetBounds(0, len(g.sectors)-1)
	return g
}

func (g *levelGenerator) Generate() {
	resourceMultipliers := []float64{
		0.4,
		0.75,
		1, // Default
		1.25,
		1.6,
	}
	g.placePlayers()
	g.placeWalls()
	g.placeCreepBases()
	g.placeCreeps()
	g.placeResources(resourceMultipliers[g.world.config.Resources])
	g.placeBoss()

	g.fillPathgrid()
}

func (g *levelGenerator) randomPos(sector gmath.Rect) gmath.Vec {
	return randomSectorPos(&g.rng, sector)
}

func (g *levelGenerator) fillPathgrid() {
	p := g.world.pathgrid

	if pathing.CellSize != wallTileSize {
		panic("update the pathgrid build algorithm")
	}

	// Traverse all relevant world objects and mark the occupied cells.

	// We're using a few assumptions here:
	// 1. Wall tiles are always grid-aligned.
	// 2. Wall tiles have the same grid size as path grid cells.
	for _, wall := range g.world.walls {
		if wall.rectShape {
			for y := wall.rect.Min.Y; y <= wall.rect.Max.Y; y += pathing.CellSize {
				for x := wall.rect.Min.X; x <= wall.rect.Max.X; x += pathing.CellSize {
					pos := gmath.Vec{X: x, Y: y}
					p.MarkCell(p.PosToCoord(pos))
				}
			}
			continue
		}
		for _, pos := range wall.points {
			p.MarkCell(p.PosToCoord(pos))
		}
	}

	// For resources we can only get an approx grid cell,
	// since resources are not grid-aligned.
	for _, essence := range g.world.essenceSources {
		p.MarkCell(p.PosToCoord(essence.pos))
	}
}

func (g *levelGenerator) placePlayers() {
	g.playerSpawn = g.world.rect.Center()
	g.createBase(g.playerSpawn, true)

	if g.world.config.SecondBase {
		g.createBase(g.playerSpawn.Add(gmath.Vec{X: 160, Y: 96}), false)
	}
}

func (g *levelGenerator) createBase(pos gmath.Vec, mainBase bool) {
	core := g.world.NewColonyCoreNode(colonyConfig{
		World:  g.world,
		Radius: 128,
		Pos:    pos,
	})
	core.priorities.SetWeight(priorityResources, 0.5)
	core.priorities.SetWeight(priorityGrowth, 0.4)
	core.priorities.SetWeight(prioritySecurity, 0.1)
	g.world.nodeRunner.AddObject(core)

	switch g.world.config.StartingResources {
	case 1:
		core.resources = maxVisualResources / 3
	case 2:
		core.resources = maxVisualResources
	}

	for i := 0; i < 5; i++ {
		a := core.NewColonyAgentNode(gamedata.WorkerAgentStats, core.pos.Add(g.rng.Offset(-20, 20)))
		g.world.nodeRunner.AddObject(a)
		a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
	}
	if mainBase {
		for _, stats := range g.world.config.ExtraDrones {
			a := core.NewColonyAgentNode(stats, core.pos.Add(g.scene.Rand().Offset(-20, 20)))
			g.world.nodeRunner.AddObject(a)
			a.AssignMode(agentModeStandby, gmath.Vec{}, nil)
		}
	}
}

func (g *levelGenerator) placeCreepsCluster(sector gmath.Rect, maxSize int, stats *creepStats) int {
	rand := &g.rng
	placed := 0
	pos := correctedPos(sector, g.randomPos(sector), 128)
	initialPos := pos
	unitPos := pos
	for i := 0; i < maxSize; i++ {
		if !posIsFree(g.world, nil, pos, 24) || pos.DistanceTo(g.playerSpawn) < 520 {
			break
		}
		creep := g.world.NewCreepNode(pos, stats)
		if stats.kind == creepCrawler {
			creep.specialModifier = crawlerGuard
		}
		g.world.nodeRunner.AddObject(creep)
		unitPos = pos
		direction := gmath.RadToVec(rand.Rad()).Mulf(32)
		if rand.Bool() {
			pos = initialPos.Add(direction)
		} else {
			pos = pos.Add(direction)
		}
		placed++
	}
	// Creep groups may have some scraps near them.
	if placed != 0 && rand.Chance(0.7) {
		numScraps := rand.IntRange(1, 2)
		for i := 0; i < numScraps; i++ {
			scrapPos := g.adjustResourcePos(gmath.RadToVec(rand.Rad()).Mulf(rand.FloatRange(64, 128)).Add(unitPos))
			if posIsFree(g.world, nil, scrapPos, 8) {
				source := g.world.NewEssenceSourceNode(scrapSource, scrapPos)
				g.world.nodeRunner.AddObject(source)
			}
		}
	}
	return placed
}

func (g *levelGenerator) adjustResourcePos(pos gmath.Vec) gmath.Vec {
	aligned := g.world.pathgrid.AlignPos(pos)
	return aligned.Add(g.rng.Offset(-10, 10))
}

func (g *levelGenerator) placeResourceCluster(sector gmath.Rect, maxSize int, kind *essenceSourceStats) int {
	rand := &g.rng
	placed := 0
	pos := correctedPos(sector, g.randomPos(sector), 196)
	initialPos := pos
	for i := 0; i < maxSize; i++ {
		pos = g.adjustResourcePos(pos)
		if !posIsFree(g.world, nil, pos, 8) {
			break
		}
		source := g.world.NewEssenceSourceNode(kind, pos)
		g.pendingResources = append(g.pendingResources, source)
		direction := gmath.RadToVec(rand.Rad()).Mulf(32)
		if rand.Bool() {
			pos = initialPos.Add(direction)
		} else {
			pos = pos.Add(direction)
		}
		placed++
	}
	return placed
}

func (g *levelGenerator) placeResources(resMultiplier float64) {
	rand := &g.rng

	worldSizeMultipliers := []float64{
		0.8,
		0.9,
		1.0,
		1.1,
	}
	multiplier := resMultiplier * worldSizeMultipliers[g.world.config.WorldSize]
	numIron := int(float64(rand.IntRange(26, 38)) * multiplier)
	numScrap := int(float64(rand.IntRange(6, 8)) * multiplier)
	numGold := int(float64(rand.IntRange(20, 28)) * multiplier)
	numCrystals := int(float64(rand.IntRange(14, 20)) * multiplier)
	numOil := int(float64(rand.IntRange(4, 6)) * multiplier)

	numRedOil := 0
	numRedCrystals := 0
	if g.world.config.EliteResources {
		numRedOil = gmath.ClampMin(int(float64(rand.IntRange(2, 3))*multiplier), 2)
		numRedCrystals = int(float64(rand.IntRange(10, 15)) * multiplier)
	}

	g.world.numRedCrystals = numRedCrystals

	g.sectorSlider.TrySetValue(rand.IntRange(0, len(g.sectors)-1))

	for numIron > 0 {
		clusterSize := rand.IntRange(2, 6)
		sector := g.sectors[g.sectorSlider.Value()]
		g.sectorSlider.Inc()
		numIron -= g.placeResourceCluster(sector, gmath.ClampMax(clusterSize, numIron), ironSource)
	}
	for numScrap > 0 {
		sector := g.sectors[g.sectorSlider.Value()]
		g.sectorSlider.Inc()
		clusterSize := 1
		if rand.Chance(0.3) {
			clusterSize = 2
		}
		kind := smallScrapSource
		roll := rand.Float()
		if roll > 0.8 {
			kind = bigScrapCreepSource
		} else if roll > 0.4 {
			kind = scrapSource
		}
		numScrap -= g.placeResourceCluster(sector, gmath.ClampMax(clusterSize, numScrap), kind)
	}
	for numGold > 0 {
		clusterSize := rand.IntRange(1, 3)
		sector := g.sectors[g.sectorSlider.Value()]
		g.sectorSlider.Inc()
		numGold -= g.placeResourceCluster(sector, gmath.ClampMax(clusterSize, numGold), goldSource)
	}
	for numOil > 0 {
		sector := g.sectors[g.sectorSlider.Value()]
		g.sectorSlider.Inc()
		numOil -= g.placeResourceCluster(sector, 1, oilSource)
	}
	for numRedOil > 0 {
		sector := g.sectors[g.sectorSlider.Value()]
		g.sectorSlider.Inc()
		numRedOil -= g.placeResourceCluster(sector, 1, redOilSource)
	}
	for numRedCrystals > 0 {
		sector := g.sectors[g.sectorSlider.Value()]
		g.sectorSlider.Inc()
		numRedCrystals -= g.placeResourceCluster(sector, 1, redCrystalSource)
	}
	for numCrystals > 0 {
		clusterSize := 1
		if rand.Chance(0.4) {
			clusterSize = 2
		}
		sector := g.sectors[g.sectorSlider.Value()]
		g.sectorSlider.Inc()
		numCrystals -= g.placeResourceCluster(sector, gmath.ClampMax(clusterSize, numCrystals), crystalSource)
	}

	// If there are no resources near the colony spawn pos,
	// place something in there.
	for _, core := range g.world.colonies {
		hasResources := xslices.ContainsWhere(g.world.essenceSources, func(source *essenceSourceNode) bool {
			// We don't count scraps as some viable starting resource.
			return source.pos.DistanceTo(core.pos) <= core.realRadius &&
				source.stats != scrapSource &&
				source.stats != smallScrapSource
		})
		if !hasResources {
			for i := 0; i < 2; i++ {
				for j := 0; j < 5; j++ {
					pos := g.adjustResourcePos(gmath.RadToVec(rand.Rad()).Mulf(80).Add(core.pos))
					if !posIsFree(g.world, nil, pos, 14) {
						continue
					}
					essence := g.world.NewEssenceSourceNode(ironSource, pos)
					g.pendingResources = append(g.pendingResources, essence)
					break
				}
			}
		}
	}

	// Now sort all resources by their Y coordinate and only
	// then add them to the scene.
	sort.Slice(g.pendingResources, func(i, j int) bool {
		return g.pendingResources[i].pos.Y < g.pendingResources[j].pos.Y
	})
	for _, source := range g.pendingResources {
		g.world.nodeRunner.AddObject(source)
	}
}

func (g *levelGenerator) placeBoss() {
	if !g.world.config.EnemyBoss {
		return
	}

	spawnLocations := []gmath.Vec{
		{X: 196, Y: 196},
		{X: g.world.width - 196, Y: 196},
		{X: 196, Y: g.world.height - 196},
		{X: g.world.width - 196, Y: g.world.height - 196},
	}
	pos := gmath.RandElem(&g.rng, spawnLocations)
	boss := g.world.NewCreepNode(pos, uberBossCreepStats)
	boss.specialDelay = g.rng.FloatRange(3*60, 4*60)
	g.world.nodeRunner.AddObject(boss)

	g.world.boss = boss
}

func (g *levelGenerator) placeCreeps() {
	if g.world.config.InitialCreeps == 0 {
		return
	}

	rand := &g.rng

	worldSizeMultipliers := []float64{
		0.7,
		0.9,
		1.1,
		1.4,
	}
	multiplier := worldSizeMultipliers[g.world.config.WorldSize]
	if g.world.config.InitialCreeps > 1 {
		multiplier *= 2
	}

	g.sectorSlider.TrySetValue(rand.IntRange(0, len(g.sectors)-1))

	numTurrets := int(math.Round(float64(rand.IntRange(4, 5)) * multiplier))
	for numTurrets > 0 {
		sector := g.sectors[g.sectorSlider.Value()]
		g.sectorSlider.Inc()
		numTurrets -= g.placeCreepsCluster(sector, 1, turretCreepStats)
	}

	numCrawlers := int(math.Round(float64(rand.IntRange(8, 12)) * multiplier))
	for numCrawlers > 0 {
		sector := g.sectors[g.sectorSlider.Value()]
		g.sectorSlider.Inc()
		stats := crawlerCreepStats
		if g.rng.Chance(0.4) {
			stats = heavyCrawlerCreepStats
		}
		numCrawlers -= g.placeCreepsCluster(sector, 1, stats)
	}

	numHowitzers := 0
	if g.world.config.InitialCreeps > 1 {
		numHowitzers = 1
	}
	for numHowitzers > 0 {
		sector := g.sectors[g.sectorSlider.Value()]
		g.sectorSlider.Inc()
		numHowitzers -= g.placeCreepsCluster(sector, 1, howitzerCreepStats)
	}
}

func (g *levelGenerator) placeCreepBases() {
	if g.world.config.NumCreepBases == 0 {
		return // Zero bases
	}
	// The bases are always located somewhere on the map boundary.
	// Bases can't be on the same border.
	pad := 128.0
	borderWidth := 440.0
	borders := []gmath.Rect{
		// top border
		{Min: gmath.Vec{X: pad, Y: pad}, Max: gmath.Vec{X: g.world.width - pad, Y: borderWidth + pad}},
		// right border
		{Min: gmath.Vec{X: g.world.width - borderWidth - pad, Y: pad}, Max: gmath.Vec{X: g.world.width - pad, Y: g.world.height - pad}},
		// left border
		{Min: gmath.Vec{X: pad, Y: pad}, Max: gmath.Vec{X: borderWidth + pad, Y: g.world.height - pad}},
		// bottom border
		{Min: gmath.Vec{X: pad, Y: g.world.height - borderWidth - pad}, Max: gmath.Vec{X: g.world.width - pad, Y: g.world.height - pad}},
	}
	gmath.Shuffle(&g.rng, borders)
	numBases := g.world.config.NumCreepBases
	for i := 0; i < numBases; i++ {
		border := borders[i]
		var basePos gmath.Vec
		for round := 0; round < 32; round++ {
			posProbe := g.randomPos(border)
			if posIsFree(g.world, nil, posProbe, 48) {
				basePos = posProbe
				break
			}
		}
		if basePos.IsZero() {
			if g.world.debugLogs {
				fmt.Println("couldn't deploy creep base", i+1)
			}
			continue
		}
		if g.world.debugLogs {
			fmt.Println("deployed a creep base", i+1, "at", basePos, "distance is", basePos.DistanceTo(g.playerSpawn))
		}
		baseRegion := gmath.Rect{
			Min: basePos.Sub(gmath.Vec{X: 96, Y: 96}),
			Max: basePos.Add(gmath.Vec{X: 96, Y: 96}),
		}
		g.placeCreepsCluster(baseRegion, 1, turretCreepStats)
		base := g.world.NewCreepNode(basePos, baseCreepStats)
		if i == 1 || i == 3 {
			base.specialDelay = (9 * 60.0) * g.rng.FloatRange(0.9, 1.1)
		} else {
			base.specialModifier = 1.0 // Initial level base
			base.specialDelay = g.rng.FloatRange(60, 120)
			base.attackDelay = g.rng.FloatRange(40, 50)
		}
		g.world.nodeRunner.AddObject(base)
	}
}

func (g *levelGenerator) placeWalls() {
	rand := &g.rng

	worldSizeMultipliers := []float64{
		0.5,
		0.75,
		1.0,
		1.4,
	}
	multiplier := worldSizeMultipliers[g.world.config.WorldSize]
	numWallClusters := int(float64(rand.IntRange(8, 10)) * multiplier)
	numMountains := int(float64(rand.IntRange(5, 9)) * multiplier)

	const (
		// A simple 1x1 wall tile (rect shape: true).
		wallPit int = iota
		// A straight line shaped wall.
		wallLine
		// Like a line, but may have branches.
		wallSpikedLine
		// A randomly drawed shape.
		wallSnake
		// A cross-like shape.
		wallCrossway
		// Like a spiked line, but more predictable.
		wallZap
	)

	shapePicker := gmath.NewRandPicker[int](rand)
	shapePicker.AddOption(wallPit, 0.05)
	shapePicker.AddOption(wallLine, 0.1)
	shapePicker.AddOption(wallSpikedLine, 0.1)
	shapePicker.AddOption(wallZap, 0.2)
	shapePicker.AddOption(wallSnake, 0.15)
	shapePicker.AddOption(wallCrossway, 0.1)

	directions := []gmath.Vec{
		{X: wallTileSize},
		{X: -wallTileSize},
		{Y: wallTileSize},
		{Y: -wallTileSize},
	}

	chooseRandDirection := func() gmath.Vec {
		roll := rand.IntRange(0, 4)
		var d gmath.Vec
		switch roll {
		case 0:
			d.X = wallTileSize
		case 1:
			d.X = -wallTileSize
		case 2:
			d.Y = wallTileSize
		default:
			d.Y = -wallTileSize
		}
		return d
	}

	reverseDirection := func(d gmath.Vec) gmath.Vec {
		reversed := d
		reversed.X = -d.X
		reversed.Y = -d.Y
		return d
	}

	rotateDirection := func(d gmath.Vec) gmath.Vec {
		rotated := d
		if rand.Bool() {
			rotated.X = d.Y
			rotated.Y = d.X
		} else {
			rotated.X = -d.Y
			rotated.Y = -d.X
		}
		return rotated
	}

	removeDuplicates := func(points []gmath.Vec) []gmath.Vec {
		set := make(map[gmath.Vec]struct{})
		for _, p := range points {
			set[p] = struct{}{}
		}
		filtered := points[:0]
		for p := range set {
			filtered = append(filtered, p)
		}
		return filtered
	}

	g.sectorSlider.TrySetValue(rand.IntRange(0, len(g.sectors)-1))
	for i := 0; i < numWallClusters; i++ {
		sector := g.sectors[g.sectorSlider.Value()]
		g.sectorSlider.Inc()

		var pos gmath.Vec
		for i := 0; i < 10; i++ {
			pos = correctedPos(sector, g.randomPos(sector), 96)
			if posIsFree(g.world, nil, pos, 48) {
				break
			}
		}
		if pos.IsZero() {
			continue
		}

		// Wall positions should be rounded to a tile size.
		{
			x := math.Floor(pos.X/wallTileSize) * wallTileSize
			y := math.Floor(pos.Y/wallTileSize) * wallTileSize
			pos = gmath.Vec{X: x + wallTileSize/2, Y: y + wallTileSize/2}
		}

		var config wallClusterConfig
		shape := shapePicker.Pick()
		switch shape {
		case wallPit:
			config.points = append(config.points, pos)

		case wallCrossway:
			config.points = append(config.points, pos)
			for _, dir := range directions {
				length := rand.IntRange(0, 3)
				currentPos := pos
				for i := 0; i < length; i++ {
					currentPos = currentPos.Add(dir)
					if !posIsFree(g.world, nil, currentPos, 48) {
						break
					}
					config.points = append(config.points, currentPos)
				}
			}

		case wallZap:
			numJoints := rand.IntRange(1, 3) + 1
			dir := chooseRandDirection()
			config.points = append(config.points, pos)
			currentPos := pos
		OuterLoop:
			for i := 0; i < numJoints; i++ {
				length := rand.IntRange(2, 4)
				for j := 0; j < length; j++ {
					currentPos = currentPos.Add(dir)
					if !posIsFree(g.world, nil, currentPos, 48) {
						break OuterLoop
					}
					config.points = append(config.points, currentPos)
					if len(config.points) == maxWallSegments {
						break OuterLoop
					}
				}
				jointPos := currentPos.Add(rotateDirection(dir))
				if !posIsFree(g.world, nil, jointPos, 48) {
					break
				}
				currentPos = jointPos
				config.points = append(config.points, jointPos)
				if len(config.points) == maxWallSegments {
					break OuterLoop
				}
			}

		case wallSnake:
			steps := rand.IntRange(3, maxWallSegments-1)
			config.points = append(config.points, pos)
			currentPos := pos
			dir := chooseRandDirection()
			for i := 0; i < steps; i++ {
				currentPos = currentPos.Add(dir)
				if !posIsFree(g.world, nil, currentPos, 48) {
					break
				}
				config.points = append(config.points, currentPos)
				roll := rand.Float()
				if roll < 0.2 {
					currentPos = pos
				} else if roll < 0.35 {
					dir = chooseRandDirection()
				} else if roll < 0.5 {
					dir = rotateDirection(dir)
				}
			}
			config.points = removeDuplicates(config.points)

		case wallLine, wallSpikedLine:
			spiked := shape == wallSpikedLine
			config.points = append(config.points, pos)
			currentPos := pos
			maxLength := 3
			if spiked {
				maxLength = 6
			}
			lengthRoll := rand.IntRange(1, maxLength)
			length := lengthRoll
			dir := chooseRandDirection()
			for length > 0 {
				currentPos = currentPos.Add(dir)
				if !posIsFree(g.world, nil, currentPos, 48) {
					break
				}
				if spiked && rand.Chance(0.5) {
					sidePos := currentPos.Add(rotateDirection(dir))
					if posIsFree(g.world, nil, sidePos, 48) {
						config.points = append(config.points, sidePos)
					}
				}
				config.points = append(config.points, currentPos)
				length--
			}
			if length > 0 {
				dir = reverseDirection(dir)
				for length > 0 {
					currentPos = currentPos.Add(dir)
					if !posIsFree(g.world, nil, currentPos, 48) {
						break
					}
					config.points = append(config.points, currentPos)
					length--
				}
			}
		}

		if len(config.points) != 0 {
			config.atlas = wallAtras{layers: landcrackAtlas}
			config.world = g.world
			wall := g.world.NewWallClusterNode(config)
			g.scene.AddObject(wall)
		}
	}

	chunkSizePicker := gmath.NewRandPicker[mountainKind](rand)
	chunkSizePicker.AddOption(mountainSmall, 0.15)
	chunkSizePicker.AddOption(mountainMedium, 0.25)
	chunkSizePicker.AddOption(mountainBig, 0.45)
	chunkSizePicker.AddOption(mountainWide, 0.1)
	chunkSizePicker.AddOption(mountainTall, 0.05)
	sideChunkSizePicker := gmath.NewRandPicker[mountainKind](rand)
	sideChunkSizePicker.AddOption(mountainSmall, 0.45)
	sideChunkSizePicker.AddOption(mountainMedium, 0.3)
	sideChunkSizePicker.AddOption(mountainBig, 0.1)
	sideChunkSizePicker.AddOption(mountainWide, 0.1)
	sideChunkSizePicker.AddOption(mountainTall, 0.05)
	for i := 0; i < numMountains; i++ {
		sector := g.sectors[g.sectorSlider.Value()]
		g.sectorSlider.Inc()

		var pos gmath.Vec
		for i := 0; i < 10; i++ {
			pos = correctedPos(sector, g.randomPos(sector), 196)
			if posIsFree(g.world, nil, pos, 80) {
				break
			}
		}
		if pos.IsZero() {
			continue
		}

		// Positions should be rounded to a tile size.
		{
			x := math.Floor(pos.X/wallTileSize) * wallTileSize
			y := math.Floor(pos.Y/wallTileSize) * wallTileSize
			pos = gmath.Vec{X: x + wallTileSize/2, Y: y + wallTileSize/2}
		}

		currentPos := pos
		mountainLength := g.rng.IntRange(3, 10)
		widthModifier := g.rng.Float()
		dir := chooseRandDirection()
		chunks := make([]wallChunk, 0, mountainLength+4)
		for j := 0; j < mountainLength; j++ {
			chunkSize := chunkSizePicker.Pick()
			chunks = append(chunks, wallChunk{
				pos:  currentPos,
				kind: chunkSize,
			})
			rotateRoll := g.rng.Float()
			if rotateRoll < 0.1 {
				dir = rotateDirection(dir)
			} else if rotateRoll < 0.45 {
				extraPos := currentPos.Add(rotateDirection(dir))
				if posIsFree(g.world, nil, extraPos, 64) {
					chunks = append(chunks, wallChunk{
						pos:  extraPos,
						kind: sideChunkSizePicker.Pick(),
					})
				}
			} else if rotateRoll < widthModifier {
				parallelDir := dir
				parallelDir.X = dir.Y
				parallelDir.Y = dir.X
				extraPos := currentPos.Add(parallelDir)
				if posIsFree(g.world, nil, extraPos, 64) {
					chunks = append(chunks, wallChunk{
						pos:  extraPos,
						kind: chunkSizePicker.Pick(),
					})
				}
			}
			currentPos = currentPos.Add(dir)
			if !posIsFree(g.world, nil, currentPos, 80) {
				break
			}
		}
		var config wallClusterConfig
		config.chunks = chunks
		config.world = g.world
		wall := g.world.NewWallClusterNode(config)
		g.scene.AddObject(wall)
	}
}
