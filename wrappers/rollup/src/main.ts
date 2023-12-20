import wasm_gz from './wasm/main.wasm.gz'

export default class AmfCustomValidator {
    private initialized: boolean = false
    private version: string = "__amf_custom_validator_version__"


    constructor() {}

    validate(profile: string, data: string): string {
        return "OK"
    }
}