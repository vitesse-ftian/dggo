package demo

import (
	"fmt"
	"testing"
)

func TestImageTaggings(t *testing.T) {
	m := ImageModel{
		Program: "/home/ftian/oss/models/tutorials/image/imagenet/classify_image.py",
		Dir:     "/home/ftian/offline/tfm",
	}

	t.Run("cat", func(t *testing.T) {
		tags := ImageTagging(m, "/home/ftian/offline/image/cat-1.jpeg")
		if tags.Err != nil {
			t.Error(tags.Err)
		} else {
			for i := 0; i < len(tags.Tags); i++ {
				fmt.Printf("CAT-%d: %s, %f\n", i, tags.Tags[i], tags.Scores[i])
			}
		}
	})
}
