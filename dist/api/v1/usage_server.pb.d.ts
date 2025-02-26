import * as fm from "../../fetch.pb";
export type GetAggregatedSummaryRequest = {
    tenant_id?: string;
    start_time?: string;
    end_time?: string;
};
export type Summary = {
    method?: string;
    total_requests?: string;
    success_requests?: string;
    failure_requests?: string;
    average_latency?: number;
};
export type AggregatedSummary = {
    summary?: Summary;
    method_summaries?: Summary[];
};
export type GetUsageDataRequestFilter = {
    api_keys?: string[];
    models?: string[];
};
export type GetUsageDataRequest = {
    start_time?: string;
    end_time?: string;
    filter?: GetUsageDataRequestFilter;
    after?: string;
    limit?: number;
};
export type UsageData = {
    data_points?: UsageDataPoint[];
};
export type UsageDataPoint = {
    user_id?: string;
    organization?: string;
    project?: string;
    api_key_id?: string;
    api_method?: string;
    status_code?: number;
    timestamp?: string;
    latency_ms?: number;
    model_id?: string;
    time_to_first_token_ms?: number;
    prompt_tokens?: number;
    completion_tokens?: number;
};
export type ListUsageDataRequest = {
    start_time?: string;
    end_time?: string;
};
export type UsageDataByGroup = {
    user_id?: string;
    api_key_id?: string;
    api_key_name?: string;
    model_id?: string;
    total_requests?: string;
    avg_latency_ms?: number;
    avg_time_to_first_token_ms?: number;
    total_prompt_tokens?: string;
    total_completion_tokens?: string;
};
export type ListUsageDataResponse = {
    usages?: UsageDataByGroup[];
};
export declare class APIUsageService {
    static GetAggregatedSummary(req: GetAggregatedSummaryRequest, initReq?: fm.InitReq): Promise<AggregatedSummary>;
    static GetUsageData(req: GetUsageDataRequest, initReq?: fm.InitReq): Promise<UsageData>;
    static ListUsageData(req: ListUsageDataRequest, initReq?: fm.InitReq): Promise<ListUsageDataResponse>;
}
