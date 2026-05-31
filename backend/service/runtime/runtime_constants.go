package runtime

// NodeLTSVersion is the Node.js LTS version pinned by the application.
const NodeLTSVersion = "v22.11.0"

// NodeDistBaseURL is the base URL for official Node.js distribution downloads.
const NodeDistBaseURL = "https://nodejs.org/dist"

// RuntimeStateFileName is the file used to persist runtime download state.
const RuntimeStateFileName = "state.json"

// NodeSubdir is the directory under data dir that holds the Node runtime tree.
const NodeSubdir = "runtime/node"
