import { describe, it } from 'mocha'
import { assert } from 'chai';
import {ProfileParser} from "../main/ProfileParser";
import {RegoGenerator} from "../main/RegoGenerator";
import * as fs from "fs";
import {OPAWrapper} from "../main/OPAWrapper";
import {reset} from "../main/VarGen";

const loadGoldenFile = (path: string) => {
    const goldenPath = path.replace(".yaml", ".rego");
    return fs.readFileSync(goldenPath).toString().trim()
}

const testProfile = async (path: string) => {
    reset();
    const parser = new ProfileParser(path);
    const profile = await parser.parse();
    const rego = new RegoGenerator(profile).generate();
    await OPAWrapper.check(rego); // let's check that this is valid Rego
    //await fs.promises.writeFile(path.replace(".yaml", ".rego"), rego);
    assert.equal(rego, loadGoldenFile(path));
}

describe("Rego generation", () => {

    it("Should generate Rego code from an AMF profile", async () => {
        return await testProfile("src/test/resources/profile1.yaml");
    });

    it("Should generate Rego code from an AMF profile with an OR rule", async () => {
        return await testProfile("src/test/resources/profile2.yaml");
    });

    it("Should generate Rego code from an AMF profile with a simple OR rule", async () => {
        return await testProfile("src/test/resources/profile3.yaml");
    });

    it("Should generate Rego code from an AMF profile with a simple NESTED rule", async () => {
        return await testProfile("src/test/resources/profile4.yaml");
    });

    /*
    it("Should generate Rego code from an AMF profile with a simple QUALIFIED rule", async () => {
        return await testProfile("src/test/resources/profile5.yaml");
    });
     */

});