package list

import (
	"fmt"
	"strings"
	"sync"

	"github.com/christian-gama/produgit/internal/git"
	"github.com/spf13/cobra"
)

var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "List all repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		var reposMu sync.Mutex
		repos := make([]string, 0)

		err := git.WalkDirs(dir, func(path string) error {
			path = strings.TrimSuffix(path, "/.git")
			reposMu.Lock()
			repos = append(repos, path)
			reposMu.Unlock()

			return nil
		})
		if err != nil {
			return err
		}

		for _, repo := range SortAndDeDuplicate(repos) {
			fmt.Println(repo)
		}

		return nil
	},
}
