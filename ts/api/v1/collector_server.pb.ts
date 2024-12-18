/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"

type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };
type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (
    keyof T extends infer K ?
      (K extends string & keyof T ? { [k in K]: T[K] } & Absent<T, K>
        : never)
    : never);
export type UsageRecord = {
  userId?: string
  tenant?: string
  organization?: string
  project?: string
  apiKeyId?: string
  apiMethod?: string
  statusCode?: number
  timestamp?: string
  latencyMs?: number
  details?: UsageDetails
}


type BaseUsageDetails = {
}

export type UsageDetails = BaseUsageDetails
  & OneOf<{ createChatCompletion: CreateChatCompletion; createCompletion: CreateCompletion }>

export type CreateChatCompletion = {
  modelId?: string
  timeToFirstTokenMs?: number
  promptTokens?: number
  completionTokens?: number
}

export type CreateCompletion = {
  modelId?: string
  timeToFirstTokenMs?: number
  promptTokens?: number
  completionTokens?: number
}

export type CreateUsageRequest = {
  records?: UsageRecord[]
}

export type Usage = {
}

export class CollectionInternalService {
  static CreateUsage(req: CreateUsageRequest, initReq?: fm.InitReq): Promise<Usage> {
    return fm.fetchReq<CreateUsageRequest, Usage>(`/llmariner.apiusage.server.v1.CollectionInternalService/CreateUsage`, {...initReq, method: "POST", body: JSON.stringify(req)})
  }
}