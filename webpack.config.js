const path = require("path");
const DefinePlugin = require("webpack").DefinePlugin;
const CleanWebpackPlugin = require("clean-webpack-plugin");
const ExtractTextPlugin = require("extract-text-webpack-plugin");
const OptimizeCssAssetsPlugin = require("optimize-css-assets-webpack-plugin");
const HtmlWebpackPlugin = require("html-webpack-plugin");

module.exports = (env = {}, opts) => {
  const DEBUG = opts.mode === "development";
  const DIST_DIR = path.resolve(__dirname, "dist");
  const JS_NAME = DEBUG ? "index.js" : "[chunkhash:10].js";
  const CSS_NAME = DEBUG ? "index.css" : "[contenthash:10].css";
  return {
    entry: path.resolve(__dirname, "ts/index/index.tsx"),
    resolve: {
      extensions: [".tsx", ".ts", ".js"],
    },
    module: {
      rules: [
        {
          test: /\.tsx?$/,
          use: "ts-loader",
          exclude: /node_modules/,
        },
        {
          test: /\.less$/,
          use: ExtractTextPlugin.extract(["css-loader", "less-loader"]),
          exclude: /node_modules/,
        },
      ],
    },
    plugins: [
      // Common plugins.
      new CleanWebpackPlugin([env.output || DIST_DIR], {allowExternal: true}),
      new DefinePlugin({"window.KNET_API_PREFIX": JSON.stringify(env.api_prefix)}),
      new ExtractTextPlugin(`static/${CSS_NAME}`),
      new HtmlWebpackPlugin({
        title: "K-pop idols network | Profiles, images and face recognition",
        favicon: path.resolve(__dirname, "ts/index/favicon.ico"),
      }),
    ].concat(DEBUG ? [
      // Development only.
      new (require("webpack-notifier")),
      new (require("webpack-livereload-plugin")),
    ] : [
      // Production only.
      new OptimizeCssAssetsPlugin(),
    ]),
    output: {
      path: env.output || DIST_DIR,
      filename: `static/${JS_NAME}`,
      publicPath: "/",
    },
  };
};
