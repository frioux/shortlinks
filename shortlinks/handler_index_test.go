package shortlinks

import (
	"testing"
)

func TestSubstitute(t *testing.T) {
	type test struct {
		path     string
		to       string
		expected string
	}

	cases := []test{
		{path: "/path", to: "https://url.com", expected: "https://url.com"},
		{path: "/path", to: "https://url.com/", expected: "https://url.com/"},
		{path: "/path/", to: "https://url.com", expected: "https://url.com"},
		{path: "/path", to: "https://url.com/q=%s", expected: "https://url.com/q="},
		{path: "/path/sub", to: "https://url.com", expected: "https://url.com"},
		{path: "/path/sub1/sub2", to: "https://url.com", expected: "https://url.com"},
		{path: "/path/sub", to: "https://url.org/q=%s", expected: "https://url.org/q=sub"},
		{path: "/path/sub1/sub2", to: "https://url.org/q=%s&t=%s", expected: "https://url.org/q=sub1/sub2&t=%s"},
		{path: "/path/sub1/sub2", to: "https://url.org/q=%s&", expected: "https://url.org/q=sub1/sub2&"},
		{path: "/path/sub1/", to: "https://url.org/q=%s&t=%s", expected: "https://url.org/q=sub1/&t=%s"},
		{path: "/path/sub1", to: "https://url.org/q=%s&t=%s", expected: "https://url.org/q=sub1&t=%s"},
		{path: "/j", to: "https://atlassian.net/browse/%s", expected: "https://atlassian.net/browse/"},
		{path: "/j/JIRA-000", to: "https://atlassian.net/browse/%s", expected: "https://atlassian.net/browse/JIRA-000"},
		{path: "/some/firstquery/secondquery", to: "https://some.url/%s/else", expected: "https://some.url/firstquery/secondquery/else"},
		{path: "/other/firstquery/secondquery", to: "https://some.url/%s/else/%s", expected: "https://some.url/firstquery/secondquery/else/%s"},
	}

	for _, test := range cases {
		s := Shortlink{To: test.to}
		_, substitutions := split(test.path)
		url := substitute(s, substitutions)
		if url != test.expected {
			t.Errorf("URL as a result of substitute did match expected,\npath: %q\nto: %q\nactual url:\t\t%q\nexpected url:\t%q", test.path, test.to, url, test.expected)
		}
	}
}
