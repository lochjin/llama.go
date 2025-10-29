package pull

import (
	"fmt"
	"github.com/Qitmeer/llama.go/api"
	"github.com/Qitmeer/llama.go/common/progress"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func PullHandler(ctx *cli.Context) error {
	var model string
	if ctx.Args().Len() > 0 {
		model = ctx.Args().Slice()[0]
	}

	client := api.DefaultClient()

	p := progress.NewProgress(os.Stderr)
	defer p.Stop()

	bars := make(map[string]*progress.Bar)

	var status string
	var spinner *progress.Spinner

	fn := func(resp api.ProgressResponse) error {
		if resp.Digest != "" {
			if resp.Completed == 0 {
				// This is the initial status update for the
				// layer, which the server sends before
				// beginning the download, for clients to
				// compute total size and prepare for
				// downloads, if needed.
				//
				// Skipping this here to avoid showing a 0%
				// progress bar, which *should* clue the user
				// into the fact that many things are being
				// downloaded and that the current active
				// download is not that last. However, in rare
				// cases it seems to be triggering to some, and
				// it isn't worth explaining, so just ignore
				// and regress to the old UI that keeps giving
				// you the "But wait, there is more!" after
				// each "100% done" bar, which is "better."
				return nil
			}

			if spinner != nil {
				spinner.Stop()
			}

			bar, ok := bars[resp.Digest]
			if !ok {
				name, isDigest := strings.CutPrefix(resp.Digest, "sha256:")
				name = strings.TrimSpace(name)
				if isDigest {
					name = name[:min(12, len(name))]
				}
				bar = progress.NewBar(fmt.Sprintf("pulling %s:", name), resp.Total, resp.Completed)
				bars[resp.Digest] = bar
				p.Add(resp.Digest, bar)
			}

			bar.Set(resp.Completed)
		} else if status != resp.Status {
			if spinner != nil {
				spinner.Stop()
			}

			status = resp.Status
			spinner = progress.NewSpinner(status)
			p.Add(status, spinner)
		}

		return nil
	}

	request := api.PullRequest{Model: model}
	return client.Pull(ctx.Context, &request, fn)
}
