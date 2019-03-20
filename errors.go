package elevengo

import "fmt"

func apiError(code int) error {
	return fmt.Errorf("api error: %d", code)
}
