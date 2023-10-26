const path = require("path");
const TerserPlugin = require('terser-webpack-plugin');
module.exports = {
    entry: './index.js',
    output: {
        path: path.join(__dirname, 'dist'),
        filename: 'main.js',
        libraryTarget: 'commonjs',
    },
    module: {
      rules: [
        {
         test: /main\.wasm\.gz/,
         type: 'asset/inline',
         generator: {
          dataUrl: content => {
            return Buffer.from(content).toString('base64');
          }
        }
        }
      ]
    },
    resolve: {
      fallback: {
        "os": require.resolve("os-browserify/browser"),
        "util": require.resolve("util/"),
        "crypto": require.resolve("crypto-browserify"),
        "buffer": require.resolve("buffer/"),
        "stream": require.resolve("stream-browserify"),
        "fs": false,
        "performance": false
      }
    },
    optimization: {
      minimize: true,
      minimizer: [new TerserPlugin({
        extractComments: false,
      })],
    },
}