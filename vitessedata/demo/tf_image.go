package demo

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type ImageTag struct {
	Tags   []string
	Scores []float32
	ErrMsg string
	Err    error
}

type ImageModel struct {
	Program string
	Dir     string
}

func ImageTagging(model ImageModel, fn string) *ImageTag {
	var ret ImageTag
	var outb, errb bytes.Buffer
	cmd := exec.Command("/usr/bin/python", model.Program, "--model_dir="+model.Dir, "--image_file="+fn)
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()

	if err == nil {
		lines := strings.Split(outb.String(), "\n")

		for _, line := range lines {
			fields := strings.Split(line, "(")
			if len(fields) != 2 {
				continue
			}
			var sc float32
			ret.Tags = append(ret.Tags, fields[0])
			fmt.Sscanf(fields[1], "score = %f)", &sc)
			ret.Scores = append(ret.Scores, sc)
		}
	} else {
		ret.ErrMsg = errb.String()
		ret.Err = err
	}

	return &ret
}
