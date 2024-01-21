import pako from "pako";
import {isBrowser} from "browser-or-node";
import compressedWasm from "../assets/main.wasm.gz";
import {loadPolyfills as loadWebPolyfills} from "./platform/web/polyfills/index";
import {loadPolyfills as loadNodePolyfills} from "./platform/node/polyfills/index";
import { Buffer } from "buffer";


class WebAssemblySingleton {

    static initialized = false;
    static go = undefined;
    static wasm = undefined;
    static wasm_gz = undefined;

    static async initialize() {
        if (this.isInitialized()) {
            return
        }
        WebAssemblySingleton.loadPolyfills(globalThis)
        WebAssemblySingleton.go = new Go();
        if(this.hasToLoadWasm()) {
            this.wasm_gz = Buffer.from(compressedWasm.split(",")[1], 'base64')
            this.wasm = pako.ungzip(this.wasm_gz)
        }

        if (WebAssembly) {
            const initRuntime = this.initializeWebAssemblyRuntime()
            const waitForInit = waitForWasmInitialization(globalThis)

            await Promise.all([initRuntime, waitForInit])
            this.initialized = true;
        } else {
            throw Error("WebAssembly is not supported in your JS environment")
        }
    }

    static hasToLoadWasm() {
        return !this.wasm_gz || !this.wasm
    }

    static isInitialized() {
        return this.initialized
    }

    static loadPolyfills(global) {
        if (isBrowser) {
            loadWebPolyfills(global)
        } else {
            loadNodePolyfills(global)
        }
    }

    static async initializeWebAssemblyRuntime() {
        const webAssembly = await WebAssembly.instantiate(this.wasm, this.go.importObject)
        this.go.run(webAssembly.instance, globalThis) // we don't await this on purpose
        return Promise.resolve()
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