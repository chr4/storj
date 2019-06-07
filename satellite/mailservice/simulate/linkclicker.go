// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information

package simulate

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/zeebo/errs"
	monkit "gopkg.in/spacemonkeygo/monkit.v2"

	"storj.io/storj/internal/post"
)

var mon = monkit.Package()

// LinkClicker is mailservice.Sender that click all links
// from html msg parts
type LinkClicker struct{}

// FromAddress return empty mail address
func (clicker *LinkClicker) FromAddress() post.Address {
	return post.Address{}
}

// SendEmail click all links from email html parts
func (clicker *LinkClicker) SendEmail(ctx context.Context, msg *post.Message) (err error) {
	defer mon.Task()(&ctx)(&err)

	// dirty way to find links without pulling in a html dependency
	regx := regexp.MustCompile(`href="([^\s])+"`)
	// collect all links
	var links []string
	for _, part := range msg.Parts {
		tags := findLinkTags(part.Content)
		for _, tag := range tags {
			href := regx.FindString(tag)
			if href == "" {
				continue
			}

			links = append(links, href[len(`href="`):len(href)-1])
		}
	}
	// click all links
	var sendError error
	for _, link := range links {
		_, err := http.Get(link)
		sendError = errs.Combine(sendError, err)
	}

	return sendError
}

func findLinkTags(body string) []string {
	var tags []string
Loop:
	for {
		stTag := strings.Index(body, "<a")
		if stTag < 0 {
			break Loop
		}

		stripped := body[stTag:]
		endTag := strings.Index(stripped, "</a>")
		if endTag < 0 {
			break Loop
		}

		offset := endTag + len("</a>") + 1
		body = stripped[offset:]
		tags = append(tags, stripped[:offset])
	}
	return tags
}
