package webapp

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func ErrAndExit(msg string) {
	app.Window().GetElementByID("header").Set("innerHTML", msg)
	app.Window().GetElementByID("header").Set(
		"className", "w-full bg-rose-700 text-center")
	log.Fatal(msg)
}

func GetDbFile() []byte {
	r, err := http.Get("web/db.json")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func Run() {
	app.Route("/", func() app.Composer {
		return &RootContainer{}
	})
	app.Route("/explore", func() app.Composer { return &ExploreContainer{} })

	app.RunWhenOnBrowser()

	if os.Getenv("BUILD_STATIC") == "true" {
		err := app.GenerateStaticWebsite(".", &app.Handler{
			Name:        "Pomiary",
			Description: "Dane z termometrów Xiaomi",
			//Resources:   app.GitHubPages("pomiary"),
			Styles: []string{
				"/web/default_styles.css",
			},
			//Scripts: []string{
			//	"https://cdn.tailwindcss.com",
			//},
			Image: "/web/temperature.png",
			Icon: app.Icon{
				Maskable: "/web/temperature.png",
				Default:  "/web/temperature.png",
				Large:    "/web/temperature.png",
				SVG:      "/web/temp.svg",
			},
			LoadingLabel: "Ładowanie...",
		})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	log.Println("Listening on :8080")
	http.Handle("/", &app.Handler{
		Name:        "Pomiary",
		Description: "Dane z termometrów Xiaomi",
		Scripts: []string{
			"https://cdn.tailwindcss.com",
		},
		Image: "/web/temperature.png",
		Icon: app.Icon{
			Maskable: "/web/temperature.png",
			Default:  "/web/temperature.png",
			Large:    "/web/temperature.png",
			SVG:      "/web/temp.svg",
		},
		LoadingLabel: "Ładowanie...",
		Styles: []string{
			"/web/default_styles.css",
		},
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
