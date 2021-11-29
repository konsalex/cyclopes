package cyclops

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/schollz/progressbar/v3"
)

func Screenshot(
	ctx context.Context,
	config *Configuration,
	pageConfig *PageConfig,
	bar *progressbar.ProgressBar) {

	serverPath, err := config.ExtractServerURL()
	if err != nil {
		FatalPrint(err)
	}

	serverURLConstructed, err := ConstructServerURL(serverPath)
	if err != nil {
		FatalPrint(err)
	}

	var buf []byte

	// Check current URL, if the previous subdirectory part of the URL
	// is the same with the current url, we reload the page on purpos
	var location string
	if err := chromedp.Run(ctx,
		chromedp.Tasks{
			chromedp.Sleep(time.Millisecond * 2000),
			chromedp.Location(&location),
		}); err != nil {
		log.Fatal(err)
	}

	currentURL, err := url.Parse(location)
	if err != nil {
		log.Fatal(err)
	}

	parsedURL, err := url.Parse(pageConfig.Path)
	if err != nil {
		log.Fatal(err)
	}

	scrolling := parsedURL.Path != currentURL.Path
	screenshotURL := serverURLConstructed + pageConfig.Path

	var filename string

	if pageConfig.Device == DESKTOP || pageConfig.Device == MOBILE {
		if err := chromedp.Run(ctx, fullScreenshot(
			screenshotURL,
			95,
			pageConfig,
			pageConfig.Device,
			scrolling,
			&buf)); err != nil {
			log.Fatal(err)
		}

		filename = fmt.Sprintf("%s-%s", string(pageConfig.Device), strings.Replace(pageConfig.Path, "/", "", -1))

		error := SaveFile(buf, config.ImagesDir, filename)
		if error != nil {
			WarningPrint("Error happened when trying to create image: " + "testing")
			WarningPrint(error)
		}
	} else {
		WarningPrint("Going for two screenshotsss")
		devices := []DEVICE{DESKTOP, MOBILE}
		for idx, device := range devices {
			/** Both Cases so two screenshots (in the second case we should not scroll) */
			if err := chromedp.Run(ctx, fullScreenshot(
				screenshotURL,
				95,
				pageConfig,
				device,
				scrolling && idx == 0,
				&buf)); err != nil {
				log.Fatal(err)
			}
			filename = fmt.Sprintf("%s-%s", string(device), strings.Replace(pageConfig.Path, "/", "", -1))
			error := SaveFile(buf, config.ImagesDir, filename)
			if error != nil {
				WarningPrint("Error happened when trying to create image: " + "testing")
				WarningPrint(error)
			}
		}
	}

	bar.Add(1)
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
	 * Scroll with steps will the bottom of the page
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

	/** Waiting for an element to appear */
	if pageConfig.WaitSelector != "" {
		SuccessPrint("Selector is :" + pageConfig.WaitSelector)
		tasks = append(tasks,
			chromedp.Tasks{
				chromedp.WaitVisible(pageConfig.WaitSelector)},
		)
	}

	tasks = append(tasks, chromedp.Tasks{
		/** User defined sleep */
		chromedp.Sleep(time.Millisecond * time.Duration(pageConfig.Delay)),
	})

	if pageConfig.Screenshot == FULLPAGE {
		tasks = append(tasks, chromedp.Tasks{
			/** User defined sleep */
			chromedp.Sleep(time.Millisecond * time.Duration(pageConfig.Delay)),
			chromedp.FullScreenshot(res, quality),
		})
	} else {
		tasks = append(tasks, chromedp.Tasks{
			/** User defined sleep */
			chromedp.Sleep(time.Millisecond * time.Duration(pageConfig.Delay)),
			chromedp.CaptureScreenshot(res),
		})
	}

	return tasks
}
