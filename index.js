#!/usr/bin/env node
const { spawn } = require("child_process");
const path = require("path");
const os = require("os");

const platform = os.platform(); 
const arch = os.arch(); 

let bin;

if (platform === "darwin") {
  if (arch === "arm64") {
    bin = path.join(__dirname, "bin", "gitbroski-darwin");
  } else {
    bin = path.join(__dirname, "bin", "gitbroski-intel-darwin");
  }
} else if (platform === "linux") {
  bin = path.join(__dirname, "bin", "gitbroski-linux");
} else if (platform === "win32") {
  bin = path.join(__dirname, "bin", "gitbroski-windows.exe");
} else {
  console.error(`Unsupported platform: ${platform} (${arch})`);
  process.exit(1);
}

const args = process.argv.slice(2);

const child = spawn(bin, args, { stdio: "inherit" });

child.on("error", (err) => {
  console.error(`Error running binary: ${bin}`);
  console.error(err.message);
  process.exit(1);
});

child.on("exit", (code) => {
  process.exit(code || 0);
});
