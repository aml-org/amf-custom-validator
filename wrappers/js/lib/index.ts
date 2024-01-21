import pako from "pako";
import {isBrowser} from "browser-or-node";
// @ts-ignore
import compressedWasm from "../assets/main.wasm.gz";
import {loadPolyfills as loadWebPolyfills} from "./platform/web/polyfills/index";
import {loadPolyfills as loadNodePolyfills} from "./platform/node/polyfills/index";
import {Buffer} from "buffer";


class WebAssemblySingleton {

    static initialized: boolean = false;
    static go: Go | undefined = undefined;
    static wasm: BufferSource | undefined = undefined;
    static wasm_gz: Buffer | undefined = undefined;

    static async initialize() {
        if (this.isInitialized()) {
            return
        }
        this.loadPolyfills(globalThis)
        this.go = new Go();

        if(this.hasToLoadWasm()) {
            this.wasm_gz = Buffer.from(compressedWasm.split(",")[1], 'base64')
            this.wasm = pako.ungzip(this.wasm_gz)
        }

        if (WebAssembly) {
            const initRuntime = this.initializeWebAssemblyRuntime(this.wasm!, this.go)
            const waitForInit = waitForWasmInitialization(globalThis)

            await Promise.all([initRuntime, waitForInit])
            this.initialized = true;
        } else {
            throw Error("WebAssembly is not supported in your JS environment")
        }
    }

    static hasToLoadWasm(): boolean {
        return !this.wasm_gz || !this.wasm
    }

    static isInitialized(): boolean {
        return this.initialized
    }

    static loadPolyfills(global: any) {
        if (isBrowser) {
            loadWebPolyfills(global)
        } else {
            loadNodePolyfills(global)
        }
    }

    static async initializeWebAssemblyRuntime(wasm: BufferSource, go: Go): Promise<void> {
        const webAssembly = await WebAssembly.instantiate(wasm, go.importObject)
        this.go?.run(webAssembly.instance, globalThis) // we don't await this on purpose
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
    validateCustomProfile(profile: string, data: string, debug: boolean) {
        return validateKernel(profile, data, debug)
    }
    validateWithReportConfiguration(profile: string, data: string, debug: boolean, reportConfig: any) {
        return validateWithReportConfigurationKernel(profile, data, debug, reportConfig);
    }
    generateRego(profile: string) {
        return runGenerateRego(profile)
    }
    normalizeInput(data: string) {
        return runNormalizeInput(data)
    }
    exit() {
        __AMF__terminateValidator()
        WebAssemblySingleton.go?.exit(0)
        WebAssemblySingleton.initialized = false;
    }
}

const validateKernel = function(profile: string, data: string, debug: boolean) {
    let before = new Date()
    const res = __AMF__validateCustomProfile(profile,data, debug);
    let after = new Date();
    if (debug) console.log("Elapsed : " + (after.getTime() - before.getTime()))
    return res;
}

const validateWithReportConfigurationKernel = function(profile: string, data: string, debug: boolean, reportConfig: any) {
    let before = new Date()
    const res = __AMF__validateCustomProfileWithConfiguration(profile,data, debug, undefined, reportConfig);
    let after = new Date();
    if (debug) console.log("Elapsed : " + (after.getTime() - before.getTime()))
    return res;
}

const runGenerateRego = function(profile: string) {
    return __AMF__generateRego(profile);
}


const runNormalizeInput = function(data: string) {
    return __AMF__normalizeInput(data);
}

const waitForWasmInitialization = (container: any) => new Promise((resolve) => {
    container.onWasmInitialized = resolve;
});