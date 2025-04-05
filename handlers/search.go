package handlers

import (
	"fmt"
	"go-games-htmx/api"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SearchDisplayItem struct {
	Type             string
	Id               int
	Link             string
	Title            string
	ShortDescription string
	ImageSmall       string
	ReleaseDate      string
	Platforms        []string
}

type ResultsPagination struct {
	PageCurrent  int
	PageNextLink string
	PageTotal    int
	TotalFound   int
}

type Results struct {
	Query        string
	ErrorMessage string
	Pagination   ResultsPagination
	Items        []SearchDisplayItem
}

func newEmptyResults(q string, errorMessage string) Results {
	return Results{
		Query:        q,
		ErrorMessage: errorMessage,
		Pagination:   ResultsPagination{},
		Items:        []SearchDisplayItem{},
	}
}

func HandleGETSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")

	q := strings.TrimSpace(r.URL.Query().Get("q"))
	page := strings.TrimSpace(r.URL.Query().Get("page"))
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if pageInt <= 0 {
		page = "1"
		pageInt = 1
	}

	var tmpl *template.Template
	if pageInt == 1 {
		tmpl = template.Must(template.New("template").Funcs(CustomTemplateFuncs).ParseFiles("components/blocks/results_items.html", "components/blocks/results.html"))
	} else {
		tmpl = template.Must(template.New("template").Funcs(CustomTemplateFuncs).ParseFiles("components/blocks/results_items.html"))
	}

	var data Results

	if len(q) < 3 {
		data = newEmptyResults(q, "Please search for more than 3 characters")

	} else {
		results, err := api.Search(q, page)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		if len(results.Results) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var items []SearchDisplayItem
		for _, r := range results.Results {

			// Skip if no description exists
			if len(strings.TrimSpace(r.ShortDescription)) == 0 {
				continue
			}

			searchDisplayItem := SearchDisplayItem{
				Type:             r.ResourceType,
				Id:               r.Id,
				Link:             fmt.Sprintf("/game?guid=%s", r.Guid),
				Title:            r.Name,
				ShortDescription: r.ShortDescription,
			}

			// If it contains the "default" API vendor image, consider it as empty
			if !strings.Contains(r.Images.Small, "-gb_default-") {
				searchDisplayItem.ImageSmall = r.Images.Small
			}

			// Format Release Date
			if r.ReleaseDate != "" {
				tmpTime, err := time.Parse("2006-01-02", r.ReleaseDate)
				if err == nil {
					searchDisplayItem.ReleaseDate = tmpTime.Format("02 January 2006")
				}
			}

			for _, rPlatform := range r.Platforms {
				searchDisplayItem.Platforms = append(searchDisplayItem.Platforms, rPlatform.Name)
			}
			sort.Strings(searchDisplayItem.Platforms)

			items = append(items, searchDisplayItem)
		}

		pageNextLink := ""

		pageNextInt := pageInt + 1
		if pageNextInt > results.Pages {
			pageNextInt = 0

		} else {
			pageNextLink = "/search?q=" + q + "&page=" + strconv.Itoa(pageNextInt)
		}

		// Build the data object
		data = Results{
			Query:        q,
			ErrorMessage: "",
			Pagination: ResultsPagination{
				PageCurrent:  pageInt,
				PageNextLink: pageNextLink,
				PageTotal:    results.Pages,
				TotalFound:   results.Total,
			},
			Items: items,
		}
	}

	var template string
	if pageInt == 1 {
		template = "results"
	} else {
		template = "results_items"
	}

	tmpl.ExecuteTemplate(w, template, data)
}
