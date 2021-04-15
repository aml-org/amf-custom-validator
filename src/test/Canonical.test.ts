import { describe, it } from 'mocha'
import { assert } from 'chai';
import {Expression} from "../main/model/Expression";
import {AndRule} from "../main/model/rules/AndRule";
import {InRule} from "../main/model/constraints/InRule";
import {Quantification, Variable} from "../main/model/Rule";
import {OrRule} from "../main/model/rules/OrRule";
import {ProfileParser} from "../main/ProfileParser";
import {Canonical} from "../main/model/CanonicalCheck";

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
describe("Canonical form", () => {
    it("Should transform simple validations into canonical form", async () => {
        const parser = new ProfileParser("src/test/resources/profile1.yaml");
        const profile = await parser.parse();
        const validations = profile.validations

        const validation1 = validations.find((violation) => violation.name == "validation1");
        const canonical = <Expression>validation1.toCanonicalForm();
        assert.equal(canonical.toString(), "validation1[VIOLATION] := ¬ ∃x : ( ( ¬In(x,'apiContract:method',[publish,subscribe]) ∧ Class(x,'apiContract:Operation') ) ∨ ( ¬MinCount(x,'apiContract:method',1) ∧ Class(x,'apiContract:Operation') ) ∨ ( ¬Pattern(x,'shacl:name','^put|post$') ∧ Class(x,'apiContract:Operation') ) )");
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