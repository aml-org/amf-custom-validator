import { describe, it } from 'mocha'
import { assert } from 'chai';
import {Expression} from "../main/model/Expression";
import {AndRule} from "../main/model/rules/AndRule";
import {InRule} from "../main/model/constraints/InRule";
import {Quantification, Variable} from "../main/model/Rule";
import {OrRule} from "../main/model/rules/OrRule";
import {ProfileParser} from "../main/ProfileParser";
import {Canonical} from "../main/model/CanonicalCheck";
import * as fs from "fs";

describe("Canonical check", () => {
    it ("Should check if a formula is canonical", () => {
        const v = new Variable("x", Quantification.ForAll);
        assert(Canonical.check(new InRule(false, v, ["a"], [""])));
        assert(Canonical.check(
            new AndRule(false).withBody([
                new InRule(false, v, ["a"], [""]),
                new InRule(false, v, ["b"], [""])
            ])
        ));
        assert(Canonical.check(
            new OrRule(false).withBody([
                new InRule(false, v, ["a"], [""]),
                new InRule(false, v, ["b"], [""])
            ])
        ));
        assert(!Canonical.check(
            new AndRule(false).withBody([
                new InRule(false, v, ["a"], [""]),
                new AndRule(false).withBody([new InRule(false, v, ["b"], [""])])
            ])
        ));
        assert(Canonical.check(
            new OrRule(false).withBody([
                new InRule(false, v, ["a"], [""]),
                new AndRule(false).withBody([new InRule(false, v, ["b"], [""])])
            ])
        ));
        assert(!Canonical.check(
            new OrRule(false).withBody([
                new InRule(false, v, ["a"], [""]),
                new AndRule(false).withBody([new InRule(false, v, ["b"], [""])]),
                new OrRule(false).withBody([new InRule(false, v, ["b"], [""])])
            ])
        ));
        assert(!Canonical.check(
            new AndRule(false).withBody([
                new OrRule(false).withBody([
                    new InRule(false, v, ["d"], [""]),
                    new InRule(false, v, ["l"], [""]),
                ]),
                new OrRule(false).withBody([
                    new InRule(false, v, ["m"], [""]),
                    new InRule(false, v, ["z"], [""]),
                ])
            ])
        ))

    });
})

const canonicalTest = async (path: string, violationName: string) => {
    const parser = new ProfileParser(path);
    const profile = await parser.parse();
    const validations = profile.validations
    const expected = await fs.promises.readFile(path.replace(".yaml", ".canonical"))
    const validation1 = validations.find((violation) => violation.name == violationName);
    const canonical = <Expression>validation1.toCanonicalForm();
    return assert.equal(canonical.toString(), expected);
}
describe("Canonical form", () => {
    it("Should transform simple validations into canonical form", async () => {
        return await canonicalTest("src/test/resources/profile1.yaml", "validation1");
    });

    it("Should transform simple validations with an OR rule into canonical form", async () => {
        return await canonicalTest("src/test/resources/profile2.yaml", "validation1");
    });

    it("Should transform simple validations with a simmple OR rule into canonical form", async () => {
        return await canonicalTest("src/test/resources/profile3.yaml", "validation1");
    });

    it("Should transform simple formulas 1", async () => {
        const v = new Variable("x", Quantification.ForAll);

        const cond = new AndRule(false).withBody([
            new InRule(false, v, ["a"], [""]),
            new OrRule(false).withBody([
                new InRule(false, v, ["c"], [""]),
                new InRule(false, v, ["d"], [""]),
            ]),
            new OrRule(false).withBody([
                new InRule(false, v, ["g"], [""]),
                new InRule(false, v, ["h"], [""]),
            ])
        ]);

        const canonical = <Expression>cond.toCanonicalForm();
        assert.equal(canonical.toString(),"( ( In(x,'c',[]) ∧ In(x,'a',[]) ∧ In(x,'g',[]) ) ∨ ( In(x,'d',[]) ∧ In(x,'a',[]) ∧ In(x,'g',[]) ) ∨ ( In(x,'c',[]) ∧ In(x,'a',[]) ∧ In(x,'h',[]) ) ∨ ( In(x,'d',[]) ∧ In(x,'a',[]) ∧ In(x,'h',[]) ) )")
    })


    it("Should transform simple formulas 2", async () => {
        const v = new Variable("x", Quantification.ForAll);

        const cond = new AndRule(false).withBody([
            new InRule(false, v, ["a"], [""]),
            new OrRule(false).withBody([
                new InRule(false, v, ["c"], [""]),
                new AndRule(false).withBody([
                    new OrRule(false).withBody([
                        new InRule(false, v, ["d"], [""]),
                        new InRule(false, v, ["l"], [""]),
                    ]),
                    new OrRule(false).withBody([
                        new InRule(false, v, ["m"], [""]),
                        new InRule(false, v, ["z"], [""]),
                    ])
                ])

            ]),
            new OrRule(false).withBody([
                new InRule(false, v, ["g"], [""]),
                new InRule(false, v, ["h"], [""]),
            ])
        ]);

        const canonical = <Expression>cond.toCanonicalForm();
        assert(canonical.toString(), "( ( In(x,'d',[]) ∧ In(x,'m',[]) ∧ In(x,'a',[]) ∧ In(x,'g',[]) ) ∨ ( In(x,'l',[]) ∧ In(x,'m',[]) ∧ In(x,'a',[]) ∧ In(x,'g',[]) ) ∨ ( In(x,'d',[]) ∧ In(x,'z',[]) ∧ In(x,'a',[]) ∧ In(x,'g',[]) ) ∨ ( In(x,'l',[]) ∧ In(x,'z',[]) ∧ In(x,'a',[]) ∧ In(x,'g',[]) ) ∨ ( In(x,'c',[]) ∧ In(x,'a',[]) ∧ In(x,'g',[]) ) ∨ ( In(x,'d',[]) ∧ In(x,'m',[]) ∧ In(x,'a',[]) ∧ In(x,'h',[]) ) ∨ ( In(x,'l',[]) ∧ In(x,'m',[]) ∧ In(x,'a',[]) ∧ In(x,'h',[]) ) ∨ ( In(x,'d',[]) ∧ In(x,'z',[]) ∧ In(x,'a',[]) ∧ In(x,'h',[]) ) ∨ ( In(x,'l',[]) ∧ In(x,'z',[]) ∧ In(x,'a',[]) ∧ In(x,'h',[]) ) ∨ ( In(x,'c',[]) ∧ In(x,'a',[]) ∧ In(x,'h',[]) ) )");
    })
});