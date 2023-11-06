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

    const stringParameter = "Nico";
    const bytes = new TextEncoder("utf8").encode(stringParameter);
    const ptr = wasmInst.exports.alloc(bytes.length);
    const mem = new Uint8Array(
        wasmInst.exports.memory.buffer, ptr, bytes.length
    );
    mem.set(new Uint8Array(bytes));

    wasmInst.exports.greet(ptr, bytes.length)
    console.log("Result")
    console.log(result)
}


module.exports.run = run
