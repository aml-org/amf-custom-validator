import {ungzip} from 'pako'
import {Buffer} from 'buffer'

require("../js/lib/wasm_exec");

export default class AmfCustomValidator {
    // constants ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    VERSION: string = '__version'
    private WASM_GZ: string = '__wasm_gz'

    // variables ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    private initialized: boolean = false
    private debug: boolean = false

    // constructors ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    constructor() {}

    // public methods ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    validate(profile: string, data: string): Promise<string> {
        if (this.initialized) {
            // @ts-ignore
            return __AMF__validateCustomProfile(profile, data, this.debug);
        } else {
            throw new Error("WASM/GO not initialized")
        }
    }

    async initialize(): Promise<void> {
        const wasm = ungzip(Buffer.from(this.WASM_GZ, 'base64'))
        const go = new globalThis.Go()
        let result = await WebAssembly
            .instantiate(wasm, go.importObject);
        go.run(result.instance);
        this.initialized = true;
    }
}