/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type RequestFilter = {
  start_timestamp?: string
  end_timestamp?: string
}

export type GetAggregatedSummaryRequest = {
  tenant_id?: string
  start_time?: string
  end_time?: string
}

export type Summary = {
  method?: string
  total_requests?: string
  success_requests?: string
  failure_requests?: string
  average_latency?: number
}

export type AggregatedSummary = {
  summary?: Summary
  method_summaries?: Summary[]
}

export type GetUsageDataRequestFilter = {
  api_keys?: string[]
  models?: string[]
}

export type GetUsageDataRequest = {
  start_time?: string
  end_time?: string
  filter?: GetUsageDataRequestFilter
  after?: string
  limit?: number
}

export type UsageData = {
  data_points?: UsageDataPoint[]
}

export type UsageDataPoint = {
  user_id?: string
  organization?: string
  project?: string
  api_key_id?: string
  api_method?: string
  status_code?: number
  timestamp?: string
  latency_ms?: number
  model_id?: string
  time_to_first_token_ms?: number
  prompt_tokens?: number
  completion_tokens?: number
}

export type ListUsageDataRequest = {
  start_time?: string
  end_time?: string
}

export type UsageDataByGroup = {
  user_id?: string
  api_key_id?: string
  api_key_name?: string
  model_id?: string
  total_requests?: string
  avg_latency_ms?: number
  avg_time_to_first_token_ms?: number
  total_prompt_tokens?: string
  total_completion_tokens?: string
}

export type ListUsageDataResponse = {
  usages?: UsageDataByGroup[]
}

export type ListModelUsageSummariesRequest = {
  filter?: RequestFilter
}

export type ListModelUsageSummariesResponseValue = {
  model_id?: string
  total_requests?: string
}

export type ListModelUsageSummariesResponseDatapoint = {
  timestamp?: string
  values?: ListModelUsageSummariesResponseValue[]
}

export type ListModelUsageSummariesResponse = {
  datapoints?: ListModelUsageSummariesResponseDatapoint[]
}

export class APIUsageService {
  static GetAggregatedSummary(req: GetAggregatedSummaryRequest, initReq?: fm.InitReq): Promise<AggregatedSummary> {
    return fm.fetchReq<GetAggregatedSummaryRequest, AggregatedSummary>(`/llmariner.apiusage.server.v1.APIUsageService/GetAggregatedSummary`, {...initReq, method: "POST", body: JSON.stringify(req)})
  }
  static GetUsageData(req: GetUsageDataRequest, initReq?: fm.InitReq): Promise<UsageData> {
    return fm.fetchReq<GetUsageDataRequest, UsageData>(`/llmariner.apiusage.server.v1.APIUsageService/GetUsageData`, {...initReq, method: "POST", body: JSON.stringify(req)})
  }
  static ListUsageData(req: ListUsageDataRequest, initReq?: fm.InitReq): Promise<ListUsageDataResponse> {
    return fm.fetchReq<ListUsageDataRequest, ListUsageDataResponse>(`/v1/api_usages?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static ListModelUsageSummaries(req: ListModelUsageSummariesRequest, initReq?: fm.InitReq): Promise<ListModelUsageSummariesResponse> {
    return fm.fetchReq<ListModelUsageSummariesRequest, ListModelUsageSummariesResponse>(`/v1/api-usage/model-usage-summaries?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
}