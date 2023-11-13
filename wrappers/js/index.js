const fs = require('fs');

require("./lib/wasm_exec");


const go = new Go()
let wasmModule, wasmInst

function removeWhiteSpace(string) {
    return string.replace(/[ \n\t]/g,'')
}
async function run(ruleset, data) {
    const source = fs.readFileSync("./lib/main.wasm")

    let report

    function assignReport(ptr, len) {
        console.log(`Assigning report: {ptr: ${ptr}, len: ${len}}`)
        const buf = new Uint8Array(wasmInst.exports.memory.buffer, ptr, len)
        report = new TextDecoder('utf8').decode(buf);
    }


    const env = {assignReport}
    wasmModule = await WebAssembly.compile(source)
    go.importObject.env = Object.assign(go.importObject.env, env)
    wasmInst = await WebAssembly.instantiate(wasmModule, go.importObject)
    go.run(wasmInst)


    const memory = new Memory(wasmInst.exports.memory.buffer, wasmInst.exports.alloc)

    const rulesetMemorySegment = memory.write(ruleset);
    const dataMemorySegment = memory.write(removeWhiteSpace(data));

    wasmInst.exports.validate(rulesetMemorySegment.pointer, rulesetMemorySegment.size, dataMemorySegment.pointer, dataMemorySegment.size)

    console.log(report)
}

class Memory {
    constructor(buffer, allocate) {
        this.buffer = buffer;
        this.allocate = allocate;
    }

    write(string) {
        const bytes = new TextEncoder("utf8").encode(string);
        const size = bytes.length;
        const startPointer = this.allocate(size);
        const allocatedMemory = new Uint8Array(this.buffer, startPointer, size);
        allocatedMemory.set(new Uint8Array(bytes));
        return new MemorySegment(startPointer, size);
    }
}

class MemorySegment {
    /**
     * @param pointer to the start of the memory segment
     * @param size of the memory segment
     */
    constructor(pointer, size) {
        this.pointer = pointer;
        this.size = size;
    }
}


module.exports.run = run
