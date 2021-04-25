import { describe, it } from 'mocha'
import { assert } from 'chai';
import {ProfileParser} from "../main/ProfileParser";
import * as fs from "fs";

const testParsing = async (profilePath: string) => {
    const parser = new ProfileParser(profilePath);
    const profile = await parser.parse();
    const expressions = profile.toString();
    //fs.writeFileSync(profilePath.replace(".yaml",".parsed"), expressions);
    const expected = await fs.promises.readFile(profilePath.replace(".yaml", ".parsed"));
    return assert.equal(expressions, expected.toString());
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

    it("Should parse a simple nested rule in a profile", async () => {
        return await testParsing("src/test/resources/profile4.yaml")
    });

    it("Should parse a simple qualified rule in a profile", async () => {
        return await testParsing("src/test/resources/profile5.yaml")
    });

    it("Should parse a simple lessThan rule in a profile", async () => {
        return await testParsing("src/test/resources/profile6.yaml")
    });

    it("Should parse complex negation rules in a profile", async () => {
        return await testParsing("src/test/resources/profile7.yaml")
    });

    it("Should parse complex property paths in a profile", async () => {
        return await testParsing("src/test/resources/profile8.yaml")
    });
});