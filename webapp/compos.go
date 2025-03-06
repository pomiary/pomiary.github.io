package webapp

import (
	"fmt"
	"log"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type InstallButton struct {
	app.Compo
	name             string
	isAppInstallable bool
}

func (b *InstallButton) OnMount(ctx app.Context) {
	b.isAppInstallable = ctx.IsAppInstallable()
}
func (b *InstallButton) OnAppInstallChange(ctx app.Context) {
	b.isAppInstallable = ctx.IsAppInstallable()
}
func (b *InstallButton) Render() app.UI {
	return app.Div().
		Body(
			app.If(b.isAppInstallable, func() app.UI {
				return app.Button().
					Text("Zainstaluj jako aplikację").
					OnClick(b.onInstallButtonClicked).
					Class("bg-green-600 hover:bg-green-800 p-2 rounded m-2")
			}),
		).Class("flex flex-col")
}
func (b *InstallButton) onInstallButtonClicked(ctx app.Context, e app.Event) {
	ctx.ShowAppInstallPrompt()
}

type LoadingWidget struct {
	app.Compo
	id      string
	visible bool
}

func (l *LoadingWidget) OnMount(ctx app.Context) {
	ctx.Handle(fmt.Sprintf("show-%s", l.id), l.Visibility)
}

func (l *LoadingWidget) Render() app.UI {
	return app.Div().Body(
		app.Div().Body(
			app.Span().
				Text("Loading...").Class("font-bold").
				Class("!absolute !-m-px !h-px !w-px !overflow-hidden !whitespace-nowrap !border-0 !p-0 ![clip:rect(0,0,0,0)]"),
		).
			Class("inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-e-transparent align-[-0.125em] text-surface motion-reduce:animate-[spin_1.5s_linear_infinite] text-sky-600"),
	).Class("fixed invisible").ID(l.id)
}

func (l *LoadingWidget) Visibility(ctx app.Context, a app.Action) {
	log.Println("handling ", a.Name)
	if !l.visible {
		app.Window().GetElementByID(l.id).Set("className", "flex flex-row justify-center items-center")
		l.visible = true
	} else {
		app.Window().GetElementByID(l.id).Set("className", "flex flex-row justify-center fixed invisible")
		l.visible = false
	}
}

type Header struct {
	app.Compo
}

func (c *Header) Render() app.UI {
	return app.Div().Body(
		app.P().Text("Pomiary").Class("text-xl").ID("header-text"),
	).Class("w-full bg-sky-900 text-center").ID("header")
}

type ExploreButton struct {
	app.Compo
}

func (b *ExploreButton) Render() app.UI {
	return app.Div().Body(
		app.Button().Text("Przeglądaj pomiary").OnClick(func(ctx app.Context, e app.Event) {
			app.Window().Set("location", "explore")
		}).Class("bg-sky-700 hover:bg-sky-800 font-bold py-2 px-4 my-2 rounded"),
	)
}

type HomeButton struct {
	app.Compo
}

func (b *HomeButton) Render() app.UI {
	return app.Div().Body(
		app.Button().Text("Powrót").OnClick(func(ctx app.Context, e app.Event) {
			app.Window().Set("location", ".")
		}).Class("bg-sky-700 hover:bg-sky-800 font-bold py-2 px-4 rounded"),
	)
}

type ThermometerContainer struct {
	app.Compo
	id          string
	roomName    string
	temperature float64
	humidity    int
	voltage     float64
	timestamp   int
}

func (c *ThermometerContainer) OnMount(ctx app.Context) {
	ctx.Async(func() {
		ctx.NewAction(fmt.Sprintf("show-%s-loading", c.id))
		m, err := LastData(c.id)
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				ErrAndExit(err.Error())
			}
			c.temperature = m.Temperature
			c.humidity = m.Humidity
			c.voltage = m.Voltage
			c.timestamp = m.Timestamp
			ctx.NewAction(fmt.Sprintf("show-%s-loading", c.id))
		})
	})
}
func (c *ThermometerContainer) Render() app.UI {
	return app.Div().Body(
		app.P().Text(c.roomName).Class("text-lg font-bold mb-4"),
		&LoadingWidget{id: fmt.Sprintf("%s-loading", c.id)},
		app.Table().Body(
			app.Tr().Body(
				app.Td().Text("Temperatura [°C]").Class("px-4 py-2 border"),
				app.Td().Text(fmt.Sprintf("%v", c.temperature)).Class("px-4 py-2 border"),
			),
			app.Tr().Body(
				app.Td().Text("Wilgotność [%]").Class("px-4 py-2 border"),
				app.Td().Text(fmt.Sprintf("%v", c.humidity)).Class("px-4 py-2 border"),
			),
			app.Tr().Body(
				app.Td().Text("Napięcie baterii [V]").Class("px-4 py-2 border"),
				app.Td().Text(fmt.Sprintf("%v", c.voltage)).Class("px-4 py-2 border"),
			),
			app.Tr().Body(
				app.Td().Text("Czas").Class("px-4 py-2 border"),
				app.Td().Text(time.Unix(int64(c.timestamp), 0).Format("2006-01-02 15:04:05")).Class(
					"px-4 py-2 border"),
			),
		).Class("w-full border-collapse border border-sky-800"),
	).Class("border border-sky-800 p-4 m-2 rounded").ID(c.id)
}

type ExploreTable struct {
	app.Compo
	measurements []Measurement
}

func (c *ExploreTable) OnMount(ctx app.Context) {
	ctx.Handle("loadMore", c.handleLoadMore)
	ctx.Async(func() {
		ctx.NewAction("show-measurements-loading")
		m, err := Data(0)
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				ErrAndExit(err.Error())
			}
			c.measurements = m
			ctx.NewAction("show-measurements-loading")
		})
	})
}

func (c *ExploreTable) Render() app.UI {
	return app.Div().Body(
		app.Table().Body(
			app.Tr().Body(
				app.Td().Text("Czas").Class("px-4 py-2 border"),
				app.Td().Text("Pokój").Class("px-4 py-2 border"),
				app.Td().Text("Temp.").Class("px-4 py-2 border"),
				app.Td().Text("Wilg.").Class("px-4 py-2 border"),
			),
			app.Range(c.measurements).Slice(func(index int) app.UI {
				return app.Tr().Body(
					app.Td().Text(time.Unix(int64(c.measurements[index].Timestamp), 0).Format("2006-01-02 15:04:05")).Class(
						"px-4 py-2 border"),
					app.Td().Text(Sensors[c.measurements[index].Id]).Class("px-4 py-2 border"),
					app.Td().Text(fmt.Sprintf("%v", c.measurements[index].Temperature)).Class("px-4 py-2 border"),
					app.Td().Text(fmt.Sprintf("%v", c.measurements[index].Humidity)).Class("px-4 py-2 border"),
				)
			}),
		).Class("border-collapse border border-sky-800"),
	)
}

func (c *ExploreTable) handleLoadMore(ctx app.Context, a app.Action) {
	log.Println("handling ", a.Name)
	var skip int
	ctx.GetState("skip", &skip)
	skip += 100
	ctx.SetState("skip", skip)
	app.Window().GetElementByID("header-text").Set("innerHTML", "Ładowanie...")
	ctx.Async(func() {
		m, err := Data(skip)
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				ErrAndExit(err.Error())
			}
			log.Println("Adding measurements...")
			c.measurements = append(c.measurements, m...)
			app.Window().GetElementByID("header-text").Set("innerHTML", "Pomiary")
			ctx.NewAction("show-measurements-loading")
		})
	})
}

type LoadMoreButton struct {
	app.Compo
}

func (c *LoadMoreButton) Render() app.UI {
	return app.Div().Body(
		app.Button().Text("Wczytaj więcej").OnClick(c.OnClick),
	).Class("bg-sky-700 hover:bg-sky-800 font-bold py-2 px-4 rounded")
}
func (c *LoadMoreButton) OnClick(ctx app.Context, e app.Event) {
	log.Println("Sending loadMore action...")
	ctx.NewAction("loadMore")
	log.Println("Sending show-measurements-loading action...")
	ctx.NewAction("show-measurements-loading")
}

type RootContainer struct {
	app.Compo
}

func (c *RootContainer) Render() app.UI {
	return app.Div().Body(
		&Header{},
		app.Br(),
		app.Div().Body(
			app.Range(Sensors).Map(func(k string) app.UI {
				return &ThermometerContainer{id: k, roomName: Sensors[k]}
			}),
			app.Br(),
			&ExploreButton{},
			&InstallButton{},
		).Class(
			"flex flex-col items-center"),
	)
}

type ExploreContainer struct {
	app.Compo
}

func (c *ExploreContainer) Render() app.UI {
	return app.Div().Body(
		&Header{},
		app.Div().Body(
			&ExploreTable{},
			app.Br(),
			&LoadMoreButton{},
			&LoadingWidget{id: "measurements-loading"},
			app.Br(),
			&HomeButton{},
		).Class("flex flex-col items-center p-4"),
	)
}
