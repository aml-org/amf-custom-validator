// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

"use strict";


import crypto from "crypto";

/**
 * Go 1.19 split the wasm_exec file into two:
 *  - a base wasm_exec file that requires to provide polyfills
 *  - a Node wasm_exec_node file that provides those Node polyfills and adds some CLI behavior
 *
 * Removed the CLI behavior stuff and left only the polyfill provision
 */

export const loadGoPolyfills = (global) => {
    global.require = require;
    global.fs = require("fs");
    global.TextEncoder = require("util").TextEncoder;
    global.TextDecoder = require("util").TextDecoder;

    if (!global.performance || !global.performance.now) {
        global.performance = {
            now() {
                const [sec, nsec] = process.hrtime();
                return sec * 1000 + nsec / 1000000; // time in milliseconds
            },
        };
    }

    const crypto = require("crypto")
    if (!global.crypto || !global.crypto.getRandomValues) {
        global.crypto = {
            getRandomValues(b) {
                crypto.randomFillSync(b);
            },
        };
    }
}