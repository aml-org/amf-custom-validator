require("../js/lib/wasm_exec")

// should cache this
const readWasm = () => {
    const pako = require("pako")
    const wasmCompressed = require("../js/lib/main.wasm.gz")
    return pako.ungzip(Buffer.from(wasmCompressed, 'base64'))
}

// should cache this
const compileWasmModule = () => WebAssembly.compile(readWasm())

// should create one instance per validate call
const createWasmInstance = (go) => compileWasmModule().then(module => WebAssembly.instantiate(module, go.importObject))

const validate = (profile, data, debug, callback) => {
    const go = new Go()
    createWasmInstance(go)
        .then(instance => {
            go.run(instance)
        })
        .then(() => {
            const result = __AMF__validateCustomProfile(profile, data, debug)
            go.exit(0)
            __AMF__terminateValidator()
            callback(result)
        })
}

module.exports.validate = validate