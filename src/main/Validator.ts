import {AmfParser} from "./model/AmfParser";
import {ProfileParser} from "./ProfileParser";
import {RegoParser} from "./RegoParser";
import {RegoGenerator} from "./RegoGenerator";
import {loadPolicy} from "@open-policy-agent/opa-wasm";
import * as fs from "fs";

export class Validator {
    private format: string;
    private mediaType: string;
    private customProfile: string;
    private customProfileWasm: string;
    private file: string;

    constructor(file: string, format: string, mediaType: string, customProfile: string, customProfileWasm?: string) {
        this.file = file;
        this.format = format;
        this.mediaType = mediaType;
        this.customProfile = customProfile;
        this.customProfileWasm = customProfileWasm;
    }

    async validate() {
        const parsedJSONLD = await new AmfParser(this.file, this.format, this.mediaType).parse();
        //console.log("=============")
        //console.log("DATA");
        //console.log(parsedJSONLD);
        const parsedProfile = await new ProfileParser(this.customProfile).parse();
        const rego = new RegoGenerator(parsedProfile).generate();
        //console.log("=============")
        //console.log("PROFILE");
        //console.log(rego);
        const wasmFile = await RegoParser.fromText(rego, parsedProfile.entrypoint());
        return await this.evalute(wasmFile, parsedJSONLD)
    }

    async evalute(wasmFile: string, parsedJSONLD: any) {
        try {
            const policyWasm = await fs.promises.readFile(wasmFile)
            const policy = await loadPolicy(policyWasm);
            policy.setData({});
            //@ts-ignore
            const result = policy.evaluate(parsedJSONLD);
            return result;
        } catch (e) {
            console.error("Error evaluating WASM policy: " + wasmFile);
            throw e;
        }
    }
}