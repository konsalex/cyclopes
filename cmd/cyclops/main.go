package cyclops

import (
	"context"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v2"
)

type DEVICE string

const DESKTOP DEVICE = "desktop"
const MOBILE DEVICE = "mobile"
const BOTH DEVICE = "both"

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

type Configuration struct {
	PageConfig `yaml:",inline"`

	// If the users want to debug
	Headless bool `yaml:"headless"`
	// Directory to save images
	ImagesDir string `yaml:"imagesDir"`
	// If multi-threading is enabled (for multiple Chrome instances?)
	Multithreading bool
	// Testing session id
	sessionID string

	Pages []PageConfig

	// If the server is enabled we should start one
	Server bool
	// The below options are needed if we should
	// spin a server to serve the build
	// The path to the build directory
	BuildDir string `yaml:"buildDir"`
	// Server's URL (used only if Server is false)
	ServerURL string `yaml:"serverURL"`
}

func Start(configPath string) {
	// YAML Unmarshal
	conf := Configuration{
		Server:         false,
		BuildDir:       "/dist",
		ImagesDir:      "./cyclopes",
		Multithreading: true,
		Headless:       false,
		sessionID:      uuid.NewString(),
	}

	// Open config file
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		FatalPrint(err)
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatal(err)
	}

	serverURL, err := conf.ExtractServerURL()
	if err != nil {
		FatalPrint(err)
	}

	SuccessPrint("Server url is: " + serverURL)
	SuccessPrint("Dev server will start: " + strconv.FormatBool(conf.Server))
	SuccessPrint("Images will be saved at: " + conf.ImagesDir)
	SuccessPrint("Session ID: " + conf.sessionID)

	/** Enable Headless mode for testing **/
	var opts []chromedp.ExecAllocatorOption
	opts = chromedp.DefaultExecAllocatorOptions[:]

	if conf.Headless {
		WarningPrint("Headless mode is enabled")
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

	// var srvr *Server
	// Source: https://stackoverflow.com/a/42533360/3247715
	if conf.Server {
		Server(conf.BuildDir)
	}

	bar := progressbar.Default(int64(len(conf.Pages)))
	for _, v := range conf.Pages {
		Screenshot(ctx, &conf, &v, bar)
	}

	SuccessPrint("Finised visual testing")
}
