import { describe, it } from 'mocha'
import { assert } from 'chai';
import {ProfileParser} from "../main/ProfileParser";
import {Expression} from "../main/model/Expression";
import * as fs from "fs";

const testNegation = async (path: string, violationNameToTest: string) => {
    const parser = new ProfileParser(path);
    const expected = await fs.promises.readFile(path.replace(".yaml", ".negated"));
    const profile = await parser.parse();
    const validations = profile.validations;
    const validation1 = validations.find((violation) => violation.name == violationNameToTest);
    const negatedValidation1 = <Expression>validation1.negation();
    return assert.equal(negatedValidation1.toString(), expected)
}

describe("Expression negation", () => {
    it("Should negate correctly a simple rule expression", async () => {
        return await testNegation("src/test/resources/profile1.yaml", "validation1");
    });

    it("Should negate correctly an OR rule expression", async () => {
        return await testNegation("src/test/resources/profile2.yaml", "validation1");
    });
});