package parser

import (
	"context"

	"github.com/chromedp/chromedp"
)

func NewBrowser() (context.Context, context.CancelFunc) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("start-maximized", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(
		context.Background(),
		opts...,
	)

	ctx, ctxCancel := chromedp.NewContext(allocCtx)

	cancel := func() {
		ctxCancel()
		allocCancel()
	}

	return ctx, cancel
}
