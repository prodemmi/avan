const assert = require("node:assert/strict");
const { readFileSync } = require("node:fs");
const vm = require("node:vm");

const homeLabel = { nodeValue: "Home" };
let toggle;
const button = { addEventListener: (_, handler) => { toggle = handler; } };
let index = -1;
const walker = {
  nextNode() { index++; return index < 1; },
  get currentNode() { return homeLabel; },
};
const document = {
  body: {},
  documentElement: {},
  createTreeWalker: () => { index = -1; return walker; },
  querySelectorAll: () => [],
  getElementById: () => button,
};
const storage = new Map();

vm.runInNewContext(readFileSync("public/admin.js", "utf8"), {
  document,
  NodeFilter: { SHOW_TEXT: 4 },
  localStorage: {
    getItem: (key) => storage.get(key),
    setItem: (key, value) => storage.set(key, value),
  },
});

assert.equal(document.documentElement.lang, "en");
toggle();
assert.equal(homeLabel.nodeValue, "خانه");
assert.equal(document.documentElement.dir, "rtl");
assert.equal(storage.get("avan-language"), "fa");
toggle();
assert.equal(homeLabel.nodeValue, "Home");
assert.equal(document.documentElement.dir, "ltr");
