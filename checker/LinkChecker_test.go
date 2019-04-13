package checker

import (
	"fmt"
	"net/url"
	"testing"
)

var HOST_URL, _ = url.Parse("www.example.com")

var linkMap = map[string]string{
	"ABSOLUTE_LINK":          "http://example.com",
	"ABSOLUTE_LINK_NOSCHEME": "//example.com",
	"NO_LINK":                "",
	"RELATIVE_LINK":          "/test",
	"RELATIVE_LINK_NOSLASH":  "test2",
	"ID_LINK":                "#top",
	"EMAIL_LINK":             "mailto:example@example.com",
	"JAVASCRIPT_LINK":        "javascript:alert('Hello');",
}

func TestCheckLinks(t *testing.T) {
	links := make([]string, 0)
	for _, link := range linkMap {
		links = append(links, link)
	}
	checkLinks(links, HOST_URL)
}

func TestGetLinks(t *testing.T) {
	node, err := getHTMLNode(mock_getHTMLBody, HOST_URL)
	if err != nil {
		t.Errorf("Failed to get *html.Node, got error %v", err)
	}

	links := getLinks(node)

	for _, link := range linkMap {
		if !contains(links, link) {
			t.Errorf("Expected to find link %v", link)
		}
	}
}

func mock_getHTMLBody(url string) Result {
	return Result{url,
		fmt.Sprintf(
			`<!doctype html>
				<html>
					<body>
						<a href=%s>Absolute Link</a>
						<a href=%s>Absolute Link (no scheme)</a>
						<a href=%s>No Link</a>
						<a href=%s>Relative Link</a>
						<a href=%s>Relative Link (no slash)</a>
						<a href=%s>Id Link</a>
						<a href=%s>Email Link</a>
						<a href=%s>Javascript Link</a>
					</body>
				</html>
			`,
			linkMap["ABSOLUTE_LINK"],
			linkMap["ABSOLUTE_LINK_NOSCHEME"],
			linkMap["NO_LINK"],
			linkMap["RELATIVE_LINK"],
			linkMap["RELATIVE_LINK_NOSLASH"],
			linkMap["ID_LINK"],
			linkMap["EMAIL_LINK"],
			linkMap["JAVASCRIPT_LINK"],
		),
	}
}

func contains(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
