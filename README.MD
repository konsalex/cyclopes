# <h1 align="center"> cyclopes </h1>

<!-- toc -->

- [Description](#description)
- [Installation](#installation)
- [How to use](#how-to-use)
- [How to configure](#how-to-configure)
- [Useful Parameters](#useful-parameters)
- [🔌 Adapters](#---adapters)
  - [Slack adapter](#slack-adapter)

<!-- tocstop -->

---

[![Release Cyclopes](https://github.com/konsalex/cyclopes/actions/workflows/release.yaml/badge.svg)](https://github.com/konsalex/cyclopes/actions/workflows/release.yaml)
![npm](https://img.shields.io/npm/v/cyclopes)
[![Gitter](https://badges.gitter.im/cyclopes-clip/community.svg)](https://gitter.im/cyclopes-clip/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

<span style="display:block;margin-top:80px" class="breaking-gap"/>
 
## Description

**Cyclopes is lightweight CLI tool to help us perform simple visual testing tasks, so we can feel more confident shipping our code.**

Visual testing with Cyclopes includes screenshotting our desired pages, but also sending the results to the tools we use (like Slack) so we can make a quick review.

Some useful use cases:

1. Visual test your GatsbyJS/NextJS etc. marketing websites. If you are not writing tests (which is often the case with marketing websites), you can use **cyclopes** to check if your website is rendered as expected.

2. Integrate visual testing with your CI, and send generated screeenshots to a dev Slack channel. The screenshots can be taken from Cypress.js or any other E2E testing framework you already use.

**Cyclopes** is only using Chrome to test as it is based in [chromedp](https://github.com/chromedp/chromedp).

> 💡 Cyclopes is not Visual Testing as a Service (Percy, Chromatic etc), but something to bridge the gap between not having a visual test, and having a complicated setup to work with.

## Installation

With `npm` or `yarn`

```bash
# npm
npm install cyclopes

# yarn
yarn add cyclopes
```

With `go get`

```bash
go get -u github.com/konsalex/cyclopes
```

Or you can download the binaries directly from [Github releases](https://github.com/konsalex/cyclopes/releases).

## How to use

The usage is fairly simple. You can run

```
cyclopes
```

which will try to find the default yaml file `cyclops.yml` in the current directory.

If you want to use a different file name, you can use

```bash
cyclopes --config=./your-config.yml
```

## How to configure

The configuration (yaml) file will start with two root steps, the `visual` and `adapters`, both which are optional. Below is an example of a yaml configuration file and below a table with the values and their possible options.

```yaml
# Directory where the images will be saved
# and/or retrieved from the channel adapters
imagesDir: "./images"

# If visual is defined, it will be the first step
# where the pages defined will be visited and screenshotted
# one by one with the defined order
visual:
  # (Optional): URL to visit. If this field does
  # not exist, builDir is required
  remoteURL: https://www.neo4j.com/
  # Directory to serve for testing.
  # Required only if remoteURL is missing
  buildDir: "./public"
  # If the testing session will be headless or not
  # for debugging purposes
  headless: true
  # Pages to visit and screenshot
  pages:
    # Path of the page
    - path: "/"
      # Type of device to screenshot
      device: "both"
      # Delay after the page is loaded
      delay: 2000
      # Wait for an element to appear
      waitSelector: "#NeoCybotCookiebotDialogBodyButtonAccept"
      # If the screenshot should be fullpage or viewport
      screenshot: "fullpage"
      # Code to be executed inside the webpage
      # This can help to remove Cookies banner for example
      code: |
        const cookiesButton = document.querySelector("#NeoCybotCookiebotDialogBodyButtonAccept")
        if (cookiesButton) {
            cookiesButton.click()
        }
    - path: "/product"
      device: "mobile"
      screenshot: "fullpage"
    - path: "/product/graph-data-science/"
      device: "desktop"
# Adapters
# Each adapter will have it's own parameters
adapters:
  slack:
    OAUTH_TOKEN: example-token
    CHANNEL_ID: example-token
```

## Useful Parameters

| parameter    | values                      | default value | required | description                                                                                     |
| ------------ | --------------------------- | ------------- | -------- | ----------------------------------------------------------------------------------------------- |
| `imagesDir`  | string                      | null          | yes      | The path where the images will be saved or retrieved from the adapters                          |
| `headless`   | bool                        | true          | no       | If our visual testing session will open the chrome browser in headless mode                     |
| `device`     | `desktop`, `mobile`, `both` | `both`        | no       | Viewport of device we should screenshot                                                         |
| `screenshot` | `fullpage`, `viewport`      | `fullpage`    | no       | Screenshotting the current viewport or fullpage screenshot                                      |
| `code`       | string                      |               | no       | Javascript code to execute inside page. It can be a Promise that will be waited to be fulfilled |

## 🔌 Adapters

Right now only Slack adapter is supported.

Soon more will be added like Trello etc. Any recommendation is welcome.

### Slack adapter

<details>
<summary>Expand</summary>

```yaml
slack:
  oauth_token: example-token
  channel_id: example-channel-id
```

or as _environment variables_

```sh
export ADAPTERS_SLACK_CHANNEL_ID=example-channel-id
export ADAPTERS_SLACK_OAUTH_TOKEN=example-token
```

</details>
