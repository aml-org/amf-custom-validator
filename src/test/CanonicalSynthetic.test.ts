import { describe, it } from 'mocha'
import { assert } from 'chai';
import {Expression} from "../main/model/Expression";
import {AndRule} from "../main/model/rules/AndRule";
import {InRule} from "../main/model/constraints/InRule";
import {AtomicRule, Quantification, Rule, Variable} from "../main/model/Rule";
import {OrRule} from "../main/model/rules/OrRule";
import {Implication} from "../main/model/Implication";

/**
 * Test Rule to geenrate statements
 */
class Pred extends AtomicRule {
    constructor(negated: boolean, variable: Variable, name : string) {
        super(negated, variable, "", [], "");
        this.name = name;
    }

    negation(): Rule {
        return new Pred(!this.negated, this.variable, this.name);
    }

    toCanonicalForm(): Rule {
        return this
    }

    toString(): string {
        let negation = ""
        if (this.negated) {
            negation = "¬"
        }
        return `${negation}${this.name}`
    }
}

/**
 * Auxiliary class to build predicates
 */
class PredBuilder {
    private variable: Variable;
    constructor(variable: Variable) {
        this.variable = variable;
    }

    pred(name: string) {
        return new Pred(false, this.variable, name);
    }

    notPred(name: string) {
        return new Pred(true, this.variable, name);
    }
}

// Common builder
const PB = new PredBuilder(new Variable("x", Quantification.ForAll));


describe("Canonical synthetic", () => {
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
    });

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
    });

    it ("Should normalize nested negations", async () => {
        const cond = new OrRule(true).withBody([
            new OrRule(true).withBody([
                new AndRule(false).withBody([
                    PB.pred("A"),
                    PB.pred("B")
                ]),
                new OrRule(false).withBody([
                    PB.pred("D")
                ])
            ])
        ]);

        const canonical = <Expression>cond.toCanonicalForm();
        assert.equal(canonical.toString(), "( ( A ∧ B ) ∨ D )");
    });

    it ("Should normalize complex negations", async () => {
        //const cond = new OrRule(true).withBody([

        //]);
        const cond = new OrRule(true).withBody([
            new AndRule(false).withBody([
                PB.pred("A"),
                PB.pred("B")
            ]),
            new OrRule(false).withBody([
                PB.pred("D"),
                PB.pred("E")
            ])
        ])
        const canonical = <Expression>cond.toCanonicalForm();
        assert.equal(canonical.toString(), "( ( ¬A ∧ ¬D ∧ ¬E ) ∨ ( ¬B ∧ ¬D ∧ ¬E ) )");
    });

    it ("Should normalize complex negations 2", async () => {
        const cond = new OrRule(true).withBody([
            new AndRule(false).withBody([
                PB.pred("A"),
                PB.pred("B"),
            ]),
            new OrRule(false).withBody([
                PB.pred("D"),
                PB.pred("E")
            ]),
            PB.pred("M")
        ]);
        const canonical = <Expression>cond.toCanonicalForm();
        assert.equal(canonical.toString(), "( ( ¬A ∧ ¬D ∧ ¬E ∧ ¬M ) ∨ ( ¬B ∧ ¬D ∧ ¬E ∧ ¬M ) )");
    });

    it("Should normalize conditionals with nested ANDs", () => {
        const cond = new Implication(false, new Variable("x", Quantification.ForAll),
            PB.pred("A"),
            new AndRule(false).withBody([
                PB.pred("B"),
                new AndRule(false).withBody([
                    PB.pred("C"),
                    PB.pred("D")
                ])
            ]));
        const canonical = <Expression>cond.toCanonicalForm();
        assert.equal(canonical.toString(), "( ¬A ∨ ( B ∧ C ∧ D ) )");
    })

    it("Should normalize nested conditions with nested ANDs", () => {
        const cond = new Implication(true, new Variable("x", Quantification.ForAll),
            PB.pred("CX"),
            new AndRule(false).withBody([
                PB.pred("MX"),
                new Implication(false , new Variable("y", Quantification.ForAll),
                    PB.pred("NXY"),
                    PB.pred("MY")
                )
            ]));
        const canonical = <Expression>cond.toCanonicalForm();
        assert.equal(canonical.toString(), "");
    })
});