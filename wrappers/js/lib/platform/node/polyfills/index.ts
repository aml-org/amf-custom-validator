import {loadPolyfills as loadCommonGoPolyfills} from "../../common/polyfills";
import {loadGoPolyfills} from "./wasm_exec_node";


export function loadPolyfills(global: any) {
    loadGoPolyfills(global)
    loadCommonGoPolyfills(global)
}