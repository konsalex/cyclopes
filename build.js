#!/usr/bin/env node
const spawn = require("child_process").spawn;
const os = require("os");
const process = require("process");

let ls;

// Run command depending on the OS
if (os.type() === "Linux")
  ls = spawn("./dist/cyclopes_linux_amd64/cyclopes", [
    process.argv.slice(2).join(" "),
  ]);
else if (os.type() === "Darwin")
  ls = spawn("./dist/cyclopes_darwin_amd64/cyclopes", [
    process.argv.slice(2).join(" "),
  ]);
else if (os.type() === "Windows_NT")
  ls = spawn("./dist/cyclopes_windows_amd64/cyclopes", [
    process.argv.slice(2).join(" "),
  ]);
else throw new Error("Unsupported OS found: " + os.type());

ls.stdout.on("data", (data) => {
  console.log(`stdout:\n${data}`);
});

ls.stderr.on("data", (data) => {
  console.log(`stderr:\n${data}`);
});

ls.on("error", (error) => {
  console.log(`error:\n${error.message}`);
});

ls.on("close", (code) => {
  console.log(`child process exited with code ${code}`);
});
