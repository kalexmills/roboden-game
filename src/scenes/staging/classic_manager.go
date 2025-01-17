package staging

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type classicManager struct {
	world *worldState

	spawnAreas []gmath.Rect

	scene *ge.Scene

	tier3spawnDelay float64
	tier3spawnRate  float64

	crawlersDelay float64
}

func newClassicManager(world *worldState) *classicManager {
	return &classicManager{
		world:          world,
		tier3spawnRate: 1,
	}
}

func (m *classicManager) Init(scene *ge.Scene) {
	m.scene = scene

	m.spawnAreas = creepSpawnAreas(m.world)

	// Start launching tier3 creeps after ~15 minutes.
	m.tier3spawnDelay = m.world.rand.FloatRange(14*60.0, 16*60.0)

	// Extra crawlers show up around the 10th minute.
	m.crawlersDelay = m.world.rand.FloatRange(8*60.0, 12*60.0)
}

func (m *classicManager) IsDisposed() bool {
	return false
}

func (m *classicManager) Update(delta float64) {
	m.tier3spawnDelay = gmath.ClampMin(m.tier3spawnDelay-delta, 0)
	if m.tier3spawnDelay == 0 {
		m.spawnTier3Creep()
	}
	m.crawlersDelay = gmath.ClampMin(m.crawlersDelay-delta, 0)
	if m.crawlersDelay == 0 {
		m.spawnCrawlers()
	}
}

func (m *classicManager) spawnCrawlers() {
	m.crawlersDelay = m.world.rand.FloatRange(50, 100)

	sector := gmath.RandElem(m.world.rand, m.spawnAreas)
	spawnPos := randomSectorPos(m.world.rand, sector)
	targetPos := correctedPos(m.world.rect, randomSectorPos(m.world.rand, sector), 520)

	numCreeps := 1
	creepStats := howitzerCreepStats
	if m.world.rand.Chance(0.65) {
		numCreeps = m.world.rand.IntRange(3, 7)
		creepStats = stealthCrawlerCreepStats
	}

	for i := 0; i < numCreeps; i++ {
		creepPos, spawnDelay := groundCreepSpawnPos(m.world, spawnPos, creepStats)
		creepTargetPos := targetPos.Add(m.world.rand.Offset(-60, 60))
		if spawnDelay > 0 {
			spawner := newCreepSpawnerNode(m.world, spawnDelay, creepPos, creepTargetPos, creepStats)
			m.world.nodeRunner.AddObject(spawner)
		} else {
			creep := m.world.NewCreepNode(creepPos, creepStats)
			m.world.nodeRunner.AddObject(creep)
			creep.SendTo(creepTargetPos)
		}
	}
}

func (m *classicManager) spawnTier3Creep() {
	m.tier3spawnRate = gmath.ClampMin(m.tier3spawnRate-0.025, 0.35)
	m.tier3spawnDelay = m.world.rand.FloatRange(55, 80) * m.tier3spawnRate

	var spawnPos gmath.Vec
	roll := m.world.rand.Float()
	if roll < 0.25 {
		spawnPos.X = m.world.width - 4
		spawnPos.Y = m.world.rand.FloatRange(0, m.world.height)
	} else if roll < 0.5 {
		spawnPos.X = m.world.rand.FloatRange(0, m.world.width)
		spawnPos.Y = m.world.height - 4
	} else if roll < 0.75 {
		spawnPos.X = 4
		spawnPos.Y = m.world.rand.FloatRange(0, m.world.height)
	} else {
		spawnPos.X = m.world.rand.FloatRange(0, m.world.width)
		spawnPos.Y = 4
	}
	spawnPos = roundedPos(spawnPos)
	stats := assaultCreepStats
	if m.world.rand.Chance(0.3) {
		stats = builderCreepStats
	}
	creep := m.world.NewCreepNode(spawnPos, stats)
	m.world.nodeRunner.AddObject(creep)
}
