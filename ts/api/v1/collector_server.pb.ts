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
  user_id?: string
  tenant?: string
  organization?: string
  project?: string
  api_key_id?: string
  api_method?: string
  status_code?: number
  timestamp?: string
  latency_ms?: number
  details?: UsageDetails
}


type BaseUsageDetails = {
}

export type UsageDetails = BaseUsageDetails
  & OneOf<{ create_chat_completion: CreateChatCompletion; create_completion: CreateCompletion; create_audio_transcription: CreateAudioTranscription }>

export type CreateChatCompletion = {
  model_id?: string
  time_to_first_token_ms?: number
  prompt_tokens?: number
  completion_tokens?: number
}

export type CreateCompletion = {
  model_id?: string
  time_to_first_token_ms?: number
  prompt_tokens?: number
  completion_tokens?: number
}

export type CreateAudioTranscription = {
  model_id?: string
  time_to_first_token_ms?: number
  input_tokens?: number
  output_tokens?: number
  total_tokens?: number
  text_tokens?: number
  audio_tokens?: number
  input_duration_seconds?: number
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