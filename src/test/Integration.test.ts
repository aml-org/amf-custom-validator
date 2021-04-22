import { describe, it } from 'mocha'
import { assert } from 'chai';

import {Validator} from "../main/Validator";
import * as fs from "fs";

const runValidation = async (example: string, profile: string) => {
    const validator = new Validator(example, "RAML 1.0", "application/yaml", profile, null)
    return await validator.validate()
};

const runTest= async (name: string) => {
    const directory = "./src/test/resources/integration/"
    const profilePath = directory + name + ".yaml";
    const positiveFilePath = directory + name + ".positive.yaml"
    const positiveJSONLDFilePath = directory + name + ".positive.jsonld"
    const negativeFilePath = directory + name + ".negative.yaml"
    const negativeJSONLDFilePath = directory + name + ".negative.jsonld"

    const positiveResult = await runValidation(positiveFilePath, profilePath);
    //await fs.promises.writeFile(positiveJSONLDFilePath, JSON.stringify(positiveResult, null, 2))
    assert(positiveResult.conforms())
    const expectedPositive = await fs.promises.readFile(positiveJSONLDFilePath);
    assert.equal(JSON.stringify(positiveResult, null, 2), expectedPositive.toString())

    const negativeResult = await runValidation(negativeFilePath, profilePath);
    //await fs.promises.writeFile(negativeJSONLDFilePath, JSON.stringify(negativeResult, null, 2))
    assert(!negativeResult.conforms())
    const expectedNegative = await fs.promises.readFile(negativeJSONLDFilePath);
    assert.equal(JSON.stringify(negativeResult, null, 2), expectedNegative.toString())
};



    const files = fs.readdirSync("./src/test/resources/integration");
    const profiles: {[name: string]: true} = {};
    files.map((f) => f.split(".")[0])/*.filter((p) => p == "profile4")*/.forEach((p) => profiles[p] = true);
    Object.keys(profiles).forEach((profile) => {
        describe("IntegrationTests", () => {
            it("Integration test " + profile, async () => {
                const result = await runTest(profile)
                return result
            })
        })
    })
