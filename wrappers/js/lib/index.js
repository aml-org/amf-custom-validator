import pako from "pako";
import {isBrowser} from "browser-or-node";
import compressedWasm from "../assets/main.wasm.gz";
import {loadPolyfills as loadWebPolyfills} from "./platform/web/polyfills/index";
import {loadPolyfills as loadNodePolyfills} from "./platform/node/polyfills/index";
import { Buffer } from "buffer";


let wasm_gz
let wasm

class WebAssemblySingleton {

    static initialized = false;
    static go = undefined;

    static async initialize() {
        if (this.isInitialized()) {
            return Promise.resolve()
        }
        WebAssemblySingleton.loadPolyfills(globalThis)
        WebAssemblySingleton.go = new Go();
        if(!wasm_gz || !wasm) {
            wasm_gz = Buffer.from(compressedWasm.split(",")[1], 'base64')
            wasm = pako.ungzip(wasm_gz)
        }

        if (WebAssembly) {
            const waitForInit = waitForWasmInitialization(globalThis)

            const initWa = WebAssembly.instantiate(wasm, WebAssemblySingleton.go.importObject)
                .then((result) => {
                    WebAssemblySingleton.go.run(result.instance, globalThis);
                }).catch(rejection => console.error(rejection));

            return Promise.all([initWa, waitForInit]).then(() => {
                WebAssemblySingleton.initialized = true;
            })
        } else {
            return Promise.reject("WebAssembly is not supported in your JS environment")
        }
    }

    static isInitialized() {
        return WebAssemblySingleton.initialized
    }

    static loadPolyfills(global) {
        if (isBrowser) {
            loadWebPolyfills(global)
        } else {
            loadNodePolyfills(global)
        }
    }
}

export class CustomValidatorFactory {
    static async create() {
        await WebAssemblySingleton.initialize()
        return new CustomValidatorSingleton()
    }
}

export class CustomValidatorSingleton {
    validateCustomProfile(profile, data, debug) {
        return validateKernel(profile, data, debug)
    }
    validateWithReportConfiguration(profile, data, debug, reportConfig) {
        return validateWithReportConfigurationKernel(profile, data, debug, reportConfig);
    }
    generateRego(profile) {
        return runGenerateRego(profile)
    }
    normalizeInput(data) {
        return runNormalizeInput(data)
    }
    exit() {
        __AMF__terminateValidator()
        WebAssemblySingleton.go.exit(0)
        WebAssemblySingleton.initialized = false;
    }
}

const validateKernel = function(profile, data, debug) {
    let before = new Date()
    const res = __AMF__validateCustomProfile(profile,data, debug);
    let after = new Date();
    if (debug) console.log("Elapsed : " + (after - before))
    return res;
}

const validateWithReportConfigurationKernel = function(profile, data, debug, reportConfig) {
    let before = new Date()
    const res = __AMF__validateCustomProfileWithConfiguration(profile,data, debug, undefined, reportConfig);
    let after = new Date();
    if (debug) console.log("Elapsed : " + (after - before))
    return res;
}

const runGenerateRego = function(profile) {
    const res = __AMF__generateRego(profile);
    return res;
}


const runNormalizeInput = function(data) {
    const res = __AMF__normalizeInput(data);
    return res;
}

const waitForWasmInitialization = (container) => new Promise((resolve) => {
    container.onWasmInitialized = resolve;
});