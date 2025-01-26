#!/usr/bin/env node
import * as cdk from "aws-cdk-lib";
import { BlogAppStack } from "../lib/blog-app-stack";

declare module "aws-cdk-lib" {
  interface StackProps {
    stage: string;
  }
}

const app = new cdk.App();
new BlogAppStack(app, "BlogAppStack", {});
