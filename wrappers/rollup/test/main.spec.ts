import {assert} from "chai";
import {describe, it} from "mocha";
import AmfCustomValidator from "../src/main";

describe("AMF Custom Validator", () => {
    it("should validate and return a report", () => {
        const validator = new AmfCustomValidator()

        const profile = "some ruleset"
        const data = "some data"

        const actualReport = validator.validate(profile, data)
        const expectedReport = "OK"

        assert.equal(actualReport, expectedReport);
    });
});