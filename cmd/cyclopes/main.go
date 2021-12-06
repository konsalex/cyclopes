package cyclopes

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/chromedp/chromedp"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/konsalex/cyclopes/cmd/cyclopes/adapters"

	"github.com/pterm/pterm"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v2"
)

type DEVICE string

const (
	DESKTOP DEVICE = "desktop"
	MOBILE  DEVICE = "mobile"
	BOTH    DEVICE = "both"
)

type SCREENSHOTSIZE string

const FULLPAGE SCREENSHOTSIZE = "fullpage"
const VIEWPORT SCREENSHOTSIZE = "viewport"

const DEFAULT_URL = "http://localhost:3000"

type PageConfig struct {
	// Relative path to visit
	Path string
	// Device to screenshot
	Device DEVICE
	// Delay in ms
	Delay int
	// Code (Javascript code to execute)
	Code string
	// Selector to wait
	WaitSelector string
	// If the screenshot should be fullpage or viewport
	Screenshot SCREENSHOTSIZE
}

type VisualTesting struct {
	Pages []PageConfig

	// If the users want to debug
	Headless *bool `yaml:"headless"`

	// Base url to visit in our session
	RemoteURL string `yaml:"remoteURL"`
	/* Website directory
	(used only if remoteURL is not
	defined to serve the static assets)
	*/
	BuildDir string `yaml:"buildDir"`
}

type Configuration struct {
	PageConfig    `yaml:",inline"`
	VisualTesting `yaml:",inline"`
	// Directory to save/retrieve images
	ImagesDir string         `yaml:"imagesDir"`
	Visual    *VisualTesting `yaml:"visual,omitempty"`
	// Testing session id
	sessionID string
	Adapters  map[string]interface{} `yaml:"adapters"`
}

func Start(configPath string) {
	// YAML Unmarshal
	conf := Configuration{
		ImagesDir: "./cyclopes",
		sessionID: uuid.NewString(),
	}

	// Read config file
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		pterm.Fatal.Println(err)
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		pterm.Fatal.Println(err)
	}

	/* ------------- Loading Adapters ----------------- */
	var adapterInstances = make(map[string]adapters.AdapterInterface)
	for key := range conf.Adapters {
		fmt.Println(key)
		a, err := adapters.NewAdapter(key, &yamlFile)
		if err != nil {
			panic(err)
		}
		adapterInstances[key] = a
	}
	// Preflight checks for all adapters to avoid unnecessary test
	// if they are not properly configured
	for _, adapter := range adapterInstances {
		err := adapter.Preflight()
		if err != nil {
			panic(err)
		}
	}

	if conf.Visual != nil {
		pterm.Info.Println("Session ID: " + color.GreenString(conf.sessionID))

		if err := CheckPath(conf.ImagesDir); err != nil {
			pterm.Fatal.Println(err)
		}
		pterm.Info.Println("Images will be saved at: " + color.GreenString(conf.ImagesDir))

		// var srvr *Server
		// Source: https://stackoverflow.com/a/42533360/3247715
		if conf.Visual.RemoteURL == "" {
			if conf.Visual.BuildDir == "" {
				pterm.Fatal.Println("You must define either remoteURL or buildDir")
			}
			Server(conf.VisualTesting.BuildDir)
		}

		/** Enable Headless mode for testing **/
		var opts []chromedp.ExecAllocatorOption
		opts = chromedp.DefaultExecAllocatorOptions[:]

		if conf.Visual.Headless != nil && !*conf.Visual.Headless {
			pterm.Warning.Println("Headless mode is enabled")
			opts = append(opts, chromedp.Flag("headless", false),
				chromedp.Flag("disable-gpu", false),
				chromedp.Flag("enable-automation", false),
			)
		}

		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		/** Initialise ChromeDp **/
		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()

		pterm.Success.Println("Starting Visual Testing session")
		bar := progressbar.Default(int64(len(conf.Visual.Pages)))
		for _, v := range conf.Visual.Pages {
			Screenshot(ctx, &conf, &v, bar)
		}
	} else {
		pterm.Warning.Println("Skipping visual testing")
	}

	// Execute all defined adapters
	for _, adapter := range adapterInstances {
		err := adapter.Execute(conf.ImagesDir)
		if err != nil {
			panic(err)
		}
	}

	pterm.Success.Println("Finised Visual Testing")
}
