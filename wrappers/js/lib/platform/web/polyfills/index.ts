import {loadPolyfills as loadCommonGoPolyfills} from "../../common/polyfills/index";
import {loadGoPolyfills} from "./wasm_exec_node.js";


export function loadPolyfills(global: object): void {
    loadGoPolyfills(global)
    loadCommonGoPolyfills(global)
}