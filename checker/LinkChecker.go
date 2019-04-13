package checker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// A Result is the result host & body
// after doing a http.get to the {host}
type Result struct {
	host string
	body string
}

// A Response is the response link & status
// after doing a http.Get to {link}
type Response struct {
	status int
	link   string
}

// A HTMLGetter is a function for getting html
// for a given url
type HTMLGetter func(url string) Result

// Check attemps to check all of the links for a given {host}
func check(host string) {
	// Remove trailing slash
	if host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}

	// Create *url.URL from host
	hostURL, err := url.Parse(host)

	// Get *.html.Node from the host
	node, err := getHTMLNode(getHTMLBody, hostURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get links from the html node
	links := getLinks(node)

	// Check the status for each link
	checkedLinks := checkLinks(links, hostURL)

	// Print out link status
	for _, res := range checkedLinks {
		fmt.Printf("[%d] - %s\n", res.status, res.link)
	}
}

// GetLinks gets the link nodes for a given *html.Node
func getLinks(node *html.Node) []string {
	links := make([]string, 0)
	getAnchorLinks(node, &links)
	return links
}

// CheckLinks takes a slice of strings (links), and attempts to do a http.Get on each
// Because the href attribute of an <a> tag can have many values
// there is a lot of logic in here to determine what we should do
// We currently handle:
// - absolute links
// - absolute links with no scheme
// - no links (empty string)
// - relative links
// - relative links with no leading slash
// - id link (#)
// - email link (mailto:)
// - javascript link (javascript:)
func checkLinks(links []string, host *url.URL) []Response {
	checkedLinks := make([]Response, 0)
	for _, a := range links {
		if len(a) > 1 {
			// URL encode spaces
			a = strings.Replace(a, " ", "%20", -1)

			// Check if the link is absolute
			_, err := url.ParseRequestURI(a)
			if err != nil {
				// Omit if link is an id anchor
				if string(a[0]) == "#" {
					continue
				}
				// Assume the link is relative (no slash)
				a = fmt.Sprintf("%s/%s", host, a)
			}

			// If the link starts with a slash
			if string(a[0]) == "/" {
				// If the link starts with // assume absolute (no scheme)
				// Else its a relative link
				if string(a[1]) == "/" {
					a = host.Scheme + ":" + a
				} else {
					a = host.Host + a
				}
			}

			resp, err := http.Get(a)
			if err != nil {
				continue
			}

			checkedLinks = append(checkedLinks, Response{resp.StatusCode, a})
		}
	}
	return checkedLinks
}

// GetAnchorLinks gets the <a> tags from a html node
// This is a recursive function
func getAnchorLinks(n *html.Node, links *[]string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				*links = append(*links, a.Val)
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getAnchorLinks(c, links)
	}
}

// GetHTMLNode gets the *.html.Node for a given host
// It will append scheme as http if not given
func getHTMLNode(getHTMLBody HTMLGetter, host *url.URL) (*html.Node, error) {
	// If no scheme, default to http
	if host.Scheme == "" {
		host.Scheme = "http"
	}

	result := getHTMLBody(host.String())
	doc, err := html.Parse(strings.NewReader(result.body))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// GetHTMLBody gets the html body and hostname for a given url
func getHTMLBody(url string) Result {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	host := resp.Request.URL.Scheme + "://" + resp.Request.URL.Host
	body, err := ioutil.ReadAll(resp.Body)
	return Result{host, string(body)}
}
