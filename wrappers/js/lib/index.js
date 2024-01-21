import pako from "pako";
import {isBrowser} from "browser-or-node";
import compressedWasm from "../assets/main.wasm.gz";
import {loadPolyfills as loadWebPolyfills} from "./platform/web/polyfills/index";
import {loadPolyfills as loadNodePolyfills} from "./platform/node/polyfills/index";
import { Buffer } from "buffer";


let wasm_gz
let wasm

let initialized = false
let go = undefined;

const loadPolyfills = (global) => {
    if (isBrowser) {
        loadWebPolyfills(global)
    } else {
        loadNodePolyfills(global)
    }
}

export const validateKernel = function(profile, data, debug) {
    let before = new Date()
    const res = __AMF__validateCustomProfile(profile,data, debug);
    let after = new Date();
    if (debug) console.log("Elapsed : " + (after - before))
    return res;
}

export const validateWithReportConfigurationKernel = function(profile, data, debug, reportConfig) {
    let before = new Date()
    const res = __AMF__validateCustomProfileWithConfiguration(profile,data, debug, undefined, reportConfig);
    let after = new Date();
    if (debug) console.log("Elapsed : " + (after - before))
    return res;
}

export const validateCustomProfile = function(profile, data, debug, cb) {
    if (initialized) {
        let res = validateKernel(profile, data, debug);
        cb(res,undefined);
    } else {
        cb(undefined,new Error("WASM/GO not initialized"))
    }
}

export const validateWithReportConfiguration = function(profile, data, debug, reportConfig, cb) {
    if (initialized) {
        let res = validateWithReportConfigurationKernel(profile, data, debug, reportConfig);
        cb(res,undefined);
    } else {
        cb(undefined,new Error("WASM/GO not initialized"))
    }
}

export const runGenerateRego = function(profile) {
    const res = __AMF__generateRego(profile);
    return res;
}

export const generateRego = function(profile, cb) {
    if (initialized) {
        let res = runGenerateRego(profile);
        cb(res,undefined);
    } else {
        cb(undefined,new Error("WASM/GO not initialized"))
    }
}

export const runNormalizeInput = function(data) {
    const res = __AMF__normalizeInput(data);
    return res;
}

export const normalizeInput = function(data, cb) {
    if (initialized) {
        let res = runNormalizeInput(data);
        cb(res,undefined);
    } else {
        cb(undefined,new Error("WASM/GO not initialized"))
    }
}

export const initialize = function(cb) {
    if (initialized === true) {
        cb(undefined)
    }
    loadPolyfills(globalThis)
    go = new Go();
    if(!wasm_gz || !wasm) {
        wasm_gz = Buffer.from(compressedWasm.split(",")[1], 'base64')
        wasm = pako.ungzip(wasm_gz)
    }
    if (WebAssembly) {
        WebAssembly.instantiate(wasm, go.importObject).then((result) => {
            go.run(result.instance, globalThis);
            setTimeout(function() {
                initialized = true;
                cb(undefined);
            }, 1000); // still need to fix this
        }).catch(rejection => console.log(rejection));
    } else {
        cb(new Error("WebAssembly is not supported in your JS environment"));
    }
}

export const exit = function() {
    if(initialized) {
        __AMF__terminateValidator()
        go.exit(0)
        initialized = false;
    }
}