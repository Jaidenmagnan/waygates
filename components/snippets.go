package components

import "fmt"

func Snippet(baseURL string, waygateID int, waygateLinkID int, ringName string) string {
	return fmt.Sprintf(`<nav class="waygate" aria-label="%s web ring">
  <a href="%s/rings/%d/%d/previous">previous</a>
  <a href="%s">waygate: %s</a>
  <a href="%s/rings/%d/%d/next">next</a>
</nav>`, ringName, baseURL, waygateID, waygateLinkID, baseURL, ringName, baseURL, waygateID, waygateLinkID)
}

func InviteURL(baseURL string, waygateID int, waygateLinkID int) string {
	return fmt.Sprintf("%s/rings/%d/%d/next", baseURL, waygateID, waygateLinkID)
}
