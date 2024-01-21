// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

"use strict";


/**
 * Go 1.19 split the wasm_exec file into two:
 *  - a base wasm_exec file that requires to provide polyfills
 *  - a Node wasm_exec_node file that provides those Node polyfills and adds some CLI behavior
 *
 * Removed the CLI behavior stuff and left only the polyfill provision
 */
export const loadGoPolyfills = (global) => {
    if (!global.require) {
        global.require = require;
    }

    if (!global.require) {
        global.require = require;
    }

    if (!global.TextEncoder) {
        global.TextEncoder = require("util").TextEncoder;
    }

    if (!global.TextDecoder) {
        global.TextDecoder = require("util").TextDecoder;
    }

    if (!global.performance) {
        global.performance = {
            now() {
                return Date.now();
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