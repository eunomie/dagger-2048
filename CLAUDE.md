# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build and Development

Always use the available MCP tools to build and test the Go code. Do not use local go commands. 

To create a local binary, run the following command:

`dagger -c '. --platform=current | binary | export dagger2048'`

Then you have access to the `dagger2048` binary file.
