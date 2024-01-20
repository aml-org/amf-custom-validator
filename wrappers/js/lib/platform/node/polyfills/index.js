import {loadPolyfills as loadGoPolyfills} from "../../common/polyfills";


export function loadPolyfills() {
    require("wrappers/js/lib/platform/node/polyfills/wasm_exec_node");
    loadGoPolyfills()
}