import {AmfParser} from "./model/AmfParser";
import {ProfileParser} from "./ProfileParser";
import {OPAWrapper} from "./OPAWrapper";
import {RegoGenerator} from "./RegoGenerator";
import {loadPolicy} from "@open-policy-agent/opa-wasm";
import * as fs from "fs";
import {Report} from "./validator/Report";

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

    async validate(): Promise<Report> {
        const parsedJSONLD = await new AmfParser(this.file, this.format, this.mediaType).parse();
        //console.log("=============")
        //console.log("DATA");
        //console.log(JSON.stringify(parsedJSONLD,null,2));
        const parsedProfile = await new ProfileParser(this.customProfile).parse();
        const rego = new RegoGenerator(parsedProfile).generate();
        //console.log("=============")
        //console.log("PROFILE");
        //console.log(rego);
        const wasmFile = await OPAWrapper.fromText(rego, parsedProfile.entrypoint());
        const evaluationReport = await this.evalute(wasmFile, parsedJSONLD)
        return new Report(evaluationReport);
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