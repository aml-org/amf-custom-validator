const fs = require('fs');

require("./lib/wasm_exec");




const go = new Go()
let wasmModule, wasmInst


async function run() {
    const source = fs.readFileSync("./lib/main.wasm")

    let result
    function assignResult(ptr, len) {
        console.log(`Assigning result: {ptr: ${ptr}, len: ${len}}`)
        const buf = new Uint8Array(wasmInst.exports.memory.buffer, ptr, len)
        result = new TextDecoder('utf8').decode(buf);
    }

    const env = {assignResult}
    wasmModule = await WebAssembly.compile(source)
    go.importObject.env = Object.assign(go.importObject.env, env)
    wasmInst = await WebAssembly.instantiate(wasmModule, go.importObject)
    go.run(wasmInst)



    wasmInst.exports.greet("nico")
    console.log("Result")
    console.log(result)
}


module.exports.run = run
