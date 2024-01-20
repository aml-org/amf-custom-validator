import {loadPolyfills as loadGoPolyfills} from "../../common/polyfills";


export function loadPolyfills() {
    require("wrappers/js/lib/platform/web/polyfills/wasm_exec_node");
    loadGoPolyfills()
}