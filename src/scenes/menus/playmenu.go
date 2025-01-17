package menus

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/roboden-game/assets"
	"github.com/quasilyte/roboden-game/controls"
	"github.com/quasilyte/roboden-game/gamedata"
	"github.com/quasilyte/roboden-game/gameui/eui"
	"github.com/quasilyte/roboden-game/session"
)

type PlayMenuController struct {
	state *session.State

	scene *ge.Scene
}

func NewPlayMenuController(state *session.State) *PlayMenuController {
	return &PlayMenuController{state: state}
}

func (c *PlayMenuController) Init(scene *ge.Scene) {
	c.scene = scene
	c.initUI()
}

func (c *PlayMenuController) Update(delta float64) {
	if c.state.MainInput.ActionIsJustPressed(controls.ActionBack) {
		c.back()
		return
	}
}

func (c *PlayMenuController) initUI() {
	uiResources := c.state.Resources.UI

	root := eui.NewAnchorContainer()
	rowContainer := eui.NewRowLayoutContainer(10, nil)
	root.AddChild(rowContainer)

	d := c.scene.Dict()

	normalFont := c.scene.Context().Loader.LoadFont(assets.FontNormal).Face

	titleLabel := eui.NewCenteredLabel(d.Get("menu.main.title")+" -> "+d.Get("menu.main.play"), normalFont)
	rowContainer.AddChild(titleLabel)

	rowContainer.AddChild(eui.NewButton(uiResources, c.scene, d.Get("menu.play.classic"), func() {
		c.scene.Context().ChangeScene(NewLobbyMenuController(c.state, gamedata.ModeClassic))
	}))

	{
		toUnlock := gamedata.ArenaModeCost - c.state.Persistent.PlayerStats.TotalScore
		label := d.Get("menu.play.arena")
		if toUnlock > 0 {
			label = fmt.Sprintf("%s: %d", d.Get("menu.play.to_unlock"), toUnlock)
		}
		b := eui.NewButton(uiResources, c.scene, label, func() {
			c.scene.Context().ChangeScene(NewLobbyMenuController(c.state, gamedata.ModeArena))
		})
		b.GetWidget().Disabled = toUnlock > 0
		rowContainer.AddChild(b)
	}

	rowContainer.AddChild(eui.NewButton(uiResources, c.scene, d.Get("menu.play.tutorial"), func() {
		c.scene.Context().ChangeScene(NewTutorialMenuController(c.state))
	}))

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}))

	rowContainer.AddChild(eui.NewButton(uiResources, c.scene, d.Get("menu.back"), func() {
		c.back()
	}))

	uiObject := eui.NewSceneObject(root)
	c.scene.AddGraphics(uiObject)
	c.scene.AddObject(uiObject)
}

func (c *PlayMenuController) back() {
	c.scene.Context().ChangeScene(NewMainMenuController(c.state))
}
