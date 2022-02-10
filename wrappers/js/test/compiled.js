const fs = require("fs");
const assert = require('assert');

describe('validator', () => {

    describe('validate compiled', () => {

        function invalidReport(report) {
            return report[0]["doc:encodes"][0]["conforms"] === false
        }

        it("should load the WASM code, compile a profile, validated with pre-compiled profile, exit", (done) => {
            const profile = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/profile.yaml").toString()
            const data = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/negative.data.jsonld").toString()

            const validator = require(__dirname + "/../index")

            function nValidations(currentIteration, maxIterations, compiledProfile, data) {
                validator.validateCompiled(compiledProfile, data, false, (r, err) => {
                    if (err) {
                        done(err);
                    } else {
                        if (currentIteration < maxIterations) {
                            nValidations(currentIteration + 1, maxIterations, compiledProfile, data)
                        } else {
                            let report = JSON.parse(r)
                            assert.ok(invalidReport(report))
                            validator.exit();
                            done();
                        }
                    }
                });
            }

            validator.initialize(() => {
                validator.compileProfile(profile, false, (compiledProfile) => {
                    nValidations(0, 10, compiledProfile, data)
                })
            })
        });
    })
})