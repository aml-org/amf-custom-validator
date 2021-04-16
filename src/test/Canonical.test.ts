import { describe, it } from 'mocha'
import { assert } from 'chai';
import {Expression} from "../main/model/Expression";
import {ProfileParser} from "../main/ProfileParser";
import * as fs from "fs";

const canonicalTest = async (path: string, violationName: string) => {
    const parser = new ProfileParser(path);
    const profile = await parser.parse();
    const validations = profile.validations
    const expected = await fs.promises.readFile(path.replace(".yaml", ".canonical"))
    const validation1 = validations.find((violation) => violation.name == violationName);
    const canonical = <Expression>validation1.toCanonicalForm();
    return assert.equal(canonical.toString(), expected.toString());
}
describe("Canonical form", () => {
    it("Should transform simple validations into canonical form", async () => {
        return await canonicalTest("src/test/resources/profile1.yaml", "validation1");
    });

    it("Should transform simple validations with an OR rule into canonical form", async () => {
        return await canonicalTest("src/test/resources/profile2.yaml", "validation1");
    });

    it("Should transform simple validations with a simple OR rule into canonical form", async () => {
        return await canonicalTest("src/test/resources/profile3.yaml", "validation1");
    });
});