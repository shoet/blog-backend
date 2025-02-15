#!/usr/bin/env node
import * as cdk from "aws-cdk-lib";
import { BlogAppStack } from "../lib/blog-app-stack";

declare module "aws-cdk-lib" {
  interface StackProps {
    stage: string;
    commitHash?: string;
  }
}

const app = new cdk.App();

const stage = process.env.STAGE;
if (!stage) {
  throw new Error("STAGE is required");
}
if (!["dev", "prod"].includes(stage)) {
  throw new Error("STAGE must be either dev or prod");
}

const commitHash = process.env.COMMIT_HASH;

console.log(`Deploying to stage: ${stage}`);

new BlogAppStack(app, `BlogAppStack-${stage}`, {
  stage: stage,
  commitHash: commitHash,
});
