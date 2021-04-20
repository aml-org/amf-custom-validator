import {Variable, VariableCardinality} from "./model/Rule";

let counter = 0;


export const genvar = (hint: string) => {
    counter++;
    return `gen_${hint}_${counter}`;
}

export const reset = () => {
    counter = 0;
}

export class VarGenerator {
    private globalVarCounter = 0

    private vars = ["x", "y", "z", "p", "q", "r", "s", "t", "u", "v", "w", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o"]

    public genExpressionVar(quantification: string, cardinality?: VariableCardinality): Variable {
        let varName;
        if (this.globalVarCounter < this.vars.length) {
            varName = this.vars[this.globalVarCounter];

        } else {
            varName = `X${this.globalVarCounter}`
        }
        const v = new Variable(varName, quantification, cardinality);
        this.globalVarCounter++;
        return v;
    }
}