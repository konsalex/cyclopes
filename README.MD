# Cyclopes

<!-- toc -->

- [🗒 Description](#-description)
- [💻 How to configure](#-how-to-configure)
- [🔌 Adapters](#-adapters)
<!-- tocstop -->

---

## 🗒 Description

Cyclopes is a CLI tool to help us manage simple visual testing tasks, so we can feel more confident shipping our code.

Visual testing with Cyclopes includes screenshotting our desired pages, but also sending the results to the tools we use (like Slack) so we can review the results.

<!-- Some useful use cases:

1. Visual test our Gatsby/NextJS etc marketing websites. If you are not writing tests (which is often the case with marketing websites), you can use this tool to check if your website is rendered as expected.

2. Perform quick sanity checks, like checking if your website is not having any broken style, and all this inside the tools of your workflow (ex. Slack). -->

<span style="display:block;margin-top:80px" class="breaking-gap"/>

## 💻 How to configure

The configuration (yaml) file will start we two root steps, the `visual` and `adapters`, both which are optional. Below is an example of a yaml configuration file and then a table with the values and their possible options.

```yaml
# Directory where the images will be saved
# or retrieved from the channel adapters
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

<span style="display:block;margin-top:80px" class="breaking-gap"/>

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
