const fs = require("fs")
const assert = require('assert')

describe('validator', () => {

    describe('validate', () => {

        function invalidReport(report) {
            return JSON.parse(report)[0]["doc:encodes"][0]["conforms"] === false
        }

        it("should load the WASM code, validate a profile, exit", () => {
            const profile = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/profile.yaml").toString()
            const data = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/negative.data.jsonld").toString()

            const validator = require("../index")

            return validator.evaluate(() => {
                let report1 = validator.validate(profile, data, false)
                let report2 = validator.validate(profile, data, false)
                assert.ok(invalidReport(report1) && invalidReport(report2))
            })
        })
    })
})