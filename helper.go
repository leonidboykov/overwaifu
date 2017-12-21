package overwaifu

import "github.com/leonidboykov/getmoe"

func hasTags(post getmoe.Post, tags []string) bool {
	for i := range tags {
		if post.HasTag(tags[i]) {
			return true
		}
	}
	return false
}
