import {loadPolyfills as loadCommonGoPolyfills} from "../../common/polyfills/index";
import {loadGoPolyfills} from "./wasm_exec_node";


export function loadPolyfills(global) {
    loadGoPolyfills(global)
    loadCommonGoPolyfills(global)
}