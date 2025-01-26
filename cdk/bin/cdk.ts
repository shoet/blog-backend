#!/usr/bin/env node
import * as cdk from "aws-cdk-lib";
import { BlogAppStack } from "../lib/blog-app-stack";

const app = new cdk.App();
new BlogAppStack(app, "BlogAppStack", {});
