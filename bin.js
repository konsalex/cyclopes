#!/usr/bin/env node
const spawn = require("child_process").spawn;
const os = require("os");
const process = require("process");
const package = require("./package.json");
const fs = require("fs");
const https = require("https");

let cmnd;
let name;
const { version } = package;
const basePath = `${__dirname}/dist`;

// If "dist" directory does not exist, create it
if (!fs.existsSync(basePath)) {
  fs.mkdir(basePath, { recursive: true }, (err) => {
    if (err) throw err;
  });
}

// Binary name depending on the OS
if (os.type() === "Linux") {
  name = `cyclopes-linux-amd64-${version}`;
} else if (os.type() === "Darwin") {
  name = `cyclopes-darwin-amd64-${version}`;
} else if (os.type() === "Windows_NT") {
  name = `cyclopes-windows-amd64-${version}.exe`;
} else throw new Error("Unsupported OS found: " + os.type());

// Binary path
const filename = `${basePath}/${name}`;

if (!fs.existsSync(filename)) {
  console.log("Binary not found, downloading...");
  // Cannot use github releases as they do not allow direct download
  const releasedBinary = `https://cyclopes.s3.eu-west-1.amazonaws.com/cyclopes/v${version}/${name}`;
  console.log(`Downloading: ${releasedBinary}`);
  download(releasedBinary, filename)
    .then(() => {
      startSubprocess(filename);
    })
    .catch((err) => {
      console.log(err);
    });
} else {
  startSubprocess(filename);
}

// https://stackoverflow.com/a/62786397/3247715
function download(url, dest) {
  return new Promise((resolve, reject) => {
    // Check file does not exist yet before hitting network
    fs.access(dest, fs.constants.F_OK, (err) => {
      if (err === null) reject("File already exists");

      const request = https.get(url, (response) => {
        if (response.statusCode === 200) {
          const file = fs.createWriteStream(dest, { mode: 0o755 });
          file.on("finish", () => resolve());
          file.on("error", (err) => {
            file.close();
            if (err.code === "EEXIST") reject("File already exists");
            else fs.unlink(dest, () => reject(err.message)); // Delete temp file
          });
          response.pipe(file);
        } else if (response.statusCode === 302 || response.statusCode === 301) {
          //Recursively follow redirects, only a 200 will resolve.
          download(response.headers.location, dest).then(() => resolve());
        } else {
          reject(
            `Server responded with ${response.statusCode}: ${response.statusMessage}`
          );
        }
      });

      request.on("error", (err) => {
        reject(err.message);
      });
    });
  });
}

function startSubprocess(filename) {
  // Execute the binary
  cmnd = spawn(filename, [process.argv.slice(2).join(" ")]);

  // Sub-process listeners
  cmnd.stdout.on("data", (data) => {
    console.log(`stdout:\n${data}`);
  });

  cmnd.stderr.on("data", (data) => {
    console.log(`stderr:\n${data}`);
  });

  cmnd.on("error", (error) => {
    console.log(`error:\n${error.message}`);
  });

  cmnd.on("close", (code) => {
    console.log(`child process exited with code ${code}`);
  });
}
