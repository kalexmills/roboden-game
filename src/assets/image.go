package assets

import (
	"runtime"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func RegisterImageResources(ctx *ge.Context, progress *float64) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageAchievementImpossible:    {Path: "image/achievement/impossible.png"},
		ImageAchievementCheapBuild10:  {Path: "image/achievement/cheapbuild10.png"},
		ImageAchievementT3Engineer:    {Path: "image/achievement/t3engineer.png"},
		ImageAchievementHighTension:   {Path: "image/achievement/hightension.png"},
		ImageAchievementSoloBase:      {Path: "image/achievement/solobase.png"},
		ImageAchievementUILess:        {Path: "image/achievement/uiless.png"},
		ImageAchievementTinyRadius:    {Path: "image/achievement/tinyradius.png"},
		ImageAchievementT1Army:        {Path: "image/achievement/t1army.png"},
		ImageAchievementGroundWin:     {Path: "image/achievement/groundwin.png"},
		ImageAchievementSpeedrunning:  {Path: "image/achievement/speedrunning.png"},
		ImageAchievementVictoryDrag:   {Path: "image/achievement/victorydrag.png"},
		ImageAchievementT3Less:        {Path: "image/achievement/t3less.png"},
		ImageAchievementTurretDamage:  {Path: "image/achievement/turretdamage.png"},
		ImageAchievementPowerOf3:      {Path: "image/achievement/powerof3.png"},
		ImageAchievementInfinite:      {Path: "image/achievement/infinite.png"},
		ImageAchievementAntiDominator: {Path: "image/achievement/antidominator.png"},
		ImageAchievementTrample:       {Path: "image/achievement/trample.png"},
		ImageAchievementNoPeeking:     {Path: "image/achievement/no_peeking.png"},

		ImageLock: {Path: "image/ui/lock.png"},

		ImageSmallShadow:      {Path: "image/shadows/small_shadow.png"},
		ImageMediumShadow:     {Path: "image/shadows/medium_shadow.png"},
		ImageBigShadow:        {Path: "image/shadows/big_shadow.png"},
		ImageUberBossShadow:   {Path: "image/shadows/uber_boss_shadow.png"},
		ImageColonyCoreShadow: {Path: "image/shadows/colony_core_shadow.png"},

		ImageCursor: {Path: "image/cursor.png"},

		ImageRadarlessButtons:  {Path: "image/ui/radarless_buttons.png"},
		ImageRadar:             {Path: "image/ui/radar.png"},
		ImageRadarWave:         {Path: "image/ui/radar_wave.png"},
		ImageRadarBossFar:      {Path: "image/ui/radar_boss_far.png"},
		ImageRadarBossNear:     {Path: "image/ui/radar_boss_near.png"},
		ImageRightPanelLayer1:  {Path: "image/ui/right_panel_layer1.png"},
		ImageRightPanelLayer2:  {Path: "image/ui/right_panel_layer2.png"},
		ImagePriorityBar:       {Path: "image/ui/priority_bar.png"},
		ImagePriorityResources: {Path: "image/ui/priority_icon_resources.png"},
		ImagePriorityGrowth:    {Path: "image/ui/priority_icon_growth.png"},
		ImagePriorityEvolution: {Path: "image/ui/priority_icon_evolution.png"},
		ImagePrioritySecurity:  {Path: "image/ui/priority_icon_security.png"},
		ImageObjectiveDisplay:  {Path: "image/ui/objective_display.png"},

		ImageFloppyYellow:     {Path: "image/ui/floppy_yellow.png"},
		ImageFloppyRed:        {Path: "image/ui/floppy_red.png"},
		ImageFloppyGreen:      {Path: "image/ui/floppy_green.png"},
		ImageFloppyBlue:       {Path: "image/ui/floppy_blue.png"},
		ImageFloppyGray:       {Path: "image/ui/floppy_gray.png"},
		ImageFloppyYellowFlip: {Path: "image/ui/floppy_yellow_flip.png", FrameWidth: 144},
		ImageFloppyRedFlip:    {Path: "image/ui/floppy_red_flip.png", FrameWidth: 144},
		ImageFloppyGreenFlip:  {Path: "image/ui/floppy_green_flip.png", FrameWidth: 144},
		ImageFloppyBlueFlip:   {Path: "image/ui/floppy_blue_flip.png", FrameWidth: 144},
		ImageFloppyGrayFlip:   {Path: "image/ui/floppy_gray_flip.png", FrameWidth: 144},

		ImageActionBuildColony:    {Path: "image/ui/action_build_colony.png"},
		ImageActionBuildTurret:    {Path: "image/ui/action_build_turret.png"},
		ImageActionAttack:         {Path: "image/ui/action_attack.png"},
		ImageActionIncreaseRadius: {Path: "image/ui/action_increase_radius.png"},
		ImageActionDecreaseRadius: {Path: "image/ui/action_decrease_radius.png"},

		ImageProjectileSmoke:            {Path: "image/effects/projectile_smoke.png", FrameWidth: 8},
		ImageStealthLaserExplosion:      {Path: "image/effects/stealth_laser_explosion.png", FrameWidth: 14},
		ImageCripplerBlasterExplosion:   {Path: "image/effects/crippler_blaster_explosion.png", FrameWidth: 8},
		ImageScoutIonExplosion:          {Path: "image/effects/scout_ion_explosion.png", FrameWidth: 5},
		ImageShockerExplosion:           {Path: "image/effects/shocker_explosion.png", FrameWidth: 8},
		ImageFighterLaserExplosion:      {Path: "image/effects/fighter_laser_explosion.png", FrameWidth: 14},
		ImageHeavyCrawlerLaserExplosion: {Path: "image/effects/heavy_crawler_laser_explosion.png", FrameWidth: 14},
		ImageSmallExplosion1:            {Path: "image/effects/small_explosion1.png", FrameWidth: 32},
		ImagePurpleExplosion:            {Path: "image/effects/purple_explosion.png", FrameWidth: 40},
		ImageVerticalExplosion:          {Path: "image/effects/vertical_explosion.png", FrameWidth: 50},
		ImageBigVerticalExplosion:       {Path: "image/effects/big_vertical_explosion.png", FrameWidth: 38},
		ImageBigExplosion:               {Path: "image/effects/big_explosion.png", FrameWidth: 64},
		ImageIonZap:                     {Path: "image/effects/ion_zap.png", FrameWidth: 28},
		ImagePurpleIonZap:               {Path: "image/effects/purple_ion_zap.png", FrameWidth: 28},
		ImageGreenZap:                   {Path: "image/effects/green_zap.png", FrameWidth: 14},
		ImageCloakWave:                  {Path: "image/effects/cloak_wave.png", FrameWidth: 28},
		ImageServantWave:                {Path: "image/effects/servant_wave.png", FrameWidth: 64},
		ImageSmokeDown:                  {Path: "image/effects/smoke_down.png", FrameWidth: 8},
		ImageSmokeSideDown:              {Path: "image/effects/smoke_side_down.png", FrameWidth: 8},
		ImageSmokeSide:                  {Path: "image/effects/smoke_side.png", FrameWidth: 15},

		ImageFactionDiode:       {Path: "image/faction_diode.png"},
		ImageColonyResourceBar1: {Path: "image/colony_resource_bar1.png"},
		ImageColonyResourceBar2: {Path: "image/colony_resource_bar2.png"},
		ImageColonyResourceBar3: {Path: "image/colony_resource_bar3.png"},
		ImageColonyCoreSelector: {Path: "image/colony_core_selector.png"},
		ImageColonyCore:         {Path: "image/colony_core.png"},
		ImageColonyCoreFlying:   {Path: "image/colony_core_flying.png"},
		ImageColonyCoreHatch:    {Path: "image/colony_core_hatch.png"},
		ImageColonyCoreDiode:    {Path: "image/colony_core_diode.png", FrameWidth: 4},

		ImageGunpointAgent:      {Path: "image/drones/gunpoint_agent.png"},
		ImageBeamtowerAgent:     {Path: "image/drones/beamtower_agent.png"},
		ImageTetherBeaconAgent:  {Path: "image/drones/tether_beacon_agent.png"},
		ImageWorkerAgent:        {Path: "image/drones/worker_agent.png", FrameWidth: 9, FrameHeight: 10},
		ImageScoutAgent:         {Path: "image/drones/scout_agent.png", FrameWidth: 11, FrameHeight: 13},
		ImageClonerAgent:        {Path: "image/drones/cloner_agent.png", FrameWidth: 13, FrameHeight: 13},
		ImageScavengerAgent:     {Path: "image/drones/scavenger_agent.png", FrameWidth: 15, FrameHeight: 12},
		ImageCourierAgent:       {Path: "image/drones/courier_agent.png", FrameWidth: 15, FrameHeight: 15},
		ImageDisintegratorAgent: {Path: "image/drones/disintegrator_agent.png", FrameWidth: 17, FrameHeight: 15},
		ImageTruckerAgent:       {Path: "image/drones/trucker_agent.png", FrameWidth: 27, FrameHeight: 22},
		ImageMarauderAgent:      {Path: "image/drones/marauder_agent.png", FrameWidth: 29, FrameHeight: 20},
		ImageMortarAgent:        {Path: "image/drones/mortar_agent.png", FrameWidth: 21, FrameHeight: 18},
		ImageCripplerAgent:      {Path: "image/drones/crippler_agent.png", FrameWidth: 15, FrameHeight: 15},
		ImageStormbringerAgent:  {Path: "image/drones/stormbringer_agent.png", FrameWidth: 21, FrameHeight: 19},
		ImageRepairAgent:        {Path: "image/drones/repair_agent.png", FrameWidth: 17, FrameHeight: 13},
		ImageAntiAirAgent:       {Path: "image/drones/antiair_agent.png", FrameWidth: 17, FrameHeight: 19},
		ImageServoAgent:         {Path: "image/drones/servo_agent.png", FrameWidth: 15, FrameHeight: 22},
		ImageRechargerAgent:     {Path: "image/drones/recharger_agent.png", FrameWidth: 17, FrameHeight: 20},
		ImageGuardianAgent:      {Path: "image/drones/guardian_agent.png", FrameWidth: 31, FrameHeight: 21},
		ImageFighterAgent:       {Path: "image/drones/fighter_agent.png", FrameWidth: 15, FrameHeight: 15},
		ImageDefenderAgent:      {Path: "image/drones/defender_agent.png", FrameWidth: 17, FrameHeight: 16},
		ImageKamikazeAgent:      {Path: "image/drones/kamikaze_agent.png", FrameWidth: 15, FrameHeight: 12},
		ImagePrismAgent:         {Path: "image/drones/prism_agent.png", FrameWidth: 15, FrameHeight: 15},
		ImageDestroyerAgent:     {Path: "image/drones/destroyer_agent.png", FrameWidth: 33, FrameHeight: 24},
		ImageRepellerAgent:      {Path: "image/drones/repeller_agent.png", FrameWidth: 15, FrameHeight: 13},
		ImageFreighterAgent:     {Path: "image/drones/freighter_agent.png", FrameWidth: 17, FrameHeight: 16},
		ImageRedminerAgent:      {Path: "image/drones/redminer_agent.png", FrameWidth: 13, FrameHeight: 18},
		ImageGeneratorAgent:     {Path: "image/drones/generator_agent.png", FrameWidth: 15, FrameHeight: 17},
		ImageSkirmisherAgent:    {Path: "image/drones/skirmisher_agent.png", FrameWidth: 19, FrameHeight: 12},

		ImageColonyDamageMask:  {Path: "image/shaders/colony_damage_mask.png"},
		ImageTurretDamageMask1: {Path: "image/shaders/turret_damage_mask1.png"},
		ImageTurretDamageMask2: {Path: "image/shaders/turret_damage_mask2.png"},
		ImageTurretDamageMask3: {Path: "image/shaders/turret_damage_mask3.png"},
		ImageTurretDamageMask4: {Path: "image/shaders/turret_damage_mask4.png"},

		ImageEssenceSourceDissolveMask:    {Path: "image/resources/essence_source_dissolve_mask.png"},
		ImageEssenceRedCrystalSource:      {Path: "image/resources/red_crystal.png"},
		ImageEssenceCrystalSource:         {Path: "image/resources/crystal_source.png", FrameWidth: 16},
		ImageEssenceGoldSource:            {Path: "image/resources/gold_source.png", FrameWidth: 28},
		ImageEssenceIronSource:            {Path: "image/resources/iron_source.png", FrameWidth: 32},
		ImageEssenceScrapSource:           {Path: "image/resources/scrap_source.png"},
		ImageEssenceSmallScrapSource:      {Path: "image/resources/small_scrap_source.png"},
		ImageEssenceScrapCreepSource:      {Path: "image/resources/scrap_source_creep.png"},
		ImageEssenceSmallScrapCreepSource: {Path: "image/resources/small_scrap_source_creep.png"},
		ImageEssenceBigScrapCreepSource:   {Path: "image/resources/big_scrap_source_creep.png"},
		ImageEssenceSource:                {Path: "image/resources/essence_source.png", FrameWidth: 32},
		ImageRedEssenceSource:             {Path: "image/resources/red_essence_source.png", FrameWidth: 32},

		ImageHowitzerCreep:       {Path: "image/creeps/howitzer_creep.png", FrameWidth: 37},
		ImageHowitzerPreparing:   {Path: "image/creeps/howitzer_preparing.png", FrameWidth: 37},
		ImageHowitzerTrunk:       {Path: "image/creeps/howitzer_trunk.png", FrameWidth: 23},
		ImageHeavyCrawlerCreep:   {Path: "image/creeps/heavy_crawler_creep.png", FrameWidth: 25},
		ImageStealthCrawlerCreep: {Path: "image/creeps/stealth_crawler_creep.png", FrameWidth: 19},
		ImageEliteCrawlerCreep:   {Path: "image/creeps/elite_crawler_creep.png", FrameWidth: 23},
		ImageCrawlerCreep:        {Path: "image/creeps/crawler_creep.png", FrameWidth: 23},
		ImagePrimitiveCreep:      {Path: "image/creeps/tier1_creep.png"},
		ImageServantCreep:        {Path: "image/creeps/servant_creep.png", FrameWidth: 15},
		ImageCreepTier2:          {Path: "image/creeps/tier2_creep.png"},
		ImageCreepTier3:          {Path: "image/creeps/tier3_creep.png", FrameWidth: 25},
		ImageCreepDominator:      {Path: "image/creeps/dominator_creep.png", FrameWidth: 23},
		ImageTurretCreep:         {Path: "image/creeps/turret_creep.png"},
		ImageUberBoss:            {Path: "image/creeps/uber_boss.png", FrameWidth: 40},
		ImageUberBossOpen:        {Path: "image/creeps/uber_boss_open.png"},
		ImageCreepBase:           {Path: "image/creeps/creep_base.png", FrameWidth: 32},
		ImageCrawlerCreepBase:    {Path: "image/creeps/crawler_base_creep.png"},
		ImageBuilderCreep:        {Path: "image/creeps/builder_creep.png", FrameWidth: 31, FrameHeight: 31},

		ImageBackgroundTiles: {Path: "image/landscape/tiles.png"},
		ImageMountainSmall:   {Path: "image/landscape/mountain_small.png", FrameWidth: 32},
		ImageMountainMedium:  {Path: "image/landscape/mountain_medium.png", FrameWidth: 48},
		ImageMountainBig:     {Path: "image/landscape/mountain_big.png", FrameWidth: 64},
		ImageMountainWide:    {Path: "image/landscape/mountain_wide.png", FrameWidth: 64},
		ImageMountainTall:    {Path: "image/landscape/mountain_tall.png", FrameWidth: 48},
		ImageLandCrack:       {Path: "image/landscape/landcrack.png", FrameWidth: 32},
		ImageLandCrack2:      {Path: "image/landscape/landcrack2.png", FrameWidth: 32},
		ImageLandCrack3:      {Path: "image/landscape/landcrack3.png", FrameWidth: 32},
		ImageLandCrack4:      {Path: "image/landscape/landcrack4.png", FrameWidth: 32},

		ImageSkirmisherProjectile:     {Path: "image/projectile/skirmisher_projectile.png"},
		ImageHeavyCrawlerProjectile:   {Path: "image/projectile/heavy_crawler_projectile.png"},
		ImageStealthCrawlerProjectile: {Path: "image/projectile/stealth_crawler_projectile.png"},
		ImageEliteCrawlerProjectile:   {Path: "image/projectile/elite_crawler_projectile.png"},
		ImageTankProjectile:           {Path: "image/projectile/tank_projectile.png"},
		ImageAssaultProjectile:        {Path: "image/projectile/assault_projectile.png"},
		ImageCripplerProjectile:       {Path: "image/projectile/crippler_projectile.png"},
		ImageScoutProjectile:          {Path: "image/projectile/scout_projectile.png"},
		ImageRepellerProjectile:       {Path: "image/projectile/repeller_projectile.png"},
		ImageGunpointProjectile:       {Path: "image/projectile/gunpoint_projectile.png"},
		ImageFighterProjectile:        {Path: "image/projectile/fighter_projectile.png"},
		ImageScavengerProjectile:      {Path: "image/projectile/scavenger_projectile.png"},
		ImageMarauderProjectile:       {Path: "image/projectile/marauder_projectile.png"},
		ImageCourierProjectile:        {Path: "image/projectile/courier_projectile.png"},
		ImageDisintegratorProjectile:  {Path: "image/projectile/disintegrator_projectile.png"},
		ImageServantProjectile:        {Path: "image/projectile/servant_projectile.png"},
		ImageWandererProjectile:       {Path: "image/projectile/wanderer_projectile.png"},
		ImageStormbringerProjectile:   {Path: "image/projectile/stormbringer_projectile.png"},
		ImageMortarProjectile:         {Path: "image/projectile/mortar_projectile.png"},
		ImageHowitzerProjectile:       {Path: "image/projectile/howitzer_projectile.png"},
		ImageHowitzerLaserProjectile:  {Path: "image/projectile/howitzer_laser_projectile.png"},
		ImageAntiAirMissile:           {Path: "image/projectile/aa_missile.png"},
		ImageMissile:                  {Path: "image/projectile/missile.png"},

		ImageStunnerLine:   {Path: "image/lines/stunner_line.png"},
		ImageBossLaserLine: {Path: "image/lines/boss_laser_line.png"},
		ImageRepairLine:    {Path: "image/lines/repair_line.png"},
		ImageRechargerLine: {Path: "image/lines/recharger_line.png"},
		ImageDefenderLine:  {Path: "image/lines/defender_line.png"},
		ImageBeamtowerLine: {Path: "image/lines/beamtower_line.png"},
		ImageTetherLine:    {Path: "image/lines/tether_line.png"},

		ImageUIButtonIdle:               {Path: "image/ebitenui/button-idle.png"},
		ImageUIButtonHover:              {Path: "image/ebitenui/button-hover.png"},
		ImageUIButtonPressed:            {Path: "image/ebitenui/button-pressed.png"},
		ImageUIButtonDisabled:           {Path: "image/ebitenui/button-disabled.png"},
		ImageUITabButtonIdle:            {Path: "image/ebitenui/tabbutton-idle.png"},
		ImageUITabButtonHover:           {Path: "image/ebitenui/tabbutton-hover.png"},
		ImageUITabButtonPressed:         {Path: "image/ebitenui/tabbutton-pressed.png"},
		ImageUITabButtonDisabled:        {Path: "image/ebitenui/tabbutton-disabled.png"},
		ImageUITextInputIdle:            {Path: "image/ebitenui/text-input-idle.png"},
		ImageUITextInputDisabled:        {Path: "image/ebitenui/text-input-disabled.png"},
		ImageUIItemButtonIdle:           {Path: "image/ebitenui/itembutton-idle.png"},
		ImageUIItemButtonHover:          {Path: "image/ebitenui/itembutton-hover.png"},
		ImageUIItemButtonPressed:        {Path: "image/ebitenui/itembutton-pressed.png"},
		ImageUIItemButtonDisabled:       {Path: "image/ebitenui/itembutton-disabled.png"},
		ImageUIAltItemButtonIdle:        {Path: "image/ebitenui/itembutton-alt-idle.png"},
		ImageUIAltItemButtonHover:       {Path: "image/ebitenui/itembutton-alt-hover.png"},
		ImageUIAltItemButtonPressed:     {Path: "image/ebitenui/itembutton-alt-pressed.png"},
		ImageUIAltItemButtonDisabled:    {Path: "image/ebitenui/itembutton-alt-disabled.png"},
		ImageUIBigItemButtonIdle:        {Path: "image/ebitenui/bigitembutton-idle.png"},
		ImageUIBigItemButtonHover:       {Path: "image/ebitenui/bigitembutton-hover.png"},
		ImageUIBigItemButtonPressed:     {Path: "image/ebitenui/bigitembutton-pressed.png"},
		ImageUIBigItemButtonDisabled:    {Path: "image/ebitenui/bigitembutton-disabled.png"},
		ImageUIAltBigItemButtonIdle:     {Path: "image/ebitenui/bigitembutton-alt-idle.png"},
		ImageUIAltBigItemButtonHover:    {Path: "image/ebitenui/bigitembutton-alt-hover.png"},
		ImageUIAltBigItemButtonPressed:  {Path: "image/ebitenui/bigitembutton-alt-pressed.png"},
		ImageUIAltBigItemButtonDisabled: {Path: "image/ebitenui/bigitembutton-alt-disabled.png"},
		ImageUIButtonSelectedIdle:       {Path: "image/ebitenui/button-selected-idle.png"},
		ImageUIButtonSelectedHover:      {Path: "image/ebitenui/button-selected-hover.png"},
		ImageUIButtonSelectedPressed:    {Path: "image/ebitenui/button-selected-pressed.png"},
		ImageUIPanelIdle:                {Path: "image/ebitenui/panel-idle.png"},
	}

	singleThread := runtime.GOMAXPROCS(-1) == 1
	progressPerItem := 1.0 / float64(len(imageResources))
	for id, res := range imageResources {
		ctx.Loader.ImageRegistry.Set(id, res)
		ctx.Loader.LoadImage(id)
		*progress += progressPerItem
		if singleThread {
			runtime.Gosched()
		}
	}
}

const (
	ImageNone resource.ImageID = iota

	ImageAchievementAntiDominator
	ImageAchievementImpossible
	ImageAchievementCheapBuild10
	ImageAchievementT3Engineer
	ImageAchievementHighTension
	ImageAchievementSoloBase
	ImageAchievementUILess
	ImageAchievementTinyRadius
	ImageAchievementT1Army
	ImageAchievementGroundWin
	ImageAchievementSpeedrunning
	ImageAchievementVictoryDrag
	ImageAchievementT3Less
	ImageAchievementTurretDamage
	ImageAchievementPowerOf3
	ImageAchievementInfinite
	ImageAchievementNoPeeking
	ImageAchievementTrample

	ImageLock

	ImageSmallShadow
	ImageMediumShadow
	ImageBigShadow

	ImageCursor

	ImageRadarlessButtons
	ImageRadar
	ImageRadarWave
	ImageRadarBossFar
	ImageRadarBossNear
	ImageRightPanelLayer1
	ImageRightPanelLayer2
	ImagePriorityBar
	ImagePriorityResources
	ImagePriorityGrowth
	ImagePriorityEvolution
	ImagePrioritySecurity
	ImageObjectiveDisplay

	ImageProjectileSmoke
	ImageStealthLaserExplosion
	ImageCripplerBlasterExplosion
	ImageScoutIonExplosion
	ImageShockerExplosion
	ImageFighterLaserExplosion
	ImageHeavyCrawlerLaserExplosion
	ImageSmallExplosion1
	ImagePurpleExplosion
	ImageVerticalExplosion
	ImageBigVerticalExplosion
	ImageBigExplosion
	ImageIonZap
	ImagePurpleIonZap
	ImageGreenZap
	ImageCloakWave
	ImageServantWave
	ImageSmokeDown
	ImageSmokeSideDown
	ImageSmokeSide

	ImageColonyDamageMask
	ImageTurretDamageMask1
	ImageTurretDamageMask2
	ImageTurretDamageMask3
	ImageTurretDamageMask4
	ImageEssenceSourceDissolveMask

	ImageFloppyYellow
	ImageFloppyRed
	ImageFloppyGreen
	ImageFloppyBlue
	ImageFloppyGray
	ImageFloppyYellowFlip
	ImageFloppyRedFlip
	ImageFloppyGreenFlip
	ImageFloppyBlueFlip
	ImageFloppyGrayFlip

	ImageActionBuildColony
	ImageActionBuildTurret
	ImageActionAttack
	ImageActionIncreaseRadius
	ImageActionDecreaseRadius

	ImageFactionDiode
	ImageUberBoss
	ImageUberBossOpen
	ImageUberBossShadow
	ImageCreepBase
	ImageColonyResourceBar1
	ImageColonyResourceBar2
	ImageColonyResourceBar3
	ImageColonyCoreSelector
	ImageColonyCore
	ImageColonyCoreFlying
	ImageColonyCoreHatch
	ImageColonyCoreDiode
	ImageColonyCoreShadow
	ImageGunpointAgent
	ImageBeamtowerAgent
	ImageTetherBeaconAgent
	ImageWorkerAgent
	ImageGeneratorAgent
	ImageSkirmisherAgent
	ImageScoutAgent
	ImageClonerAgent
	ImageScavengerAgent
	ImageCourierAgent
	ImageDisintegratorAgent
	ImageMarauderAgent
	ImageTruckerAgent
	ImageStormbringerAgent
	ImageAntiAirAgent
	ImagePrismAgent
	ImageRepairAgent
	ImageRechargerAgent
	ImageGuardianAgent
	ImageFighterAgent
	ImageDefenderAgent
	ImageKamikazeAgent
	ImageDestroyerAgent
	ImageCripplerAgent
	ImageMortarAgent
	ImageRedminerAgent
	ImageRepellerAgent
	ImageServoAgent
	ImageFreighterAgent
	ImageEssenceRedCrystalSource
	ImageEssenceCrystalSource
	ImageEssenceGoldSource
	ImageEssenceIronSource
	ImageEssenceScrapSource
	ImageEssenceSmallScrapSource
	ImageEssenceScrapCreepSource
	ImageEssenceSmallScrapCreepSource
	ImageEssenceBigScrapCreepSource
	ImageEssenceSource
	ImageRedEssenceSource
	ImageCrawlerCreep
	ImageEliteCrawlerCreep
	ImageHeavyCrawlerCreep
	ImageHowitzerCreep
	ImageHowitzerPreparing
	ImageHowitzerTrunk
	ImageStealthCrawlerCreep
	ImagePrimitiveCreep
	ImageServantCreep
	ImageCreepTier2
	ImageCreepTier3
	ImageCreepDominator
	ImageTurretCreep
	ImageBuilderCreep
	ImageCrawlerCreepBase

	ImageBackgroundTiles
	ImageMountainSmall
	ImageMountainMedium
	ImageMountainBig
	ImageMountainTall
	ImageMountainWide
	ImageLandCrack
	ImageLandCrack2
	ImageLandCrack3
	ImageLandCrack4

	ImageSkirmisherProjectile
	ImageHeavyCrawlerProjectile
	ImageStealthCrawlerProjectile
	ImageEliteCrawlerProjectile
	ImageTankProjectile
	ImageAssaultProjectile
	ImageCripplerProjectile
	ImageScoutProjectile
	ImageRepellerProjectile
	ImageGunpointProjectile
	ImageFighterProjectile
	ImageScavengerProjectile
	ImageMarauderProjectile
	ImageCourierProjectile
	ImageDisintegratorProjectile
	ImageWandererProjectile
	ImageServantProjectile
	ImageStormbringerProjectile
	ImageMortarProjectile
	ImageHowitzerProjectile
	ImageHowitzerLaserProjectile
	ImageAntiAirMissile
	ImageMissile

	ImageBossLaserLine
	ImageStunnerLine
	ImageRepairLine
	ImageRechargerLine
	ImageDefenderLine
	ImageBeamtowerLine
	ImageTetherLine

	ImageUIButtonIdle
	ImageUIButtonHover
	ImageUIButtonPressed
	ImageUIButtonDisabled
	ImageUITabButtonIdle
	ImageUITabButtonHover
	ImageUITabButtonPressed
	ImageUITabButtonDisabled
	ImageUITextInputIdle
	ImageUITextInputDisabled
	ImageUIItemButtonIdle
	ImageUIItemButtonHover
	ImageUIItemButtonPressed
	ImageUIItemButtonDisabled
	ImageUIAltItemButtonIdle
	ImageUIAltItemButtonHover
	ImageUIAltItemButtonPressed
	ImageUIAltItemButtonDisabled
	ImageUIBigItemButtonIdle
	ImageUIBigItemButtonHover
	ImageUIBigItemButtonPressed
	ImageUIBigItemButtonDisabled
	ImageUIAltBigItemButtonIdle
	ImageUIAltBigItemButtonHover
	ImageUIAltBigItemButtonPressed
	ImageUIAltBigItemButtonDisabled
	ImageUIButtonSelectedIdle
	ImageUIButtonSelectedHover
	ImageUIButtonSelectedPressed
	ImageUIPanelIdle
)
