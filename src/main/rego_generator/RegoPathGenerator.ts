import * as md5 from "md5";
import {AndPath, OrPath, Property, PropertyPath} from "../profile_parser/PathParser";

export interface RegoPathResult {
    rego: string[],
    pathVariables: string[],
    variable: string,
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
            // AND we traverse each component and merge the results one after the other and generating a single result
            return this.traverseAnd(<AndPath>path, counter, rego, pathVariables, paths);
        //@ts-ignore
        } else if (path.or != null) {
            // OR, we need to branch the lookup in the path and generate two results for each branch, duplicating the accumulator
            return this.traverseOr(<OrPath>path, counter, rego, pathVariables, paths);
        } else {
            // SIMPLE property, we just add it to the current path
            return this.traverseProperty(<Property>path, counter, rego, pathVariables, paths);
        }
    }

    private traverseAnd(path: AndPath, counter: number, rego: string[], pathVariables: string[], paths: string[]): RegoPathResultInternal[] {
        const first = path.and[0];
        const last = path.and[1];
        const firstParsed = this.traverse(first, counter, rego, pathVariables, paths);
        const acc: RegoPathResultInternal[] = [];
        firstParsed.forEach((pathResult) => {
            this.traverse(last, pathResult.counter, pathResult.rego, paths,pathResult.pathVariables).forEach((secondPathResult) => {
                acc.push({
                    rego: pathResult.rego.concat(secondPathResult.rego),
                    pathVariables: pathResult.pathVariables.concat(secondPathResult.pathVariables),
                    paths: pathResult.paths.concat(secondPathResult.paths),
                    counter: secondPathResult.counter,
                    variable: secondPathResult.variable
                });
            });
        });
        return acc;
    }

    private traverseOr(path: OrPath, counter: number, rego: string[], pathVariables: string[], paths: string[]): RegoPathResultInternal[] {
        const first = path.or[0];
        const last = path.or[1];
        const firstParsed = this.traverse(first, counter, rego, pathVariables, paths);
        const acc: RegoPathResultInternal[] = [];
        firstParsed.forEach((pathResult) => {
            const newCounter = this.branchingCounter++ // We are generating a newly unique counter for the properties in this path
            this.traverse(last, newCounter, pathResult.rego, pathResult.pathVariables, pathResult.paths).forEach((pr) => acc.push(pr));
        });
        return acc;
    }

    private traverseProperty(path: Property, counter: number, rego: string[], pathVariables: string[], paths: string[]): RegoPathResultInternal[] {
        const idx = counter === 0 ? pathVariables.length : `${pathVariables.length}_${counter}`
        let binding = this.variable + "_" + idx + "_" + this.id + "_" + this.hint;
        const previousBinding = pathVariables[pathVariables.length-1];
        rego.push(`${binding} = ${previousBinding}["${path.iri}"]`)
        return [{
            rego: rego,
            pathVariables: pathVariables.concat([binding]),
            paths: paths.concat([path.iri]),
            counter: counter,
            variable: binding
        }]
    }

    public generatePropertyValues(): RegoPathResult {
        const paths = this.traverse(this.path, this.branchingCounter, [], [this.variable], []).map((result) => {
            result.rego.pop();
            const binding = result.variable;
            const previous_binding = result.pathVariables[result.pathVariables.length-2] || this.variable;
            const nextPath = result.paths[result.paths.length-1];
            result.rego.push(`${binding} = ${previous_binding}["${nextPath}"]`)
            return result;
        });

        if (paths.length === 1) {
            return paths[0];
        }  else {
            return this.accumulatePaths(paths);
        }
    }

    public generateNodeArray(): RegoPathResult {
        const paths = this.traverse(this.path, this.branchingCounter, [], [this.variable], []).map((result) => {
            result.rego.pop();
            const binding = result.variable;
            const previous_binding = result.pathVariables[result.pathVariables.length-2] || this.variable;
            const nextPath = result.paths[result.paths.length-1];
            result.rego.push(`nested_nodes[${binding}] with data.nodes as ${previous_binding}["${nextPath}"]`)
            return result;
        });

        if (paths.length === 1) {
            return paths[0];
        }  else {
            return this.accumulatePaths(paths)
        }
    }


    public generatePropertyArray(): RegoPathResult {

        const paths = this.traverse(this.path, this.branchingCounter, [], [this.variable], []).map((result) => {
            result.rego.pop();
            const binding = result.variable;
            const previous_binding = result.pathVariables[result.pathVariables.length-2] || this.variable;
            const nextPath = result.paths[result.paths.length-1];
            result.rego.push(`${binding} = object.get(${previous_binding},"${nextPath}",[])`)
            return result;
        });

        if (paths.length === 1) {
            return paths[0];
        }  else {
            return this.accumulatePaths(paths);
        }
    }


    /**
     * Since there are many alternative paths to reach the nodes, we need to provide a single
     * collection of nodes of the rest of the checks.
     * @param paths
     * @private
     */
    private accumulatePaths(paths: RegoPathResultInternal[]) {
        const variables = paths.map((result) => result.variable);
        const rego = [];
        paths.map((p) => p.rego.forEach((line) => rego.push(line)))
        this.branchingCounter++;
        const variable =  this.variable + "_final_" + this.branchingCounter + "_" + this.id + "_" + this.hint;
        for (let i=0; i<variables.length; i++) {
            const v = variables[i];
            if (i === 0) {
                rego.push(`${variable}_${i} = array.concat(${v},[])`);
            } else if (i === variables.length-1) {
                rego.push(`${variable} = array.concat(${v},${variables[i-1]})`);
            } else {
                rego.push(`${variable}_${i} = array.concat(${v},${variables[i-1]})`);
            }
        }
        variables.push(variable);
        return {
            rego: rego,
            pathVariables: variables,
            paths: paths[paths.length-1].paths,
            counter: this.branchingCounter,
            variable: variable
        };
    }

}