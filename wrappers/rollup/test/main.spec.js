const fs = require('fs')
const path = require('path')
const assert = require('assert');
const AmfCustomValidator = require('../dist/main')

describe('AMF Custom Validator', () => {
    it('should validate and return a report', () => {
        const validator = new AmfCustomValidator()

        // Paths
        const basePath = path.join('..', '..', 'test', 'data', 'integration', 'profile1')
        const profilePath = path.join(basePath, 'profile.yaml')
        const dataPath = path.join(basePath, 'negative.data.jsonld')

        // Read profile & data
        const profile = fs.readFileSync(profilePath, 'utf-8')
        const data = fs.readFileSync(dataPath, 'utf-8')

        // Validate
        validator
            .initialize()
            .then(validator.validate(profile, data))
            .then((report) => {
                const expectedReport = 'OK'
                assert.equal(report, expectedReport);
                done();
            })
    });
});
