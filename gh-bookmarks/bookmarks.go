package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/google/go-github/v38/github"
)

func createBookmarksHTML(repos []*github.Repository) {
	// Create an HTML bookmark file from the list of repos
	htmlContent := generateHTMLBookmarks(repos)

	// Write the HTML content to a file
	err := writeToFile("bookmarks.html", htmlContent)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
	}
}

func generateHTMLBookmarks(repos []*github.Repository) string {
	// Group repositories by owner
	groupedRepos := make(map[string][]*github.Repository)

	// map the list of repo URLs to the owner URL
	// doing it here is simpler than using text/template
	for _, repo := range repos {
		groupedRepos[strings.TrimPrefix(*repo.Owner.HTMLURL, "https://")] = append(groupedRepos[strings.TrimPrefix(*repo.Owner.HTMLURL, "https://")], repo)
	}

	// Create an HTML template
	tmpl := `
<!DOCTYPE NETSCAPE-Bookmark-file-1>
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
{{- range $folder, $repos := .}}
    <DT><H3>{{$folder}}</H3>
    <DL><p>
	{{- range $repo := $repos}}
        <DT><A HREF="{{$repo.HTMLURL}}" ADD_DATE="0" LAST_VISIT="0" LAST_MODIFIED="0">{{$repo.FullName}}</A>
    {{- end}}
    </DL><p>
{{- end}}
</DL><p>
`

	// Execute the template to generate the HTML content
	var htmlContent strings.Builder
	t := template.Must(template.New("bookmark").Parse(tmpl))
	err := t.Execute(&htmlContent, groupedRepos)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
	}

	return htmlContent.String()
}

func writeToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
