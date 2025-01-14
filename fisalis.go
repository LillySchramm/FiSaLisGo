package fisalisgo

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	FISALIS_URL  = "https://www.finanz-sanktionsliste.de/fisalis/"
	SEARCH_FIELD = "txtSearch"
	CMD_FIELD    = "cmdSearch"
)

type Document struct {
	Id   string
	Link string
	Date time.Time

	Names         []string
	BirthDates    []time.Time
	Orgs          []string
	DocumentNames []string
}
type Result struct {
	Id string

	Match       float64
	Description string
	Documents   []Document
}

func Search(ctx context.Context, text string) ([]Result, error) {
	str, err := request(ctx, text)
	if err != nil {
		return nil, err
	}

	return parse(str)
}

func request(ctx context.Context, text string) (*string, error) {
	form := url.Values{}
	form.Add(SEARCH_FIELD, text)
	form.Add(CMD_FIELD, "🔍")

	req, err := http.NewRequestWithContext(ctx, "POST", FISALIS_URL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "FiSaLisGo (https://github.com/LillySchramm/FiSaLisGo)")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	str := string(body)
	return &str, nil
}

func parse(str *string) ([]Result, error) {
	// Yes, I know this is not the best way to parse HTML, but the absolute garbage of HTML
	// that is returned by the website makes it impossible to use a proper HTML parser.
	// But hey, at least it did not change since at least a century, so it should be fine forever.

	body := *str
	body = strings.Split(body, "<form action=\"?\" method=\"post\">")[1]
	body = strings.Split(body, "</form>")[0]
	body = strings.Split(body, "Ergebnisse anzeigen -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -->")[1]
	body = strings.ReplaceAll(body, "<em>", "")
	body = strings.ReplaceAll(body, "</em>", "")

	parts := strings.Split(body, "<hr />")
	parts = parts[1:]

	results := make([]Result, len(parts))
	for i, part := range parts {
		results[i] = Result{}

		part = strings.ReplaceAll(part, "<h3><span style='color:", "")
		part = strings.SplitN(part, ">", 2)[1]

		matchStr := strings.SplitN(part, "%", 2)[0]
		match, err := strconv.ParseFloat(matchStr, 64)
		if err != nil {
			return nil, err
		}
		results[i].Match = match

		part := strings.Split(part, "</span>: (")[1]
		id, part := strings.SplitN(part, ")", 2)[0], strings.SplitN(part, ")", 2)[1]
		part = strings.Trim(part, " ")

		results[i].Id = id

		part = strings.SplitN(part, "</h3>", 2)[1]

		if !strings.HasPrefix(part, "<p><a") {
			results[i].Description, part = strings.SplitN(part, "</p>", 2)[0], strings.SplitN(part, "</p>", 2)[1]
			results[i].Description = stripHtmlTags(results[i].Description)
		}

		documentParts := strings.Split(part, "<p><a href='")
		results[i].Documents = make([]Document, len(documentParts)-1)
		for j, documentStr := range documentParts[1:] {
			document := Document{}
			document.BirthDates = []time.Time{}
			document.Names = []string{}
			document.Orgs = []string{}

			link, documentStr := strings.SplitN(documentStr, "'", 2)[0], strings.SplitN(documentStr, "'", 2)[1]
			document.Link = link

			id, documentStr := strings.SplitN(documentStr, "</a><small>/", 2)[0], strings.SplitN(documentStr, "</a><small>/", 2)[1]
			id = strings.Split(id, ">")[1]
			document.Id = id

			dateStr, documentStr := strings.SplitN(documentStr, ":", 2)[0], strings.SplitN(documentStr, ":", 2)[1]
			document.Date, err = time.Parse("02.01.2006", dateStr)
			if err != nil {
				return nil, err
			}

			documentStr = stripHtmlTags(documentStr)
			documentStr = strings.ReplaceAll(documentStr, "\n", "")
			documentStr = strings.Trim(documentStr, " ")

			infoParts := strings.Split(documentStr, ",")
			for _, infoPart := range infoParts {
				infoPart = strings.Trim(infoPart, " ")
				infoPart = strings.TrimSuffix(infoPart, "</p>")
				if strings.HasPrefix(infoPart, "Name:") {
					infoPart = strings.TrimPrefix(infoPart, "Name: ")
					document.Names = append(document.Names, infoPart)
				} else if strings.HasPrefix(infoPart, "Geboren:") {
					infoPart = strings.TrimPrefix(infoPart, "Geboren: ")
					infoPart = strings.Trim(infoPart, ".")
					if !strings.Contains(infoPart, ".") {
						infoPart = "01.01." + infoPart
					}
					birthDate, err := time.Parse("02.01.2006", infoPart)
					if err != nil {
						return nil, err
					}
					document.BirthDates = append(document.BirthDates, birthDate)
				} else if strings.HasPrefix(infoPart, "Orga:") {
					infoPart = strings.TrimPrefix(infoPart, "Orga: ")
					document.Orgs = append(document.Orgs, infoPart)
				} else if strings.HasPrefix(infoPart, "Dokument:") {
					infoPart = strings.TrimPrefix(infoPart, "Dokument: ")
					document.DocumentNames = append(document.DocumentNames, infoPart)
				}
			}

			results[i].Documents[j] = document
		}
	}

	return results, nil
}
