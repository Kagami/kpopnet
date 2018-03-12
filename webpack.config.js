const path = require("path");
const LiveReloadPlugin = require("webpack-livereload-plugin");
const CleanWebpackPlugin = require("clean-webpack-plugin");
const ExtractTextPlugin = require("extract-text-webpack-plugin");
const HtmlWebpackPlugin = require("html-webpack-plugin");

const DEBUG = process.env.NODE_ENV !== "production";
const JS_NAME = DEBUG ? "index.js" : "[chunkhash:10].js";
const CSS_NAME = DEBUG ? "index.css" : "[contenthash:10].css";

module.exports = {
  entry: "./ts/index/index.tsx",
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
        test: /\.css$/,
        use: ExtractTextPlugin.extract({use: "css-loader"}),
        exclude: /node_modules/,
      },
    ],
  },
  plugins: [
    new LiveReloadPlugin(),
    new CleanWebpackPlugin(["dist"]),
    new ExtractTextPlugin(`static/${CSS_NAME}`),
    new HtmlWebpackPlugin({
      title: "K-pop idols network | Profiles, images and face recognition",
      favicon: "ts/index/favicon.ico",
    }),
  ],
  output: {
    path: path.resolve(__dirname, "dist"),
    filename: `static/${JS_NAME}`,
    publicPath: "/",
  },
};
