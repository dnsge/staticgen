# staticgen
staticgen is a tool that generates file bundles and watches directories for changes.

## Configuration
Project configurations are specified in the `staticgen.json` file and specify which files to include, where to place output files, and generation scripts.

Here is an example configuration:
```json
{
  "output": "./dist",
  "clean": true,
  "include": [
    "./src/index.html",
    "./src/js",
    "./src/images"
  ],
  "scripts": [
    {
      "name": "Compile SASS",
      "build": ["sass", "--no-source-map", "./src/scss/:./dist/css"],
      "watch": ["sass", "-c", "--watch", "./src/scss/:./dist/css"]
    }
  ]
}
```
 - `output` specifies the output directory for included files
 - `clean` specifies whether to clean the output directory before running
 - `include` lists files and directories that will be copied to the output directory
 - `scripts` specifies scripts (command-line programs) to run on build and watch

When running `./staticgen`, the files and directories from `include` are copied into `output` and the `build` portion of each script is run. 

If you want to automatically update the files (as is common during development), run `./staticgen --watch`. This flag enables watching of everything specified in `include` and will run each script's `watch` command in parallel.

## Inspiration
This project was inspired by my want to streamline the development of a static HTML website without dealing with the complexity of webpack, among other things.