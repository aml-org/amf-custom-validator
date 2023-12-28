import resolve from '@rollup/plugin-node-resolve'
import commonjs from '@rollup/plugin-commonjs'
import replace from '@rollup/plugin-replace'
import fs from 'fs'
import pkg from './package.json' assert {type: 'json'} // Warning: importing JSON modules is an experimental JS feature

export default [
    {
        input: 'src/main.js',
        output: {
            name: 'main',
            file: 'dist/main.js',
            format: 'umd'
        },
        plugins: [
            resolve({
                preferBuiltins: false // prefer node libraries over dependencies
            }),
            commonjs(),
            replace({
                __version: pkg.version,
                __wasm_gz: () => {
                    return fs.readFileSync('./src/wasm/main.wasm.gz', {encoding: 'base64'})
                },
                preventAssignment: true
            })
        ]
    }
];