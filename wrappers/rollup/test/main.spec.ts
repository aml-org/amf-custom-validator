import {assert} from 'chai';
import {describe, it} from 'mocha';
const acv = require('../dist/main')

describe('AMF Custom Validator', () => {
    it('should validate and return a report', () => {
        const validator = new acv.AmfCustomValidator()

        const profile = 'some ruleset'
        const data = 'some data'

        const actualReport = validator.validate(profile, data)
        const expectedReport = 'OK'

        assert.equal(actualReport, expectedReport);
    });
});