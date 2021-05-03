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
	links := make([]string, 0)
	getAnchorLinks(node, &links)

	// Check the status for each link
	parsedLinks := parseLinks(links, hostURL)
	ch := make(chan Response)

	for _, link := range parsedLinks {
		go checkLink(link, ch)
	}

	results := make(map[int]int)
	for range parsedLinks {
		res := <-ch
		if _, ok := results[res.status]; ok {
			results[res.status]++
		} else {
			results[res.status] = 1
		}

		fmt.Printf("[%d] - %s\n", res.status, res.link)
	}

	fmt.Println(results)
}

func checkLink(link string, ch chan<- Response) {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
	}
	ch <- Response{resp.StatusCode, link}
}

// ParseLinks takes a slice of strings (anchor tag hrefs), and converts them into fully qualified links
// Because the href attribute of an <a> tag can have many values
// there is a lot of logic in here to determine what we should do
// Handles:
// - absolute links
// - absolute links with no scheme
// - no links (empty string)
// - relative links
// - relative links with no leading slash
// - id link (#)
// - email link (mailto:)
// - javascript link (javascript:)
func parseLinks(links []string, host *url.URL) []string {
	parsedLinks := make([]string, 0)

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
					a = host.Scheme + "://" + host.Host + a
				}
			}

			if strings.HasPrefix(a, "mailto:") || strings.HasPrefix(a, "javascript:") {
				continue
			}

			parsedLinks = append(parsedLinks, a)
		}
	}
	return parsedLinks
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
