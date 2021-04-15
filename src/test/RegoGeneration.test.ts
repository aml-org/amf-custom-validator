import { describe, it } from 'mocha'
import { assert } from 'chai';
import {ProfileParser} from "../main/ProfileParser";
import {RegoGenerator} from "../main/RegoGenerator";
import * as fs from "fs";
import {RegoParser} from "../main/RegoParser";

const loadGoldenFile = (path: string) => {
    const goldenPath = path.replace(".yaml", ".rego");
    return fs.readFileSync(goldenPath).toString().trim()
}

const testProfile = async (path: string) => {
    const parser = new ProfileParser(path);
    const profile = await parser.parse();
    const rego = new RegoGenerator(profile).generate();
    await RegoParser.check(rego); // let's check that this is valid Rego
    assert.equal(rego, loadGoldenFile(path));
}

describe("Rego generation", () => {

    it("Should generate Rego code from an AMF profile", async () => {
        return await testProfile("src/test/resources/profile1.yaml");
    });

});