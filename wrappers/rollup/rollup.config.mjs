import resolve from '@rollup/plugin-node-resolve'
import commonjs from '@rollup/plugin-commonjs'
import typescript from '@rollup/plugin-typescript'

export default [
    {
        input: 'src/main.ts',
        output: {
            name: 'main',
            file: 'dist/main.js',
            format: 'umd'
        },
        plugins: [
            resolve(),
            commonjs(),
            typescript()
        ]
    }
];
