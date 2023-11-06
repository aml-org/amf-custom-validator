const fs = require('fs');

require("./lib/wasm_exec");




const go = new Go()
let wasmModule, wasmInst


async function run() {
    const source = fs.readFileSync("./lib/main.wasm")

    function log(ptr, size) {
        console.log(memToString(ptr, size))
    }

    function memToString(ptr, len) {
        const buf = new Uint8Array(wasmInst.exports.memory.buffer, ptr, len)
        return new TextDecoder('utf8').decode(buf)
    }

    const env = {log}
    wasmModule = await WebAssembly.compile(source)
    go.importObject.env = Object.assign(go.importObject.env, env)
    wasmInst = await WebAssembly.instantiate(wasmModule, go.importObject)
    go.run(wasmInst)



    wasmInst.exports.greet()
}


module.exports.run = run
