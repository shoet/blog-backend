#!/usr/bin/env node
import * as cdk from "aws-cdk-lib";
import { BlogBackendAppStack } from "../lib/blog-backend-app-stack";
import { BlogCDNStack } from "../lib/blog-cdk-stack";

declare module "aws-cdk-lib" {
  interface StackProps {
    stage: string;
    commitHash?: string;
  }
}

const app = new cdk.App();

const stage = app.node.tryGetContext("STAGE");
if (!stage) {
  throw new Error("STAGE is required");
}
if (!["dev", "prod"].includes(stage)) {
  throw new Error("STAGE must be either dev or prod");
}

const commitHash = app.node.tryGetContext("COMMIT_HASH");

type StackType = "app" | "cdn";
const stackType: StackType = app.node.tryGetContext("TYPE") || "app";

console.log(`Deploying to stage: ${stage}`);

console.log(`Deploying [${stackType}] stack`);
switch (stackType) {
  case "app":
    new BlogBackendAppStack(app, `BlogBackendAppStack-${stage}`, {
      stage: stage,
      commitHash: commitHash,
    });
    break;
  case "cdn":
    new BlogCDNStack(app, `BlogCDNStack-${stage}`, {
      stage: stage,
      commitHash: commitHash,
    });
    break;
  default:
    console.log(`Not found stack type: ${stackType}`);
    break;
}
