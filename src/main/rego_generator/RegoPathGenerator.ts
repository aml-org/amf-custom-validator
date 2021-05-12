import * as md5 from "md5";
import {AndPath, OrPath, Property, PropertyPath} from "../profile_parser/PathParser";
import {genvar} from "../VarGen";

export interface RegoPathResult {
    rego: string[],
    rule: string,
    variable: string
}

interface RegoPathResultInternal {
    rego: string[],
    pathVariables: string[],
    paths: string[],
    variable: string,
    counter: number
}

export class RegoPathGenerator {

    private path: PropertyPath;
    private variable: string;
    private id: any;
    private hint: string;
    private branchingCounter: number;

    constructor(path: PropertyPath, variable: string, hint: string) {
        this.id = md5(path.source)
        this.path = path;
        this.variable = variable;
        this.hint = hint;
        this.branchingCounter = 0;
    }

    /**
     * This method is the common traversal method for all paths
     * @param path
     * @param counter
     * @param rego
     * @param pathVariables
     * @param paths
     */
    private traverse(path: PropertyPath, counter: number, rego: string[], pathVariables: string[], paths: string[]): RegoPathResultInternal[] {
        //@ts-ignore
        if (path.and != null) {
            // AND we traverse each component and merge the results one after the other to generate a single result
            return this.traverseAnd(<AndPath>path, counter, rego, pathVariables, paths);
        //@ts-ignore
        } else if (path.or != null) {
            // OR, we need to branch the path traversal and generate two results for each branch, duplicating the accumulator
            return this.traverseOr(<OrPath>path, counter, rego, pathVariables, paths);
        } else {
            // SIMPLE property, we just add it to the current path
            return this.traverseProperty(<Property>path, counter, rego, pathVariables, paths);
        }
    }

    private traverseAnd(path: AndPath, counter: number, rego: string[], pathVariables: string[], paths: string[]): RegoPathResultInternal[] {
        const first = path.and.shift();
        const firstParsed = this.traverse(first, counter, rego, pathVariables, paths);

        if (path.and.length > 0) {
            const acc: RegoPathResultInternal[] = [];
            firstParsed.forEach((pathResult) => {
                const next = {and: path.and.concat([])} // clone the path
                this.traverse(next, pathResult.counter, pathResult.rego, pathResult.pathVariables, pathResult.paths).forEach((secondPathResult) => {
                    acc.push(secondPathResult);
                });
            });

            return acc;
        } else {
            return firstParsed;
        }
    }

    private traverseOr(path: OrPath, counter: number, rego: string[], pathVariables: string[], paths: string[]): RegoPathResultInternal[] {
        let acc = [];
        path.or.forEach((pathElement) => {
            let counter = this.branchingCounter++;  // We are generating a newly unique counter for the properties in this path
            const parsed = this.traverse(pathElement, counter, rego.concat([]), pathVariables.concat([]), paths.concat([]));
            acc = acc.concat(parsed)
        });

        return acc;
    }

    private traverseProperty(path: Property, counter: number, rego: string[], pathVariables: string[], paths: string[]): RegoPathResultInternal[] {
        const idx = counter === 0 ? pathVariables.length : `${pathVariables.length}_${counter}`
        let binding = this.variable + "_" + idx + "_" + this.id + "_" + this.hint;
        const previousBinding = pathVariables[pathVariables.length-1];
        if (pathVariables.length === 0) {
            rego.push(`${this.variable} = data.sourceNode`)
            rego.push(`tmp_${binding} = nested_nodes with data.nodes as ${this.variable}["${path.iri}"]`)
            rego.push(`${binding} = tmp_${binding}[_][_]`)
        } else {
            rego.push(`tmp_${binding} = nested_nodes with data.nodes as ${previousBinding}["${path.iri}"]`)
            rego.push(`${binding} = tmp_${binding}[_][_]`)
        }
        return [{
            rego: rego,
            pathVariables: pathVariables.concat([binding]),
            paths: paths.concat([path.iri]),
            counter: counter,
            variable: binding
        }]
    }

    public generatePropertyValues(): RegoPathResult {
        const paths = this.traverse(this.path, this.branchingCounter, [], [], []).map((result) => {
            // remove the last comprehension
            if (result.rego.length > 2) {
                result.rego.pop();
                result.rego.pop();
            }
            const previous_binding = result.pathVariables[result.pathVariables.length-2] || this.variable;
            const nextPath = result.paths[result.paths.length-1];
            result.rego.push(`nodes_tmp = ${previous_binding}["${nextPath}"]`); // return value or array of values
            result.rego.push(`nodes_tmp2 = nodes_array with data.nodes as nodes_tmp`) // we make sure we got an array
            result.rego.push(`nodes = nodes_tmp2[_]`) // iterate through each element of the array to return int wrapped in the result
            return result;
        });

        return this.accumulatePaths(paths);
    }

    public generateNodeArray(): RegoPathResult {
        const paths = this.traverse(this.path, this.branchingCounter, [], [], []).map((result) => {
            result.rego.pop();
            result.rego.pop();
            const binding = result.variable;
            const previous_binding = result.pathVariables[result.pathVariables.length-2] || this.variable;
            const nextPath = result.paths[result.paths.length-1];
            result.rego.push(`tmp_${binding} = nested_nodes with data.nodes as ${previous_binding}["${nextPath}"]`)
            result.rego.push(`${binding} = tmp_${binding}[_][_]`) // iterate up to the individual level so I can return each individual element in the rule array
            result.rego.push(`nodes = ${binding}`)
            return result;
        });

        return this.accumulatePaths(paths)
    }


    public generatePropertyArray(): RegoPathResult {

        const paths = this.traverse(this.path, this.branchingCounter, [], [], []).map((result) => {
            result.rego.pop();
            result.rego.pop();
            const previous_binding = result.pathVariables[result.pathVariables.length-2] || this.variable;
            const nextPath = result.paths[result.paths.length-1];
            result.rego.push(`nodes_tmp = object.get(${previous_binding},"${nextPath}",[])`)
            result.rego.push(`nodes_tmp2 = nodes_array with data.nodes as nodes_tmp`) // this returns and array
            result.rego.push(`nodes = nodes_tmp2[_]`) // I need to iterate to each element in the array so it can be wrapped in the rule result
            return result;
        });

        return this.accumulatePaths(paths);
    }


    /**
     * Since there are many alternative paths to reach the nodes, we need to provide a single
     * collection of nodes of the rest of the checks.
     * @param paths
     * @private
     */
    private accumulatePaths(paths: RegoPathResultInternal[]): RegoPathResult {
        const variables = paths.map((result) => result.variable);
        // Let's generate a rule that will return the flat list of nodes in the path
        // If there are more than one paths (because of ORs) a rule with multiple clauses
        // will be generated and the final list of nodes will be the UNION of all the clauses
        const rego = [];
        const ruleName = genvar("path_rule");
        paths.map((p, i) => {
            if (i === 0 ) {
                rego.push(`${ruleName}[nodes] {`);
            } else {
                rego.push("} {")
            }
            p.rego.forEach((line) => {
                rego.push("  " + line);
            })
        });
        rego.push("}")

        return {
            rego: rego,
            rule: ruleName,
            variable: this.variable
        };
    }

}