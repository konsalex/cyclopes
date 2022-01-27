package cyclopes

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/pterm/pterm"
)

// Screenshot takes a screenshot
// for a specific PageConfig
func Screenshot(
	ctx context.Context,
	config *Configuration,
	pageConfig *PageConfig) {

	serverPath := config.ExtractServerURL()

	serverURLConstructed, err := ConstructServerURL(serverPath)
	if err != nil {
		pterm.Fatal.Println(err)
	}

	var buf []byte

	// Check current URL, if the previous subdirectory part of the URL
	// is the same with the current url, we reload the page on purpose
	// Example:
	// Previous session: could be "/"
	// Current session: could be "/#about"
	var location string
	if err := chromedp.Run(ctx,
		chromedp.Tasks{
			chromedp.Sleep(time.Millisecond * 500),
			chromedp.Location(&location),
		}); err != nil {
		pterm.Fatal.Println(err)
	}

	currentURL, err := url.Parse(location)
	if err != nil {
		pterm.Error.Println(err)
	}

	parsedURL, err := url.Parse(pageConfig.Path)
	if err != nil {
		pterm.Error.Println(err)
	}

	scrolling := parsedURL.Path != currentURL.Path
	screenshotURL := serverURLConstructed + pageConfig.Path

	var filename string

	devices := []DEVICE{DESKTOP, MOBILE}
	for idx, device := range devices {
		if pageConfig.Device != BOTH && pageConfig.Device != device {
			continue
		}

		// Timeout in 20 seconds of running the task
		ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*20))
		defer cancel()

		/** Both Cases so two screenshots (in the second case we should not scroll) */
		pterm.Info.Printfln("[%s]: Screenshotting %s", device, screenshotURL)
		if err := chromedp.Run(ctx, fullScreenshot(
			screenshotURL,
			95,
			pageConfig,
			device,
			scrolling && idx == 0,
			&buf)); err != nil {
			pterm.Fatal.Println(err)
		}
		filename = fmt.Sprintf("%s-%s", string(device), strings.Replace(pageConfig.Path, "/", "", -1))
		error := SaveFile(buf, config.ImagesDir, filename)
		if error != nil {
			pterm.Warning.Println(err)
		}
	}
}

func fullScreenshot(
	urlstr string,
	quality int,
	pageConfig *PageConfig,
	screen DEVICE,
	scrolling bool,
	res *[]byte) chromedp.Tasks {

	// width, height
	deviceDim := make([]int64, 2)

	if screen == MOBILE {
		deviceDim[0], deviceDim[1] = 500, 1000
	} else {
		deviceDim[0], deviceDim[1] = 2560, 1600
	}

	tasks := chromedp.Tasks{
		chromedp.EmulateViewport(deviceDim[0], deviceDim[1]),
		chromedp.Navigate(urlstr),
		/** Sleep half second */
		chromedp.Sleep(time.Millisecond * 500),
	}

	/**
	 * Scroll with steps till the bottom of the page
	 * to force images to download
	 * (Example Gatsby)
	 */
	if scrolling && pageConfig.Screenshot == FULLPAGE {
		tasks = append(tasks,
			chromedp.Tasks{
				chromedp.Evaluate(`
					const scrollToBottomStepper = async () => {
					const height = window.visualViewport.height
					let start = height;	
				
					while(start<document.body.scrollHeight){
						window.scrollTo(0, start)
						start = start + height
						await new Promise(resolve => setTimeout(resolve, 200))
					}
					
					// Scroll to top
					window.scrollTo(0, 0)
					await new Promise(resolve => setTimeout(resolve, 200))
					}
				
					scrollToBottomStepper()
				`,
					nil, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
						return p.WithAwaitPromise(true)
					}),
			})
	}

	/*
		Waiting for an element to appear
		Important to state that first we wait for the element to appear
		and then we execute any JS code
	*/
	if pageConfig.WaitSelector != "" {
		tasks = append(tasks,
			chromedp.Tasks{
				chromedp.WaitVisible(pageConfig.WaitSelector)},
		)
	}

	/** Execute JS Code inside page */
	if pageConfig.Code != "" {
		tasks = append(tasks,
			chromedp.Tasks{
				chromedp.Evaluate(pageConfig.Code,
					nil, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
						return p.WithAwaitPromise(true)
					}),
			},
		)
	}

	tasks = append(tasks, chromedp.Tasks{
		/** User defined sleep */
		chromedp.Sleep(time.Millisecond * time.Duration(pageConfig.Delay)),
	})

	if pageConfig.Screenshot == VIEWPORT {
		tasks = append(tasks, chromedp.Tasks{
			/** User defined sleep */
			chromedp.Sleep(time.Millisecond * time.Duration(pageConfig.Delay)),
			chromedp.CaptureScreenshot(res),
		})
	} else {
		tasks = append(tasks, chromedp.Tasks{
			/** User defined sleep */
			chromedp.Sleep(time.Millisecond * time.Duration(pageConfig.Delay)),
			chromedp.FullScreenshot(res, quality),
		})

	}

	return tasks
}
