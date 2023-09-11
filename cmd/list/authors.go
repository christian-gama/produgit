package list

import (
	"fmt"
	"sort"
	"sync"

	"github.com/christian-gama/produgit/internal/git"
	"github.com/spf13/cobra"
)

var authorCmd = &cobra.Command{
	Use:   "author",
	Short: "List authors of all repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		var authorsMu sync.Mutex
		authors := make([]string, 0)

		err := git.WalkDirs(dir, func(path string) error {
			a, err := git.ListAllAuthors(path)
			if err != nil {
				return err
			}

			authorsMu.Lock()
			authors = append(authors, a...)
			authorsMu.Unlock()

			return nil
		})
		if err != nil {
			return err
		}

		for _, author := range SortAndDeDuplicate(authors) {
			fmt.Println(author)
		}

		return nil
	},
}

func SortAndDeDuplicate(authors []string) []string {
	unique := make(map[string]bool)
	for _, author := range authors {
		unique[author] = true
	}

	var sorted []string
	for author := range unique {
		sorted = append(sorted, author)
	}

	sort.Strings(sorted)

	return sorted
}
