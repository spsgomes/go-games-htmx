package handlers

import (
	"fmt"
	"go-games-htmx/api"
	"html/template"
	"net/http"
	"sort"
	"strings"
	"time"
)

type GameDisplayItem struct {
	Type             string
	Id               int
	Title            string
	ShortDescription string
	Description      template.HTML
	ImageLarge       string
	ReleaseDate      string
	GameRating       []string
	Platforms        []string
	Developers       []string
	Genres           []string
	Publishers       []string
	SimilarGames     []string
}

func HandleGETGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")

	tmpl := template.Must(template.New("template").Funcs(CustomTemplateFuncs).ParseFiles("components/pages/game.html", "components/base.html"))

	guid := r.URL.Query().Get("guid")

	if guid == "" {
		Handle404(w, r)
		return
	}

	var display GameDisplayItem

	results, err := api.Game(guid)
	if err != nil {
		fmt.Println(err)

		Handle404(w, r)
		return
	}

	if results.Error != "OK" {
		fmt.Println("results.Error: " + results.Error)

		Handle404(w, r)
		return
	}

	result := results.Results

	display = GameDisplayItem{
		Type:             "game",
		Id:               result.Id,
		Title:            result.Name,
		ShortDescription: result.ShortDescription,
	}

	description := strings.TrimSpace(result.Description)

	const GB_WEBSITE = "https://www.giantbomb.com"

	description = strings.ReplaceAll(description, "href=\"", "target=\"_blank\" rel=\"noopener noreferrer\" href=\""+GB_WEBSITE)
	description = strings.ReplaceAll(description, GB_WEBSITE+"https", "https")
	description = strings.ReplaceAll(description, "data-src", "src")

	display.Description = template.HTML(description)

	// If it contains the "default" API vendor image, consider it as empty
	if !strings.Contains(result.Images.Large, "-gb_default-") {
		display.ImageLarge = result.Images.Large
	}

	// Format Release Date
	if result.ReleaseDate != "" {
		tmpTime, err := time.Parse("2006-01-02", result.ReleaseDate)
		if err == nil {
			display.ReleaseDate = tmpTime.Format("02 January 2006")
		}
	}

	// Build Game Ratings
	for _, val := range result.GameRating {
		display.GameRating = append(display.GameRating, val.Name)
	}

	// Build Platforms
	for _, val := range result.Platforms {
		display.Platforms = append(display.Platforms, val.Name)
	}
	sort.Strings(display.Platforms)

	// Build Developers
	for _, val := range result.Developers {
		display.Developers = append(display.Developers, val.Name)
	}

	// Build Genres
	for _, val := range result.Genres {
		display.Genres = append(display.Genres, val.Name)
	}

	// Build Publishers
	for _, val := range result.Publishers {
		display.Publishers = append(display.Publishers, val.Name)
	}

	// Build SimilarGames
	for _, val := range result.SimilarGames {
		display.SimilarGames = append(display.SimilarGames, val.Name)
	}

	tmpl.ExecuteTemplate(w, "base", display)
}
