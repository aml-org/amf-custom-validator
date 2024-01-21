import {loadGoPolyfills} from "./wasm_exec";

export const loadPolyfills = (global) => {
    return loadGoPolyfills(global)
}