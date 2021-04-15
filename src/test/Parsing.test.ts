import { describe, it } from 'mocha'
import { assert } from 'chai';
import {ProfileParser} from "../main/ProfileParser";

describe("Profile parsing", () => {
    it("Should parse a simple validation profile", async () => {
        const parser = new ProfileParser("src/test/resources/profile1.yaml");
        const profile = await parser.parse();
        const expressions = profile.toString();
        assert.equal(expressions, "validation1[VIOLATION] :=  ∀x : ( Class(x,'apiContract:Operation') → ( In(x,'apiContract:method',[publish,subscribe]) ∧ MinCount(x,'apiContract:method',1) ∧ Pattern(x,'shacl:name','^put|post$') ) )")
    });
});