import { describe, it } from 'mocha'
import { assert } from 'chai';
import {ProfileParser} from "../main/ProfileParser";
import {Expression} from "../main/model/Expression";

describe("Expression negation", () => {
    it("Should negate correctly any rule expression", async () => {
        const parser = new ProfileParser("src/test/resources/profile1.yaml");
        const profile = await parser.parse();
        const validations = profile.validations

        const validation1 = validations.find((violation) => violation.name == "validation1");
        const negatedValidation1 = <Expression>validation1.negation();
        assert.equal(negatedValidation1.toString(), "validation1[VIOLATION] := ¬ ∃x : ( Class(x,'apiContract:Operation') ∧ ( ¬In(x,'apiContract:method',[publish,subscribe]) ∨ ¬MinCount(x,'apiContract:method',1) ∨ ¬Pattern(x,'shacl:name','^put|post$') ) )")
    });
});