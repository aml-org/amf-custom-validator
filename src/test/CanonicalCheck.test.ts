import { describe, it } from 'mocha'
import { assert } from 'chai';
import {AndRule} from "../main/model/rules/AndRule";
import {InRule} from "../main/model/constraints/InRule";
import {Quantification, Variable} from "../main/model/Rule";
import {OrRule} from "../main/model/rules/OrRule";
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
});