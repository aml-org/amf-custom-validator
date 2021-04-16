import { describe, it } from 'mocha'
import { assert } from 'chai';
import {ProfileParser} from "../main/ProfileParser";
import * as fs from "fs";

const testParsing = async (profilePath: string) => {
    const expected = await fs.promises.readFile(profilePath.replace(".yaml", ".parsed"));
    const parser = new ProfileParser(profilePath);
    const profile = await parser.parse();
    const expressions = profile.toString();
    return assert.equal(expressions, expected);
}
describe("Profile parsing", () => {
    it("Should parse a simple validation profile", async () => {
        return await testParsing("src/test/resources/profile1.yaml")
    });

    it("Should parse an or validation profile", async () => {
        return await testParsing("src/test/resources/profile2.yaml")
    });

    it("Should parse a simple or validation profile", async () => {
        return await testParsing("src/test/resources/profile3.yaml")
    });
});