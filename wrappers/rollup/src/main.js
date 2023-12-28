import {ungzip} from 'pako'
import {Buffer} from 'buffer'

export default class AmfCustomValidator {
    // constants ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    VERSION = '__version'
    WASM_GZ = '__wasm_gz'

    // variables ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    initialized = false
    debug = false

    // constructors ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    constructor() {}

    // public methods ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    async validate(profile, data) {
        if (this.initialized) {
            // @ts-ignore
            return __AMF__validateCustomProfile(profile, data, this.debug);
        } else {
            throw new Error("WASM/GO not initialized")
        }
    }

    async initialize() {
        const wasm = ungzip(Buffer.from(this.WASM_GZ, 'base64'))
        const go = new Go()
        let result = await WebAssembly
            .instantiate(wasm, go.importObject);
        go.run(result.instance);
        this.initialized = true;
    }
}
