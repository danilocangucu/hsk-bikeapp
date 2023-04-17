package tests

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/PuerkitoBio/goquery"

	_ "github.com/mattn/go-sqlite3"
)

func TestIndexGet(t *testing.T) {
	// Get the absolute path of the index.html file
	templatePath, err := filepath.Abs("../index.html")
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	templ := template.Must(template.ParseFiles(templatePath))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := templ.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("http.Get: %v", err)
	}
	defer resp.Body.Close()

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("resp.StatusCode = %d; want %d", got, want)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Fatalf("failed to create document from response body: %v", err)
	}

	title := doc.Find("title").Text()
	if got, want := title, "Document"; got != want {
		t.Errorf("title = %q; want %q", got, want)
	}

	cssFile := "src/styles.css"
	if doc.Find("link[href='"+cssFile+"']").Length() <= 0 {
		t.Errorf("%s not found in response body", cssFile)
	}

	jsFile := "src/App.js"
	if doc.Find("script[src='"+jsFile+"']").Length() <= 0 {
		t.Errorf("%s not found in response body", jsFile)
	}

	if doc.Find("#stations-list").Length() <= 0 {
		t.Errorf("#stations-list not found in response body")
	}

	fontFile := "https://fonts.googleapis.com/css2?family=Anton&family=Roboto:wght@300&display=swap"
	if doc.Find("link[href='"+fontFile+"']").Length() <= 0 {
		t.Errorf("%s not found in response body", fontFile)
	}

	videoFile := "src/video.mp4"
	if doc.Find("source[src='"+videoFile+"']").Length() <= 0 {
		t.Errorf("%s not found in response body", videoFile)
	}

	sectionTitle := "STATIONS"
	section := doc.Find("div.section:has(h1:contains('" + sectionTitle + "'))")
	if section.Length() <= 0 {
		t.Errorf("section with title %s not found in response body", sectionTitle)
	}

	journeysTitle := doc.Find("div.journeys.section h1")
	if journeysTitle.Length() <= 0 {
		t.Errorf("journeys section title not found in response body")
	} else if got, want := journeysTitle.Text(), "JOURNEYS"; got != want {
		t.Errorf("journeys section title = %q; want %q", got, want)
	}

}
