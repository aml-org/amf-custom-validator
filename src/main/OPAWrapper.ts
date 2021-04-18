import { exec } from 'child_process';
import * as os from 'os'
import * as uuid from 'uuid';
import * as fs from "fs";

export class OPAWrapper {

    private path: string;
    private entrypoint: string;

    constructor(path: string, entrypoint: string) {
        this.path = path;
        this.entrypoint = entrypoint;
    }

    parse(): Promise<string> {
        return new Promise<string>((resolve, reject) => {
            const tmpFile = `${os.tmpdir()}/opa_wasm${uuid.v4().replace(/-/g, "")}`
            const command = `./opa build -t wasm -e ${this.entrypoint} ${this.path} -o ${tmpFile} --capabilities opa_capabilities.json`
            exec(command, (error, stdout, stderr) => {
                if (error) {
                    console.log(`error: ${error.message}`);
                    reject(error);
                }
                const tmpDir = tmpFile + "_dir"
                fs.mkdirSync(tmpDir);
                const tarCommand = `tar -zxvf ${tmpFile} --directory ${tmpDir}`;
                exec(tarCommand, (error, stdout, stderr) => {
                    if (error) {
                        console.log(`error: ${error.message}`);
                        reject(error);
                    }
                    resolve(tmpDir + "/policy.wasm");
                });
            });
        })
    }

    static async fromText(text: string, entrypoint: string): Promise<string> {
        const tmpFile = `${os.tmpdir()}/rego${uuid.v4().replace(/-/g, "")}`;
        fs.writeFileSync(tmpFile, text);
        try {
            const parsed = await  new OPAWrapper(tmpFile, entrypoint).parse();
            fs.unlinkSync(tmpFile);
            return parsed;
        } catch (e) {
            console.error("Error parsing generated rego script:")
            console.error(text)
            throw e;
        }
    }

    static async check(text: string): Promise<string> {
        return new Promise<string>((resolve, reject) => {
            const tmpFile = `${os.tmpdir()}/rego${uuid.v4().replace(/-/g, "")}`;
            fs.writeFileSync(tmpFile, text);
            const command = `./opa check ${tmpFile}`
            exec(command, (error, stdout, stderr) => {
                fs.unlinkSync(tmpFile);
                if (error) {
                    console.log(`error: ${error.message}`);
                    reject(error);
                }
                if (stderr) {
                    reject(new Error(stderr))
                }
                resolve(stdout);
            });
        })
    }

}