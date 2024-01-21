import {loadGoPolyfills} from "./wasm_exec";

export const loadPolyfills = (global: any) => {
    return loadGoPolyfills(global)
}