// rollup.config.mjs
export default {
    input: 'lib/index.js',
    output: {
        file: 'dist/bundle.js',
        format: 'umd',
        name: 'amf-custom-validator',
    },
};