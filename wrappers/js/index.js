require("./lib/wasm_exec")
const fs = require("fs")
const pako = require("pako")

let go = new Go()
let wasmCompilation = compileWasm()

function compileWasm() {
    if (WebAssembly) {
        let wasmGz = fs.readFileSync(__dirname + "/lib/main.wasm.gz")
        let wasm = pako.ungzip(wasmGz)
        return WebAssembly.compile(wasm)
    } else {
        throw new Error("WebAssembly is not supported in your JS environment")
    }
}

function validate(profile, data, debug) {
    return global._acv_validate(profile, data, debug)
}

function evaluate(callback) {
    return wasmCompilation
        .then(module => WebAssembly.instantiate(module, go.importObject))
        .then(instance => {
            global._acv_user_callback = callback
            return go.run(instance)
        })
}

module.exports.evaluate = evaluate
module.exports.validate = validate
