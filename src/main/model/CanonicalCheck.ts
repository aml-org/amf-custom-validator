import {AtomicRule, Rule} from "./Rule";
import {OrRule} from "./rules/OrRule";
import {AndRule} from "./rules/AndRule";

export class Canonical {
    static check(rule: Rule): boolean {
        if (rule instanceof AtomicRule) {
            return true;
        } else if (rule instanceof OrRule) {
            for (let i=0; i<rule.body.length; i++) {
                const e = rule.body[i];
                if ((e instanceof OrRule) || !Canonical.check(e)) {
                    return false;
                }
            }
            return true;
        } else if (rule instanceof  AndRule) {
            for (let i=0; i<rule.body.length; i++) {
                const e = rule.body[i];
                if (!(e instanceof AtomicRule)) {
                    return false;
                }
            }
            return true;
        } else {
            return false;
        }
    }
}
