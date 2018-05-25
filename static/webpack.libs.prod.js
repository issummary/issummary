module.exports = {
  mode: 'production',
  entry: './src/libs.ts',
  output: {
    filename: 'libs.bundle.js',
    path: __dirname + '/dist'
  },

  resolve: {
    // Add '.ts' and '.tsx' as resolvable extensions.
    extensions: ['.ts', '.tsx', '.js', '.json']
  },

  node: {
    __dirname: false,
    __filename: false
  },

  module: {
    rules: [
      // All files with a '.ts' or '.tsx' extension will be handled by 'awesome-typescript-loader'.
      { test: /\.tsx?$/, loader: 'ts-loader' }
    ]
  }
};
