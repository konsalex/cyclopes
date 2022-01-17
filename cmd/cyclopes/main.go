package cyclopes

import (
	"context"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/konsalex/cyclopes/cmd/cyclopes/adapters"
	"github.com/spf13/viper"

	"github.com/pterm/pterm"
	"github.com/schollz/progressbar/v3"
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

	// If users want to debug
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
	// Initialise Configuration struct
	conf := Configuration{
		sessionID: uuid.NewString(),
	}

	// Unmarshal YAML with Viper
	viper.AutomaticEnv()
	viper.SetConfigName(configPath)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		pterm.Fatal.Printfln("Fatal error config file: %s", err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		pterm.Fatal.Printfln("Unable to parse the config file, %s", err)
	}

	// Throw if ImagesDir is not defined
	if conf.ImagesDir == "" {
		pterm.Error.Println("imagesDir is not defined")
		panic("imagesDir is not defined")
	}

	/* ------------- Loading Adapters ----------------- */
	/*
	  We allow adapters to be empty arrays if the variables
	  are stored as env variables.
	  Check `example-configs/novisual.yml` for more info.
	*/
	var adapterInstances = make(map[string]adapters.AdapterInterface)
	pterm.Info.Println("Loading adapters")
	for key := range conf.Adapters {
		pterm.Info.Println("Adapter: " + key)
		a, err := adapters.NewAdapter(key)
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

	pterm.Info.Println("Session ID: " + color.GreenString(conf.sessionID))

	if conf.Visual != nil {
		if err := CheckPath(conf.ImagesDir); err != nil {
			pterm.Fatal.Println(err)
		}
		pterm.Info.Println("Images will be saved at: " + color.GreenString(conf.ImagesDir))

		// Source: https://stackoverflow.com/a/42533360/3247715
		if conf.Visual.RemoteURL == "" {
			if conf.Visual.BuildDir == "" {
				pterm.Fatal.Println("You must define either remoteURL or buildDir")
			}
			pterm.Info.Println("Serving static assets from: " + color.GreenString(conf.Visual.BuildDir))
			Server(conf.Visual.BuildDir)
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
		/*
			chromedp handles context unexpectly,
			so we need to run an empty task first
			https://github.com/chromedp/chromedp/issues/513
		*/
		if err := chromedp.Run(ctx); err != nil {
			pterm.Error.Println(err)
			panic(err)
		}

		pterm.Success.Println("Starting Visual Testing session")
		bar := progressbar.Default(int64(len(conf.Visual.Pages)))
		for _, v := range conf.Visual.Pages {
			Screenshot(ctx, &conf, &v)
			bar.Add(1)
		}
	} else {
		pterm.Info.Println("Skipping visual testing")
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
